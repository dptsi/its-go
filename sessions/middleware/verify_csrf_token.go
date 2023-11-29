package middleware

import (
	"bitbucket.org/dptsi/base-go-libraries/app/errors"
	"bitbucket.org/dptsi/base-go-libraries/sessions"
	"github.com/gin-gonic/gin"
)

var errInvalidCSRFToken = errors.NewForbidden(errors.ForbiddenParam{
	Message: "invalid_csrf_token",
})
var methodsWithoutCSRFToken = []string{"GET", "HEAD", "OPTIONS"}

func VerifyCSRFToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sess := sessions.Default(ctx)
		sessionCSRFToken := sess.CSRFToken()
		requestCSRFToken := ctx.Request.Header.Get("X-CSRF-TOKEN")

		// Skip CSRF token verification for some methods
		for _, method := range methodsWithoutCSRFToken {
			if ctx.Request.Method == method {
				ctx.Next()
				return
			}
		}

		if sessionCSRFToken == "" || sessionCSRFToken != requestCSRFToken {
			ctx.Error(errInvalidCSRFToken)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
