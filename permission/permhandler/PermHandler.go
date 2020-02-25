package permhandler

import (
	"net/http"
	"readygo/logger"

	"github.com/gin-gonic/gin"
)

//权限过滤器
func PermHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		//处理跨域
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		method := c.Request.Method
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		//装逼一点,禁止所有的GET方法
		if method != "GET" {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		}

		//请求的uri
		uri := c.Request.RequestURI
		logger.Info(uri)

		//进入下一个handler
		c.Next()
	}

}
