package test

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"readygo/cache"
	"readygo/permission/permservice"
	"readygo/permission/permstruct"

	"gitee.com/chunanyong/zorm"
	_ "github.com/go-sql-driver/mysql"
)

var dbDao *zorm.DBDao

// 性别枚举
type Gander int

func (g Gander) String() string {
	return []string{"Male", "Female", "Bisexual"}[g]
}

const (
	// 男的
	Male = iota
	// 女的
	Female
	// 跨性别
	Bisexual
)

var ctx = context.Background()

func init() {
	dataSourceConfig := zorm.DataSourceConfig{
		DSN: "root:root@tcp(127.0.0.1:3306)/readygo?charset=utf8&parseTime=true",
		// DriverName 数据库驱动名称:mysql,postgres,oci8,sqlserver,sqlite3,go_ibm_db,clickhouse,dm,kingbase,aci,taosSql|taosRestful 和Dialect对应
		DriverName: "mysql",
		// Dialect 数据库方言:mysql,postgresql,oracle,mssql,sqlite,db2,clickhouse,dm,kingbase,shentong,tdengine 和 DriverName 对应
		Dialect: "mysql",
		// MaxOpenConns 数据库最大连接数 默认50
		MaxOpenConns: 50,
		// MaxIdleConns 数据库最大空闲连接数 默认50
		MaxIdleConns: 50,
		// ConnMaxLifetimeSecond 连接存活秒时间. 默认600(10分钟)后连接被销毁重建.避免数据库主动断开连接,造成死连接.MySQL默认wait_timeout 28800秒(8小时)
		ConnMaxLifetimeSecond: 600,
		// SlowSQLMillis 慢sql的时间阈值,单位毫秒.小于0是禁用SQL语句输出;等于0是只输出SQL语句,不计算执行时间;大于0是计算SQL执行时间,并且>=SlowSQLMillis值
		SlowSQLMillis: 0,
	}
	dbDao, _ = zorm.NewDBDao(&dataSourceConfig)

	//	cache.NewMemeryCacheManager()
	cache.NewRedisClient(ctx, &cache.RedisConfig{
		Addr: "127.0.0.1:6379",
	})
}

func initDate() {
}

func TestQuey(t *testing.T) {
	finder := zorm.NewSelectFinder(permstruct.UserStructTableName)
	finder.Append(" order by id ")
	page := zorm.NewPage()

	var users []permstruct.UserStruct

	background := context.Background()

	err := zorm.Query(background, finder, &users, page)

	fmt.Println(users)

	if err != nil {
		// 标记测试失败
		t.Errorf("TestQuey错误:%v", err)
	}
	t.Log("总条数:", page.TotalCount)
	t.Log(users)
}

func TestNull(t *testing.T) {
	finder := zorm.NewFinder()

	finder.Append("select * from t_user limit 1")

	queryMap, err := zorm.QueryRowMap(nil, finder)
	if err != nil {
		t.Errorf("TestNull：%v", err)
	}

	t.Log(queryMap)
}

func TestCount(t *testing.T) {
	finder := zorm.NewFinder()

	finder.Append("select count(*) as c from ").Append(permstruct.WxCpconfigStructTableName)

	queryMap, err := zorm.QueryRowMap(context.Background(), finder)
	if err != nil {
		t.Errorf("TestCount错误：%v", err)
	}

	t.Log(queryMap)
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {
		var u permstruct.UserStruct
		//
		u.UserName = "zyf"
		u.CreateTime = time.Now()
		u.UpdateTime = time.Now()

		u.Sex = "男" + string(id)
		// u.Active = 2/0
		incr, _ := cache.RedisINCR(ctx, "permstruct.UserStruct")

		u.Id = strconv.Itoa(int(incr.(int64)))

		_, e2 := zorm.Insert(ctx, &u)
		if e2 != nil {
			// 标记测试失败
			// t.Errorf("TestTrancSave错误:%v", e2)
			return nil, e2
		}

		finder := zorm.NewSelectFinder(permstruct.UserStructTableName).Append(" where id = ?", "1583077877688617000")

		_, ee := zorm.QueryRow(ctx, finder, &u)

		fmt.Println(ee)

		// u.UserName = u.UserName + "test" + string(id)
		u.UserName = strconv.Itoa(id)

		u.UserType = id

		_, e3 := zorm.Update(ctx, &u)
		if e3 != nil {
			// 标记测试失败
			// t.Errorf("TestTrancUpdate错误:%v", e3)
			return nil, e3
		}

		return nil, nil
	})
}

func TestTranc(t *testing.T) {
	var wg sync.WaitGroup

	begin := time.Now()

	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	subTime := time.Now().Sub(begin)
	t.Log(subTime)

	wg.Wait()
}

func TestPrem(t *testing.T) {
	finder, err := permservice.WrapOrgIdFinderByPrivateOrgRoleId(context.Background(), "r_10001", "u_10001")
	if err != nil {
		t.Error(err)
	}
	t.Log(finder.GetSQL())
}

func TestUpdateNotZero(t *testing.T) {
	user := permstruct.UserStruct{}
	user.Id = "1583077877688617000"
	user.UserName = "abc"

	zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {
		zorm.UpdateNotZeroValue(ctx, &user)

		return nil, nil
	})
}
