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
	dbName      = "readygo"
	packageName = "code"
)

func init() {
	baseDaoConfig := orm.DataSourceConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   dbName,
		UserName: "root",
		PassWord: "root",
		DBType:   orm.DBType_MYSQL,
	}

	baseDao, _ = orm.NewBaseDao(&baseDaoConfig)
}

func main() {
	code("t_user")

	tableNames := selectAllTable()
	for _, tableName := range tableNames {
		code(tableName)
	}

}

//生成代码
func code(tableName string) {

	info := selectTableColumn(tableName)

	structFileName := "./code/" + info["structName"].(string) + "Struct.go"
	f, err := os.Create(structFileName)

	//w := bufio.NewWriter(f) // 创建新的 Writer 对象
	defer func() {
		f.Close()

	}()

	t, err := template.ParseFiles("./templates/struct.txt")
	if err != nil {
		fmt.Println(err)
	}

	t.Execute(f, info)
}

//获取所有的表名
func selectAllTable() []string {
	finder := orm.NewFinder()
	finder.Append("select table_name from information_schema.TABLES where  TABLE_SCHEMA =?", dbName)
	tableNames := []string{}
	baseDao.QueryStructList(nil, finder, &tableNames, nil)
	return tableNames
}

//根据表名查询字段信息和主键名称
func selectTableColumn(tableName string) map[string]interface{} {
	tableComment := ""
	finder := orm.NewFinder()
	finder.Append("select table_comment from information_schema.TABLES where  TABLE_SCHEMA =? and TABLE_Name=? ", dbName, tableName)
	baseDao.QueryStruct(nil, finder, &tableComment)

	finder2 := orm.NewFinder()
	// select * from information_schema.COLUMNS where table_schema ='readygo' and table_name='t_user';
	finder2.Append("select COLUMN_NAME,DATA_TYPE,IS_NULLABLE,COLUMN_COMMENT from information_schema.COLUMNS where  TABLE_SCHEMA =? and TABLE_NAME=? order by ORDINAL_POSITION asc", dbName, tableName)

	maps, _ := baseDao.QueryMapList(nil, finder2, nil)

	for _, m := range maps {
		dataType := m["DATA_TYPE"].(string)
		if dataType == "varchar" {
			dataType = "string"
		} else if dataType == "datetime" {
			dataType = "time.Time"
		} else if dataType == "bigint" {
			dataType = "int64"
		}
		m["DATA_TYPE"] = dataType
		m["field"] = camelCaseName(m["COLUMN_NAME"].(string))
	}

	finderPK := orm.NewFinder()
	finderPK.Append("SELECT column_name FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA=? and  table_name=? AND constraint_name=?", dbName, tableName, "PRIMARY")
	pkName := ""
	baseDao.QueryStruct(nil, finderPK, &pkName)
	info := make(map[string]interface{})
	info["columns"] = maps
	info["pkName"] = pkName
	info["tableName"] = tableName
	info["structName"] = camelCaseName(tableName)
	info["packageName"] = packageName
	info["tableComment"] = tableComment
	return info
}

//首字母大写
func capitalize(str string) string {
	str = strings.ToUpper(string(str[0:1])) + string(str[1:])
	return str
}

//驼峰
func camelCaseName(tableName string) string {
	tableName = strings.Replace(tableName, "t_", "", 1)
	names := strings.Split(tableName, "_")
	structName := ""
	for _, name := range names {
		structName = structName + capitalize(name)
	}

	return structName

}
