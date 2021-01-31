package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

func RegisterValidation(e *echo.Echo) {
	e.Validator = &customValidator{
		validate: validator.New(),
	}
}

type customValidator struct {
	validate *validator.Validate
}

func (c *customValidator) Validate(i interface{}) error {
	return c.validate.Struct(i)
}
