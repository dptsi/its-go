package errors

import "net/http"

type UnauthorizedParam struct {
	// Code is the status code of the error
	//
	// Default to HTTP status code 401
	Code int
}

// Unauthorized is an error that occurs when the request is unauthorized
type Unauthorized struct {
	code int
}

func NewUnauthorized(param UnauthorizedParam) Unauthorized {
	if param.Code == 0 {
		param.Code = http.StatusUnauthorized
	}
	return Unauthorized{param.Code}
}

func (e Unauthorized) Error() string {
	return "unauthorized"
}
