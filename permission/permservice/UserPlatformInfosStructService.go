package permservice

import (
	"context"
	"errors"
	"fmt"
	"readygo/permission/permstruct"

	"gitee.com/chunanyong/logger"

	"gitee.com/chunanyong/zorm"
)

//SaveUserPlatformInfosStruct 保存用户平台信息表
//如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
//如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
//如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
//但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveUserPlatformInfosStruct(ctx context.Context, userPlatformInfosStruct *permstruct.UserPlatformInfosStruct) error {

	// userPlatformInfosStruct对象指针不能为空
	if userPlatformInfosStruct == nil {
		return errors.New("userPlatformInfosStruct对象指针不能为空")
	}
	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveUserPlatformInfosStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

		//事务下的业务代码开始

		//赋值主键Id
		if len(userPlatformInfosStruct.Id) < 1 {
			userPlatformInfosStruct.Id = zorm.FuncGenerateStringID(ctx)
		}

		_, errSaveUserPlatformInfosStruct := zorm.Insert(ctx, userPlatformInfosStruct)

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
//如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
//如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
//如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
//但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateUserPlatformInfosStruct(ctx context.Context, userPlatformInfosStruct *permstruct.UserPlatformInfosStruct) error {

	// userPlatformInfosStruct对象指针或主键Id不能为空
	if userPlatformInfosStruct == nil || len(userPlatformInfosStruct.Id) < 1 {
		return errors.New("userPlatformInfosStruct对象指针或主键Id不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateUserPlatformInfosStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

		//事务下的业务代码开始
		_, errUpdateUserPlatformInfosStruct := zorm.Update(ctx, userPlatformInfosStruct)

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

//DeleteUserPlatformInfosStructById 根据Id删除用户平台信息表
//如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
//如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
//如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
//但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteUserPlatformInfosStructById(ctx context.Context, id string) error {

	//id不能为空
	if len(id) < 1 {
		return errors.New("id不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteUserPlatformInfosStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

		//事务下的业务代码开始
		finder := zorm.NewDeleteFinder(permstruct.UserPlatformInfosStructTableName).Append(" WHERE id=?", id)
		_, errDeleteUserPlatformInfosStruct := zorm.UpdateFinder(ctx, finder)

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
//ctx中如果没有dbConnection,则会使用默认的datasource进行无事务查询
func FindUserPlatformInfosStructById(ctx context.Context, id string) (*permstruct.UserPlatformInfosStruct, error) {
	//id不能为空
	if len(id) < 1 {
		return nil, errors.New("id不能为空")
	}

	//根据Id查询
	finder := zorm.NewSelectFinder(permstruct.UserPlatformInfosStructTableName).Append(" WHERE id=?", id)
	userPlatformInfosStruct := permstruct.UserPlatformInfosStruct{}
	_, errFindUserPlatformInfosStructById := zorm.QueryRow(ctx, finder, &userPlatformInfosStruct)

	//记录错误
	if errFindUserPlatformInfosStructById != nil {
		errFindUserPlatformInfosStructById := fmt.Errorf("permservice.FindUserPlatformInfosStructById错误:%w", errFindUserPlatformInfosStructById)
		logger.Error(errFindUserPlatformInfosStructById)
		return nil, errFindUserPlatformInfosStructById
	}

	return &userPlatformInfosStruct, nil

}

//FindUserPlatformInfosStructList 根据Finder查询用户平台信息表列表
//ctx中如果没有dbConnection,则会使用默认的datasource进行无事务查询
func FindUserPlatformInfosStructList(ctx context.Context, finder *zorm.Finder, page *zorm.Page) ([]permstruct.UserPlatformInfosStruct, error) {

	//finder不能为空
	if finder == nil {
		return nil, errors.New("finder不能为空")
	}

	userPlatformInfosStructList := make([]permstruct.UserPlatformInfosStruct, 0)
	errFindUserPlatformInfosStructList := zorm.Query(ctx, finder, &userPlatformInfosStructList, page)

	//记录错误
	if errFindUserPlatformInfosStructList != nil {
		errFindUserPlatformInfosStructList := fmt.Errorf("permservice.FindUserPlatformInfosStructList错误:%w", errFindUserPlatformInfosStructList)
		logger.Error(errFindUserPlatformInfosStructList)
		return nil, errFindUserPlatformInfosStructList
	}

	return userPlatformInfosStructList, nil
}
