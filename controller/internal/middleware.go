package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"go-api-samp/model/dto"
	"go-api-samp/util/scope"
	"net/http"
)

func getContext(eCtx echo.Context) context.Context {
	if ctx := eCtx.Get(string(scope.RequestIDContextKey)); ctx != nil {
		return ctx.(context.Context)
	}
	return context.Background()
}

func Handler (next echo.HandlerFunc) echo.HandlerFunc{
	return func (eCtx echo.Context) error {
		token := eCtx.Request().Header.Get("Auth-Token")
		if token != "auth-token" {
			return eCtx.JSON(http.StatusUnauthorized, &dto.ErrorResponse{
				Message: "Unauthorized request",
			})
		}
		return next(eCtx)
	}
}

func SetContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		ctx := eCtx.Request().Context()
		id, _ := uuid.NewRandom()
		ctx = scope.SetRequestID(ctx, id.String())
		eCtx.Set(string(scope.RequestIDContextKey), ctx)
		return next(eCtx)
	}
}

func ErrorHandler(err error, eCtx echo.Context) {

}
