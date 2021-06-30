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
	"readygo/permission/permhandler"
	"readygo/permission/permutil"

	"gitee.com/chunanyong/zorm"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Creates a router without any middleware by default
var GinEngine *gin.Engine = gin.New()

//初始化
func init() {
	// 初始化Gin引擎
	initGinEngine()

	//初始化DBDao
	dbDaoConfig := zorm.DataSourceConfig{
		DSN:        "root:root@tcp(127.0.0.1:3306)/readygo?charset=utf8&parseTime=true",
		DriverName: "mysql",
		PrintSQL:   true,
	}
	_, _ = zorm.NewDBDao(&dbDaoConfig)
	cache.NewMemeryCacheManager()

	permutil.NewJWEConfig("permission/permcert/private.pem", "readygo", 0)
}

// initGinEngine 初始化Gin引擎
func initGinEngine() {

	//GinEngine = gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	GinEngine.Use(ginext.GinLogger())
	//r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	GinEngine.Use(ginext.GinRecovery())
	//r.Use(gin.Recovery())

	//加载自定义的权限过滤器
	GinEngine.Use(permhandler.PermHandler())

	//css js等静态文件
	GinEngine.Static("/assets", "./assets")
	GinEngine.StaticFS("/more_static", http.Dir("my_file_system"))
	GinEngine.StaticFile("/favicon.ico", "./resources/favicon.ico")
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.

// @host 127.0.0.1:8080
// @BasePath /
func main() {

	// swagger
	GinEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	GinEngine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	GinEngine.GET("/login", api.Login)

	GinEngine.GET("/system/menu/tree", func(c *gin.Context) {
		user, err := permhandler.GetCurrentUserFromContext(c.Request.Context())
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

	GinEngine.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
