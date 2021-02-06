// +build e2e

package e2e

import (
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

/*
 * 結合テスト(e2e)なので実際に起動しているサーバにアクセスする
 * アクセスするサーバは環境変数より読み込む
 * ライブラリを利用
 * https://github.com/go-resty/resty
 */

const (
	authHeader = "Auth-Token"
	token = "auth-token"
	E2E_ENDPOINT = "E2E_ENDPOINT"
)

var (
	url = ""
)

func TestE2E(t *testing.T) {
	t.Helper()
	if os.Getenv(E2E_ENDPOINT) == "" {
		t.Logf("not set E2E_ENDPOINT env")
		t.FailNow()
	}
	url = os.Getenv(E2E_ENDPOINT)

	tests := []struct{
		name string
		test func(t *testing.T)
	}{
		{
			name: "ヘルスチェック",
			test: testHealthCheck,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, test.test)
	}
}

func testHealthCheck(t *testing.T) {
	tests := []struct{
		name string
		path string
	}{
		{
			name: "ヘルスチェック",
			path: "/health",
		},
	}

	client := resty.New()
	client.RetryCount = 3
	client.RetryWaitTime = 1 * time.Second
	// client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	for _, test := range tests {
		tp := test
		t.Run(tp.name, func(t *testing.T) {
			res, err := client.R().
				SetHeader(authHeader, token).
				SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
				Get(url + tp.path)

			if err != nil {
				t.Error(err)
			}

			t.Log(string(res.Body()))
			assert.Equal(t, http.StatusOK, res.StatusCode())
		})
	}
}
