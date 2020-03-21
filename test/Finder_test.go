package test

import (
	"fmt"
	"testing"

	"gitee.com/chunanyong/zorm"
)

func init() {

	dataSourceConfig := zorm.DataSourceConfig{
		DSN:        "root:root@tcp(127.0.0.1:3306)/readygo",
		DriverName: "mysql",
		PrintSQL:   true,
	}
	baseDao, _ = zorm.NewBaseDao(&dataSourceConfig)
}

func TestAppend(t *testing.T) {
	finder := zorm.NewFinder()
	finder.Append("SELECT * FROM t_user ")
	fmt.Println(finder.GetSQL())
}

func TestNewSelectFinder(t *testing.T) {
	finder := zorm.NewSelectFinder("t_user", "id")
	//finder.Append("SELECT * FROM t_user ")
	fmt.Println(finder.GetSQL())
}
