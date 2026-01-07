package validator

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func HandleRequestError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			switch e.Tag() {
			case "required":
				return fmt.Sprintf("%s is require", e.Field())
			case "email":
				return fmt.Sprintf("%s is not a valid email", e.Field())
			case "min":
				return fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
			case "max":
				return fmt.Sprintf("%s cannot exceed %s characters", e.Field(), e.Param())
			case "len":
				return fmt.Sprintf("%s must have exactly %s characters", e.Field(), e.Param())
			case "numeric":
				return fmt.Sprintf("%s must be a number", e.Field())
			case "uuid4":
				return fmt.Sprintf("%s must be a valid UUID v4", e.Field())
			case "oneof":
				return fmt.Sprintf("%s must have the value: %s", e.Field(), e.Param())
			default:
				return fmt.Sprintf("%s is not valid", e.Field())
			}
		}
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return fmt.Sprintf("%s must be a %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String())
	}

	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return fmt.Sprintf("Invalid JSON at byte %d", syntaxError.Offset)
	}

	if err != nil {
		return err.Error()
	}

	return "Invalid request"
}