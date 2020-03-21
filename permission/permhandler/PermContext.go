/*
 * @Author: your name
 * @Date: 2020-03-11 22:35:04
 * @LastEditTime: 2020-03-12 12:50:04
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\permission\permhandler\PermContext.go
 */
package permhandler

import (
	"context"
	"errors"
	"readygo/permission/permstruct"
)

// WrapCurrentUserKey 上下文Key类型
type WrapCurrentUserKey string

// CurrentUserKey 上下文Key
var CurrentUserKey = WrapCurrentUserKey("currentUserKey")

// SetCurrentUserToCtx 将当前用户信息添加到上下文
func setCurrentUserToCtx(c context.Context, userInfo interface{}) (context.Context, error) {

	if c == nil {
		return nil, errors.New("context不能为nil")
	}
	c = context.WithValue(c, CurrentUserKey, userInfo)
	return c, nil
}

// GetCurrentUserFromCtx 将当前上下文获取用户信息
func getCurrentUserFromCtx(c context.Context) (interface{}, error) {
	if c == nil {
		return nil, errors.New("context不能为nil")
	}
	userInfo := c.Value(CurrentUserKey)
	return userInfo, nil
}

// SetCurrentUser 设置当前登录用户到上下文
func SetCurrentUser(c context.Context, user permstruct.UserVOStruct) (context.Context, error) {
	return setCurrentUserToCtx(c, user)
}

// GetCurrentUser 从上下文获取登录用户信息
func GetCurrentUser(c context.Context) (permstruct.UserVOStruct, error) {
	var user permstruct.UserVOStruct
	ctxUser, error := getCurrentUserFromCtx(c)
	if error != nil {
		return user, error
	}
	if user, ok := ctxUser.(permstruct.UserVOStruct); ok {
		return user, nil
	}
	return user, errors.New("没有用户的登录信息")
}