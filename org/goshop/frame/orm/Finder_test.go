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
