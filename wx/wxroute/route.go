package wxroute

import (
	"readygo/wx/wxapi"

	"github.com/gin-gonic/gin"
)

// RegisterWXRoute 路由配置
func RegisterWXRoute(r *gin.Engine) {
	//GinEngine := gin.Default()
	//r := ginext.GinEngine()
	// 路由
	v1 := r.Group("/api/wx/v1")
	{
		v1.POST("ping", wxapi.Ping)

		v1.POST("WxMaCode2Session", wxapi.WxMaCode2Session)
		v1.POST("WxPayUnifiedOrder", wxapi.WxPayUnifiedOrder)
		v1.POST("WxPayNotifyPay", wxapi.WxPayNotifyPay)

		v1.POST("WxPayAppSign", wxapi.WxPayAppSign)
		v1.GET("WxMaSubscribeMessageSend", wxapi.WxMaSubscribeMessageSend)

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

}
