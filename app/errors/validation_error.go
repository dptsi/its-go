package errors

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// ValidationErrorData is a struct that contains data for validation error
type ValidationErrorData struct {
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// ErrValidation is an error that occurs when the request is invalid
var ErrValidation = errors.New("validation_error")

// GetValidationErrors is a function that returns ValidationErrorData
func GetValidationErrors(errs validator.ValidationErrors) map[string]ValidationErrorData {
	data := make(map[string]ValidationErrorData)
	for _, e := range errs {
		data[e.Field()] = ValidationErrorData{
			Tag:   e.Tag(),
			Param: e.Param(),
		}
	}

	return data
}
