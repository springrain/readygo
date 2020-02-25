package main

import (
	"net/http"
	"readygo/ginext"
	"readygo/orm"
	"readygo/permission/permhandler"

	"github.com/gin-gonic/gin"
)

//初始化BaseDao
func init() {
	baseDaoConfig := orm.DataSourceConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "readygo",
		UserName: "root",
		PassWord: "root",
		DBType:   "mysql",
	}
	_, _ = orm.NewBaseDao(&baseDaoConfig)
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
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
