package middleware

import "bitbucket.org/dptsi/go-framework/app/errors"

var unauthorizedError = errors.NewUnauthorized(errors.UnauthorizedParam{})
