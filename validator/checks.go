package validator

import (
	"math"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// RegisterChecks реистрирует кастомные проверки тегов структур
func RegisterChecks(v *validator.Validate) {
	_ = v.RegisterValidation("multipleof", validateMultipleOf)
	_ = v.RegisterValidation("host_port", validateHostPort)

}

// проверка что делится на 10
func validateMultipleOf(fl validator.FieldLevel) bool {
	param := strings.TrimSpace(fl.Param())
	if param == "" {
		return false
	}
	denom, errU := strconv.ParseUint(param, 10, 64)
	if errU != nil || denom == 0 {
		return false
	}
	field := fl.Field()

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int()%int64(denom) == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint()%denom == 0
	case reflect.Float32, reflect.Float64:
		val := field.Float()
		div := float64(denom)
		return math.Abs((val/div)-math.Round(val/div)) < 0.01
	default:
		return false
	}
}

func validateHostPort(fl validator.FieldLevel) bool {
	param := strings.TrimSpace(fl.Param())
	if param == "" {
		return false
	}
	_, _, err := net.SplitHostPort(param)
	return err == nil
}
