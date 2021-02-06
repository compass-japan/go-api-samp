package controller

import (
	"github.com/labstack/echo/v4"
	"go-api-samp/model/dto"
	"go-api-samp/model/errors"
	service "go-api-samp/service/interface"
	"net/http"
)

/*
 * Controller Handler
 * golangではcontextを渡して、DBや外部APIアクセスのタイムアウトを検知しリソースを開放できるようにする
 * https://tip.golang.org/pkg/context/
 */

type APIController struct {
	APIService service.APIService
}

func (c *APIController) RegisterHandler(eCtx echo.Context) error {
	ctx := getContext(eCtx)

	req := &dto.RegisterRequest{}
	if err := eCtx.Bind(req); err != nil {
		return errors.InvalidRequestError(err)
	}

	if err := eCtx.Validate(req); err != nil {
		return errors.InvalidRequestParamError(err)
	}

	if err := c.APIService.Register(ctx, req); err != nil {
		return err
	}

	return eCtx.JSON(http.StatusOK, &dto.RegisterResponse{
		Message: "weather registered",
	})
}

func (c *APIController) GetWeatherHandler(eCtx echo.Context) error {
	ctx := getContext(eCtx)

	req := &dto.GetWeatherRequest{}
	if err := eCtx.Bind(req); err != nil {
		return errors.InvalidRequestError(err)
	}

	if err := eCtx.Validate(req); err != nil {
		return errors.InvalidRequestParamError(err)
	}

	response, err := c.APIService.GetWeather(ctx, req)
	if err != nil {
		return err
	}

	return eCtx.JSON(http.StatusOK, response)
}

func (c *APIController) GetAPIDataHandler(eCtx echo.Context) error {
	ctx := getContext(eCtx)

	w, err := c.APIService.GetAPIData(ctx)
	if err != nil {
		return err
	}

	return eCtx.JSON(http.StatusOK, w)
}
