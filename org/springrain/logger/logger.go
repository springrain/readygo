package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//  logger
var logger *zap.Logger

const appName = "goshop"

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

//获取日志级别
func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

//初始化日志
func init() {
	fileName := "./logs/" + appName + ".log"
	level := getLoggerLevel("debug")
	hook := lumberjack.Logger{
		Filename:   fileName, // 日志文件路径
		MaxSize:    128,      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,       // 日志文件最多保存多少个备份
		MaxAge:     7,        // 文件最多保存多少天
		Compress:   true,     // 是否压缩
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("app", appName))
	//因为对zap做了封装,在记录日志时,跳过本次封装,这样才能记录到真实的异常信息文件和行数,不然一直显示是本文件的方法行数
	skip := zap.AddCallerSkip(1)
	// 构造日志
	logger = zap.New(core, caller, development, filed, skip)

	//logger = logger.Sugar()
}

//隔离引用
type LogField = zap.Field

func Debug(msg string, fields ...LogField) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...LogField) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...LogField) {
	logger.Warn(msg, fields...)
}

func Error(err error, fields ...LogField) {

	logger.Error(err.Error(), fields...)
}

//func DPanic(msg string, fields ...LogField) {
//	logger.DPanic(msg, fields...)
//}

func Panic(msg error, fields ...LogField) {
	logger.Panic(msg.Error(), fields...)
}

func Fatal(msg string, fields ...LogField) {
	logger.Fatal(msg, fields...)
}

func String(key string, val string) LogField {
	return zap.String(key, val)
}

func Int(key string, val int) LogField {
	return zap.Int(key, val)
}
func Int32(key string, val int32) LogField {
	return zap.Int32(key, val)
}
func Int64(key string, val int64) LogField {
	return zap.Int64(key, val)
}
func Bool(key string, val bool) LogField {
	return zap.Bool(key, val)
}

func Duration(key string, val time.Duration) LogField {
	return zap.Duration(key, val)
}
func Time(key string, val time.Time) LogField {
	return zap.Time(key, val)
}
func Float32(key string, val float32) LogField {
	return zap.Float32(key, val)
}
func Float64(key string, val float64) LogField {

	return zap.Float64(key, val)
}

func ByteString(key string, val []byte) LogField {

	return zap.ByteString(key, val)
}
