package controller

import (
	"github.com/labstack/echo"
	controller "go-api-samp/controller/internal"
	"net/http"
)

func RegisterRoute(e *echo.Echo) {
	e.Use(controller.SetContext())

	e.GET("/health", func(eCtx echo.Context) error {
		return eCtx.String(http.StatusOK, "{ Status: OK }")
	})
}
