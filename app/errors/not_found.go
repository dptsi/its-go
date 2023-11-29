package errors

import "net/http"

type NotFoundParam struct {
	// Code is the status code of the error
	//
	// Default to HTTP status code 404
	Code int

	// Msg is the message of the error
	//
	// Default to "not_found"
	Msg string
}

// NotFound is an error that occurs when the resource is not found
type NotFound struct {
	code int
	msg  string
}

func NewNotFound(param NotFoundParam) NotFound {
	if param.Code == 0 {
		param.Code = http.StatusNotFound
	}
	if param.Msg == "" {
		param.Msg = "not_found"
	}
	return NotFound{param.Code, param.Msg}
}

func (e NotFound) Code() int {
	return e.code
}

func (e NotFound) Error() string {
	return e.msg
}
