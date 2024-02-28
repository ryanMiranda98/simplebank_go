package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/ryanMiranda98/simplebank/util"
)

var ValidCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// Check if currency is supported
		return util.IsSupportedCurrency(currency)
	}
	return false
}
