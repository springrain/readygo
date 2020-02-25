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
		//放行所有OPTIONS方法
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		//请求的uri
		uri := c.Request.RequestURI
		logger.Info(uri)
		c.Next()
	}

}
