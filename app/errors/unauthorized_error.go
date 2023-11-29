package errors

import "net/http"

type UnauthorizedErrorParam struct {
	// Code is the status code of the error
	//
	// Default to HTTP status code 401
	Code int
}

// UnauthorizedError is an error that occurs when the request is unauthorized
type UnauthorizedError struct {
	code int
}

func NewUnauthorizedError(param UnauthorizedErrorParam) UnauthorizedError {
	if param.Code == 0 {
		param.Code = http.StatusUnauthorized
	}
	return UnauthorizedError{param.Code}
}

func (e UnauthorizedError) Error() string {
	return "unauthorized"
}
