package utils

import (

	"github.com/go-playground/validator"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
}

func FormatValidationErrors(err error) map[string]string {
	errorMessages := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errorMessages[e.Field()] = "This field is required"
			case "email":
				errorMessages[e.Field()] = "Invalid email format"
			case "endswith":
				errorMessages[e.Field()] = "Email must end with @vitstudent.ac.in"
			case "url":
				errorMessages[e.Field()] = "Invalid URL format"
			case "len":
				errorMessages[e.Field()] = "Invalid length"
			}
		}
	}

	return errorMessages
}
