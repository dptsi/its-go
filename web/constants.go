package web

const (
	validationError     string = "validation_error"
	internalServerError string = "internal_server_error"
	forbiddenError      string = "forbidden"
	unauthorizedError   string = "unauthorized"
)

var statusCode = map[string]int{
	validationError:     9000,
	unauthorizedError:   9001,
	forbiddenError:      9003,
	internalServerError: 9005,
}
