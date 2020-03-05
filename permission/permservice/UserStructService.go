package permservice

import (
	"errors"
	"fmt"
	"readygo/cache"
	"readygo/logger"
	"readygo/permission/permstruct"
	"readygo/zorm"
)

//SaveUserStruct 保存用户
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveUserStruct(dbConnection *zorm.DBConnection, userStruct *permstruct.UserStruct) error {

	// userStruct对象指针不能为空
	if userStruct == nil {
		return errors.New("userStruct对象指针不能为空")
	}
	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveUserStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始

		//赋值主键Id
		if len(userStruct.Id) < 1 {
			userStruct.Id = zorm.GenerateStringID()
		}

		errSaveUserStruct := zorm.SaveStruct(dbConnection, userStruct)

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
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateUserStruct(dbConnection *zorm.DBConnection, userStruct *permstruct.UserStruct) error {

	// userStruct对象指针或主键Id不能为空
	if userStruct == nil || len(userStruct.Id) < 1 {
		return errors.New("userStruct对象指针或主键Id不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateUserStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		errUpdateUserStruct := zorm.UpdateStruct(dbConnection, userStruct)

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

//DeleteUserStructById 根据Id删除用户
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteUserStructById(dbConnection *zorm.DBConnection, id string) error {

	//id不能为空
	if len(id) < 1 {
		return errors.New("id不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteUserStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		finder := zorm.NewDeleteFinder(permstruct.UserStructTableName).Append(" WHERE id=?", id)
		errDeleteUserStruct := zorm.UpdateFinder(dbConnection, finder)

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
//dbConnection如果为nil,则会使用默认的datasource进行无事务查询
func FindUserStructById(dbConnection *zorm.DBConnection, id string) (*permstruct.UserStruct, error) {
	//id不能为空
	if len(id) < 1 {
		return nil, errors.New("id不能为空")
	}
	userStruct := permstruct.UserStruct{}
	cacheKey := "FindUserStructById_" + id
	cache.GetFromCache(qxCacheKey, cacheKey, &userStruct)
	if len(userStruct.Id) > 0 { //缓存存在
		return &userStruct, nil
	}

	//根据Id查询
	finder := zorm.NewSelectFinder(permstruct.UserStructTableName).Append(" WHERE id=?", id)
	errFindUserStructById := zorm.QueryStruct(dbConnection, finder, &userStruct)

	//记录错误
	if errFindUserStructById != nil {
		errFindUserStructById := fmt.Errorf("permservice.FindUserStructById错误:%w", errFindUserStructById)
		logger.Error(errFindUserStructById)
		return nil, errFindUserStructById
	}

	//放入缓存
	cache.PutToCache(qxCacheKey, cacheKey, userStruct)

	return &userStruct, nil

}

//FindUserStructList 根据Finder查询用户列表
//dbConnection如果为nil,则会使用默认的datasource进行无事务查询
func FindUserStructList(dbConnection *zorm.DBConnection, finder *zorm.Finder, page *zorm.Page) ([]permstruct.UserStruct, error) {

	//finder不能为空
	if finder == nil {
		return nil, errors.New("finder不能为空")
	}

	userStructList := make([]permstruct.UserStruct, 0)
	errFindUserStructList := zorm.QueryStructList(dbConnection, finder, &userStructList, page)

	//记录错误
	if errFindUserStructList != nil {
		errFindUserStructList := fmt.Errorf("permservice.FindUserStructList错误:%w", errFindUserStructList)
		logger.Error(errFindUserStructList)
		return nil, errFindUserStructList
	}

	return userStructList, nil
}
