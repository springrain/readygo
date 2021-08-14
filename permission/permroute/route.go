package permroute

import (
	"readygo/ginext"
	"readygo/permission/permapi"
)

func init() {
	NewRouter()
}

// NewRouter 路由配置
func NewRouter() {
	r := ginext.GinEngine()

	// 路由
	v1 := r.Group("/api/perm/v1")
	{
		v1.POST("ping", permapi.Ping)

		v1.POST("menulist", permapi.QueryMenu)

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
