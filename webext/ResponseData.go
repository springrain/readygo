package webext

import (
	"readygo/config"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"gitee.com/chunanyong/zorm"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// Creates a router without any middleware by default
var webEngine = server.Default(server.WithHostPorts(":" + strconv.Itoa(config.Cfg.Server.Port)), server.WithBasePath(config.Cfg.Server.BasePath + "/"))

func WebEngine() *server.Hertz {
	return webEngine
}

// SetContextPath 设置项目名前缀contextPath,因为gin暂时不支持直接修改RouterGroup的basePath,使用unsafe.Pointer修改
// 需要在路由初始化前调用
func SetContextPath(contextPath string) {
	if contextPath == "" {
		return
	}
	if !strings.HasSuffix(contextPath, "/") {
		contextPath = contextPath + "/"
	}
	// 获取引擎
	h := WebEngine()
	// 因为Engine匿名注入了RouterGroup,所以直接获取Engine的反射对象
	engine := reflect.ValueOf(h).Elem()
	// 获取RouterGroup的basePath属性反射值对象
	basePath := engine.FieldByName("basePath")
	// 获取basePath的UnsafeAddr
	p := unsafe.Pointer(basePath.UnsafeAddr())
	// 重新赋值basePath的反射值,NewAt默认返回的是指针,使用Elem获取反射值对象
	basePath = reflect.NewAt(basePath.Type(), p).Elem()
	// 设置反射值
	basePath.Set(reflect.ValueOf(contextPath))
}

// 三位数错误编码为复用http原本含义
// 五位数错误编码为应用自定义错误
// 五开头的五位数错误编码为服务器端错误，比如数据库操作失败
// 四开头的五位数错误编码为客户端错误，有时候是客户端代码写错了，有时候是用户操作错误
const (
	// CodeCheckLogin 未登录
	CodeCheckLogin = 401
	// CodeNoRightErr 未授权访问
	CodeNoRightErr = 403
	// CodeDBError 数据库操作失败
	CodeDBError = 50001
	// CodeEncryptError 加密失败
	CodeEncryptError = 50002
	// CodeParamErr 各种奇奇怪怪的参数错误
	CodeParamErr = 40001
)

// ErrorResult 响应错误
type ErrorResult struct {
	Error ErrorItem `json:"error"` // 错误项
}

// ErrorItem 响应错误项
type ErrorItem struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
}

// ResponseData 返回数据包装器
type ResponseData struct {
	// 业务状态代码 // 异常 1, 成功 0,默认成功0,业务代码见说明
	StatusCode int `json:"statusCode"`
	// HttpCode http的状态码
	// HttpCode int `json:"httpCode,omitempty"`
	// 返回数据
	Data interface{} `json:"data,omitempty"`
	// 返回的信息内容,配合StatusCode
	Message string `json:"message,omitempty"`
	// 扩展的map,用于处理返回多个值的情况
	ExtMap map[string]interface{} `json:"extMap,omitempty"`
	// 列表的分页对象
	Page zorm.Page `json:"page,omitempty"`
	// 查询条件的struct回传
	QueryStruct interface{} `json:"queryStruct,omitempty"`
	ERR         error       // 响应错误
}

// TrackedErrorResponse 有追踪信息的错误响应
type TrackedErrorResponse struct {
	ResponseData
	TrackID string `json:"track_id"`
}

// CheckLogin 检查登录
func CheckLogin() ResponseData {
	return ResponseData{
		StatusCode: CodeCheckLogin,
		Message:    "未登录",
	}
}

// SuccessReponseData 通用成功处理
func SuccessReponseData(data interface{}, message string) ResponseData {
	res := ResponseData{
		Data:       data,
		StatusCode: 200,
		Message:    message,
	}
	return res
}

/*
	ErrorReponseData 通用错误处理
	err 错误信息，可传可不传
*/
func ErrorReponseData(errCode int, message string, err ...error) ResponseData {
	res := ResponseData{
		StatusCode: errCode,
		Message:    message,
	}
	/*
		// 生产环境隐藏底层报错
		if err != nil && gin.Mode() != gin.ReleaseMode {
			res.Message = err.Error()
		}
	*/
	return res
}

// DBErr 数据库操作失败
func DBErr(message string, err error) ResponseData {
	if message == "" {
		message = "数据库操作失败"
	}
	return ErrorReponseData(CodeDBError, message, err)
}

// ParamErr 各种参数错误
func ParamErr(message string, err error) ResponseData {
	if message == "" {
		message = "参数错误"
	}
	return ErrorReponseData(CodeParamErr, message, err)
}
