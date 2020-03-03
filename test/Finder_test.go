package test

import (
	"fmt"
	"readygo/orm"
	"testing"
)



func init() {

	dataSourceConfig := orm.DataSourceConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "readygo",
		UserName: "root",
		PassWord: "root",
		DBType:   "mysql",
	}
	baseDao, _ = orm.NewBaseDao(&dataSourceConfig)
}

func TestAppend(t *testing.T) {
	finder := orm.NewFinder()
	finder.Append("SELECT * FROM t_user ")
	fmt.Println(finder.GetSQL())
}

func TestNewSelectFinder(t *testing.T) {
	finder := orm.NewSelectFinder("t_user", "id")
	//finder.Append("SELECT * FROM t_user ")
	fmt.Println(finder.GetSQL())
}
