package permapi

import (
	"context"
	"readygo/webext"

	"github.com/cloudwego/hertz/pkg/app"
)

// Ping 状态检查页面
func Ping(ctx context.Context, c *app.RequestContext) {
	c.JSON(200, webext.ResponseData{
		StatusCode: 0,
		Message:    "Pong",
	})
}
