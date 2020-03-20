package wxapi

import (
	"gitee.com/chunanyong/gowe"
	"gitee.com/chunanyong/zorm"
	"github.com/gin-gonic/gin"
	"readygo/ginext/Ginserializer"
	"readygo/wx/wxstruct"
	"strings"
)

var WXPay *wxstruct.WxPayConfig

func init()  {

	WXPay = &wxstruct.WxPayConfig{
		Id:"test",
		AppId:"xxx",
		Secret:"xxx",
		MchID:"xxx",
		Key:"xxxx",
	}


}

func WxPayUnifiedOrder(c *gin.Context)  {

	openid := c.Query("openid")

	body := &gowe.WxPayUnifiedOrderBody{
		Body:"人参果",
		OutTradeNo:strings.Replace(zorm.GenerateStringID(),"-","",-1),// zorm.GenerateStringID(),
		TotalFee: gowe.ServiceTypeNormalDomestic,
		SpbillCreateIP:"127.0.0.1",
		NotifyUrl:"http://www.qq.com",
		OpenId:openid,
		TradeType:gowe.TradeTypeApplet,
	}




	order, err := gowe.WxPayUnifiedOrder(WXPay, body)

	if err != nil {
		c.JSON(505, Ginserializer.Response{
			Code: 1,
			Msg:  err.Error(),
		})
	}else{
		c.JSON(200, Ginserializer.Response{
			Code: 0,
			Data:  order,
		})
	}

}


