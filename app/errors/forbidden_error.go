package errors

import "net/http"

type ForbiddenErrorParam struct {
	// Code is the status code of the error
	//
	// Default to HTTP status code 403
	Code int

	// Message is the message of the error
	//
	// Default to "forbidden"
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

// ForbiddenError is an error that occurs when the request is forbidden
type ForbiddenError struct {
	msg     string
	details string
}

func NewForbiddenError(param ForbiddenErrorParam) ForbiddenError {
	if param.Code == 0 {
		param.Code = http.StatusForbidden
	}
	if param.Message == "" {
		param.Message = "forbidden"
	}
	if param.RemoveDetailsInProduction {
		param.Details = ""
	}
	return ForbiddenError{param.Message, param.Details}
}

func (e ForbiddenError) Error() string {
	return e.msg
}

func (e ForbiddenError) Details() string {
	return e.details
}
