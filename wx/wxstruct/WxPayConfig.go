package wxstruct

type WxPayConfig struct {
	Id     string
	AppId  string
	Secret string
	MchID  string
	Key    string
}

func (w WxPayConfig) GetId() string {
	return w.Id
}

func (w WxPayConfig) GetAppId() string {
	return w.AppId
}

func (w WxPayConfig) GetAccessToken() string {
	panic("implement me")
}

func (w WxPayConfig) GetSecret() string {
	return w.Secret
}

func (w WxPayConfig) GetCertificateFile() string {
	return "../cert/apiclient_cert.pem"
}

func (w WxPayConfig) GetMchId() string {
	return w.MchID
}

func (w WxPayConfig) GetSubAppId() string {
	panic("implement me")
}

func (w WxPayConfig) GetSubMchId() string {
	panic("implement me")
}

func (w WxPayConfig) GetAPIKey() string {
	return w.Key
}

func (w WxPayConfig) GetNotifyUrl() string {
	return ""
}

func (w WxPayConfig) GetSignType() string {
	return "MD5"
}

func (w WxPayConfig) GetServiceType() int {
	return 1
}

func (w WxPayConfig) IsProd() bool {
	return true
}

func (w WxPayConfig) IsMch() bool {
	return false
}
