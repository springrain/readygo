package wxstruct




type WxConfig struct {
	Id string
	AppId string
	AccessToken string
	Secret string


}

func (w WxConfig) GetId() string {
	return w.Id
}

func (w WxConfig) GetAppId() string {
	return w.AppId
}

func (w WxConfig) GetAccessToken() string {
	return w.AccessToken
}

func (w WxConfig) GetSecret() string {
	return w.Secret
}
