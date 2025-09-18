// permission/handler/login_handler.go
package permhandler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"readygo/permission/permservice"
	"readygo/permission/permstruct"
)

// LoginHandler 是 Hertz 的 HTTP 处理函数
func LoginHandler(ctx context.Context, c *app.RequestContext) {
    // 1. 定义请求对象
    var req permstruct.UserStruct

    // 2. 解析 JSON 请求体
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(consts.StatusBadRequest, map[string]any{
            "error": "参数错误: " + err.Error(),
        })
        return
    }

    // 3. 调用业务函数 Login
    resp, err := permservice.Login(ctx, &req)
    if err != nil {
        c.JSON(consts.StatusUnauthorized, map[string]any{
            "error": err.Error(),
        })
        return
    }

    // 4. 返回成功响应
    c.JSON(consts.StatusOK, resp)
}