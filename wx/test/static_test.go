package test

import (
	"gitee.com/chunanyong/zorm"
	"strings"
	"testing"
)

func TestUUID(t *testing.T)  {

	uuid := zorm.GenerateStringID()
	strings.Replace(uuid,"-","",-1)



	t.Log(len(strings.Replace(uuid,"-","",-1)))

}
