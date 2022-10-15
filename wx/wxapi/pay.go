package wxapi

import (
	"context"
	"fmt"
	"os"
	"readygo/webext"
	"readygo/wx/wxstruct"
	"strings"

	"gitee.com/chunanyong/gowe"
	"gitee.com/chunanyong/zorm"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/joho/godotenv"
)

var WXPay *wxstruct.WxPayConfig

func init() {
	godotenv.Load()

	WXPay = &wxstruct.WxPayConfig{
		AppId:  os.Getenv("WXPayAppId"),
		Secret: os.Getenv("WXPaySecret"),
		MchID:  os.Getenv("WXPayMchID"),
		Key:    os.Getenv("WXPayKey"),
	}

}

//支付结果通知
func WxPayNotifyPay(ctx context.Context, c *app.RequestContext) {
	//
	//var body gowe.WxPayNotifyPayBody
	//c.Bind(&body)

	body := c.Request.Body()

	gowe.WxPayNotifyPay(WXPay, body, func(wxPayNotifyPayBody gowe.WxPayNotifyPayBody) error {

		fmt.Println(wxPayNotifyPayBody)

		return nil
	})

}

//统一下单
func WxPayAppSign(ctx context.Context, c *app.RequestContext) {

	body := make(map[string]string, 0)
	c.Bind(&body)

	fmt.Println(body)

	paySign := gowe.WxPayMaSign(WXPay.GetAppId(), body["nonceStr"], body["packages"], body["signType"], body["timeStamp"], WXPay.GetAPIKey())

	body["paySign"] = paySign

	c.JSON(200, webext.ResponseData{
		StatusCode: 0,
		Data:       body,
	})
}

func WxPayUnifiedOrder(ctx context.Context, c *app.RequestContext) {

	openid := c.Query("openid")

	body := &gowe.WxPayUnifiedOrderBody{
		Body:       "人参果",
		OutTradeNo: strings.Replace(zorm.FuncGenerateStringID(ctx), "-", "", -1),
		TotalFee:   gowe.ServiceTypeNormalDomestic,

		SpbillCreateIP: "127.0.0.1",
		NotifyUrl:      "http://www.qq.com",
		OpenId:         openid,
		TradeType:      gowe.TradeTypeMiniApp,
	}

	order, err := gowe.WxPayUnifiedOrder(WXPay, body)

	if err != nil {
		c.JSON(505, webext.ResponseData{
			StatusCode: 1,
			Message:    err.Error(),
		})
	} else {
		c.JSON(200, webext.ResponseData{
			StatusCode: 0,
			Data:       order,
		})
	}

}
