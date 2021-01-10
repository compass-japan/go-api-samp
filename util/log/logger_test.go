package log

import (
	"github.com/stretchr/testify/assert"
	"go-api-samp/util/config"
	"testing"
)

func TestLogPrioritySuccess(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		prior    priority
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
