package permservice

import (
	"context"
	"errors"
	"fmt"

	"readygo/cache"
	"readygo/permission/permstruct"

	"gitee.com/chunanyong/zorm"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// SaveRoleStruct 保存角色
// 如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
// 如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
// 如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
// 但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveRoleStruct(ctx context.Context, roleStruct *permstruct.RoleStruct) error {
	// roleStruct对象指针不能为空
	if roleStruct == nil {
		return errors.New("roleStruct对象指针不能为空")
	}
	// 匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveRoleStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		// 事务下的业务代码开始

		// 赋值主键Id
		if len(roleStruct.Id) < 1 {
			roleStruct.Id = zorm.FuncGenerateStringID(ctx)
		}

		_, errSaveRoleStruct := zorm.Insert(ctx, roleStruct)

		if errSaveRoleStruct != nil {
			return nil, errSaveRoleStruct
		}

		return nil, nil
		// 事务下的业务代码结束
	})

	// 记录错误
	if errSaveRoleStruct != nil {
		errSaveRoleStruct := fmt.Errorf("permservice.SaveRoleStruct错误:%w", errSaveRoleStruct)
		hlog.Error(errSaveRoleStruct)
		return errSaveRoleStruct
	}

	return nil
}

// UpdateRoleStruct 更新角色
// 如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
// 如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
// 如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
// 但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateRoleStruct(ctx context.Context, roleStruct *permstruct.RoleStruct) error {
	// roleStruct对象指针或主键Id不能为空
	if roleStruct == nil || len(roleStruct.Id) < 1 {
		return errors.New("roleStruct对象指针或主键Id不能为空")
	}

	// 匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateRoleStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		// 事务下的业务代码开始
		_, errUpdateRoleStruct := zorm.Update(ctx, roleStruct)

		if errUpdateRoleStruct != nil {
			return nil, errUpdateRoleStruct
		}

		return nil, nil
		// 事务下的业务代码结束
	})

	// 记录错误
	if errUpdateRoleStruct != nil {
		errUpdateRoleStruct := fmt.Errorf("permservice.UpdateRoleStruct错误:%w", errUpdateRoleStruct)
		hlog.Error(errUpdateRoleStruct)
		return errUpdateRoleStruct
	}

	// 清除缓存
	go cache.EvictKey(ctx, baseInfoCacheKey, "FindRoleStructById_"+roleStruct.Id)

	return nil
}

// DeleteRoleStructById 根据Id删除角色
// 如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
// 如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
// 如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
// 但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteRoleStructById(ctx context.Context, id string) error {
	// id不能为空
	if len(id) < 1 {
		return errors.New("id不能为空")
	}

	// 匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteRoleStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		// 事务下的业务代码开始
		finder := zorm.NewDeleteFinder(permstruct.RoleStructTableName).Append(" WHERE id=?", id)
		_, errDeleteRoleStruct := zorm.UpdateFinder(ctx, finder)

		if errDeleteRoleStruct != nil {
			return nil, errDeleteRoleStruct
		}

		return nil, nil
		// 事务下的业务代码结束
	})

	// 记录错误
	if errDeleteRoleStruct != nil {
		errDeleteRoleStruct := fmt.Errorf("permservice.DeleteRoleStruct错误:%w", errDeleteRoleStruct)
		hlog.Error(errDeleteRoleStruct)
		return errDeleteRoleStruct
	}

	// 清除缓存
	go cache.EvictKey(ctx, baseInfoCacheKey, "FindRoleStructById_"+id)

	return nil
}

// FindRoleStructById 根据Id查询角色信息
// ctx中如果没有dbConnection,则会使用默认的datasource进行无事务查询
func FindRoleStructById(ctx context.Context, id string) (*permstruct.RoleStruct, error) {
	// id不能为空
	if len(id) < 1 {
		return nil, errors.New("id不能为空")
	}

	roleStruct := permstruct.RoleStruct{}

	cacheKey := "FindRoleStructById_" + id
	cache.GetFromCache(ctx, baseInfoCacheKey, cacheKey, &roleStruct)
	if len(roleStruct.Id) > 0 { // 如果缓存中存在
		return &roleStruct, nil
	}
	// 根据Id查询
	finder := zorm.NewSelectFinder(permstruct.RoleStructTableName).Append(" WHERE id=?", id)

	_, errFindRoleStructById := zorm.QueryRow(ctx, finder, &roleStruct)

	// 记录错误
	if errFindRoleStructById != nil {
		errFindRoleStructById := fmt.Errorf("permservice.FindRoleStructById错误:%w", errFindRoleStructById)
		hlog.Error(errFindRoleStructById)
		return nil, errFindRoleStructById
	}

	// 放入缓存
	cache.PutToCache(ctx, baseInfoCacheKey, cacheKey, roleStruct)
	return &roleStruct, nil
}

// FindRoleStructList 根据Finder查询角色列表
// ctx中如果没有dbConnection,则会使用默认的datasource进行无事务查询
func FindRoleStructList(ctx context.Context, finder *zorm.Finder, page *zorm.Page) ([]permstruct.RoleStruct, error) {
	// finder不能为空
	if finder == nil {
		return nil, errors.New("finder不能为空")
	}

	roleStructList := make([]permstruct.RoleStruct, 0)
	errFindRoleStructList := zorm.Query(ctx, finder, &roleStructList, page)

	// 记录错误
	if errFindRoleStructList != nil {
		errFindRoleStructList := fmt.Errorf("permservice.FindRoleStructList错误:%w", errFindRoleStructList)
		hlog.Error(errFindRoleStructList)
		return nil, errFindRoleStructList
	}

	return roleStructList, nil
}
