/*
 * @Author: your name
 * @Date: 2020-03-11 22:35:04
 * @LastEditTime: 2020-03-12 12:50:04
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\permission\permstruct\PermContextUserVO.go
 */
package permstruct

import (
	"context"
	"errors"
)

// wrapCurrentUserKey 上下文Key类型
type wrapCurrentUserKey string

// currentUserKey 上下文Key
var currentUserKey = wrapCurrentUserKey("currentUserKey")

// SetCurrentUserToCtx 将当前用户信息添加到上下文
func setCurrentUserToCtx(c context.Context, userInfo interface{}) (context.Context, error) {

	if c == nil {
		return nil, errors.New("context不能为nil")
	}
	c = context.WithValue(c, currentUserKey, userInfo)
	return c, nil
}

// GetCurrentUserFromCtx 将当前上下文获取用户信息
func getCurrentUserFromCtx(ctx context.Context) (interface{}, error) {
	if ctx == nil {
		return nil, errors.New("context不能为nil")
	}
	userInfo := ctx.Value(currentUserKey)
	return userInfo, nil
}

// BindContextCurrentUser 设置当前登录用户到上下文
func BindContextCurrentUser(ctx context.Context, userVO UserVOStruct) (context.Context, error) {
	return setCurrentUserToCtx(ctx, userVO)
}

// GetCurrentUserFromContext 从上下文获取登录用户信息
func GetCurrentUserFromContext(ctx context.Context) (UserVOStruct, error) {
	var user UserVOStruct
	ctxUser, error := getCurrentUserFromCtx(ctx)
	if error != nil {
		return user, error
	}
	if user, ok := ctxUser.(UserVOStruct); ok {
		return user, nil
	}
	return user, errors.New("没有用户的登录信息")
}
