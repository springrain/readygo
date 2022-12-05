package webext

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gitee.com/chunanyong/zorm"
	"github.com/cloudwego/hertz/pkg/app"
)

// 定义上下文中的键
const (
	prefix     = "fxs"
	ReqBodyKey = prefix + "/req-body"
	ResBodyKey = prefix + "/res-body"
)

// GetToken 获取用户令牌
func GetToken(h *app.RequestContext) string {
	var token string
	auth := string(h.GetHeader("Authorization"))
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

// ParseParamID Parse path id
func ParseParamID(h *app.RequestContext, key string) uint64 {
	val := h.Param(key)
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

// GetBody Get request body
func GetBody(h *app.RequestContext) []byte {
	if v, ok := h.Get(ReqBodyKey); ok {
		if b, ok := v.([]byte); ok {
			return b
		}
	}
	return nil
}

// ParseJSON 解析请求JSON
func ParseJSON(h *app.RequestContext, obj interface{}) error {
	if err := h.Bind(obj); err != nil {
		return err
	}
	return nil
}

// ParseQuery 解析Query参数
func ParseQuery(h *app.RequestContext, obj interface{}) error {
	if err := h.Bind(obj); err != nil {
		return err
	}
	return nil
}

// ParseForm 解析Form请求
func ParseForm(h *app.RequestContext, obj interface{}) error {
	if err := h.Bind(obj); err != nil {
		return err
	}
	return nil
}

// ResList 响应列表数据
func ResList(h *app.RequestContext, v interface{}) {
	ResSuccess(h, ResponseData{Data: v})
}

// ResPage 响应分页数据
func ResPage(h *app.RequestContext, v interface{}, pr zorm.Page) {
	list := ResponseData{
		Data: v,
		Page: pr,
	}
	ResSuccess(h, list)
}

// ResSuccess 响应成功
func ResSuccess(h *app.RequestContext, v interface{}) {
	ResJSON(h, http.StatusOK, v)
}

// ResJSON 响应JSON数据
func ResJSON(h *app.RequestContext, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	h.Set(ResBodyKey, buf)
	h.Data(status, "application/json; charset=utf-8", buf)
	h.Abort()
}

func ResError(h *app.RequestContext, err error) {
	ResJSON(h, 0, err.Error())
}
