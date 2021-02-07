// +build e2e

package e2e

import (
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go-api-samp/model/dto"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
	"encoding/json"
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
	client *resty.Client = nil
)

func TestE2E(t *testing.T) {
	t.Helper()
	if os.Getenv(E2E_ENDPOINT) == "" {
		t.Logf("not set E2E_ENDPOINT env")
		t.FailNow()
	}
	url = os.Getenv(E2E_ENDPOINT)

	c := resty.New()
	c.SetTimeout(5 * time.Second)
	client = c

	tests := []struct{
		name string
		test func(t *testing.T)
	}{
		{
			name: "ヘルスチェック",
			test: testHealthCheck,
		},
		{
			name: "外部API",
			test: testExAPI,
		},
		{
			name: "Get Weather",
			test: testGetWeather,
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
				Get(url + "/health")

			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, http.StatusOK, res.StatusCode())
			assert.Equal(t, "{ Status: OK }", res.String())
		})
	}
}

func testExAPI(t *testing.T) {
	tests := []struct{
		name string
	}{
		{
			name: "正常系",
		},
	}

	for _, test := range tests {
		tp := test
		t.Run(tp.name, func(t *testing.T) {
			res, err := client.R().
				SetHeader(authHeader, token).
				SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
				Get(url + "/get/apidata")

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode())
			t.Log(res.String())
			assert.NotEmpty(t, res.String())
		})
	}
}

func testGetWeather(t *testing.T) {
	var get = func(locationId int, date string, isErr bool, status int, expected *dto.GetWeatherResponse) func() {
		return func(){
			res, err := client.R().
				SetHeader(authHeader, token).
				SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
				Get(url + "/get/" + strconv.Itoa(locationId) + "/" + date)

			if isErr {
				assert.Equal(t, status, res.StatusCode())
			}
			if !isErr {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, res.StatusCode())
				b, _ := json.Marshal(expected)
				assert.JSONEq(t, string(b), res.String())
			}
			t.Log(locationId, date, expected)
		}
	}
	var register = func(isErr bool, status int, request *dto.RegisterRequest) func() {
		return func(){
			res, err := client.R().
				SetHeader(authHeader, token).
				SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
				SetBody(request).
				Post(url + "/register")

			if isErr {
				assert.Equal(t, status, res.StatusCode())
			}
			if !isErr {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, res.StatusCode())
				assert.JSONEq(t, `{"message": "weather registered"}`, res.String())
			}
			t.Log(request)
		}
	}
	tests := []struct{
		name string
		functions []interface{}
	}{
		{
			name: "register(insert), get, register(update), get",
			functions: []interface{}{
				register(false, http.StatusOK, &dto.RegisterRequest{
					LocationId: 1,
					Date: "20001010",
					Weather: 3,
					Comment: "register(insert) comment",
				}),
				get(1, "20001010", false, http.StatusOK, &dto.GetWeatherResponse{
					Location: "新宿",
					Date: "20001010",
					Weather: "Rainy",
					Comment: "register(insert) comment",
				}),
				register(false, http.StatusOK, &dto.RegisterRequest{
					LocationId: 1,
					Date: "20001010",
					Weather: 1,
					Comment: "register(update) comment",
				}),
				get(1, "20001010", false, http.StatusOK, &dto.GetWeatherResponse{
					Location: "新宿",
					Date: "20001010",
					Weather: "Sunny",
					Comment: "register(update) comment",
				}),
			},
		},
		{
			name: "register error",
			functions: []interface{}{
				register(true, http.StatusBadRequest, &dto.RegisterRequest{
					LocationId: 0, // not required
					Date: "20001010",
					Weather: 1,
					Comment: "no register",
				}),
				register(true, http.StatusInternalServerError, &dto.RegisterRequest{
					LocationId: 999, // no location
					Date: "20001010",
					Weather: 1,
					Comment: "no register",
				}),
				register(true, http.StatusBadRequest, &dto.RegisterRequest{
					LocationId: 1,
					Date: "2000101", // len != 8
					Weather: 1,
					Comment: "no register",
				}),
				register(true, http.StatusBadRequest, &dto.RegisterRequest{
					LocationId: 1,
					Date: "20001010",
					Weather: 0, // not required
					Comment: "no register",
				}),
			},
		},
		{
			name: "get error",
			functions: []interface{}{
				get(0, "20001010", true, http.StatusBadRequest,nil),
				get(999, "20001010", true, http.StatusInternalServerError, nil),
				get(1, "2000101", true, http.StatusBadRequest, nil),
			},
		},
	}

	for _, test := range tests {
		tp := test
		t.Run(tp.name, func(t *testing.T) {
			for _, f := range tp.functions {
				v, _ := f.(func())
				v()
			}
		})
	}
}
