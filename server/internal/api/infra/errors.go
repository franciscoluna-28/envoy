package infra

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func MapValidationErrors(err error) []ValidationError {
	var errs []ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, f := range validationErrs {
			errs = append(errs, ValidationError{
				Field:   f.Field(),
				Message: getErrorMessage(f),
			})
		}
	}
	return errs
}

func getErrorMessage(f validator.FieldError) string {
	switch f.Tag() {
	case "required":
		return "This field is mandatory"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", f.Param())
	default:
		return "Invalid value"
	}
}
