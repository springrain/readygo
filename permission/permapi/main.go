package permapi

import (
	"github.com/gin-gonic/gin"
	"readygo/serializer"
)

// Ping 状态检查页面
func Ping(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "Pong",
	})
}
