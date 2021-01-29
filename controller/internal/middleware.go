package controller

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"go-api-samp/model/dto"
	"go-api-samp/model/errors"
	"go-api-samp/util/log"
	"go-api-samp/util/scope"
	"net/http"
	"strings"
)

func getContext(eCtx echo.Context) context.Context {
	if ctx := eCtx.Get(fmt.Sprint(scope.RequestIDContextKey)); ctx != nil {
		return ctx.(context.Context)
	}
	return context.Background()
}

func HeaderHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		token := eCtx.Request().Header.Get("Auth-Token")
		if strings.ToLower(token) != "auth-token" {
			return eCtx.JSON(http.StatusUnauthorized, &dto.ErrorResponse{
				Message: errors.Application.UnauthorizedError(nil).Message(),
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
		eCtx.Set(fmt.Sprint(scope.RequestIDContextKey), ctx)
		return next(eCtx)
	}
}

func ErrorHandler(err error, eCtx echo.Context) {
	if e, ok := err.(*echo.HTTPError); ok {
		switch e.Code {
		case http.StatusNotFound:
			err = errors.Application.HttpRouteNotFoundError(e)
		case http.StatusMethodNotAllowed:
			err = errors.Application.HttpMethodNotAllowedError(e)
		default:
			err = errors.Application.InternalServerError(e)
		}
	}

	e, ok := err.(errors.ApplicationError)
	if !ok {
		e = errors.Application.InternalServerError(e)
	}

	if e.LogIgnorable() {
		logger := log.GetLogger()
		logger.Info(getContext(eCtx), "error footprint: %v", e.Causes())
	}

	eCtx.JSON(e.StatusCode(), &dto.ErrorResponse{
		Message: e.Message(),
	})
}
