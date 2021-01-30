package controller

import (
	"context"
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
	if ctx := eCtx.Get(scope.RequestIDContextKey); ctx != nil {
		return ctx.(context.Context)
	}
	return context.Background()
}

func HeaderHandler() echo.MiddlewareFunc {

	const (
		headerKey  = "Auth-Token"
		fixedValue = "auth-token"
	)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(eCtx echo.Context) error {
			token := eCtx.Request().Header.Get(headerKey)
			if strings.ToLower(token) != fixedValue {
				return eCtx.JSON(http.StatusUnauthorized, &dto.ErrorResponse{
					Message: errors.Application.UnauthorizedError(nil).Message(),
				})
			}
			return next(eCtx)
		}
	}
}

func SetRequestID() echo.MiddlewareFunc {

	generator := func() string {
		id, _ := uuid.NewRandom()
		return id.String()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(eCtx echo.Context) error {
			req := eCtx.Request()
			res := eCtx.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = generator()
			}
			ctx := eCtx.Request().Context()
			ctx = scope.SetRequestID(ctx, rid)
			eCtx.Set(scope.RequestIDContextKey, ctx)
			res.Header().Set(echo.HeaderXRequestID, rid)
			return next(eCtx)
		}
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
