package test

import (
	"strings"
	"testing"

	"gitee.com/chunanyong/zorm"
)

func TestUUID(t *testing.T) {

	uuid := zorm.FuncGenerateStringID()
	strings.Replace(uuid, "-", "", -1)

	t.Log(len(strings.Replace(uuid, "-", "", -1)))

}
