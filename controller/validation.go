package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

func RegisterValidation(e *echo.Echo) {
	e.Validator = &customEchoValidator{
		validate: validator.New(),
	}
}

type customEchoValidator struct {
	validate *validator.Validate
}

func (c *customEchoValidator) Validate(i interface{}) error {
	return c.validate.Struct(i)
}
