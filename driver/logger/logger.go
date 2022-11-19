package logger

import (
	"context"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	// 设置日志格式为json格式
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	logrus.SetOutput(os.Stdout)

	// 设置日志级别为warn以上
	logrus.SetLevel(logrus.WarnLevel)
}

func withLogFields(ctx context.Context, kv map[string]interface{}) logrus.Fields {
	var result logrus.Fields = map[string]interface{}{}
	for j := range kv {
		result[j] = kv[j]
	}
	pc, file, line, _ := runtime.Caller(2)
	method := runtime.FuncForPC(pc).Name()
	result["file"] = file
	result["line"] = line
	result["method"] = method
	return result
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logrus.WithFields(withLogFields(ctx, nil)).Infof(format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithFields(withLogFields(ctx, nil)).Warnf(format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithFields(withLogFields(ctx, nil)).Errorf(format, args...)
}
