package main

import (
	"goshop/org/springrain/orm"
	"testing"
	"time"
)

var baseDao *orm.BaseDao

func initDatabase() {

	dataSourceConfig := orm.DataSourceConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "goshop",
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
	baseDao.Transaction(func(session *orm.Session) (interface{}, error) {
		baseDao.SaveStruct(session, &user)
		return nil, nil
	})

}

func TestAdd(t *testing.T) {

	initDatabase()

	//initDate()

	table := orm.NewEntityMap("user")

	table.Set("id", 11)

	baseDao.Transaction(func(session *orm.Session) (interface{}, error) {
		e1 := baseDao.SaveMap(session, &table)
		if e1 != nil {
			return nil, e1
		}
		var l Language
		l.Id = 11

		l.Name = "englist"
		e2 := baseDao.SaveStruct(session, &l)
		if e2 != nil {
			return nil, e2
		}
		return nil, nil
	})

}
