// Package errors contains all the errors that can be returned by the application
// This error handled and sent to the client by the framework
//
// You can pass this error to the framework by using the following code:
//
//		ctx.Error(errors.ErrRecordNotFound)
//	 ctx.Abort()
//	 return
package errors

import "net/http"

// AggregateVersionMismatchParam is a struct that contains parameters for AggregateVersionMismatch
type AggregateVersionMismatchParam struct {
	// Code is the status code of the error
	//
	// Default to HTTP status code 409
	Code int

	// Msg is the message of the error
	//
	// Default to "aggregate_version_mismatch"
	Msg string
}

// AggregateVersionMismatch is an error that occurs
// when the aggregate version in database is not the same
// as the aggregate version in the command
type AggregateVersionMismatch struct {
	code int
	msg  string
}

func NewAggregateVersionMismatch(param AggregateVersionMismatchParam) AggregateVersionMismatch {
	if param.Code == 0 {
		param.Code = http.StatusConflict
	}
	if param.Msg == "" {
		param.Msg = "aggregate_version_mismatch"
	}
	return AggregateVersionMismatch{param.Code, param.Msg}
}

func (e AggregateVersionMismatch) Code() int {
	return e.code
}

func (e AggregateVersionMismatch) Error() string {
	return e.msg
}
