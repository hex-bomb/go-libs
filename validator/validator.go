package validator

import govalidator "github.com/go-playground/validator/v10"

var Validator *govalidator.Validate

func init() {
	Validator = govalidator.New()

	RegisterChecks(Validator)
}

func Get() *govalidator.Validate {
	return Validator
}
