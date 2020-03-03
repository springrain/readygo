package permservice

import (
	"errors"
	"fmt"
	"readygo/logger"
	"readygo/permission/permstruct"
	"readygo/zorm"
)

//SaveMenuStruct 保存菜单
//如果入参session为nil,使用defaultDao开启事务并最后提交.如果入参session没有事务,调用session.begin()开启事务并最后提交.如果入参session有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveMenuStruct(session *zorm.Session, menuStruct *permstruct.MenuStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveMenuStruct := zorm.Transaction(session, func(session *zorm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errSaveMenuStruct := zorm.SaveStruct(session, menuStruct)

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
//如果入参session为nil,使用defaultDao开启事务并最后提交.如果入参session没有事务,调用session.begin()开启事务并最后提交.如果入参session有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateMenuStruct(session *zorm.Session, menuStruct *permstruct.MenuStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateMenuStruct := zorm.Transaction(session, func(session *zorm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errUpdateMenuStruct := zorm.UpdateStruct(session, menuStruct)

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
//如果入参session为nil,使用defaultDao开启事务并最后提交.如果入参session没有事务,调用session.begin()开启事务并最后提交.如果入参session有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteMenuStruct(session *zorm.Session, menuStruct *permstruct.MenuStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteMenuStruct := zorm.Transaction(session, func(session *zorm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errDeleteMenuStruct := zorm.DeleteStruct(session, menuStruct)

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
//session如果为nil,则会使用默认的datasource进行无事务查询
func FindMenuStructById(session *zorm.Session, id string) (*permstruct.MenuStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := zorm.NewSelectFinder(" WHERE id=?", id)
	menuStruct := permstruct.MenuStruct{}
	errFindMenuStructById := zorm.QueryStruct(session, finder, &menuStruct)

	//记录错误
	if errFindMenuStructById != nil {
		errFindMenuStructById := fmt.Errorf("permservice.FindMenuStructById错误:%w", errFindMenuStructById)
		logger.Error(errFindMenuStructById)
		return nil, errFindMenuStructById
	}

	return &menuStruct, nil

}

//FindMenuStructList 根据Finder查询菜单列表
//session如果为nil,则会使用默认的datasource进行无事务查询
func FindMenuStructList(session *zorm.Session, finder *zorm.Finder, page *zorm.Page) ([]permstruct.MenuStruct, error) {
	menuStructList := make([]permstruct.MenuStruct, 0)
	errFindMenuStructList := zorm.QueryStructList(session, finder, &menuStructList, page)

	//记录错误
	if errFindMenuStructList != nil {
		errFindMenuStructList := fmt.Errorf("permservice.FindMenuStructList错误:%w", errFindMenuStructList)
		logger.Error(errFindMenuStructList)
		return nil, errFindMenuStructList
	}

	return menuStructList, nil
}
