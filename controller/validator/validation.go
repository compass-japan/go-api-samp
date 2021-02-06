package validator

import (
	"github.com/go-playground/validator/v10"
	"time"
)

/*
 * カスタムバリデーター
 */

const (
	layout = "20060102"
)

// after len == 8 alphanum validator
func IsDateFormat(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	if _, err := time.Parse(layout, str); err != nil {
		return false
	}
	return true
}
