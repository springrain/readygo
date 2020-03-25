package wxapi

import (
	"fmt"
	"gitee.com/chunanyong/gowe"
	"gitee.com/chunanyong/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"readygo/ginext/Ginserializer"
	"readygo/wx/wxstruct"
)



var WX *wxstruct.WxConfig

func init()  {
	godotenv.Load()
	WX = &wxstruct.WxConfig{
 		AppId:os.Getenv("APPID"),
		Secret:os.Getenv("SECRET"),
	}
}


// Ping 状态检查页面
func Ping(c *gin.Context) {
	c.JSON(200, Ginserializer.Response{
		Code: 0,
		Msg:  "Pong",
	})
}

func WxMaSubscribeMessageSend(c *gin.Context){


	token, _ := gowe.GetAccessToken(WX)
	WX.AccessToken = token.AccessToken

	var params  gowe.WxMaTemplateMsgSendBody
	params.Page = "index"
	params.Touser =os.Getenv("OPENID")
	params.MiniprogramState = "developer"
	params.Lang="zh_CN"
	params.TemplateId ="xszovj1cMuDhp-SClpKYdnt5flFjtB4TN-mt3CrhrFE"
	params.AddData("thing1","恭喜您下单获得一笔新的奖励")
	params.AddData("amount2","1.2元")
	params.AddData("thing3","每月20日统一结算")
	params.AddData("thing4","请查看详细信息")

	send, err := gowe.WxMaSubscribeMessageSend(WX, &params)

	if err != nil{
		fmt.Println(err.Error())
	}

	fmt.Println(send)

	c.JSON(200, Ginserializer.Response{
		Code: 0,
		Data: send,
	})
}

func WxMaCode2Session(c *gin.Context)  {

	code := c.Query("jsCode")
	logger.Info(code)



	session, err := gowe.WxMaCode2Session(WX, code)
	logger.Info(session.OpenId)

	if err != nil {
		c.JSON(505, Ginserializer.Response{
			Code: 1,
			Msg:  err.Error(),
		})
	}else{
		fmt.Println(session)

		c.JSON(200, Ginserializer.Response{
			Code: 0,
			Msg: session.OpenId,
		})
	}


}

