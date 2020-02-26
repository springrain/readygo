package test

import (
	"readygo/orm"
	"readygo/permission/permstruct"
	"testing"
	"time"
)

var baseDao *orm.BaseDao



//性别枚举
type Gander int

func (g Gander)String() string{
	return []string{"Male","Female","Bisexual"}[g]
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

	var user User
	user.CreatedAt = time.Now()
	user.Id = 3
	orm.Transaction(func(session *orm.Session) (interface{}, error) {
		orm.SaveStruct(session, &user)
		return nil, nil
	})

}

func TestQuey(t *testing.T)  {



}

func TestTranc(t *testing.T) {

	/*
		table := orm.NewEntityMap("user")

		table.Set("id", 11)
		e1 := baseDao.SaveMap(session, &table)

		if e1 != nil {
			return nil, e1
		}
	*/

	orm.Transaction(func(session *orm.Session) (interface{}, error) {

		var u permstruct.UserStruct

		u.UserName = "zyf"
		u.CreateTime = time.Now()

		u.Sex = "男"



 		e2 := orm.SaveStruct(session, &u)
		if e2 != nil {
			return nil, e2
		}


		return nil, nil
	})

}
