package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	validator2 "go-api-samp/controller/validator"
)

/*
 * バリデータの作成
 * https://github.com/go-playground/validator
 * https://pkg.go.dev/gopkg.in/go-playground/validator.v10
 */

func RegisterValidation(e *echo.Echo) {
	v := validator.New()
	v.RegisterValidation("dateformat", validator2.IsDateFormat)
	e.Validator = &customValidator{
		validate: v,
	}
}

type customValidator struct {
	validate *validator.Validate
}

func (c *customValidator) Validate(i interface{}) error {
	return c.validate.Struct(i)
}
