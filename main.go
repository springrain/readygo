package main

import (
	"fmt"
	"goshop/org/springrain/orm"
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

	finder := orm.NewSelectFinder(user.GetTableName(), "id,account")
	finder.Append(" WHERE id=?", "id")

	baseDao.DeleteStruct(&user)
	baseDao.SaveStruct(&user)

	user = shop.User2{}

	baseDao.QueryStruct(finder, &user)
	fmt.Println(user.Account)
	/*
		user.Account = "update"
		baseDao.UpdateStruct(&user)
		baseDao.QueryStruct(finder, &user)

		userMap := orm.NewEntityMap("t_user")

		userMap.Set("id", "mapId")
		userMap.Set("account", "mapAccount")
		baseDao.SaveMap(&userMap)
		userMap.Set("account", "213")
		baseDao.UpdateMap(&userMap)
		baseDao.QueryStruct(finder, &user)

		finder2 := orm.NewUpdateFinder(user.GetTableName())
		finder2.Append("acc")
		finder2.Append("ount=?", "adad")
		baseDao.UpdateFinder(finder2)

		baseDao.QueryStruct(finder, &user)
	*/

}
