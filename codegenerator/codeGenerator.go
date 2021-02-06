package codegenerator

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"gitee.com/chunanyong/zorm"
	_ "github.com/go-sql-driver/mysql"
)

var dbDao *zorm.DBDao

const (
	dbName             = "readygo"
	packageName        = "permstruct"
	servicePackageName = "permservice"
)

func init() {
	dbDaoConfig := zorm.DataSourceConfig{
		DSN:        "root:root@tcp(127.0.0.1:3306)/readygo?charset=utf8&parseTime=true",
		DriverName: "mysql",
		PrintSQL:   true,
	}

	dbDao, _ = zorm.NewDBDao(&dbDaoConfig)
}

//生成代码
func code(tableName string) {
	ctx := context.Background()

	info := selectTableColumn(ctx, tableName)

	//创建目录
	os.MkdirAll("./code/struct", os.ModePerm)
	os.MkdirAll("./code/service", os.ModePerm)

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
	finder := zorm.NewFinder()
	finder.Append("select table_name from information_schema.TABLES where  TABLE_SCHEMA =?", dbName)
	tableNames := []string{}
	zorm.Query(nil, finder, &tableNames, nil)
	return tableNames
}

//根据表名查询字段信息和主键名称
func selectTableColumn(ctx context.Context, tableName string) map[string]interface{} {

	info := make(map[string]interface{})

	tableComment := ""
	finder := zorm.NewFinder()
	finder.Append("select table_comment from information_schema.TABLES where  TABLE_SCHEMA =? and TABLE_Name=? ", dbName, tableName)
	zorm.QueryRow(ctx, finder, &tableComment)

	//查找主键
	finderPK := zorm.NewFinder()
	finderPK.Append("SELECT column_name FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA=? and  table_name=? AND constraint_name=?", dbName, tableName, "PRIMARY")
	pkName := ""
	zorm.QueryRow(ctx, finderPK, &pkName)

	finder2 := zorm.NewFinder()
	// select * from information_schema.COLUMNS where table_schema ='readygo' and table_name='t_user';
	finder2.Append("select COLUMN_NAME,DATA_TYPE,IS_NULLABLE,COLUMN_COMMENT from information_schema.COLUMNS where  TABLE_SCHEMA =? and TABLE_NAME=? and COLUMN_NAME not like ?  order by ORDINAL_POSITION asc", dbName, tableName, "bak%")

	maps, _ := zorm.QueryMap(ctx, finder2, nil)

	for _, m := range maps {
		dataType := m["DATA_TYPE"].(string)
		dataType = strings.ToUpper(dataType)

		nullable := m["IS_NULLABLE"].(string)
		nullable = strings.ToUpper(nullable)

		if dataType == "VARCHAR" || dataType == "NVARCHAR" || dataType == "TEXT" {
			//if nullable == "YES" {
			//	dataType = "sql.NullString"
			//} else {
			dataType = "string"
			//}

		} else if dataType == "DATETIME" || dataType == "TIMESTAMP" {
			//if nullable == "YES" {
			//	dataType = "sql.NullTime"
			//} else {
			dataType = "time.Time"
			//}

		} else if dataType == "INT" {
			//if nullable == "YES" {
			//	dataType = "sql.NullInt32"
			//} else {
			dataType = "int"
			//}

		} else if dataType == "BIGINT" {
			//if nullable == "YES" {
			//	dataType = "sql.NullInt64"
			//} else {
			dataType = "int64"
			//}

		} else if dataType == "FLOAT" {
			//if nullable == "YES" {
			//	dataType = "sql.NullFloat64"
			//} else {
			dataType = "float32"
			//}

		} else if dataType == "DOUBLE" {
			//if nullable == "YES" {
			//	dataType = "sql.NullFloat64"
			//} else {
			dataType = "float64"
			//}

		} else if dataType == "DECIMAL" {
			dataType = "decimal.Decimal"
		}

		m["DATA_TYPE"] = dataType
		fieldName := camelCaseName(m["COLUMN_NAME"].(string))
		m["field"] = fieldName

		//设置主键的struct属性名称
		if m["COLUMN_NAME"].(string) == pkName {
			info["pkField"] = fieldName
		}

	}

	info["columns"] = maps
	info["pkName"] = pkName
	info["tableName"] = tableName
	structName := tableName
	if strings.HasPrefix(structName, "t_") {
		structName = structName[2:]
	}
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
