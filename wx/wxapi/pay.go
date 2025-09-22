package wxapi

import (
	"context"
	"fmt"
	"os"
	"strings"

	"readygo/webext"
	"readygo/wx/wxstruct"

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

// 支付结果通知
func WxPayNotifyPay(ctx context.Context, c *app.RequestContext) {
	// 1. 验证签名
	fmt.Println("微信回调请求到达----------")

	// 从请求头获取签名相关字段
	timestamp := string(c.Request.Header.Peek("Wechatpay-Timestamp"))
	nonce := string(c.Request.Header.Peek("Wechatpay-Nonce"))
	signature := string(c.Request.Header.Peek("Wechatpay-Signature"))
	serial := string(c.Request.Header.Peek("Wechatpay-Serial"))

	if timestamp == "" || nonce == "" || signature == "" || serial == "" {
		fmt.Println("缺少签名头----------")
		c.JSON(400, map[string]string{"message": "缺少签名头"})
		return
	}

	body, err := c.Body()
	if err != nil {
		fmt.Println("读取请求体失败----------", err.Error())
		c.JSON(400, map[string]string{"message": "读取请求体失败"})
		return
	}
	// 通常不需要再手动调用 SetBody，但如果你后续处理需要，可以设置
	c.Request.SetBody(body)
	//验签
	err = gowe.VerifyWechatSignature(ctx, WXPay, timestamp, nonce, signature, serial, body)
	if err != nil {
		fmt.Printf("签名验证失败: %v", err)
		c.JSON(400, map[string]string{"message": fmt.Sprintf("签名验证失败: %v", err)})
		return
	}

	//回调
	callback := gowe.WechatPayCallback(ctx, WXPay, body)
	if callback.Code == 0 {
		fmt.Println("微信回调成功----------")
		// 7. 返回成功
		c.JSON(200, map[string]string{"message": "Success"})
	} else {
		fmt.Println("微信回调失败----------")
		// 7. 返回成功
		c.JSON(400, map[string]string{"message": "error"})
	}
}

// 统一下单
func WxPayAppSign(ctx context.Context, c *app.RequestContext) {
	body := make(map[string]string, 0)
	c.Bind(&body)

	fmt.Println(body)

	paySign := gowe.WxPayMaSign(WXPay.GetAppId(ctx), body["nonceStr"], body["packages"], body["signType"], body["timeStamp"], WXPay.GetAPIKey(ctx))

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

	order, err := gowe.WxPayUnifiedOrder(ctx, WXPay, body)

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
