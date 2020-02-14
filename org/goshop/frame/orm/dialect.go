package orm

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func wrapsavesql(dbType DBTYPE, entity IBaseEntity, columns []reflect.StructField, values []interface{}) (string, error) {
	if dbType == DBType_MySQL {
		return wrapsavesql_mysql(entity, columns, values)
	}
	return "", errors.New("不支持的数据库")
}

//mysql的save语句拼接
func wrapsavesql_mysql(entity IBaseEntity, columns []reflect.StructField, values []interface{}) (string, error) {

	//SQL语句的构造器
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("INSERT INTO ")
	sqlBuilder.WriteString(entity.GetTableName())
	sqlBuilder.WriteString("(")

	//SQL语句中,VALUES(?,?,...)语句的构造器
	var valueSQLBuilder strings.Builder
	valueSQLBuilder.WriteString(" VALUES (")

	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	for i := 0; i < len(columns); i++ {
		field := columns[i]
		if strings.EqualFold(field.Name, entity.GetPkName()) { //如果是主键
			pkKind := field.Type.Kind()

			if !(pkKind == reflect.String || pkKind == reflect.Int) { //只支持字符串和int类型的主键
				return "", errors.New("不支持的主键类型")
			}

			pkValue := values[i]
			if (pkKind == reflect.String) && (pkValue.(string) == "") { //主键是字符串类型,并且值为"",赋值id
				values[i] = id
				//如果是数字类型,并且值为0,需要从数组中删除掉主键的信息,让数据库自己生成
			} else if (pkKind == reflect.Int) && (pkValue.(int) == 0) {
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

	fmt.Println(rebind(sqlstr))

	return sqlstr, nil

}

func rebind(query string) string {
	//switch bindType {
	//case QUESTION, UNKNOWN:
	//	return query
	//}

	// Add space enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(query)+10)

	var i, j int

	for i = strings.Index(query, "?"); i != -1; i = strings.Index(query, "?") {
		rqb = append(rqb, query[:i]...)

		//switch bindType {
		//case DOLLAR:
		rqb = append(rqb, '$')
		//case NAMED:
		//	rqb = append(rqb, ':', 'a', 'r', 'g')
		//case AT:
		//	rqb = append(rqb, '@', 'p')
		//}

		j++
		rqb = strconv.AppendInt(rqb, int64(j), 10)

		query = query[i+1:]
	}

	return string(append(rqb, query...))
}
