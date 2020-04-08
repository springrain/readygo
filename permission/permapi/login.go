package permapi

import (
	"fmt"
	"readygo/util/captcha"

	"github.com/gin-gonic/gin"
)

// Login 登录
func Login(c *gin.Context) {

}

//Captcha 获取验证码
func Captcha(c *gin.Context) {
	//获取base64的图片验证码
	base64, err := captcha.NewCaptchaB64string(60, 45, "ab34")
	if err != nil {
		return
	}
	fmt.Println(base64)
}
