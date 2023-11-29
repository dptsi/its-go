package errors

import "net/http"

type NotFoundErrorParam struct {
	// Code is the status code of the error
	//
	// Default to HTTP status code 404
	Code int

	// Msg is the message of the error
	//
	// Default to "not_found"
	Msg string
}

// NotFoundError is an error that occurs when the resource is not found
type NotFoundError struct {
	code int
	msg  string
}

func NewNotFoundError(param NotFoundErrorParam) NotFoundError {
	if param.Code == 0 {
		param.Code = http.StatusNotFound
	}
	if param.Msg == "" {
		param.Msg = "not_found"
	}
	return NotFoundError{param.Code, param.Msg}
}

func (e NotFoundError) Code() int {
	return e.code
}

func (e NotFoundError) Error() string {
	return e.msg
}
