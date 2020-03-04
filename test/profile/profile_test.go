package profile

import (
	"fmt"
	"readygo/zorm"
	"testing"
)


var baseDao *zorm.BaseDao

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


func TestQuery(t *testing.T) {

	finder := zorm.NewFinder()

	finder.Append("select * from t_user limit 1")

	for i := 0; i< 10000; i++{

		queryMap, err := zorm.QueryMap(nil, finder)

		if err != nil {
			t.Errorf("TestNull：%v", err)
		}

		fmt.Println(queryMap)
	}


	//ok      readygo/test/profile    5.735s


}

