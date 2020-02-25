package main

import (
	"fmt"
	"readygo/orm"
)

var baseDao *orm.BaseDao

const dbName = "readygo"

func init() {
	baseDaoConfig := orm.DataSourceConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   dbName,
		UserName: "root",
		PassWord: "123456789",
		DBType:   orm.DBType_MYSQL,
	}

	baseDao, _ = orm.NewBaseDao(&baseDaoConfig)
}

func main() {
	//selectAllTable()
	columns := selectTableColumn("t_user")
	fmt.Println(columns)


}

//获取所有的表名
func selectAllTable() []string {
	finder := orm.NewFinder()
	finder.Append("select table_name from information_schema.TABLES where  TABLE_SCHEMA =?", dbName)
	tableNames := []string{}
	baseDao.QueryStructList(nil, finder, &tableNames, nil)
	return tableNames
}

func selectTableColumn(tableName string) []map[string]interface{} {
	finder := orm.NewFinder()

	// select * from information_schema.COLUMNS where table_schema ='readygo' and table_name='t_user';
	finder.Append("select TABLE_NAME,COLUMN_NAME,DATA_TYPE,IS_NULLABLE from information_schema.COLUMNS where  TABLE_SCHEMA =? and TABLE_NAME=? order by ORDINAL_POSITION asc", dbName, tableName)

	maps, _ := baseDao.QueryMapList(nil, finder, nil)

	finderPK := orm.NewFinder()
	finderPK.Append("SELECT column_name FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA=? and  table_name=? AND constraint_name=?", dbName, tableName, "PRIMARY")
	pkName := ""
	baseDao.QueryStruct(nil, finderPK, &pkName)

	return maps
}
