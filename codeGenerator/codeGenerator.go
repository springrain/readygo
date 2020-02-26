package main

import (
	"fmt"
	"os"
	"readygo/orm"
	"strings"
	"text/template"
)

var baseDao *orm.BaseDao

const (
	dbName             = "readygo"
	packageName        = "permstruct"
	servicePackageName = "permservice"
)

func init() {
	baseDaoConfig := orm.DataSourceConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   dbName,
		UserName: "root",
		PassWord: "root",
		DBType:   "mysql",
	}

	baseDao, _ = orm.NewBaseDao(&baseDaoConfig)
}

func main() {
	//code("t_user")

	tableNames := selectAllTable()
	for _, tableName := range tableNames {
		code(tableName)
	}

}

//生成代码
func code(tableName string) {

	info := selectTableColumn(tableName)

	structFileName := "./code/struct/" + info["structName"].(string) + ".go"
	serviceFileName := "./code/service/" + info["structName"].(string) + "Service.go"
	structFile, _ := os.Create(structFileName)
	serviceFile, _ := os.Create(serviceFileName)

	//w := bufio.NewWriter(f) // 创建新的 Writer 对象
	defer func() {
		structFile.Close()
		serviceFile.Close()
	}()

	structTemplate, err1 := template.ParseFiles("./templates/struct.txt")
	if err1 != nil {
		fmt.Println(err1)
	}
	structTemplate.Execute(structFile, info)

	serviceTemplate, err2 := template.ParseFiles("./templates/service.txt")
	if err2 != nil {
		fmt.Println(err2)
	}
	serviceTemplate.Execute(serviceFile, info)

}

//获取所有的表名
func selectAllTable() []string {
	finder := orm.NewFinder()
	finder.Append("select table_name from information_schema.TABLES where  TABLE_SCHEMA =?", dbName)
	tableNames := []string{}
	orm.QueryStructList(nil, finder, &tableNames, nil)
	return tableNames
}

//根据表名查询字段信息和主键名称
func selectTableColumn(tableName string) map[string]interface{} {
	tableComment := ""
	finder := orm.NewFinder()
	finder.Append("select table_comment from information_schema.TABLES where  TABLE_SCHEMA =? and TABLE_Name=? ", dbName, tableName)
	orm.QueryStruct(nil, finder, &tableComment)

	finder2 := orm.NewFinder()
	// select * from information_schema.COLUMNS where table_schema ='readygo' and table_name='t_user';
	finder2.Append("select COLUMN_NAME,DATA_TYPE,IS_NULLABLE,COLUMN_COMMENT from information_schema.COLUMNS where  TABLE_SCHEMA =? and TABLE_NAME=? and COLUMN_NAME not like ?  order by ORDINAL_POSITION asc", dbName, tableName, "bak%")

	maps, _ := orm.QueryMapList(nil, finder2, nil)

	for _, m := range maps {
		dataType := m["DATA_TYPE"].(string)
		dataType = strings.ToUpper(dataType)

		nullable := m["IS_NULLABLE"].(string)
		nullable = strings.ToUpper(nullable)

		if dataType == "VARCHAR" || dataType == "NVARCHAR" || dataType == "TEXT" {
			if nullable == "YES" {
				dataType = "sql.NullString"
			} else {
				dataType = "string"
			}

		} else if dataType == "DATETIME" || dataType == "TIMESTAMP" {
			if nullable == "YES" {
				dataType = "sql.NullTime"
			} else {
				dataType = "time.Time"
			}

		} else if dataType == "INT" {
			if nullable == "YES" {
				dataType = "sql.NullInt32"
			} else {
				dataType = "int"
			}

		} else if dataType == "BIGINT" {
			if nullable == "YES" {
				dataType = "sql.NullInt64"
			} else {
				dataType = "int64"
			}

		} else if dataType == "FLOAT" {
			if nullable == "YES" {
				dataType = "sql.NullFloat64"
			} else {
				dataType = "float32"
			}

		} else if dataType == "DOUBLE" {
			if nullable == "YES" {
				dataType = "sql.NullFloat64"
			} else {
				dataType = "float64"
			}

		}
		m["DATA_TYPE"] = dataType
		m["field"] = camelCaseName(m["COLUMN_NAME"].(string))
	}

	finderPK := orm.NewFinder()
	finderPK.Append("SELECT column_name FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA=? and  table_name=? AND constraint_name=?", dbName, tableName, "PRIMARY")
	pkName := ""
	orm.QueryStruct(nil, finderPK, &pkName)
	info := make(map[string]interface{})
	info["columns"] = maps
	info["pkName"] = pkName
	info["tableName"] = tableName
	structName := tableName
	structName = strings.Replace(structName, "t_", "", 1)
	structName = camelCaseName(structName) + "Struct"
	info["structName"] = structName
	info["pname"] = firstToLower(structName)
	info["packageName"] = packageName
	info["tableComment"] = tableComment
	info["servicePackageName"] = servicePackageName
	return info
}

//首字母大写
func firstToUpper(str string) string {
	str = strings.ToUpper(string(str[0:1])) + string(str[1:])
	return str
}

//首字母小写
func firstToLower(str string) string {
	str = strings.ToLower(string(str[0:1])) + string(str[1:])
	return str
}

//驼峰
func camelCaseName(name string) string {
	names := strings.Split(name, "_")
	structName := ""
	for _, name := range names {
		structName = structName + firstToUpper(name)
	}

	return structName

}
