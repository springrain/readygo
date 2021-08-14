/*
 * @Author: your name
 * @Date: 2020-02-25 23:00:00
 * @LastEditTime: 2021-03-10 17:42:33
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\main.go
 */
package main

import (
	"net/http"
	"readygo/api"
	"readygo/apistruct"
	"readygo/cache"
	"readygo/ginext"
	"readygo/permission/permapi"
	"readygo/permission/permhandler"
	"readygo/permission/permroute"
	"readygo/permission/permstruct"
	"readygo/permission/permutil"
	"readygo/wx/wxroute"

	"gitee.com/chunanyong/zorm"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//初始化
func init() {

	//初始化DBDao
	dbDaoConfig := zorm.DataSourceConfig{
		DSN:        "root:a123456@tcp(192.168.31.165:3306)/readygo?charset=utf8&parseTime=true",
		DriverName: "mysql",
		PrintSQL:   true,
	}
	_, _ = zorm.NewDBDao(&dbDaoConfig)
	cache.NewMemeryCacheManager()

	permutil.NewJWEConfig("permission/permcert/private.pem", "readygo", 0)

	//初始化initGinEngine
	initGinEngine()
}

// initGinEngine 初始化Gin引擎
func initGinEngine() {

	//获取引擎
	r := ginext.GinEngine()
	//设置前缀,需要在路由初始化前调用
	ginext.SetContextPath("/readygo/")

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(ginext.GinLogger())
	//r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(ginext.GinRecovery())
	//r.Use(gin.Recovery())

	//加载自定义的权限过滤器
	r.Use(permhandler.PermHandler())

	//css js等静态文件
	r.Static("/assets", "./assets")
	r.StaticFS("/more_static", http.Dir("my_file_system"))
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	r.GET("/login", api.Login)
	r.POST("/Captcha", permapi.Captcha)

	r.GET("/system/menu/tree", func(c *gin.Context) {

		user, err := permstruct.GetCurrentUserFromContext(c.Request.Context())
		// token := c.GetHeader(JWTTokenName)
		// userid, err := permutil.GetInfoFromToken(token, &user)
		if err == nil {
			c.JSON(http.StatusOK, apistruct.ResponseBodyModel{
				Status:  200,
				Message: "",
				Data:    gin.H{"userid": user.UserId, "extInfo": user},
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, apistruct.ResponseBodyModel{
				Status:  500,
				Message: err.Error(),
				Data:    gin.H{"msg": err.Error()},
			})
		}
	})

	permroute.RegisterPermRoute(r)
	wxroute.RegisterWXRoute(r)

}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.

// @host 127.0.0.1:8080
// @BasePath /
func main() {
	r := ginext.GinEngine()
	r.Run(":7080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
