package valid

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func Struct(i any) error {
	return validate.Struct(i)
}
