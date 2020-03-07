package permservice

import (
	"context"
	"errors"
	"fmt"
	"readygo/cache"
	"readygo/permission/permstruct"

	"gitee.com/chunanyong/logger"

	"gitee.com/chunanyong/zorm"
)

//SaveUserStruct 保存用户
//如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
//如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
//如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
//但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveUserStruct(ctx context.Context, userStruct *permstruct.UserStruct) error {

	// userStruct对象指针不能为空
	if userStruct == nil {
		return errors.New("userStruct对象指针不能为空")
	}
	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveUserStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

		//事务下的业务代码开始

		//赋值主键Id
		if len(userStruct.Id) < 1 {
			userStruct.Id = zorm.GenerateStringID()
		}

		errSaveUserStruct := zorm.SaveStruct(ctx, userStruct)

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
//如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
//如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
//如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
//但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateUserStruct(ctx context.Context, userStruct *permstruct.UserStruct) error {

	// userStruct对象指针或主键Id不能为空
	if userStruct == nil || len(userStruct.Id) < 1 {
		return errors.New("userStruct对象指针或主键Id不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateUserStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

		//事务下的业务代码开始
		errUpdateUserStruct := zorm.UpdateStruct(ctx, userStruct)

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
	//清理缓存
	cache.EvictKey(baseInfoCacheKey, "FindUserStructById_"+userStruct.Id)
	return nil
}

//DeleteUserStructById 根据Id删除用户
//如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
//如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
//如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
//但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteUserStructById(ctx context.Context, id string) error {

	//id不能为空
	if len(id) < 1 {
		return errors.New("id不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteUserStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

		//事务下的业务代码开始
		finder := zorm.NewDeleteFinder(permstruct.UserStructTableName).Append(" WHERE id=?", id)
		errDeleteUserStruct := zorm.UpdateFinder(ctx, finder)

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

	//清理缓存
	cache.EvictKey(baseInfoCacheKey, "FindUserStructById_"+id)

	return nil
}

//FindUserStructById 根据Id查询用户信息
//ctx中如果没有dbConnection,则会使用默认的datasource进行无事务查询
func FindUserStructById(ctx context.Context, id string) (*permstruct.UserStruct, error) {
	//id不能为空
	if len(id) < 1 {
		return nil, errors.New("id不能为空")
	}
	userStruct := permstruct.UserStruct{}
	cacheKey := "FindUserStructById_" + id
	cache.GetFromCache(baseInfoCacheKey, cacheKey, &userStruct)
	if len(userStruct.Id) > 0 { //缓存存在
		return &userStruct, nil
	}

	//根据Id查询
	finder := zorm.NewSelectFinder(permstruct.UserStructTableName).Append(" WHERE id=?", id)
	errFindUserStructById := zorm.QueryStruct(ctx, finder, &userStruct)

	//记录错误
	if errFindUserStructById != nil {
		errFindUserStructById := fmt.Errorf("permservice.FindUserStructById错误:%w", errFindUserStructById)
		logger.Error(errFindUserStructById)
		return nil, errFindUserStructById
	}

	//放入缓存
	cache.PutToCache(baseInfoCacheKey, cacheKey, userStruct)

	return &userStruct, nil

}

//FindUserStructList 根据Finder查询用户列表
//ctx中如果没有dbConnection,则会使用默认的datasource进行无事务查询
func FindUserStructList(ctx context.Context, finder *zorm.Finder, page *zorm.Page) ([]permstruct.UserStruct, error) {

	//finder不能为空
	if finder == nil {
		return nil, errors.New("finder不能为空")
	}

	userStructList := make([]permstruct.UserStruct, 0)
	errFindUserStructList := zorm.QueryStructList(ctx, finder, &userStructList, page)

	//记录错误
	if errFindUserStructList != nil {
		errFindUserStructList := fmt.Errorf("permservice.FindUserStructList错误:%w", errFindUserStructList)
		logger.Error(errFindUserStructList)
		return nil, errFindUserStructList
	}

	return userStructList, nil
}
