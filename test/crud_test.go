package test

import (
	"fmt"
	"readygo/orm"
	"readygo/permission/permstruct"
	"testing"
	"time"
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

	finder := orm.NewFinder()
	finder.Append("select * from ").Append(permstruct.UserStructTableName)

	page := orm.NewPage()

	var users = make([]permstruct.UserStruct, 1)

	err := orm.QueryStructList(nil, finder, users, &page)

	if err != nil {
		fmt.Println(err)

	}

	fmt.Println(users)

}

func TestTranc(t *testing.T) {

	orm.Transaction(nil, func(session *orm.Session) (interface{}, error) {

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
