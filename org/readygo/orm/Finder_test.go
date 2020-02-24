package orm

import (
	"fmt"
	"testing"
)

func TestAppend(t *testing.T) {
	finder := NewFinder()
	finder.Append("SELECT * FROM t_user ")
	fmt.Println(finder.GetSQL())
}

func TestNewSelectFinder(t *testing.T) {
	finder := NewSelectFinder("t_user", "id")
	//finder.Append("SELECT * FROM t_user ")
	fmt.Println(finder.GetSQL())
}

func TestOrderBy(t *testing.T) {
	sqlstr := "select * FROM  t_user ORDER           by id asc "
	locOrderBy := findOrderByIndex(sqlstr)
	sqlstr = sqlstr[:locOrderBy[0]]
	fmt.Println(sqlstr)

}
