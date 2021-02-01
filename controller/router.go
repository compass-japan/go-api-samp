package controller

import (
	"github.com/labstack/echo/v4"
	controller "go-api-samp/controller/internal"
	"go-api-samp/service"
	"net/http"
)

func RegisterRoute(e *echo.Echo, provider service.Provider) {

	s := provider.GetAPIService()
	c := &controller.APIController{
		APIService: s,
	}

	e.Use(controller.HeaderHandler(), controller.SetRequestID())

	e.GET("/health", func(eCtx echo.Context) error {
		return eCtx.String(http.StatusOK, "{ Status: OK }")
	})

	e.POST("/register", c.RegisterHandler)
	e.GET("/get/:locationId/:date", c.GetWeatherHandler)
	e.GET("/get/apidata", c.GetAPIDataHandler)

	e.HTTPErrorHandler = controller.ErrorHandler
}
