package wxapi

import (
	"context"
	"fmt"
	"os"

	"readygo/webext"
	"readygo/wx/wxstruct"

	"gitee.com/chunanyong/gowe"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/joho/godotenv"
)

var WX *wxstruct.WxConfig

func init() {
	godotenv.Load()
	WX = &wxstruct.WxConfig{
		AppId:  os.Getenv("APPID"),
		Secret: os.Getenv("SECRET"),
	}
}

// Ping 状态检查页面
func Ping(ctx context.Context, c *app.RequestContext) {
	c.JSON(200, webext.ResponseData{
		StatusCode: 0,
		Message:    "Pong",
	})
}

// 订阅消息
func WxMaSubscribeMessageSend(ctx context.Context, c *app.RequestContext) {
	token, _ := gowe.GetAccessToken(ctx, WX)
	WX.AccessToken = token.AccessToken

	var params gowe.WxMaTemplateMsgSendBody
	params.Page = "index"
	params.Touser = os.Getenv("OPENID")
	params.MiniprogramState = "developer"
	params.Lang = "zh_CN"
	params.TemplateId = "xszovj1cMuDhp-SClpKYdnt5flFjtB4TN-mt3CrhrFE"
	params.AddData("thing1", "恭喜您下单获得一笔新的奖励")
	params.AddData("amount2", "1.2元")
	params.AddData("thing3", "每月20日统一结算")
	params.AddData("thing4", "请查看详细信息")

	send, err := gowe.WxMaSubscribeMessageSend(ctx, WX, &params)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(send)
	c.JSON(200, webext.ResponseData{
		StatusCode: 0,
		Data:       send,
	})
}

// 登录凭证校验
func WxMaCode2Session(ctx context.Context, c *app.RequestContext) {
	code := c.Query("jsCode")
	hlog.CtxInfof(ctx, code)

	session, err := gowe.WxMaCode2Session(ctx, WX, code)
	hlog.CtxInfof(ctx, session.OpenId)

	if err != nil {
		c.JSON(505, webext.ResponseData{
			StatusCode: 1,
			Message:    err.Error(),
		})
	} else {
		fmt.Println(session)

		c.JSON(200, webext.ResponseData{
			StatusCode: 0,
			Message:    session.OpenId,
		})
	}
}
