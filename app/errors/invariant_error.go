package errors

import "net/http"

// InvariantErrorParam is a struct that contains parameters for InvariantError
type InvariantErrorParam struct {
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

	// RemoveDetailsInProduction is a flag that indicates whether to remove details in production
	//
	// Default to true
	RemoveDetailsInProduction bool
}

// InvariantError is an error that occurs when invariant is violated
type InvariantError struct {
	code    int
	message string
	details string
}

func NewInvariantError(param InvariantErrorParam) InvariantError {
	if param.Code == 0 {
		param.Code = http.StatusBadRequest
	}
	if param.Message == "" {
		param.Message = "invariant_error"
	}
	if param.RemoveDetailsInProduction {
		param.Details = ""
	}
	return InvariantError{param.Code, param.Message, param.Details}
}

func (e InvariantError) Code() int {
	return e.code
}

func (e InvariantError) Message() string {
	return e.message
}

func (e InvariantError) Details() string {
	return e.details
}

func (e InvariantError) Error() string {
	return e.message
}
