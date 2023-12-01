package middleware

import "bitbucket.org/dptsi/base-go-libraries/app/errors"

var unauthorizedError = errors.NewUnauthorized(errors.UnauthorizedParam{})
