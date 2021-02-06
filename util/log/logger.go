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
 * ログレベルでのメソッドを用意(Error, Warn, Info, Debug)
 * 監視のためにError,Warn・Info,Debugでそれぞれstderr, stdoutへと出力を向ける
 */

func GetLogger() Logger {
	return logger
}

func NewLogger(config *config.LogConfig) {
	logger = &appLogger{
		config: config,
	}
}

type lvl uint8

const (
	DEBUG lvl = iota
	INFO
	WARN
	ERROR
)

var (
	logger *appLogger = &appLogger{
		config: &config.LogConfig{
			Level: "DEBUG",
		},
	}
	levels = []string{
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
	}
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
)

type Logger interface {
	Error(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
	Debug(context.Context, string, ...interface{})
}

type appLogger struct {
	config *config.LogConfig
}

func (l *appLogger) Error(ctx context.Context, s string, v ...interface{}) {
	l.log(ctx, ERROR, stderr, s, v...)
}

func (l *appLogger) Warn(ctx context.Context, s string, v ...interface{}) {
	l.log(ctx, WARN, stderr, s, v...)
}

func (l *appLogger) Info(ctx context.Context, s string, v ...interface{}) {
	l.log(ctx, INFO, stdout, s, v...)
}

func (l *appLogger) Debug(ctx context.Context, s string, v ...interface{}) {
	l.log(ctx, DEBUG, stdout, s, v...)
}

const (
	LOG_FORMAT  = "[%s] requestID: %s, time: %s, message: %s\n"
	TIME_FORMAT = "2006/01/02 15:04:05" // go特有のフォーマット指定
)

func (l *appLogger) log(ctx context.Context, level lvl, dest io.Writer, s string, v ...interface{}) {
	if l.skipLevel(level) {
		return
	}

	m := fmt.Sprintf(s, v...)
	log := fmt.Sprintf(LOG_FORMAT, levels[level], scope.GetRequestID(ctx),
		time.Now().Format(TIME_FORMAT), m)

	_, err := fmt.Fprintf(dest, log)
	if err != nil {
		fmt.Fprintf(os.Stderr, log)
	}
}

func (l *appLogger) skipLevel(level lvl) bool {
	for i := ERROR; level < i; i-- {
		if strings.ToUpper(l.config.Level) == levels[i] {
			return true
		}
	}
	return false
}
