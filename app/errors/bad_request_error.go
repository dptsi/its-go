package errors

import "net/http"

// BadRequestErrorParam is a struct that contains parameters for BadRequestError
type BadRequestErrorParam struct {
	// Code is the status code of the error
	//
	// Default to HTTP status code 400
	Code int

	// Message is the message of the error
	//
	// Default to "bad_request"
	Message string

	// Data is the additional data of the error
	//
	// Default to nil
	Data map[string]interface{}
}

// BadRequestError is an error that occurs when the request is invalid
// and cannot be processed
type BadRequestError struct {
	code    int
	message string
	data    map[string]interface{}
}

func NewBadRequestError(param BadRequestErrorParam) BadRequestError {
	if param.Code == 0 {
		param.Code = http.StatusBadRequest
	}
	if param.Message == "" {
		param.Message = "bad_request"
	}

	return BadRequestError{
		code:    param.Code,
		message: param.Message,
		data:    param.Data,
	}
}

func (e BadRequestError) Code() int {
	return e.code
}

func (e BadRequestError) Message() string {
	return e.message
}

func (e BadRequestError) Data() map[string]interface{} {
	return e.data
}

func (e BadRequestError) Error() string {
	return e.message
}
