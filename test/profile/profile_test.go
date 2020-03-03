package profile

import (
	"fmt"
	"readygo/orm"
	"testing"
)


var baseDao *orm.BaseDao

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


func TestQuery(t *testing.T) {

	finder := orm.NewFinder()

	finder.Append("select * from t_user limit 1")

	for i := 0; i< 10000; i++{

		queryMap, err := orm.QueryMap(nil, finder)

		if err != nil {
			t.Errorf("TestNull：%v", err)
		}

		fmt.Println(queryMap)
	}


	//ok      readygo/test/profile    5.735s


}


