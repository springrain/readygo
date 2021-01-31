package profile

import (
	"context"
	"fmt"
	"testing"

	"gitee.com/chunanyong/zorm"
)

var dbDao *zorm.DBDao

func init() {

	dataSourceConfig := zorm.DataSourceConfig{
		DSN:        "root:root@tcp(127.0.0.1:3306)/readygo?charset=utf8&parseTime=true",
		DriverName: "mysql",
		PrintSQL:   true,
	}
	dbDao, _ = zorm.NewDBDao(&dataSourceConfig)
}

func TestQuery(t *testing.T) {

	finder := zorm.NewFinder()

	finder.Append("select * from t_user limit 1")

	for i := 0; i < 10000; i++ {

		ctx := context.Background()

		queryMap, err := zorm.QueryRowMap(ctx, finder)

		if err != nil {
			t.Errorf("TestNull：%v", err)
		}

		fmt.Println(queryMap)
	}

	//ok      readygo/test/profile    5.735s

}
