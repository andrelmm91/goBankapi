package api

import (
	"simplebank/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(FieldLevel validator.FieldLevel) bool {
	if currency, ok := FieldLevel.Field().Interface().(string); ok {
		// check currency is supported
		return util.IsSupportedCurrency(currency)
	}
	return false
}

var validRole validator.Func = func(FieldLevel validator.FieldLevel) bool {
	if role, ok := FieldLevel.Field().Interface().(string); ok {
		// check currency is supported
		return util.IsSupportedRole(role)
	}
	return false
}
