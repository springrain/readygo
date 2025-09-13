package permservice

import (
	"context"
	"errors"
	"fmt"

	"readygo/cache"
	"readygo/config"
	"readygo/permission/permstruct"
	"readygo/util"
	"readygo/webext"

	"gitee.com/chunanyong/zorm"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// SaveUserStruct 保存用户
// 如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
// 如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
// 如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
// 但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveUserStruct(ctx context.Context, userStruct *permstruct.UserStruct) error {
	// userStruct对象指针不能为空
	if userStruct == nil {
		return errors.New("userStruct对象指针不能为空")
	}
	// 匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveUserStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		// 事务下的业务代码开始

		// 赋值主键Id
		if len(userStruct.Id) < 1 {
			userStruct.Id = zorm.FuncGenerateStringID(ctx)
		}

		_, errSaveUserStruct := zorm.Insert(ctx, userStruct)

		if errSaveUserStruct != nil {
			return nil, errSaveUserStruct
		}

		return nil, nil
		// 事务下的业务代码结束
	})

	// 记录错误
	if errSaveUserStruct != nil {
		errSaveUserStruct := fmt.Errorf("permservice.SaveUserStruct错误:%w", errSaveUserStruct)
		hlog.Error(errSaveUserStruct)
		return errSaveUserStruct
	}

	return nil
}

// UpdateUserStruct 更新用户
// 如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
// 如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
// 如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
// 但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateUserStruct(ctx context.Context, userStruct *permstruct.UserStruct) error {
	// userStruct对象指针或主键Id不能为空
	if userStruct == nil || len(userStruct.Id) < 1 {
		return errors.New("userStruct对象指针或主键Id不能为空")
	}

	// 匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateUserStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		// 事务下的业务代码开始
		_, errUpdateUserStruct := zorm.Update(ctx, userStruct)

		if errUpdateUserStruct != nil {
			return nil, errUpdateUserStruct
		}

		return nil, nil
		// 事务下的业务代码结束
	})

	// 记录错误
	if errUpdateUserStruct != nil {
		errUpdateUserStruct := fmt.Errorf("permservice.UpdateUserStruct错误:%w", errUpdateUserStruct)
		hlog.Error(errUpdateUserStruct)
		return errUpdateUserStruct
	}
	// 清理缓存
	cache.EvictKey(ctx, baseInfoCacheKey, "FindUserStructById_"+userStruct.Id)
	return nil
}

// DeleteUserStructById 根据Id删除用户
// 如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
// 如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
// 如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
// 但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteUserStructById(ctx context.Context, id string) error {
	// id不能为空
	if len(id) < 1 {
		return errors.New("id不能为空")
	}

	// 匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteUserStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		// 事务下的业务代码开始
		finder := zorm.NewDeleteFinder(permstruct.UserStructTableName).Append(" WHERE id=?", id)
		_, errDeleteUserStruct := zorm.UpdateFinder(ctx, finder)

		if errDeleteUserStruct != nil {
			return nil, errDeleteUserStruct
		}

		return nil, nil
		// 事务下的业务代码结束
	})

	// 记录错误
	if errDeleteUserStruct != nil {
		errDeleteUserStruct := fmt.Errorf("permservice.DeleteUserStruct错误:%w", errDeleteUserStruct)
		hlog.Error(errDeleteUserStruct)
		return errDeleteUserStruct
	}

	// 清理缓存
	cache.EvictKey(ctx, baseInfoCacheKey, "FindUserStructById_"+id)

	return nil
}

