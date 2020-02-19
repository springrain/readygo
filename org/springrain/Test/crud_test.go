package main

import (
	"goshop/org/springrain/orm"
	"testing"
	"time"
)

var  baseDao *orm.BaseDao

func initDatabase()  {


	dataSourceConfig := orm.DataSourceConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "goshop",
		UserName: "root",
		PassWord: "123456789",
		DBType:   orm.DBType_MYSQL,
	}
	baseDao, _ = orm.NewBaseDao(&dataSourceConfig)
}

func initDate()  {

	var user User
	user.CreatedAt = time.Now()
	user.Id = 3

	baseDao.SaveStruct(&user)



	
}


func TestAdd(t *testing.T) {

	initDatabase()

	//initDate()

	table := orm.NewEntityMap("user")

	table.Set("id",4)

	baseDao.SaveMap(&table)


	var l Language
	l.Id = 9

	l.Name = "englist"
	baseDao.SaveStruct(&l)

}
