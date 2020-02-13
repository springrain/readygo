package main

import (
	"goshop/org/goshop/frame/orm"
	"goshop/shop"
)

func main() {

	dataSourceConfig := orm.DataSourceConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "shop",
		UserName: "root",
		PassWord: "root",
		DBType:   orm.DBType_MySQL,
	}
	baseDao, _ := orm.NewBaseDao(&dataSourceConfig)

	user := shop.User{
		Id:      "id",
		Account: "user1_username",
	}

	baseDao.Save(&user)
	baseDao.Query("select id,account from t_user")

}
