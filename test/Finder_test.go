package test

import (
	"fmt"
	"testing"

	"gitee.com/chunanyong/zorm"
)

func TestAppend(t *testing.T) {
	finder := zorm.NewFinder()
	finder.Append("SELECT * FROM t_user ")
	fmt.Println(finder.GetSQL())
}

func TestNewSelectFinder(t *testing.T) {
	finder := zorm.NewSelectFinder("t_user", "id")
	// finder.Append("SELECT * FROM t_user ")
	fmt.Println(finder.GetSQL())
}
