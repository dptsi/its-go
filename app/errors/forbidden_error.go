package errors

import "net/http"

type ForbiddenParam struct {
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

// Forbidden is an error that occurs when the request is forbidden
type Forbidden struct {
	msg     string
	details string
}

func NewForbidden(param ForbiddenParam) Forbidden {
	if param.Code == 0 {
		param.Code = http.StatusForbidden
	}
	if param.Message == "" {
		param.Message = "forbidden"
	}
	if param.RemoveDetailsInProduction {
		param.Details = ""
	}
	return Forbidden{param.Message, param.Details}
}

func (e Forbidden) Error() string {
	return e.msg
}

func (e Forbidden) Details() string {
	return e.details
}
