package wxapi

import (
	"gitee.com/chunanyong/gowe"
	"gitee.com/chunanyong/zorm"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"readygo/ginext/Ginserializer"
	"readygo/wx/wxstruct"
	"strings"
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

func WxPayNotifyPay(c *gin.Context){

}

func WxPayUnifiedOrder(c *gin.Context) {

	openid := c.Query("openid")

	body := &gowe.WxPayUnifiedOrderBody{
		Body:           "人参果",
		OutTradeNo:     strings.Replace(zorm.GenerateStringID(), "-", "", -1), // zorm.GenerateStringID(),
		TotalFee:       gowe.ServiceTypeNormalDomestic,

		SpbillCreateIP: "127.0.0.1",
		NotifyUrl:      "http://www.qq.com",
		OpenId:         openid,
		TradeType:      gowe.TradeTypeMiniApp,
	}

	order, err := gowe.WxPayUnifiedOrder(WXPay, body)

	if err != nil {
		c.JSON(505, Ginserializer.Response{
			Code: 1,
			Msg:  err.Error(),
		})
	} else {
		c.JSON(200, Ginserializer.Response{
			Code: 0,
			Data: order,
		})
	}

}
