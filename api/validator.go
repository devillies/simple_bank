package api

import (
	"github.com/devillies/simple_bank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {

	if curr, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(curr)
	}
	return false
}
