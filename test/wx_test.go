package test

import (
	"context"
	"os"
	"testing"

	"gitee.com/chunanyong/gowe"
)

type wxconfig struct {
	id          string
	appid       string
	accesstoken string
	secret      string
}

var wx = &wxconfig{
	appid:  os.Getenv("APPID"),
	secret: os.Getenv("SECRET"),
}

func (w wxconfig) GetId(ctx context.Context) string {
	return w.id
}

func (w wxconfig) GetAppId(ctx context.Context) string {
	return w.appid
}

func (w wxconfig) GetAccessToken(ctx context.Context) string {
	return w.accesstoken
}

func (w wxconfig) GetSecret(ctx context.Context) string {
	return w.secret
}

func TestGetAccessToken(t *testing.T) {
	token, err := gowe.GetAccessToken(ctx, wx)
	if err != nil {
		t.Log("error:", err)
	}

	t.Log("token:", token)
}

func TestGetJsTicket(t *testing.T) {
	token, err := gowe.GetAccessToken(ctx, wx)
	if err != nil {
		t.Log("error:", err)
	}

	wx.accesstoken = token.AccessToken

	ticket, err := gowe.GetJsTicket(ctx, wx)
	if err != nil {
		t.Log("error:", err)
	}

	t.Log("ticket:", ticket)
}

func TestGetCardTicket(t *testing.T) {
	token, err := gowe.GetAccessToken(ctx, wx)
	if err != nil {
		t.Log("error:", err)
	}

	wx.accesstoken = token.AccessToken

	ticket, err := gowe.GetCardTicket(ctx, wx)
	if err != nil {
		t.Log("error:", err)
	}

	t.Log("ticket:", ticket)
}
