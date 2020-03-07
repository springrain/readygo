package permapi

import (
	"readygo/ginext/Ginserializer"

	"github.com/gin-gonic/gin"
)

// Ping 状态检查页面
func Ping(c *gin.Context) {
	c.JSON(200, Ginserializer.Response{
		Code: 0,
		Msg:  "Pong",
	})
}
