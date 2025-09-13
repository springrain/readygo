package wxstruct

import "context"

type WxPayConfig struct {
	Id     string
	AppId  string
	Secret string
	MchID  string
	Key    string
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

func (w WxPayConfig) GetMchId(ctx context.Context) string {
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

func (w WxPayConfig) GetNotifyUrl(ctx context.Context) string {
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
