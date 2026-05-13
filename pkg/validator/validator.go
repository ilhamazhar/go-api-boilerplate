package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Validate(s any) []FieldError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil
	}

	var errs []FieldError
	for _, e := range ve {
		errs = append(errs, FieldError{
			Field:   e.Field(),
			Message: message(e),
		})
	}
	return errs
}

func message(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "eqfield":
		return fmt.Sprintf("%s must be equal to %s", e.Field(), e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("%s must be at least %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s is not valid", e.Field())
	}
}
