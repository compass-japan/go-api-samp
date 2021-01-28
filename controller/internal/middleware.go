package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"go-api-samp/util/scope"
)

func getContext(eCtx echo.Context) context.Context {
	if ctx := eCtx.Get(string(scope.RequestIDContextKey)); ctx != nil {
		return ctx.(context.Context)
	}
	return context.Background()
}

func SetContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(eCtx echo.Context) error {
			ctx := eCtx.Request().Context()
			id, _ := uuid.NewRandom()
			ctx = scope.SetRequestID(ctx, id.String())
			eCtx.Set(string(scope.RequestIDContextKey), ctx)
			return next(eCtx)
		}
	}
}
