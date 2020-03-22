package main

import (
	"readygo/permission/permroute"

	"gitee.com/chunanyong/zorm"
)

//初始化BaseDao
func init() {

	baseDaoConfig := zorm.DataSourceConfig{
		DSN:        "root:root@tcp(127.0.0.1:3306)/readygo?charset=utf8&parseTime=true",
		DriverName: "mysql",
		PrintSQL:   true,
	}
	_, _ = zorm.NewBaseDao(&baseDaoConfig)
}

func main() {

	r := permroute.NewRouter()
	r.Run(":3001")

}
