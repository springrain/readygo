package permservice

import (
	"errors"
	"fmt"
	"readygo/logger"
	"readygo/orm"
	"readygo/permission/permstruct"
)

//SaveUserStruct 保存用户
//如果入参session为nil或者没事务,则会使用默认的datasource的开启事务并最后提交.如果session有事务,则只使用,不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveUserStruct(session *orm.Session, userStruct *permstruct.UserStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, saveUserStructErr := orm.Transaction(session, func(session *orm.Session) (interface{}, error) {

		//事务下的业务代码开始

		saveUserStructErr := orm.SaveStruct(session, userStruct)
		if saveUserStructErr != nil {
			return nil, saveUserStructErr
		}
		return nil, nil

		//事务下的业务代码结束

	})

	//记录错误
	if saveUserStructErr != nil {
		saveUserStructErr := fmt.Errorf("permservice.SaveUserStruct错误:%w", saveUserStructErr)
		logger.Error(saveUserStructErr)
		return saveUserStructErr
	}

	return nil
}

//UpdateUserStruct 更新用户
//如果入参session为nil或者没事务,则会使用默认的datasource的开启事务并最后提交.如果session有事务,则只使用,不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateUserStruct(session *orm.Session, userStruct *permstruct.UserStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, updateUserStructErr := orm.Transaction(session, func(session *orm.Session) (interface{}, error) {

		//事务下的业务代码开始

		updateUserStructErr := orm.UpdateStruct(session, userStruct)
		if updateUserStructErr != nil {
			return nil, updateUserStructErr
		}

		return nil, nil

		//事务下的业务代码结束

	})

	//记录错误
	if updateUserStructErr != nil {
		updateUserStructErr := fmt.Errorf("permservice.UpdateUserStruct错误:%w", updateUserStructErr)
		logger.Error(updateUserStructErr)
		return updateUserStructErr
	}

	return nil
}

//DeleteUserStruct 删除用户
//如果入参session为nil或者没事务,则会使用默认的datasource的开启事务并最后提交.如果session有事务,则只使用,不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteUserStruct(session *orm.Session, userStruct *permstruct.UserStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, deleteUserStructErr := orm.Transaction(session, func(session *orm.Session) (interface{}, error) {

		//事务下的业务代码开始

		deleteUserStructErr := orm.DeleteStruct(session, userStruct)
		if deleteUserStructErr != nil {
			return nil, deleteUserStructErr
		}

		return nil, nil

		//事务下的业务代码结束

	})

	//记录错误
	if deleteUserStructErr != nil {
		deleteUserStructErr := fmt.Errorf("permservice.DeleteUserStruct错误:%w", deleteUserStructErr)
		logger.Error(deleteUserStructErr)
		return deleteUserStructErr
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
	findUserStructByIdErr := orm.QueryStruct(session, finder, &userStruct)

	//记录错误
	if findUserStructByIdErr != nil {
		findUserStructByIdErr := fmt.Errorf("permservice.FindUserStructById错误:%w", findUserStructByIdErr)
		logger.Error(findUserStructByIdErr)
		return nil, findUserStructByIdErr
	}

	return &userStruct, nil

}

//FindUserStructList 根据Finder查询用户列表
//session如果为nil,则会使用默认的datasource进行无事务查询
func FindUserStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.UserStruct, error) {
	userStructList := make([]permstruct.UserStruct, 0)
	findUserStructListErr := orm.QueryStructList(session, finder, &userStructList, page)

	//记录错误
	if findUserStructListErr != nil {
		findUserStructListErr := fmt.Errorf("permservice.FindUserStructList错误:%w", findUserStructListErr)
		logger.Error(findUserStructListErr)
		return nil, findUserStructListErr
	}

	return userStructList, nil
}
