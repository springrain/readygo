package wxstruct

import "context"

type WxPayConfig struct {
	Id                   string
	AppId                string
	Secret               string
	MchID                string
	MchAPIv3Key          string
	Key                  string
	CertSerialNo         string
	PrivateKey           string
	WechatPayCertificate string
}

func (w WxPayConfig) GetId(ctx context.Context) string {
	return w.Id
}

func (w WxPayConfig) GetAppId(ctx context.Context) string {
	return w.AppId
}

func (w WxPayConfig) GetAccessToken(ctx context.Context) string {
	panic("implement me")
}

func (w WxPayConfig) GetSecret(ctx context.Context) string {
	return w.Secret
}

func (w WxPayConfig) GetCertificateFile(ctx context.Context) string {
	return "../cert/apiclient_cert.pem"
}

// 商户证书序列号
func (w WxPayConfig) GetCertSerialNo(ctx context.Context) string {
	return w.CertSerialNo
}

// 商户API v3密钥 (32字节)
func (w WxPayConfig) GetMchAPIv3Key(ctx context.Context) string {
	return w.MchAPIv3Key
}

func (w WxPayConfig) GetMchID(ctx context.Context) string {
	return w.MchID
}

func (w WxPayConfig) GetSubAppId(ctx context.Context) string {
	panic("implement me")
}

func (w WxPayConfig) GetSubMchId(ctx context.Context) string {
	panic("implement me")
}

func (w WxPayConfig) GetAPIKey(ctx context.Context) string {
	return w.Key
}

// 商户API私钥
func (w WxPayConfig) GetPrivateKey(ctx context.Context) string {
	return w.PrivateKey
}

// 微信支付平台证书（示例，实际应从微信平台下载并定期更新）
func (w WxPayConfig) GetWechatPayCertificate(ctx context.Context) string {
	return w.WechatPayCertificate

}

func (w WxPayConfig) GetNotifyURL(ctx context.Context) string {
	return ""
}

func (w WxPayConfig) GetSignType(ctx context.Context) string {
	return "MD5"
}

func (w WxPayConfig) GetServiceType(ctx context.Context) int {
	return 1
}

func (w WxPayConfig) IsProd(ctx context.Context) bool {
	return true
}

func (w WxPayConfig) IsMch(ctx context.Context) bool {
	return false
}
