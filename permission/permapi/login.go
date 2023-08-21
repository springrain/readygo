package permapi

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// Login 登录
func Login(ctx context.Context, c *app.RequestContext) {
}

/*
// Captcha 获取验证码
func Captcha(ctx context.Context, c *app.RequestContext) {
	// 获取base64的图片验证码
	captchaKey, base64, err := newCaptchaKeyB64()
	if err != nil {
		return
	}
	data := make(map[string]string)
	data["captchaKey"] = captchaKey
	data["captchaImg"] = base64
	res := webext.SuccessReponseData(data, "")

	c.JSON(200, res)
}


// 生成验证码的key和base64的值
func newCaptchaKeyB64() (string, string, error) {
	// 获取一个uuid字符串
	key := zorm.FuncGenerateStringID(nil)

	// 10秒为计算单位,计算当前值
	tenSeconds := time.Now().Unix() / 10
	str := key + strconv.FormatInt(tenSeconds, 10)
	md5str := util.GenerateMD5(str)
	castr := md5str[:4]
	base64, err := captcha.NewCaptchaB64string(60, 45, castr)
	return key, base64, err
}

// verifyCaptcha 校验验证码是否准确
func verifyCaptcha(captchaKey string, captchaValue string) bool {
	// 10秒为计算单位,计算当前值
	tenSeconds := time.Now().Unix() / 10
	// 当前的10秒
	str := captchaKey + strconv.FormatInt(tenSeconds, 10)
	md5str := util.GenerateMD5(str)
	castr := md5str[:4]
	if castr == captchaValue {
		return true
	}
	// 上一个10秒
	str = captchaKey + strconv.FormatInt(tenSeconds-1, 10)
	md5str = util.GenerateMD5(str)
	castr = md5str[:4]
	if castr == captchaValue {
		return true
	}
	return false
}
*/
