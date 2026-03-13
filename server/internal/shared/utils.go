package shared

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func ParseAndValidate(r *http.Request, v *validator.Validate, data any) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return fmt.Errorf("invalid request body")
	}

	return v.Struct(data)
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteValidationError(w http.ResponseWriter, err error) {
	validationErrors := make(map[string]string)

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			switch e.Tag() {
			case "required":
				validationErrors[e.Field()] = "This field is required"
			case "email":
				validationErrors[e.Field()] = "Please enter a valid email address"
			case "min":
				validationErrors[e.Field()] = fmt.Sprintf("Must be at least %s characters", e.Param())
			case "max":
				validationErrors[e.Field()] = fmt.Sprintf("Must be at most %s characters", e.Param())
			case "len":
				validationErrors[e.Field()] = fmt.Sprintf("Must be exactly %s characters", e.Param())
			case "numeric":
				validationErrors[e.Field()] = "Must be a number"
			case "alpha":
				validationErrors[e.Field()] = "Must contain only letters"
			case "alphanum":
				validationErrors[e.Field()] = "Must contain only letters and numbers"
			default:
				validationErrors[e.Field()] = "Invalid value for this field"
			}
		}
	}

	WriteJSON(w, http.StatusBadRequest, ErrorResponse{
		Message: "validation failed",
		Errors:  validationErrors,
	})
}
