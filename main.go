package main

import (
	"fmt"
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
	baseDao.Save(&user)
	baseDao.Query("select id,account from t_user")
	user.Account = "update"
	baseDao.Update(&user, false)
	baseDao.Query("select id,account from t_user")

	userMap := orm.NewEntityMap("t_user")

	userMap.Set("id", "mapId")
	userMap.Set("account", "mapAccount")
	baseDao.SaveMap(&userMap)
	userMap.Set("account", "213")
	baseDao.UpdateMap(&userMap)
	baseDao.Query("select id,account from t_user")

	finder := orm.NewFinder()
	finder.Append("SELECT * sfsdf ")
	fmt.Println(finder.GetSQL())

}
