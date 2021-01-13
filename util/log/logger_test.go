package log

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"go-api-samp/util/config"
	"testing"
	"time"
)

func TestLogPriority(t *testing.T) {
	tests := []struct {
		name       string
		logLevel   string
		prior      priority
		isOutEmpty bool
	}{
		{
			name:     "ERROR priority",
			logLevel: "ERROR",
			prior:    priorityError,
		},
		{
			name:     "WARNING priority",
			logLevel: "WARNING",
			prior:    priorityWarning,
		},
		{
			name:     "INFO priority",
			logLevel: "INFO",
			prior:    priorityInfo,
		},
		{
			name:     "DEBUG priority",
			logLevel: "DEBUG",
			prior:    priorityDebug,
		},
		{
			name:     "default priority",
			logLevel: "def",
			prior:    priorityInfo,
		},
		{
			name:     "lower case",
			logLevel: "error",
			prior:    priorityError,
		},
	}

	t.Parallel()
	for _, test := range tests {
		tp := test
		t.Run(tp.name, func(t *testing.T) {
			logger := &appLogger{
				Config: &config.LogConfig{
					Level: tp.logLevel,
				},
			}
			assert.Equal(t, tp.prior, logger.getPriority())
		})
	}
}

func TestLogOutput(t *testing.T) {
	tests := []struct {
		name         string
		logLevelConf string
		isOutEmpty   bool
	}{
		{
			name:         "ERROR higher than INFO",
			logLevelConf: "ERROR",
			isOutEmpty:   true,
		},
		{
			name:         "WARNING higher than INFO",
			logLevelConf: "WARNING",
			isOutEmpty:   true,
		},
		{
			name:         "INFO equal INFO",
			logLevelConf: "INFO",
			isOutEmpty:   false,
		},
		{
			name:         "DEBUG lower than INFO",
			logLevelConf: "DEBUG",
			isOutEmpty:   false,
		},
	}

	t.Parallel()
	for _, test := range tests {
		tp := test
		t.Run(tp.name, func(t *testing.T) {
			logger := &appLogger{
				Config: &config.LogConfig{
					Level: tp.logLevelConf,
				},
			}

			// INFO固定でテスト
			b := &bytes.Buffer{}
			tempDest := logInfo.dest
			logInfo.dest = b
			defer func() {
				logInfo.dest = tempDest
			}()

			logger.Info(context.Background(), "log test")
			assert.Equal(t, tp.isOutEmpty, b.String() == "", "unexpected output: %s", time.Now())
		})
	}
}
