/*
 * @Author: your name
 * @Date: 2020-02-25 23:00:00
 * @LastEditTime: 2020-02-27 17:13:53
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\main.go
 */
package main

import (
	"net/http"
	"readygo/ginext"
	"readygo/permission/permhandler"
	"readygo/utility/jwe"

	"gitee.com/chunanyong/zorm"

	"github.com/gin-gonic/gin"
)

//初始化BaseDao
func init() {

	baseDaoConfig := zorm.DataSourceConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "readygo",
		UserName: "root",
		PassWord: "root",
		DBType:   "mysql",
	}
	_, _ = zorm.NewBaseDao(&baseDaoConfig)
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
		user := jwe.TokenUser{
			Name:  "readygo",
			Role:  "admin",
			Group: 0,
		}
		token, err := jwe.CreateToken("1001", user)
		if err == nil {
			c.JSON(200, gin.H{"result": "OK", "token": token})
		} else {
			c.JSON(500, gin.H{"result": "Error", "msg": err.Error()})
		}
	})

	r.GET("/userInfo", func(c *gin.Context) {
		user := jwe.TokenUser{}
		token := c.GetHeader("READYGOTOKEN")
		userid, err := jwe.GetInfoFromToken(token, &user)
		if err == nil {
			c.JSON(200, gin.H{"result": "OK", "userid": userid, "extInfo": user})
		} else {
			c.JSON(500, gin.H{"result": "Error", "msg": err.Error()})
		}
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
