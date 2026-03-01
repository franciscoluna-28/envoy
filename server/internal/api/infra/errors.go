package infra

import (
	"fmt"

	"github.com/franciscoluna/envoy/server/internal/shared"
	"github.com/go-playground/validator/v10"
)

func MapValidationErrors(err error) []shared.ValidationError {
	var errs []shared.ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, f := range validationErrs {
			errs = append(errs, shared.ValidationError{
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
