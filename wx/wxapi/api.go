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

