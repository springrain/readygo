package permroute

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterPermRoute 路由配置
func RegisterPermRoute(r *server.Hertz) {
	// r := webext.WebEngine()

	// 路由
	v1 := r.Group("/api/v1")

	//// 用户登录
	//v1.POST("user/register", UserRegister)
	//
	//// 用户登录
	v1.POST("login", nil)
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
