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
	"context"
	"net/http"
	"os"
	"readygo/config"

	"readygo/cache"
	"readygo/permission/permhandler"
	"readygo/permission/permroute"
	"readygo/permission/permstruct"
	"readygo/webext"
	"readygo/wx/wxroute"

	"gitee.com/chunanyong/zorm"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	_ "github.com/go-sql-driver/mysql"
)

// 初始化
func init() {
	// 初始化DBDao
	dbDaoConfig := zorm.DataSourceConfig{
		DSN: config.Cfg.Database.DSN,
		// DriverName 数据库驱动名称:mysql,postgres,oci8,sqlserver,sqlite3,go_ibm_db,clickhouse,dm,kingbase,aci,taosSql|taosRestful 和Dialect对应
		DriverName: config.Cfg.Database.DriverName,
		// Dialect 数据库方言:mysql,postgresql,oracle,mssql,sqlite,db2,clickhouse,dm,kingbase,shentong,tdengine 和 DriverName 对应
		Dialect: config.Cfg.Database.Dialect,
		// MaxOpenConns 数据库最大连接数 默认50
		MaxOpenConns: config.Cfg.Database.MaxOpenConns,
		// MaxIdleConns 数据库最大空闲连接数 默认50
		MaxIdleConns: config.Cfg.Database.MaxIdleConns,
		// ConnMaxLifetimeSecond 连接存活秒时间. 默认600(10分钟)后连接被销毁重建.避免数据库主动断开连接,造成死连接.MySQL默认wait_timeout 28800秒(8小时)
		ConnMaxLifetimeSecond: config.Cfg.Database.ConnMaxLifetimeSecond,
		// SlowSQLMillis 慢sql的时间阈值,单位毫秒.小于0是禁用SQL语句输出;等于0是只输出SQL语句,不计算执行时间;大于0是计算SQL执行时间,并且>=SlowSQLMillis值
		SlowSQLMillis: config.Cfg.Database.SlowSQLMillis,
	}
	_, _ = zorm.NewDBDao(&dbDaoConfig)
	
	// 初始化redis
	cache.NewRedisClient(context.Background(), &cache.RedisConfig{
		Addr: config.Cfg.Redis.Addr,
		Password: config.Cfg.Redis.Password,
		PoolSize: config.Cfg.Redis.PoolSize,
		MinIdleConns: config.Cfg.Redis.MinIdleConns,
	})

	cache.NewMemeryCacheManager()
	// 初始化initWebEngine
	initWebEngine()
}

// initWebEngine 初始化Web引擎
func initWebEngine() {
	// 获取引擎
	h := webext.WebEngine()
	// 设置前缀,需要在路由初始化前调用
	//webext.SetContextPath("/readygo/")

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	h.Use(webext.WebLogger())
	// r.Use(gin.Logger())

	f, err := os.OpenFile("./logs/readygo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	hlog.SetOutput(f)

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	// h.Use(webext.WebRecovery())

	// css js等静态文件
	h.Static("/assets", "./assets")
	h.StaticFS("/more_static", &app.FS{Root: "my_file_system", GenerateIndexPages: true})
	h.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// --------------------------
    // 第一步：先注册不需要权限的路由
    // --------------------------
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, utils.H{"hello": "world"})
	})

	h.GET("/login", permhandler.LoginHandler)
	//h.POST("/Captcha", permapi.Captcha)

	h.GET("/system/menu/tree", func(ctx context.Context, c *app.RequestContext) {
		user, err := permstruct.GetCurrentUserFromContext(ctx)
		// token := c.GetHeader(JWTTokenName)
		// userid, err := permutil.GetInfoFromToken(token, &user)
		if err == nil {
			c.JSON(http.StatusOK, webext.ResponseData{
				StatusCode: 0,
				Message:    "",
				Data:       utils.H{"userid": user.UserId, "extInfo": user},
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, webext.ResponseData{
				StatusCode: 1,
				Message:    err.Error(),
				Data:       utils.H{"msg": err.Error()},
			})
		}
	})

    // --------------------------
    // 第二步：再注册权限中间件
    // --------------------------
	h.Use(permhandler.PermHandler())

	// --------------------------
    // 第三步：最后注册需要权限的路由
    // --------------------------
	permroute.RegisterPermRoute(h)
	wxroute.RegisterWXRoute(h)
}

func main() {
	h := webext.WebEngine()
	h.Spin()
}
