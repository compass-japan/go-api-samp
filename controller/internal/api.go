package controller

import (
	"github.com/labstack/echo"
	"go-api-samp/model/dto"
	service "go-api-samp/service/interface"
	"net/http"
)

type APIController struct {
	APIService service.APIService
}

func (c *APIController) RegisterHandler(eCtx echo.Context) error {
	return eCtx.JSON(http.StatusOK, &dto.RegisterResponse{
		Message: "weather registered",
	})
}

func (c *APIController) GetWeatherHandler(eCtx echo.Context) error {
	return eCtx.JSON(http.StatusOK, &dto.GetWeatherResponse{
		Date:    "20200101",
		Weather: "sunny",
		Comment: "dammy",
	})
}

func (c *APIController) GetAPIDataHandler(eCtx echo.Context) error {
	return eCtx.JSON(http.StatusOK, &dto.ApiDataResponse{
		Date:    "20200202",
		Weather: "cloudy",
	})
}
