package permservice

import (
	"errors"
	"fmt"
	"readygo/logger"
	"readygo/orm"
	"readygo/permission/permstruct"
)

//SaveUserStruct 保存用户
//如果入参session为nil,使用defaultDao开启事务并最后提交.如果入参session没有事务,调用session.begin()开启事务并最后提交.如果入参session有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveUserStruct(session *orm.Session, userStruct *permstruct.UserStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveUserStruct := orm.Transaction(session, func(session *orm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errSaveUserStruct := orm.SaveStruct(session, userStruct)

		if errSaveUserStruct != nil {
			return nil, errSaveUserStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errSaveUserStruct != nil {
		errSaveUserStruct := fmt.Errorf("permservice.SaveUserStruct错误:%w", errSaveUserStruct)
		logger.Error(errSaveUserStruct)
		return errSaveUserStruct
	}

	return nil
}

//UpdateUserStruct 更新用户
//如果入参session为nil,使用defaultDao开启事务并最后提交.如果入参session没有事务,调用session.begin()开启事务并最后提交.如果入参session有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateUserStruct(session *orm.Session, userStruct *permstruct.UserStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateUserStruct := orm.Transaction(session, func(session *orm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errUpdateUserStruct := orm.UpdateStruct(session, userStruct)

		if errUpdateUserStruct != nil {
			return nil, errUpdateUserStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errUpdateUserStruct != nil {
		errUpdateUserStruct := fmt.Errorf("permservice.UpdateUserStruct错误:%w", errUpdateUserStruct)
		logger.Error(errUpdateUserStruct)
		return errUpdateUserStruct
	}

	return nil
}

//DeleteUserStruct 删除用户
//如果入参session为nil,使用defaultDao开启事务并最后提交.如果入参session没有事务,调用session.begin()开启事务并最后提交.如果入参session有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteUserStruct(session *orm.Session, userStruct *permstruct.UserStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteUserStruct := orm.Transaction(session, func(session *orm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errDeleteUserStruct := orm.DeleteStruct(session, userStruct)

		if errDeleteUserStruct != nil {
			return nil, errDeleteUserStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errDeleteUserStruct != nil {
		errDeleteUserStruct := fmt.Errorf("permservice.DeleteUserStruct错误:%w", errDeleteUserStruct)
		logger.Error(errDeleteUserStruct)
		return errDeleteUserStruct
	}

	return nil
}

//FindUserStructById 根据Id查询用户信息
//session如果为nil,则会使用默认的datasource进行无事务查询
func FindUserStructById(session *orm.Session, id string) (*permstruct.UserStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	userStruct := permstruct.UserStruct{}
	errFindUserStructById := orm.QueryStruct(session, finder, &userStruct)

	//记录错误
	if errFindUserStructById != nil {
		errFindUserStructById := fmt.Errorf("permservice.FindUserStructById错误:%w", errFindUserStructById)
		logger.Error(errFindUserStructById)
		return nil, errFindUserStructById
	}

	return &userStruct, nil

}

//FindUserStructList 根据Finder查询用户列表
//session如果为nil,则会使用默认的datasource进行无事务查询
func FindUserStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.UserStruct, error) {
	userStructList := make([]permstruct.UserStruct, 0)
	errFindUserStructList := orm.QueryStructList(session, finder, &userStructList, page)

	//记录错误
	if errFindUserStructList != nil {
		errFindUserStructList := fmt.Errorf("permservice.FindUserStructList错误:%w", errFindUserStructList)
		logger.Error(errFindUserStructList)
		return nil, errFindUserStructList
	}

	return userStructList, nil
}
