package api

import (
	"simplebank/utils"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	c := fieldLevel.Field().String()
	return utils.IsSupportedCurrency(c)
}
