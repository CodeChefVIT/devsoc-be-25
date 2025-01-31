package utils

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
}

func FormatValidationErrors(err error) string {

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				return fmt.Sprintf("%s field is required", e.Field())
			case "email":
				return fmt.Sprintf("%s field is invalid email format", e.Field())
			case "endswith":
				return fmt.Sprintf("%s field must end with @vitstudent.ac.in", e.Field())
			case "url":
				return fmt.Sprintf("%s field is invalid URL format", e.Field())
			case "len":
				return fmt.Sprintf("%s field is invalid length", e.Field())
			case "alphanum":
				return fmt.Sprintf("%s field must contain only letters or numbers", e.Field())
			}
		}
	}

	return "validation issue :- contact us on discord"
}

func ValidateAlphaNum(str string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	return regex.MatchString(str)
}
