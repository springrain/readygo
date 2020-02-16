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

	finder := orm.NewSelectFinder(user.GetTableName())

	baseDao.Delete(&user)
	baseDao.Save(&user)
	baseDao.Query(finder, &user)
	user.Account = "update"
	baseDao.Update(&user)
	baseDao.Query(finder, &user)

	userMap := orm.NewEntityMap("t_user")

	userMap.Set("id", "mapId")
	userMap.Set("account", "mapAccount")
	baseDao.SaveMap(&userMap)
	userMap.Set("account", "213")
	baseDao.UpdateMap(&userMap)
	baseDao.Query(finder, &user)

	finder2 := orm.NewUpdateFinder(user.GetTableName())
	finder2.Append("acc")
	finder2.Append("ount=?", "adad")
	baseDao.UpdateFinder(finder2)

	baseDao.Query(finder, &user)

}
