package errors

import "net/http"

// InvariantParam is a struct that contains parameters for InvariantError
type InvariantParam struct {
	// Code is the status code of the error
	//
	// Default to HTTP status code 400
	Code int

	// Message is the message of the error
	//
	// Default to "invariant_error"
	Message string

	// Details is the details of the error
	//
	// Default to ""
	Details string
}

// Invariant is an error that occurs when invariant is violated
type Invariant struct {
	code    int
	message string
	details string
}

func NewInvariant(param InvariantParam) Invariant {
	if param.Code == 0 {
		param.Code = http.StatusBadRequest
	}
	if param.Message == "" {
		param.Message = "invariant_error"
	}
	return Invariant{param.Code, param.Message, param.Details}
}

func (e Invariant) Code() int {
	return e.code
}

func (e Invariant) Message() string {
	return e.message
}

func (e Invariant) Details() string {
	return e.details
}

func (e Invariant) Error() string {
	return e.message
}
