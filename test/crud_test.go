package test

import (
	"readygo/orm"
	"testing"
	"time"
)

var baseDao *orm.BaseDao

func init() {

	dataSourceConfig := orm.DataSourceConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "readygo",
		UserName: "root",
		PassWord: "root",
		DBType:   orm.DBType_MYSQL,
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

		var l Language
		l.Id = 22

		l.Name = "englist"
		e2 := orm.SaveStruct(session, &l)
		if e2 != nil {
			return nil, e2
		}
		return nil, nil
	})

}
