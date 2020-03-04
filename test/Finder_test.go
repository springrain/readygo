package test

import (
	"fmt"
	"readygo/zorm"
	"testing"
)

func init() {

	dataSourceConfig := zorm.DataSourceConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "readygo",
		UserName: "root",
		PassWord: "root",
		DBType:   "mysql",
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
