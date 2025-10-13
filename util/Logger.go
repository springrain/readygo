package util

import (
	"context"
	"io"
	"os"

	"gitee.com/chunanyong/zorm"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// levelError 日志级别,使用变量,提升日志级别的优先级
var levelError = func() hlog.Level {
	hlog.SetLevel(hlog.LevelError)
	return hlog.LevelError
}()

// InitLog 初始化日志文件
func InitLog() *os.File {
	f, err := os.OpenFile("./readygo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	// https://github.com/cloudwego/hertz/issues/292
	//defer f.Close()
	fileWriter := io.MultiWriter(f, os.Stdout)
	hlog.SetOutput(fileWriter)
	hlog.SetSilentMode(true)
	//使用变量,提升日志级别的优先级
	hlog.SetLevel(levelError)
	zorm.FuncLogError = FuncLogError
	zorm.FuncLogPanic = FuncLogPanic
	zorm.FuncPrintSQL = FuncPrintSQL
	return f
}

// LogCallDepth 记录日志调用层级,用于定位到业务层代码
// Log Call Depth Record the log call level, used to locate the business layer code
var LogCallDepth = 4

// FuncLogError 记录error日志.NewDBDao方法里的异常,ctx为nil,扩展时请注意
// FuncLogError Record error log
var FuncLogError func(ctx context.Context, err error) = defaultLogError

// FuncLogPanic  记录panic日志,默认使用"defaultLogError"实现
// FuncLogPanic Record panic log, using "defaultLogError" by default
var FuncLogPanic func(ctx context.Context, err error) = defaultLogPanic

// FuncPrintSQL 打印sql语句,参数和执行时间,小于0是禁用日志输出;等于0是只输出日志,不计算SQ执行时间;大于0是计算执行时间,并且大于指定值
// FuncPrintSQL Print sql statement and parameters
var FuncPrintSQL func(ctx context.Context, sqlstr string, args []interface{}, execSQLMillis int64) = defaultPrintSQL

func defaultLogError(ctx context.Context, err error) {
	hlog.Error(err)
}

func defaultLogPanic(ctx context.Context, err error) {
	defaultLogError(ctx, err)
}

func defaultPrintSQL(ctx context.Context, sqlstr string, args []interface{}, execSQLMillis int64) {
	if args != nil {

		hlog.Errorf("sql:", sqlstr, ",args:", args, ",execSQLMillis:", execSQLMillis)
		//log.Output(LogCallDepth, fmt.Sprintln("sql:", sqlstr, ",args:", args, ",execSQLMillis:", execSQLMillis))
	} else {
		hlog.Errorf("sql:", sqlstr, ",args: [] ", ",execSQLMillis:", execSQLMillis)
		//log.Output(LogCallDepth, fmt.Sprintln("sql:", sqlstr, ",args: [] ", ",execSQLMillis:", execSQLMillis))
	}
}
