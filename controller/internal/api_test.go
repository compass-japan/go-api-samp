package controller

import (
	"bytes"
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	validator2 "go-api-samp/controller/validator"
	"go-api-samp/model/dto"
	"go-api-samp/model/errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	authHeaderKey = "Auth-Token"
	authToken = "auth-token"
)

var (
	err = errors2.New("error test")
)

func TestController(t *testing.T) {
	tests := []struct {
		name string
		test func(e *echo.Echo) func(t *testing.T)
	}{
		{
			name: "RegisterSuccess",
			test: testRegisterSuccess,
		},
		{
			name: "RegisterError",
			test: testRegisterRequestError,
		},
		{
			name: "GetWeatherSuccess",
			test: testGetWeatherSuccess,
		},
		{
			name: "GetWeatherError",
			test: testGetWeatherError,
		},
		{
			name: "GetAPIDataSuccess",
			test: testGetAPIDataSuccess,
		},
		{
			name: "GetAPIDataError",
			test: testGetAPIDataError,
		},
	}

	e := echo.New()
	e.Validator = &mockValidator{}

	for _, test := range tests {
		test := test
		t.Run(test.name, test.test(e))
	}
}

func testRegisterSuccess(e *echo.Echo) func(t *testing.T) {
	reqBody := fmt.Sprintf(`{"location_id":%d,"date":"%s","weather":%d,"comment":"%s"}`,
		1,"20200101",1, "comment")
	return func(t *testing.T) {
		tests := []struct{
			name string
			mock *mockAPIService
		}{
			{
				name: "正常系",
				mock: &mockAPIService{
					regErr: nil,
				},
			},
		}

		for _, test := range tests {
			tp := test
			t.Run(tp.name, func(t *testing.T) {
				res := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
				req.Header.Add(authHeaderKey, authToken)
				req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

				ctx := e.NewContext(req, res)

				c := &APIController{
					APIService: tp.mock,
				}

				b, _ := json.Marshal(&dto.RegisterResponse{
					Message: "weather registered",
				})

				err := HeaderHandler()(SetRequestID()(c.RegisterHandler))(ctx)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.JSONEq(t, string(b), res.Body.String())
			})
		}
	}
}

func testRegisterRequestError(e *echo.Echo) func(t *testing.T) {
	base := &dto.RegisterRequest{
		LocationId: 1,
		Date: "20200101",
		Weather: 1,
		Comment: "",
	}
	return func(t *testing.T) {
		tests := []struct{
			name      string
			req       *dto.RegisterRequest
			mock *mockAPIService
			isBindErr bool
			appErr    func(cause error) errors.ApplicationError
		}{
			{
				name: "bindエラー",
				req: func() *dto.RegisterRequest {
					m := copyRegisterRequest(base)
					m.LocationId = 0
					return m
				}(),
				isBindErr: true,
				appErr:    errors.InvalidRequestError,
			},
			{
				name: "location_id不正(未指定==0)",
				req: func() *dto.RegisterRequest {
					m := copyRegisterRequest(base)
					m.LocationId = 0
					return m
				}(),
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "date不正(未指定)",
				req: func() *dto.RegisterRequest {
					m := copyRegisterRequest(base)
					m.Date = ""
					return m
				}(),
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "date不正(length == 7)",
				req: func() *dto.RegisterRequest {
					m := copyRegisterRequest(base)
					m.Date = "2020111"
					return m
				}(),
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "date不正(length == 9)",
				req: func() *dto.RegisterRequest {
					m := copyRegisterRequest(base)
					m.Date = "202011100"
					return m
				}(),
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "date不正(記号)",
				req: func() *dto.RegisterRequest {
					m := copyRegisterRequest(base)
					m.Date = "202*0101"
					return m
				}(),
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "date不正(月)",
				req: func() *dto.RegisterRequest {
					m := copyRegisterRequest(base)
					m.Date = "20201301"
					return m
				}(),
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "date不正(日)",
				req: func() *dto.RegisterRequest {
					m := copyRegisterRequest(base)
					m.Date = "20200230"
					return m
				}(),
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "weather不正(未指定=0)",
				req: func() *dto.RegisterRequest {
					m := copyRegisterRequest(base)
					m.Weather = 0
					return m
				}(),
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "comment不正(len > 100)",
				req: func() *dto.RegisterRequest {
					m := copyRegisterRequest(base)
					m.Comment = strings.Repeat("a", 101)
					return m
				}(),
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "サービス層のエラー",
				req: func() *dto.RegisterRequest {
					m := copyRegisterRequest(base)
					return m
				}(),
				mock: &mockAPIService{
					regErr: err,
				},
				appErr: errors.InternalServerError,
			},
		}

		for _ ,test := range tests {
			tp := test
			t.Run(tp.name, func(t *testing.T) {
				b, _ := json.Marshal(tp.req)
				res := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/get/apidata", bytes.NewReader(b))
				req.Header.Add(authHeaderKey, authToken)
				if !tp.isBindErr {
					req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
				}

				ctx := e.NewContext(req, res)
				c := &APIController{
					APIService: tp.mock,
				}

				err := HeaderHandler()(SetRequestID()(c.RegisterHandler))(ctx)
				assert.Error(t, err)
				ErrorHandler(err, ctx)
				assert.Equal(t, tp.appErr(err).StatusCode(), res.Code)
				response, _ := json.Marshal(&dto.ErrorResponse{Message: tp.appErr(err).Message()})
				assert.JSONEq(t, string(response), res.Body.String())
			})
		}
	}
}

