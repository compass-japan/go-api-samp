package log

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"go-api-samp/util/config"
	"os"
	"testing"
)

/*
 * ログ出力スキップ判定とレベル設定ごとのログ出力のテスト
 */

func TestLogLevel(t *testing.T) {
	tests := []struct {
		name         string
		confLogLevel string
		level        []lvl
		isSkip       bool
	}{
		{
			name:         "conf: ERROR, out: ERROR",
			confLogLevel: "ERROR",
			level:        []lvl{ERROR},
			isSkip:       false,
		},
		{
			name:         "conf: ERROR, out: lte WARN",
			confLogLevel: "ERROR",
			level:        []lvl{WARN, INFO, DEBUG},
			isSkip:       true,
		},
		{
			name:         "conf: WARN, out: gte WARN",
			confLogLevel: "WARN",
			level:        []lvl{ERROR, WARN},
			isSkip:       false,
		},
		{
			name:         "conf: WARN, out: lte INFO",
			confLogLevel: "WARN",
			level:        []lvl{INFO, DEBUG},
			isSkip:       true,
		},
		{
			name:         "conf: INFO, out: gte INFO",
			confLogLevel: "INFO",
			level:        []lvl{ERROR, WARN, INFO},
			isSkip:       false,
		},
		{
			name:         "conf: INFO, out: lte DEBUG",
			confLogLevel: "INFO",
			level:        []lvl{DEBUG},
			isSkip:       true,
		},
		{
			name:         "conf: DEBUG, out: ALL",
			confLogLevel: "DEBUG",
			level:        []lvl{ERROR, WARN, INFO, DEBUG},
			isSkip:       false,
		},
		{
			name:         "small case conf: info, out: gte INFO",
			confLogLevel: "info",
			level:        []lvl{ERROR, WARN, INFO},
			isSkip:       false,
		},
		{
			name:         "small case conf: info, out: lte DEBUG",
			confLogLevel: "info",
			level:        []lvl{DEBUG},
			isSkip:       true,
		},
		{
			name:         "invalid level string conf: DEBUGG, out: ALL",
			confLogLevel: "DEBUGG",
			level:        []lvl{ERROR, WARN, INFO, DEBUG},
			isSkip:       false,
		},
	}

	for _, test := range tests {
		tp := test
		t.Run(tp.name, func(t *testing.T) {
			a := &appLogger{
				config: &config.LogConfig{
					Level: tp.confLogLevel,
				},
			}

			for _, v := range tp.level {
				assert.Equal(t, tp.isSkip, a.skipLevel(v), "conf:%s, level:%s", tp.confLogLevel, levels[v])
				func() {
					b := &bytes.Buffer{}
					stderr = b
					stdout = b
					defer func() {
						stderr = os.Stderr
						stdout = os.Stdout
					}()
					s := "log test"
					ctx := context.Background()
					switch v {
					case ERROR:
						a.Error(ctx, s)
					case WARN:
						a.Warn(ctx, s)
					case INFO:
						a.Info(ctx, s)
					case DEBUG:
						a.Debug(ctx, s)
					}
					assert.Equal(t, tp.isSkip, b.String() == "", "unexpected output")
				}()
			}
		})
	}
}
