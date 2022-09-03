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
		DSN: "root:root@tcp(127.0.0.1:3306)/readygo?charset=utf8&parseTime=true",
		//DriverName 数据库驱动名称:mysql,postgres,oci8,sqlserver,sqlite3,go_ibm_db,clickhouse,dm,kingbase,aci,taosSql|taosRestful 和Dialect对应
		DriverName: "mysql",
		//Dialect 数据库方言:mysql,postgresql,oracle,mssql,sqlite,db2,clickhouse,dm,kingbase,shentong,tdengine 和 DriverName 对应
		Dialect: "mysql",
		//MaxOpenConns 数据库最大连接数 默认50
		MaxOpenConns: 50,
		//MaxIdleConns 数据库最大空闲连接数 默认50
		MaxIdleConns: 50,
		//ConnMaxLifetimeSecond 连接存活秒时间. 默认600(10分钟)后连接被销毁重建.避免数据库主动断开连接,造成死连接.MySQL默认wait_timeout 28800秒(8小时)
		ConnMaxLifetimeSecond: 600,
		//SlowSQLMillis 慢sql的时间阈值,单位毫秒.小于0是禁用SQL语句输出;等于0是只输出SQL语句,不计算执行时间;大于0是计算SQL执行时间,并且>=SlowSQLMillis值
		SlowSQLMillis: 0,
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