func testGetWeatherSuccess(e *echo.Echo) func(t *testing.T) {
	return func(t *testing.T) {
		tests := []struct{
			name string
			mock *mockAPIService
		}{
			{
				name: "正常系",
				mock: &mockAPIService{
					wResponse: &dto.GetWeatherResponse{
						Location: "新宿",
						Date: "20200101",
						Weather: "Sunny",
						Comment: "test comment",
					},
				},
			},
		}

		for _, test := range tests {
			tp := test
			t.Run(tp.name, func(t *testing.T) {
				res := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Add(authHeaderKey, authToken)

				ctx := e.NewContext(req, res)
				ctx.SetPath("/get/:locationId/:date")
				ctx.SetParamNames("locationId", "date")
				ctx.SetParamValues("1", "20200101")

				c:= &APIController{
					APIService: tp.mock,
				}

				b, _ := json.Marshal(tp.mock.wResponse)
				err := HeaderHandler()(SetRequestID()(c.GetWeatherHandler))(ctx)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.JSONEq(t, string(b), res.Body.String())
			})
		}
	}
}

func testGetWeatherError(e *echo.Echo) func(t *testing.T) {
	return func(t *testing.T) {
		tests := []struct{
			name string
			mock *mockAPIService
			paramLocationId string
			paramDate string
			bindErr bool
			appErr func(cause error) errors.ApplicationError
		}{
			{
				name: "bind エラー",
				paramLocationId: "1",
				paramDate: "20200101",
				bindErr: true,
				appErr: errors.InvalidRequestError,
			},
			{
				name: "locationId不正(空)",
				paramLocationId: "",
				paramDate: "20200101",
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "locationId不正(0)",
				paramLocationId: "0",
				paramDate: "20200101",
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "date不正(未指定)",
				paramLocationId: "1",
				paramDate: "",
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "date不正(length == 7)",
				paramLocationId: "1",
				paramDate: "2020010",
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "date不正(length == 9)",
				paramLocationId: "1",
				paramDate: "202001011",
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "date不正(記号)",
				paramLocationId: "1",
				paramDate: "202*0101",
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "date不正(月)",
				paramLocationId: "1",
				paramDate: "20201301",
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "date不正(日)",
				paramLocationId: "1",
				paramDate: "20200230",
				appErr: errors.InvalidRequestParamError,
			},
			{
				name: "サービス層のエラー",
				paramLocationId: "1",
				paramDate: "20200101",
				mock: &mockAPIService{
					wErr: err,
				},
				appErr: errors.InternalServerError,
			},
		}

		for _, test := range tests {
			tp := test

			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tp.bindErr {
				req = httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte("aa")))
			}
			req.Header.Add(authHeaderKey, authToken)

			ctx := e.NewContext(req, res)
			ctx.SetPath("/get/:locationId/:date")
			ctx.SetParamNames("locationId", "date")
			ctx.SetParamValues(tp.paramLocationId, tp.paramDate)

			c:= &APIController{
				APIService: tp.mock,
			}

			err := HeaderHandler()(SetRequestID()(c.GetWeatherHandler))(ctx)
			assert.Error(t, err)
			ErrorHandler(err, ctx)
			assert.Equal(t, tp.appErr(err).StatusCode(), res.Code)
			response, _ := json.Marshal(&dto.ErrorResponse{Message: tp.appErr(err).Message()})
			assert.JSONEq(t, string(response), res.Body.String())
		}
	}
}

