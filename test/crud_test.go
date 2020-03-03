package test

import (
	"fmt"
	"readygo/orm"
	"readygo/permission/permstruct"
	"strconv"
	"sync"
	"testing"
)

var baseDao *orm.BaseDao

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

func initDate() {

}

func TestQuey(t *testing.T) {

	finder := orm.NewSelectFinder(permstruct.UserStructTableName)

	page := orm.NewPage()

	var users []permstruct.UserStruct

	err := orm.QueryStructList(nil, finder, &users, &page)

	if err != nil {
		//标记测试失败
		t.Errorf("TestQuey错误:%v", err)
	}
	fmt.Println("总条数:", page.TotalCount)
	fmt.Println(users)

}
func TestNull(t *testing.T) {
	finder := orm.NewFinder()

	finder.Append("select * from t_user limit 1")



	queryMap, err := orm.QueryMap(nil, finder)

	if err != nil {
		t.Errorf("TestNull：%v", err)
	}

	fmt.Println(queryMap)
}

func TestCount(t *testing.T) {
	finder := orm.NewFinder()

	finder.Append("select count(*) as c from ").Append(permstruct.WxCpconfigStructTableName)

	queryMap, err := orm.QueryMap(nil, finder)

	if err != nil {
		t.Errorf("TestCount错误：%v", err)
	}

	fmt.Println(queryMap)
}

func worker(id int, wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println(id)

	orm.Transaction(nil, func(session *orm.Session) (interface{}, error) {

		var u permstruct.UserStruct
		//
		//u.UserName = "zyf"
		//u.CreateTime = time.Now()
		//u.Sex = "男"+string(id)
		////u.Active = 2/0
		//
		//e2 := orm.SaveStruct(session, &u)
		//if e2 != nil {
		//	//标记测试失败
		//	//t.Errorf("TestTrancSave错误:%v", e2)
		//	return nil, e2
		//}

		finder := orm.NewSelectFinder(permstruct.UserStructTableName).Append(" where id = ?", "1583077877688617000")

		orm.QueryStruct(session, finder, &u)

		//u.UserName = u.UserName + "test" + string(id)
		u.UserName = strconv.Itoa(id)

		u.UserType = id

		e3 := orm.UpdateStruct(session, &u)
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

	for i := 1; i <= 6; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()

}
