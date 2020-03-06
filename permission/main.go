package main

import (
	"readygo/permission/permroute"
	"readygo/zorm"
)

//初始化BaseDao
func init() {

	baseDaoConfig := zorm.DataSourceConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "readygo",
		UserName: "root",
		PassWord: "root",
		DBType:   "mysql",
	}
	_, _ = zorm.NewBaseDao(&baseDaoConfig)
}


func main() {

	r := permroute.NewRouter()
	r.Run(":3001")


}