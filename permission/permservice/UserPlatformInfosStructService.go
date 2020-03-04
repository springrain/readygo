package permservice

import (
	"errors"
	"fmt"
	"readygo/logger"
	"readygo/permission/permstruct"
	"readygo/zorm"
)

//SaveUserPlatformInfosStruct 保存用户平台信息表
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveUserPlatformInfosStruct(dbConnection *zorm.DBConnection, userPlatformInfosStruct *permstruct.UserPlatformInfosStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveUserPlatformInfosStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		errSaveUserPlatformInfosStruct := zorm.SaveStruct(dbConnection, userPlatformInfosStruct)

		if errSaveUserPlatformInfosStruct != nil {
			return nil, errSaveUserPlatformInfosStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errSaveUserPlatformInfosStruct != nil {
		errSaveUserPlatformInfosStruct := fmt.Errorf("permservice.SaveUserPlatformInfosStruct错误:%w", errSaveUserPlatformInfosStruct)
		logger.Error(errSaveUserPlatformInfosStruct)
		return errSaveUserPlatformInfosStruct
	}

	return nil
}

//UpdateUserPlatformInfosStruct 更新用户平台信息表
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateUserPlatformInfosStruct(dbConnection *zorm.DBConnection, userPlatformInfosStruct *permstruct.UserPlatformInfosStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateUserPlatformInfosStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		errUpdateUserPlatformInfosStruct := zorm.UpdateStruct(dbConnection, userPlatformInfosStruct)

		if errUpdateUserPlatformInfosStruct != nil {
			return nil, errUpdateUserPlatformInfosStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errUpdateUserPlatformInfosStruct != nil {
		errUpdateUserPlatformInfosStruct := fmt.Errorf("permservice.UpdateUserPlatformInfosStruct错误:%w", errUpdateUserPlatformInfosStruct)
		logger.Error(errUpdateUserPlatformInfosStruct)
		return errUpdateUserPlatformInfosStruct
	}

	return nil
}

//DeleteUserPlatformInfosStruct 删除用户平台信息表
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteUserPlatformInfosStruct(dbConnection *zorm.DBConnection, userPlatformInfosStruct *permstruct.UserPlatformInfosStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteUserPlatformInfosStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		errDeleteUserPlatformInfosStruct := zorm.DeleteStruct(dbConnection, userPlatformInfosStruct)

		if errDeleteUserPlatformInfosStruct != nil {
			return nil, errDeleteUserPlatformInfosStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errDeleteUserPlatformInfosStruct != nil {
		errDeleteUserPlatformInfosStruct := fmt.Errorf("permservice.DeleteUserPlatformInfosStruct错误:%w", errDeleteUserPlatformInfosStruct)
		logger.Error(errDeleteUserPlatformInfosStruct)
		return errDeleteUserPlatformInfosStruct
	}

	return nil
}

//FindUserPlatformInfosStructById 根据Id查询用户平台信息表信息
//dbConnection如果为nil,则会使用默认的datasource进行无事务查询
func FindUserPlatformInfosStructById(dbConnection *zorm.DBConnection, id string) (*permstruct.UserPlatformInfosStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := zorm.NewSelectFinder(" WHERE id=?", id)
	userPlatformInfosStruct := permstruct.UserPlatformInfosStruct{}
	errFindUserPlatformInfosStructById := zorm.QueryStruct(dbConnection, finder, &userPlatformInfosStruct)

	//记录错误
	if errFindUserPlatformInfosStructById != nil {
		errFindUserPlatformInfosStructById := fmt.Errorf("permservice.FindUserPlatformInfosStructById错误:%w", errFindUserPlatformInfosStructById)
		logger.Error(errFindUserPlatformInfosStructById)
		return nil, errFindUserPlatformInfosStructById
	}

	return &userPlatformInfosStruct, nil

}

//FindUserPlatformInfosStructList 根据Finder查询用户平台信息表列表
//dbConnection如果为nil,则会使用默认的datasource进行无事务查询
func FindUserPlatformInfosStructList(dbConnection *zorm.DBConnection, finder *zorm.Finder, page *zorm.Page) ([]permstruct.UserPlatformInfosStruct, error) {
	userPlatformInfosStructList := make([]permstruct.UserPlatformInfosStruct, 0)
	errFindUserPlatformInfosStructList := zorm.QueryStructList(dbConnection, finder, &userPlatformInfosStructList, page)

	//记录错误
	if errFindUserPlatformInfosStructList != nil {
		errFindUserPlatformInfosStructList := fmt.Errorf("permservice.FindUserPlatformInfosStructList错误:%w", errFindUserPlatformInfosStructList)
		logger.Error(errFindUserPlatformInfosStructList)
		return nil, errFindUserPlatformInfosStructList
	}

	return userPlatformInfosStructList, nil
}
