package permservice

import (
	"errors"
	"fmt"
	"readygo/logger"
	"readygo/permission/permstruct"
	"readygo/zorm"
)

//SaveMenuStruct 保存菜单
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveMenuStruct(dbConnection *zorm.DBConnection, menuStruct *permstruct.MenuStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveMenuStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		errSaveMenuStruct := zorm.SaveStruct(dbConnection, menuStruct)

		if errSaveMenuStruct != nil {
			return nil, errSaveMenuStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errSaveMenuStruct != nil {
		errSaveMenuStruct := fmt.Errorf("permservice.SaveMenuStruct错误:%w", errSaveMenuStruct)
		logger.Error(errSaveMenuStruct)
		return errSaveMenuStruct
	}

	return nil
}

//UpdateMenuStruct 更新菜单
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateMenuStruct(dbConnection *zorm.DBConnection, menuStruct *permstruct.MenuStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateMenuStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		errUpdateMenuStruct := zorm.UpdateStruct(dbConnection, menuStruct)

		if errUpdateMenuStruct != nil {
			return nil, errUpdateMenuStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errUpdateMenuStruct != nil {
		errUpdateMenuStruct := fmt.Errorf("permservice.UpdateMenuStruct错误:%w", errUpdateMenuStruct)
		logger.Error(errUpdateMenuStruct)
		return errUpdateMenuStruct
	}

	return nil
}

//DeleteMenuStruct 删除菜单
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteMenuStruct(dbConnection *zorm.DBConnection, menuStruct *permstruct.MenuStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteMenuStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		errDeleteMenuStruct := zorm.DeleteStruct(dbConnection, menuStruct)

		if errDeleteMenuStruct != nil {
			return nil, errDeleteMenuStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errDeleteMenuStruct != nil {
		errDeleteMenuStruct := fmt.Errorf("permservice.DeleteMenuStruct错误:%w", errDeleteMenuStruct)
		logger.Error(errDeleteMenuStruct)
		return errDeleteMenuStruct
	}

	return nil
}

//FindMenuStructById 根据Id查询菜单信息
//dbConnection如果为nil,则会使用默认的datasource进行无事务查询
func FindMenuStructById(dbConnection *zorm.DBConnection, id string) (*permstruct.MenuStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := zorm.NewSelectFinder(" WHERE id=?", id)
	menuStruct := permstruct.MenuStruct{}
	errFindMenuStructById := zorm.QueryStruct(dbConnection, finder, &menuStruct)

	//记录错误
	if errFindMenuStructById != nil {
		errFindMenuStructById := fmt.Errorf("permservice.FindMenuStructById错误:%w", errFindMenuStructById)
		logger.Error(errFindMenuStructById)
		return nil, errFindMenuStructById
	}

	return &menuStruct, nil

}

//FindMenuStructList 根据Finder查询菜单列表
//dbConnection如果为nil,则会使用默认的datasource进行无事务查询
func FindMenuStructList(dbConnection *zorm.DBConnection, finder *zorm.Finder, page *zorm.Page) ([]permstruct.MenuStruct, error) {
	menuStructList := make([]permstruct.MenuStruct, 0)
	errFindMenuStructList := zorm.QueryStructList(dbConnection, finder, &menuStructList, page)

	//记录错误
	if errFindMenuStructList != nil {
		errFindMenuStructList := fmt.Errorf("permservice.FindMenuStructList错误:%w", errFindMenuStructList)
		logger.Error(errFindMenuStructList)
		return nil, errFindMenuStructList
	}

	return menuStructList, nil
}

//FindMenuByPid 根据pid查询所有的子菜单
func FindMenuByPid(dbConnection *zorm.DBConnection, pid string, page *zorm.Page) ([]string, error) {

	f_select := zorm.NewSelectFinder(permstruct.AliPayconfigStructTableName, "id").Append(" WHERE active=1 ")

	if len(pid) > 0 { // pid不是根节点
		menu, errById := FindMenuStructById(dbConnection, pid)
		if errById != nil {
			return nil, errById
		}

		if menu.Comcode == "" { //没有编码,错误数据
			return nil, errors.New("Comcode为空,错误数据,pid:" + pid)
		}
		f_select.Append(" and comcode like ? ", menu.Comcode+"%")
	}

	f_select.Append(" order by sortno desc ")
	menuIds := make([]string, 0)
	errQueryList := zorm.QueryStructList(dbConnection, f_select, &menuIds, page)
	if errQueryList != nil {
		return menuIds, errQueryList
	}

	return menuIds, nil
}
