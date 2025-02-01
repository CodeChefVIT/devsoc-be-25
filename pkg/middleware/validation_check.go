package middleware

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func isValidName(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[A-Za-z0-9\s]+$`).MatchString(fl.Field().String())
}

func init() {
	validate.RegisterValidation("isValidName", isValidName)
}

func TrimSpaces(input interface{}) error {
	val := reflect.ValueOf(input)

	if val.Kind() != reflect.Ptr {
		return errors.New("invalid request: expected a pointer")
	}

	elem := val.Elem()

	trimStruct(elem)

	return nil
}

func trimStruct(item reflect.Value) {
	for i := 0; i < item.NumField(); i++ {
		field := item.Field(i)
		if field.Kind() == reflect.String {
			field.SetString(strings.TrimSpace(field.String()))
		}
	}
}