func testGetAPIDataSuccess(e *echo.Echo) func(t *testing.T) {
	return func(t *testing.T) {
		tests := []struct{
			name string
			mock *mockAPIService
		}{
			{
				name: "正常系",
				mock: &mockAPIService{
					apiResponse: &dto.ExApiResponse{},
					apiErr: nil,
				},
			},
		}

		for _, test := range tests {
			tp := test
			t.Run(tp.name, func(t *testing.T) {
				res := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/get/apidata", nil)
				req.Header.Add(authHeaderKey, authToken)

				ctx := e.NewContext(req, res)

				c := &APIController{
					APIService: tp.mock,
				}
				b, _ := json.Marshal(tp.mock.apiResponse)

				err := HeaderHandler()(SetRequestID()(c.GetAPIDataHandler))(ctx)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.JSONEq(t, string(b), res.Body.String())
			})
		}
	}
}

func testGetAPIDataError(e *echo.Echo) func(t *testing.T) {
	return func(t *testing.T) {
		tests := []struct{
			name string
			mock *mockAPIService
		}{
			{
				name: "サービス error",
				mock: &mockAPIService{
					apiResponse: nil,
					apiErr: err,
				},
			},
		}

		for _, test := range tests {
			tp := test
			t.Run(tp.name, func(t *testing.T) {
				res := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/get/apidata", nil)
				req.Header.Add(authHeaderKey, authToken)

				ctx := e.NewContext(req, res)

				c := &APIController{
					APIService: tp.mock,
				}

				err := HeaderHandler()(SetRequestID()(c.GetAPIDataHandler))(ctx)
				assert.Error(t, err)
				ErrorHandler(err, ctx)
				assert.Equal(t, http.StatusInternalServerError, res.Code)
				response, _ := json.Marshal(&dto.ErrorResponse{Message: errors.InternalServerError(err).Message()})
				assert.JSONEq(t, string(response), res.Body.String())
			})
		}
	}
}

type mockAPIService struct {
	regErr error
	wResponse *dto.GetWeatherResponse
	wErr error
	apiResponse *dto.ExApiResponse
	apiErr error
}

func (m *mockAPIService) Register(ctx context.Context, payload *dto.RegisterRequest) error {
	return m.regErr
}

func (m *mockAPIService) GetWeather(ctx context.Context, payload *dto.GetWeatherRequest) (*dto.GetWeatherResponse, error) {
	return m.wResponse, m.wErr
}

func (m *mockAPIService) GetAPIData(ctx context.Context) (*dto.ExApiResponse, error) {
	return m.apiResponse, m.apiErr
}


type mockValidator struct {
}


func (m *mockValidator) Validate(i interface{}) error {
	v := validator.New()
	v.RegisterValidation("dateformat", validator2.IsDateFormat)
	return v.Struct(i)
}

func copyRegisterRequest(base *dto.RegisterRequest) *dto.RegisterRequest {
	return &dto.RegisterRequest{
		LocationId: base.LocationId,
		Date: base.Date,
		Weather: base.Weather,
		Comment: base.Comment,
	}
}
