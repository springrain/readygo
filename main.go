/*
 * @Author: your name
 * @Date: 2020-02-25 23:00:00
 * @LastEditTime: 2020-03-12 19:22:35
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\main.go
 */
package main

import (
	"net/http"
	"readygo/apistruct"
	"readygo/cache"
	"readygo/ginext"
	"readygo/permission/permhandler"
	"readygo/permission/permutility/jwe"

	"gitee.com/chunanyong/zorm"

	"github.com/gin-gonic/gin"
)

//初始化BaseDao
func init() {

	baseDaoConfig := zorm.DataSourceConfig{
		DSN:        "root:root@tcp(127.0.0.1:3306)/readygo?charset=utf8&parseTime=true",
		DriverName: "mysql",
		PrintSQL:   true,
	}
	_, _ = zorm.NewBaseDao(&baseDaoConfig)
	cache.NewMemeryCacheManager()
}
func main() {

	// Creates a router without any middleware by default
	r := gin.New()

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

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	r.GET("/login", func(c *gin.Context) {
		// user := permhandler.TokenUser{
		// 	Id:    "1001",
		// 	Name:  "readygo",
		// 	Role:  "admin",
		// 	Group: 0,
		// }
		//ctx, _ := permhandler.SetCurrentUser(c.Request.Context(), user)
		//c.Request = c.Request.WithContext(ctx)
		//cuser, error := permhandler.GetCurrentUser(c.Request.Context())
		// if error == nil {
		// 	fmt.Println(cuser)
		// }
		token, err := jwe.CreateToken("u_10001", nil)
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
	})

	r.GET("/system/menu/tree", func(c *gin.Context) {
		user, err := permhandler.GetCurrentUser(c.Request.Context())
		// token := c.GetHeader("READYGOTOKEN")
		// userid, err := jwe.GetInfoFromToken(token, &user)
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

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
