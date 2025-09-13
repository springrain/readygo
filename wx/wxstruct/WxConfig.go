package wxstruct

import "context"

type WxConfig struct {
	Id          string
	AppId       string
	AccessToken string
	Secret      string
}

func (w WxConfig) GetId(ctx context.Context) string {
	return w.Id
}

func (w WxConfig) GetAppId(ctx context.Context) string {
	return w.AppId
}

func (w WxConfig) GetAccessToken(ctx context.Context) string {
	return w.AccessToken
}

func (w WxConfig) GetSecret(ctx context.Context) string {
	return w.Secret
}
