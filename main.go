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
		DBType:   orm.DBType_MYSQL,
	}
	baseDao, _ := orm.NewBaseDao(&dataSourceConfig)

	user := shop.User2{
		Id:      "id",
		Account: "test",
	}
	baseDao.Delete(&user)
	//baseDao.Save(&user)
	baseDao.Query("select id,account from t_user")

	baseDao.Update(&user, false)

}
