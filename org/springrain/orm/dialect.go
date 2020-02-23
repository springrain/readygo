package orm

import (
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//数据库连接字符串
func wrapDBDSN(config *DataSourceConfig) (string, error) {
	if config == nil {
		return "", nil
	}
	if config.DBType == DBType_MYSQL {
		//username:password@tcp(127.0.0.1:3306)/dbName
		dsn := config.UserName + ":" + config.PassWord + "@tcp(" + config.Host + ":" + strconv.Itoa(config.Port) + ")/" + config.DBName + "?charset=utf8&loc=Asia%2FShanghai&parseTime=true"
		return dsn, nil
	}

	return "", errors.New("不支持的数据库")
}

//包装基础的SQL语句
func wrapSQL(dbType DBTYPE, sqlstr string) (string, error) {
	if dbType == DBType_MYSQL || dbType == DBType_UNKNOWN {
		return sqlstr, nil
	}
	//根据数据库类型,调整SQL变量符号,例如?,? $1,$2这样的
	sqlstr = rebind(dbType, sqlstr)
	return sqlstr, nil
}

//包装分页的SQL语句
func wrapPageSQL(dbType DBTYPE, sqlstr string, page *Page) (string, error) {

	var sqlbuilder strings.Builder
	sqlbuilder.WriteString(sqlstr)
	if dbType == DBType_MYSQL || dbType == DBType_UNKNOWN { //MySQL数据库
		sqlbuilder.WriteString(" limit ")
		sqlbuilder.WriteString(strconv.Itoa(page.PageSize * (page.PageNo - 1)))
		sqlbuilder.WriteString(",")
		sqlbuilder.WriteString(strconv.Itoa(page.PageSize))

	} else if dbType == DBType_POSTGRESQL { //postgresql
		sqlbuilder.WriteString(" limit ")
		sqlbuilder.WriteString(strconv.Itoa(page.PageSize))
		sqlbuilder.WriteString(" offset ")
		sqlbuilder.WriteString(strconv.Itoa(page.PageSize * (page.PageNo - 1)))
	} else if dbType == DBType_MSSQL { //mssql
		//先不写啦
		//bug(chunanyong) 还需要其他的数据库分页语句
	}
	sqlstr = sqlbuilder.String()
	return wrapSQL(dbType, sqlstr)
}

//包装保存Struct语句.返回语句,是否自增,错误信息
//数组传递,如果外部方法有调用append的逻辑,传递指针,因为append会破坏指针引用
func wrapSaveStructSQL(dbType DBTYPE, entity IEntityStruct, columns *[]reflect.StructField, values *[]interface{}) (string, bool, error) {

	//是否自增,默认false
	autoIncrement := false
	//SQL语句的构造器
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("INSERT INTO ")
	sqlBuilder.WriteString(entity.GetTableName())
	sqlBuilder.WriteString("(")

	//SQL语句中,VALUES(?,?,...)语句的构造器
	var valueSQLBuilder strings.Builder
	valueSQLBuilder.WriteString(" VALUES (")

	for i := 0; i < len(*columns); i++ {
		field := (*columns)[i]
		fieldName, e := entityPKFieldName(entity)
		if e != nil {
			return "", autoIncrement, e
		}
		if field.Name == fieldName { //如果是主键
			pkKind := field.Type.Kind()

			if !(pkKind == reflect.String || pkKind == reflect.Int) { //只支持字符串和int类型的主键
				return "", autoIncrement, errors.New("不支持的主键类型")
			}
			//主键的值
			pkValue := (*values)[i]
			if len(entity.GetPkSequence()) > 0 { //如果是主键序列
				//拼接字符串
				sqlBuilder.WriteString(field.Tag.Get(tagColumnName))
				sqlBuilder.WriteString(",")
				valueSQLBuilder.WriteString(entity.GetPkSequence())
				valueSQLBuilder.WriteString(",")
				//去掉这一列,后续不再处理
				*columns = append((*columns)[:i], (*columns)[i+1:]...)
				*values = append((*values)[:i], (*values)[i+1:]...)
				i = i - 1
				continue

			} else if (pkKind == reflect.String) && (pkValue.(string) == "") { //主键是字符串类型,并且值为"",赋值id
				id := generateStringID()
				(*values)[i] = id
				//给对象主键赋值
				v := reflect.ValueOf(entity).Elem()
				v.FieldByName(field.Name).Set(reflect.ValueOf(id))
				//如果是数字类型,并且值为0,认为是数据库自增,从数组中删除掉主键的信息,让数据库自己生成
			} else if (pkKind == reflect.Int) && (pkValue.(int) == 0) {
				autoIncrement = true
				//去掉这一列,后续不再处理
				*columns = append((*columns)[:i], (*columns)[i+1:]...)
				*values = append((*values)[:i], (*values)[i+1:]...)
				i = i - 1
				continue
			}
		}
		//拼接字符串
		sqlBuilder.WriteString(field.Tag.Get(tagColumnName))
		sqlBuilder.WriteString(",")
		valueSQLBuilder.WriteString("?,")

	}
	//去掉字符串最后的 , 号
	sqlstr := sqlBuilder.String()
	if len(sqlstr) > 0 {
		sqlstr = sqlstr[:len(sqlstr)-1]
	}
	valuestr := valueSQLBuilder.String()
	if len(valuestr) > 0 {
		valuestr = valuestr[:len(valuestr)-1]
	}
	sqlstr = sqlstr + ")" + valuestr + ")"
	savesql, err := wrapSQL(dbType, sqlstr)
	return savesql, autoIncrement, err

}

//包装更新Struct语句
func wrapUpdateStructSQL(dbType DBTYPE, entity IEntityStruct, columns []reflect.StructField, values []interface{}, onlyupdatenotnull bool) (string, error) {

	//SQL语句的构造器
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("UPDATE ")
	sqlBuilder.WriteString(entity.GetTableName())
	sqlBuilder.WriteString(" SET ")

	//主键的值
	var pkValue interface{}

	for i := 0; i < len(columns); i++ {
		field := columns[i]

		fieldName, e := entityPKFieldName(entity)
		if e != nil {
			return "", e
		}

		if field.Name == fieldName { //如果是主键
			pkValue = values[i]
			//去掉这一列,最后处理主键
			columns = append(columns[:i], columns[i+1:]...)
			values = append(values[:i], values[i+1:]...)
			i = i - 1
			continue
		}

		//只更新不为nil的字段
		if onlyupdatenotnull && (values[i] == nil) {
			//去掉这一列,不再处理
			columns = append(columns[:i], columns[i+1:]...)
			values = append(values[:i], values[i+1:]...)
			i = i - 1
			continue

		}

		sqlBuilder.WriteString(field.Tag.Get(tagColumnName))
		sqlBuilder.WriteString("=?,")

	}
	//主键的值是最后一个
	values = append(values, pkValue)
	//去掉字符串最后的 , 号
	sqlstr := sqlBuilder.String()
	sqlstr = sqlstr[:len(sqlstr)-1]

	sqlstr = sqlstr + " WHERE " + entity.GetPKColumnName() + "=?"

	return wrapSQL(dbType, sqlstr)
}

//包装删除Struct语句
func wrapDeleteStructSQL(dbType DBTYPE, entity IEntityStruct) (string, error) {

	//SQL语句的构造器
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("DELETE FROM ")
	sqlBuilder.WriteString(entity.GetTableName())
	sqlBuilder.WriteString(" WHERE ")
	sqlBuilder.WriteString(entity.GetPKColumnName())
	sqlBuilder.WriteString("=?")
	sqlstr := sqlBuilder.String()

	return wrapSQL(dbType, sqlstr)

}

//包装保存Map语句,Map因为没有字段属性,无法完成Id的类型判断和赋值,需要确保Map的值是完整的.
func wrapSaveMapSQL(dbType DBTYPE, entity IEntityMap) (string, []interface{}, error) {

	dbFieldMap := entity.GetDBFieldMap()
	if len(dbFieldMap) < 1 {
		return "", nil, errors.New("GetDBFieldMap()返回值不能为空")
	}
	//SQL对应的参数
	values := []interface{}{}

	//SQL语句的构造器
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("INSERT INTO ")
	sqlBuilder.WriteString(entity.GetTableName())
	sqlBuilder.WriteString("(")

	//SQL语句中,VALUES(?,?,...)语句的构造器
	var valueSQLBuilder strings.Builder
	valueSQLBuilder.WriteString(" VALUES (")

	for k, v := range dbFieldMap {
		//拼接字符串
		sqlBuilder.WriteString(k)
		sqlBuilder.WriteString(",")
		valueSQLBuilder.WriteString("?,")
		values = append(values, v)
	}
	//去掉字符串最后的 , 号
	sqlstr := sqlBuilder.String()
	if len(sqlstr) > 0 {
		sqlstr = sqlstr[:len(sqlstr)-1]
	}
	valuestr := valueSQLBuilder.String()
	if len(valuestr) > 0 {
		valuestr = valuestr[:len(valuestr)-1]
	}
	sqlstr = sqlstr + ")" + valuestr + ")"

	var e error
	sqlstr, e = wrapSQL(dbType, sqlstr)
	if e != nil {
		return "", nil, e
	}
	return sqlstr, values, nil
}

//包装Map更新语句,Map因为没有字段属性,无法完成Id的类型判断和赋值,需要确保Map的值是完整的.
func wrapUpdateMapSQL(dbType DBTYPE, entity IEntityMap) (string, []interface{}, error) {
	dbFieldMap := entity.GetDBFieldMap()
	if len(dbFieldMap) < 1 {
		return "", nil, errors.New("GetDBFieldMap()返回值不能为空")
	}
	//SQL语句的构造器
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("UPDATE ")
	sqlBuilder.WriteString(entity.GetTableName())
	sqlBuilder.WriteString(" SET ")

	//SQL对应的参数
	values := []interface{}{}
	//主键名称
	var pkValue interface{}

	for k, v := range dbFieldMap {

		if k == entity.GetPKColumnName() { //如果是主键
			pkValue = v
			continue
		}

		//拼接字符串
		sqlBuilder.WriteString(k)
		sqlBuilder.WriteString("=?,")
		values = append(values, v)
	}
	//主键的值是最后一个
	values = append(values, pkValue)
	//去掉字符串最后的 , 号
	sqlstr := sqlBuilder.String()
	sqlstr = sqlstr[:len(sqlstr)-1]

	sqlstr = sqlstr + " WHERE " + entity.GetPKColumnName() + "=?"

	var e error
	sqlstr, e = wrapSQL(dbType, sqlstr)
	if e != nil {
		return "", nil, e
	}
	return sqlstr, values, nil
}

//封装查询语句
func wrapQuerySQL(dbType DBTYPE, finder *Finder, page *Page) (string, error) {

	//获取到没有page的sql的语句
	sqlstr, err := finder.GetSQL()
	if err != nil {
		return "", err
	}
	if page == nil {
		sqlstr, err = wrapSQL(dbType, sqlstr)
	} else {
		sqlstr, err = wrapPageSQL(dbType, sqlstr, page)
	}

	if err != nil {
		return "", err
	}
	return sqlstr, err
}

//根据数据库类型,调整SQL变量符号,例如?,? $1,$2这样的
func rebind(dbType DBTYPE, query string) string {

	// Add space enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(query)+10)

	var i, j int

	for i = strings.Index(query, "?"); i != -1; i = strings.Index(query, "?") {
		rqb = append(rqb, query[:i]...)

		if dbType == DBType_POSTGRESQL {
			rqb = append(rqb, '$')
		} else if dbType == DBType_MSSQL {
			rqb = append(rqb, '@', 'p')
		}
		j++
		rqb = strconv.AppendInt(rqb, int64(j), 10)

		query = query[i+1:]
	}

	return string(append(rqb, query...))
}

//查询order by在sql中出现的开始位置和结束位置
var orderByExpr = "\\s+(order)\\s+(by)+\\s"
var orderByRegexp, _ = regexp.Compile(orderByExpr)

//查询order by在sql中出现的开始位置和结束位置
func findOrderByIndex(strsql string) []int {
	loc := orderByRegexp.FindStringIndex(strings.ToLower(strsql))
	return loc
}

//查询group by在sql中出现的开始位置和结束位置
var groupByExpr = "\\s+(group)\\s+(by)+\\s"
var groupByRegexp, _ = regexp.Compile(groupByExpr)

//查询group by在sql中出现的开始位置和结束位置
func findGroupByIndex(strsql string) []int {
	loc := groupByRegexp.FindStringIndex(strings.ToLower(strsql))
	return loc
}

//查询 from 在sql中出现的开始位置和结束位置
var fromExpr = "\\s+(from)+\\s"
var fromRegexp, _ = regexp.Compile(fromExpr)

//查询from在sql中出现的开始位置和结束位置
func findFromIndex(strsql string) []int {
	loc := fromRegexp.FindStringIndex(strings.ToLower(strsql))
	return loc
}

//生成主键字符串
func generateStringID() string {
	pk := strconv.FormatInt(time.Now().UnixNano(), 10)
	return pk
}

func converValueColumnType(v interface{}, columnType *sql.ColumnType) interface{} {

	if v == nil {
		return nil
	}

	//如果是字节数组
	value, ok := v.([]byte)
	if !ok { //转化失败
		return v
	}
	if len(value) < 1 { //值为空,为nil
		return nil
	}

	//获取数据库类型,自己对应golang的基础类型
	databaseTypeName := columnType.DatabaseTypeName()
	//如果是字符串
	if databaseTypeName == "VARCHAR" || databaseTypeName == "NVARCHAR" || databaseTypeName == "TEXT" {
		return String(v)
	} else if databaseTypeName == "INT" { //如果是INT
		return Int(v)
	} else if databaseTypeName == "BIGINT" { //如果是BIGINT
		return Int64(v)
	} else if databaseTypeName == "FLOAT" { //如果是FLOAT
		return Float32(v)
	} else if databaseTypeName == "DOUBLE" { //如果是DOUBLE
		return Float64(v)
	} else if databaseTypeName == "DATETIME" || databaseTypeName == "TIMESTAMP" { //如果是DATETIME
		return Time(v, "2006-01-02 15:04:05", time.Local)
	}
	//其他类型以后再写.....

	return nil
}
