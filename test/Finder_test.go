package test

import (
	"fmt"
	"readygo/orm"
	"testing"
)

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
