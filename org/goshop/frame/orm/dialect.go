package orm

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//包装基础的SQL语句
func wrapSQL(dbType DBTYPE, sqlstr string) (string, error) {
	if dbType == DBType_MYSQL || dbType == DBType_UNKNOWN {
		return sqlstr, nil
	}
	//根据数据库类型,调整SQL变量符号,例如?,? $1,$2这样的
	sqlstr = rebind(dbType, sqlstr)
	return sqlstr, nil
}

//包装保存Struct语句
func wrapSaveStructSQL(dbType DBTYPE, entity IEntityStruct, columns []reflect.StructField, values []interface{}) (string, error) {

	//SQL语句的构造器
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("INSERT INTO ")
	sqlBuilder.WriteString(entity.GetTableName())
	sqlBuilder.WriteString("(")

	//SQL语句中,VALUES(?,?,...)语句的构造器
	var valueSQLBuilder strings.Builder
	valueSQLBuilder.WriteString(" VALUES (")

	for i := 0; i < len(columns); i++ {
		field := columns[i]
		if field.Name == entityPKFieldName(entity) { //如果是主键
			pkKind := field.Type.Kind()

			if !(pkKind == reflect.String || pkKind == reflect.Int) { //只支持字符串和int类型的主键
				return "", errors.New("不支持的主键类型")
			}
			//主键的值
			pkValue := values[i]
			if len(entity.GetPkSequence()) > 0 { //如果是主键序列
				//拼接字符串
				sqlBuilder.WriteString(field.Tag.Get(tagColumnName))
				sqlBuilder.WriteString(",")
				valueSQLBuilder.WriteString(entity.GetPkSequence())
				valueSQLBuilder.WriteString(",")
				//去掉这一列,后续不再处理
				columns = append(columns[:i], columns[i+1:]...)
				values = append(values[:i], values[i+1:]...)
				i = i - 1
				continue

			} else if (pkKind == reflect.String) && (pkValue.(string) == "") { //主键是字符串类型,并且值为"",赋值id
				id := strconv.FormatInt(time.Now().UnixNano(), 10)
				values[i] = id
				//给对象主键赋值
				v := reflect.ValueOf(entity).Elem()
				v.FieldByName(field.Name).Set(reflect.ValueOf(id))
				//如果是数字类型,并且值为0,需要从数组中删除掉主键的信息,让数据库自己生成
			} else if (pkKind == reflect.Int) && (pkValue.(int) == 0) {
				//去掉这一列,后续不再处理
				columns = append(columns[:i], columns[i+1:]...)
				values = append(values[:i], values[i+1:]...)
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

	if dbType == DBType_MYSQL || dbType == DBType_UNKNOWN {
		return sqlstr, nil
	}
	//根据数据库类型,调整SQL变量符号,例如?,? $1,$2这样的
	sqlstr = rebind(dbType, sqlstr)
	return sqlstr, nil

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
		if field.Name == entityPKFieldName(entity) { //如果是主键
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

	if dbType == DBType_MYSQL || dbType == DBType_UNKNOWN {
		return sqlstr, nil
	}
	//根据数据库类型,调整SQL变量符号,例如?,? $1,$2这样的
	sqlstr = rebind(dbType, sqlstr)
	return sqlstr, nil
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

	if dbType == DBType_MYSQL || dbType == DBType_UNKNOWN {
		return sqlstr, nil
	}
	//根据数据库类型,调整SQL变量符号,例如?,? $1,$2这样的
	sqlstr = rebind(dbType, sqlstr)
	return sqlstr, nil

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

	if dbType == DBType_MYSQL || dbType == DBType_UNKNOWN {
		return sqlstr, values, nil
	}
	//根据数据库类型,调整SQL变量符号,例如?,? $1,$2这样的
	sqlstr = rebind(dbType, sqlstr)
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

	if dbType == DBType_MYSQL || dbType == DBType_UNKNOWN {
		return sqlstr, values, nil
	}
	//根据数据库类型,调整SQL变量符号,例如?,? $1,$2这样的
	sqlstr = rebind(dbType, sqlstr)
	return sqlstr, values, nil
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
