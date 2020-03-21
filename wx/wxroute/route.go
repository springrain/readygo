package wxroute

import (
	_ "gitee.com/chunanyong/gowe"
	"github.com/gin-gonic/gin"
	"readygo/wx/wxapi"
 )

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 路由
	v1 := r.Group("/api/wx/v1")
	{
		v1.POST("ping", wxapi.Ping)

		v1.POST("WxMaCode2Session", wxapi.WxMaCode2Session)
		v1.GET("WxPayUnifiedOrder",wxapi.WxPayUnifiedOrder)




		//// 用户登录
		//v1.POST("user/register", UserRegister)
		//
		//// 用户登录
		//v1.POST("user/login", UserLogin)
		//
		//v1.GET("user/demo",Demo)
		//
		//// 需要登录保护的
		//auth := v1.Group("")
		//auth.Use(middleware.AuthRequired())
		//{
		//	// User Routing
		//	auth.GET("user/me", UserMe)
		//	auth.DELETE("user/logout", UserLogout)
		//}
	}


	return r
}