package test

import (
	"context"
	"fmt"
	"readygo/cache"
	"readygo/permission/permservice"
	"readygo/permission/permstruct"
	"strconv"
	"sync"
	"testing"
	"time"

	"gitee.com/chunanyong/zorm"
)

var baseDao *zorm.BaseDao

//性别枚举
type Gander int

func (g Gander) String() string {
	return []string{"Male", "Female", "Bisexual"}[g]
}

const (
	//男的
	Male = iota
	//女的
	Female
	//跨性别
	Bisexual
)

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

//	cache.NewMemeryCacheManager()
	cache.NewRedisClient(&cache.RedisConfig{
		Addr:         "127.0.0.1:6379",
	})
}

func initDate() {

}

func TestQuey(t *testing.T) {

	finder := zorm.NewSelectFinder(permstruct.UserStructTableName)

	page := zorm.NewPage()

	var users []permstruct.UserStruct

	err := zorm.QueryStructList(context.Background(), finder, &users, &page)

	if err != nil {
		//标记测试失败
		t.Errorf("TestQuey错误:%v", err)
	}
	fmt.Println("总条数:", page.TotalCount)
	fmt.Println(users)

}
func TestNull(t *testing.T) {
	finder := zorm.NewFinder()

	finder.Append("select * from t_user limit 1")

	queryMap, err := zorm.QueryMap(nil, finder)

	if err != nil {
		t.Errorf("TestNull：%v", err)
	}

	fmt.Println(queryMap)
}

func TestCount(t *testing.T) {
	finder := zorm.NewFinder()

	finder.Append("select count(*) as c from ").Append(permstruct.WxCpconfigStructTableName)

	queryMap, err := zorm.QueryMap(context.Background(), finder)

	if err != nil {
		t.Errorf("TestCount错误：%v", err)
	}

	fmt.Println(queryMap)
}

func worker(id int, wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println(id)

	zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {

		var u permstruct.UserStruct
		//
		u.UserName = "zyf"
		u.CreateTime = time.Now()
		u.Sex = "男"+string(id)
		//u.Active = 2/0
		incr, _ := cache.RedisINCR("permstruct.UserStruct")


		fmt.Println(incr)
		u.Id = strconv.Itoa(int(incr.(int64)))

		e2 := zorm.SaveStruct(ctx, &u)
		if e2 != nil {
			//标记测试失败
			//t.Errorf("TestTrancSave错误:%v", e2)
			return nil, e2
		}

		finder := zorm.NewSelectFinder(permstruct.UserStructTableName).Append(" where id = ?", "1583077877688617000")

		zorm.QueryStruct(ctx, finder, &u)

		//u.UserName = u.UserName + "test" + string(id)
		u.UserName = strconv.Itoa(id)

		u.UserType = id

		e3 := zorm.UpdateStruct(ctx, &u)
		if e3 != nil {
			//标记测试失败
			//t.Errorf("TestTrancUpdate错误:%v", e3)
			return nil, e3
		}

		return nil, nil
	})
}

func TestTranc(t *testing.T) {

	var wg sync.WaitGroup

	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()

}

func TestPrem(t *testing.T) {
	finder, err := permservice.WrapOrgIdFinderByPrivateOrgRoleId(context.Background(), "r_10001", "u_10001")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(finder.GetSQL())

}

func TestUpdateNotZero(t *testing.T) {
	user := permstruct.UserStruct{}
	user.Id = "1583077877688617000"
	user.UserName = "abc"

	zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {
		zorm.UpdateStructNotZeroValue(ctx, &user)

		return nil, nil
	})

}