// FindUserStructById 根据Id查询用户信息
// ctx中如果没有dbConnection,则会使用默认的datasource进行无事务查询
func FindUserStructById(ctx context.Context, id string) (*permstruct.UserStruct, error) {
	// id不能为空
	if len(id) < 1 {
		return nil, errors.New("id不能为空")
	}
	userStruct := permstruct.UserStruct{}
	cacheKey := "FindUserStructById_" + id
	cache.GetFromCache(ctx, baseInfoCacheKey, cacheKey, &userStruct)
	if len(userStruct.Id) > 0 { // 缓存存在
		return &userStruct, nil
	}

	// 根据Id查询
	finder := zorm.NewSelectFinder(permstruct.UserStructTableName).Append(" WHERE id=?", id)
	_, errFindUserStructById := zorm.QueryRow(ctx, finder, &userStruct)

	// 记录错误
	if errFindUserStructById != nil {
		errFindUserStructById := fmt.Errorf("permservice.FindUserStructById错误:%w", errFindUserStructById)
		hlog.Error(errFindUserStructById)
		return nil, errFindUserStructById
	}

	// 放入缓存
	cache.PutToCache(ctx, baseInfoCacheKey, cacheKey, userStruct)

	return &userStruct, nil
}

// FindUserStructList 根据Finder查询用户列表
// ctx中如果没有dbConnection,则会使用默认的datasource进行无事务查询
func FindUserStructList(ctx context.Context, finder *zorm.Finder, page *zorm.Page) ([]permstruct.UserStruct, error) {
	// finder不能为空
	if finder == nil {
		return nil, errors.New("finder不能为空")
	}

	userStructList := make([]permstruct.UserStruct, 0)
	errFindUserStructList := zorm.Query(ctx, finder, &userStructList, page)

	// 记录错误
	if errFindUserStructList != nil {
		errFindUserStructList := fmt.Errorf("permservice.FindUserStructList错误:%w", errFindUserStructList)
		hlog.Error(errFindUserStructList)
		return nil, errFindUserStructList
	}

	return userStructList, nil
}

func FindUserVOStructByUserId(ctx context.Context, userId string) (permstruct.UserVOStruct, error) {
	userVO := permstruct.UserVOStruct{}
	if userId == "" {
		return userVO, errors.New("userId不能为空")
	}
	userStruct, err := FindUserStructById(ctx, userId)
	if err != nil {
		return userVO, err
	}
	userVO.UserId = userStruct.Id
	userVO.Account = userStruct.Account
	userVO.Email = userStruct.Email
	userVO.UserName = userStruct.UserName
	userVO.UserType = userStruct.UserType
	userVO.Active = userStruct.Active
	userVO.PrivateOrgRoleId = ""
	return userVO, nil
}


// 用户登录
func Login(ctx context.Context, userStruct *permstruct.UserStruct) (interface{}, error) {
	// 参数校验
   	if userStruct == nil {
		return webext.ErrorReponseData(500, "登录请求不能为空"), nil
	}
	if userStruct.Account == "" {
		return webext.ErrorReponseData(500, "账号不能为空"), nil
	}
	if userStruct.Password == "" {
		return webext.ErrorReponseData(500, "密码不能为空"), nil
	}

	// 初始化用户实体，用于接收数据库查询结果
	var user permstruct.UserStruct

	// 根据Account查询
	finder := zorm.NewSelectFinder(permstruct.UserStructTableName).Append(" WHERE account=?", userStruct.Account)
	hasRecord, err := zorm.QueryRow(ctx, finder, &user)
	if err != nil {
		detailedErr := fmt.Errorf("permservice.Login - 查询数据库失败: %w", err)
		hlog.Error(detailedErr)
		return webext.ErrorReponseData(500, "查询数据库失败", detailedErr), nil
	}

	// 如果没有查到记录
	if !hasRecord {
		return webext.ErrorReponseData(500, "账号或密码错误"), nil
	}

	// 核对密码
	if userStruct.Password != user.Password {
		return webext.ErrorReponseData(500, "账号或密码错误"), nil
	}
	// 产生 Token
	token, err := util.NewJWTToken(user.Id)
	if err != nil {
		return webext.ErrorReponseData(500, "生成Token失败", err), nil
	}

	// 登录成功
	response := map[string]interface{}{
		config.Cfg.Jwt.TokenName: map[string]string{
			"access_token": token,
			"token_type":   "Bearer",
		},
		"user": map[string]string{
			"userId":   user.Id,
			"username": user.UserName,
		},
	}

	return webext.SuccessReponseData(response, "登录成功"), nil
}
