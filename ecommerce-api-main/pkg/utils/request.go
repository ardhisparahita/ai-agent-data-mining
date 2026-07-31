package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationStruct(data any) error {
	return Validate.Struct(data)
}

func ValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		errors[field] = e.Tag()
	}

	return errors
}
