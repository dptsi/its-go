package errors

import "net/http"

// BadRequestParam is a struct that contains parameters for BadRequestError
type BadRequestParam struct {
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

// BadRequest is an error that occurs when the request is invalid
// and cannot be processed
type BadRequest struct {
	code    int
	message string
	data    map[string]interface{}
}

func NewBadRequest(param BadRequestParam) BadRequest {
	if param.Code == 0 {
		param.Code = http.StatusBadRequest
	}
	if param.Message == "" {
		param.Message = "bad_request"
	}

	return BadRequest{
		code:    param.Code,
		message: param.Message,
		data:    param.Data,
	}
}

func (e BadRequest) Code() int {
	return e.code
}

func (e BadRequest) Message() string {
	return e.message
}

func (e BadRequest) Data() map[string]interface{} {
	return e.data
}

func (e BadRequest) Error() string {
	return e.message
}
