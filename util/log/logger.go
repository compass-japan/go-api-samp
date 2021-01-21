package log

import (
	"context"
	"fmt"
	"go-api-samp/util/config"
	"go-api-samp/util/scope"
	"io"
	"os"
	"strings"
	"time"
)

/*
 * ログ定義
 * 環境ごとのpriorityによりログ出力レベルを変える
 * ログレベルでのメソッドを用意(Error, Warning, Info, Debug)
 */

type priority uint32

const (
	priorityError priority = iota
	priorityWarning
	priorityInfo
	priorityDebug
)

const (
	labelError   = "ERROR"
	labelWarning = "WARNING"
	labelInfo    = "INFO"
	labelDebug   = "DEBUG"
)

func GetLogger() Logger {
	return &appLogger{config.Log}
}

type appLogger struct {
	Config *config.LogConfig
}

type logConf struct {
	dest     io.Writer
	level    string
	priority priority
}

var (
	logError   = &logConf{dest: os.Stderr, level: labelError, priority: priorityError}
	logWarning = &logConf{dest: os.Stderr, level: labelWarning, priority: priorityWarning}
	logInfo    = &logConf{dest: os.Stdout, level: labelInfo, priority: priorityInfo}
	logDebug   = &logConf{dest: os.Stdout, level: labelDebug, priority: priorityDebug}
	timeFormat = "2006/01/02 15:04:05"
)

type Logger interface {
	Error(context.Context, string, ...interface{})
	Warning(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
	Debug(context.Context, string, ...interface{})
}

func (l *appLogger) Error(ctx context.Context, s string, v ...interface{}) {
	l.log(ctx, logError, s, v...)
}
func (l *appLogger) Warning(ctx context.Context, s string, v ...interface{}) {
	l.log(ctx, logWarning, s, v...)
}
func (l *appLogger) Info(ctx context.Context, s string, v ...interface{}) {
	l.log(ctx, logInfo, s, v...)
}
func (l *appLogger) Debug(ctx context.Context, s string, v ...interface{}) {
	l.log(ctx, logDebug, s, v...)
}

func (l *appLogger) log(ctx context.Context, conf *logConf, s string, v ...interface{}) {
	if l.getPriority() < conf.priority {
		return
	}

	m := fmt.Sprintf(s, v...)
	log := fmt.Sprintf("[%s] requestID: %s, time: %s, message: %s\n",
		conf.level, scope.GetRequestID(ctx), time.Now().Format(timeFormat), m)

	_, err := fmt.Fprintf(conf.dest, log)
	if err != nil {
		fmt.Fprintf(os.Stderr, log)
	}
}

func (l *appLogger) getPriority() priority {
	switch strings.ToUpper(l.Config.Level) {
	case labelError:
		return priorityError
	case labelWarning:
		return priorityWarning
	case labelInfo:
		return priorityInfo
	case labelDebug:
		return priorityDebug
	default:
		return priorityInfo
	}
}
