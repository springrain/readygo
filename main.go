package main

import (
	"fmt"
	"goshop/org/springrain/orm"
	"goshop/shop"
	"reflect"
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

	/*
		finder4 := orm.NewSelectFinder("t_user", "*")
		maps, err := baseDao.QueryMapList(finder4, nil)
		fmt.Println(maps, err)
		users := []shop.User2{}
		finder3 := orm.NewSelectFinder("t_user", "*")
		baseDao.QueryStructList(finder3, &users, nil)
		fmt.Println(users)
	*/
	user5 := shop.User2{}
	finder5 := orm.NewSelectFinder(user5.GetTableName()).Append(" WHERE id=? and id in (?) and id=?", "id", &[]string{"id", "abc", "sfsdf"}, "id")
	baseDao.QueryStruct(finder5, &user5)
	fmt.Println(user5)
	fmt.Println(reflect.TypeOf([]byte{}))

	/*


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
			fmt.Println(user.Account)
	*/

}
