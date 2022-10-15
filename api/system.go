/*
 * @Author: your name
 * @Date: 2021-03-10 17:02:24
 * @LastEditTime: 2021-03-10 17:45:17
 * @LastEditors: Please set LastEditors
 * @Description: 系统接口
 * @FilePath: \readygo\api\system.go
 */

package api

import (
	"context"
	"readygo/apistruct"
	"readygo/permission/permutil"

	_ "readygo/docs" // swagger使用

	"github.com/cloudwego/hertz/pkg/app"
)

// Login 登录方法
// @Summary 接口概要说明
// @Description 接口详细描述信息
// @Tags 登陆
// @Param id path int true "ID"
// @Param name query string false "name"
// @Success 200 {string} string "ok"
// @Router /test/{id} [get]    //路由信息，一定要写上
func Login(ctx context.Context, c *app.RequestContext) {
	token, err := permutil.JWECreateToken("u_10001", nil)
	if err == nil {
		c.JSON(200, apistruct.ResponseBodyModel{
			Status:  200,
			Message: "",
			Data:    token,
		})
	} else {
		c.JSON(500, apistruct.ResponseBodyModel{
			Status:  500,
			Message: err.Error(),
			Data:    "",
		})
	}
}
