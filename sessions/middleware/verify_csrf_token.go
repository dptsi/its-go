package middleware

import (
	"bitbucket.org/dptsi/base-go-libraries/app/errors"
	"bitbucket.org/dptsi/base-go-libraries/sessions"
	"github.com/gin-gonic/gin"
)

var errInvalidCSRFToken = errors.NewForbidden(errors.ForbiddenParam{
	Message: "invalid_csrf_token",
	Details: "Ambil token CSRF dari cookie dan masukkan ke header X-CSRF-TOKEN",
})
var methodsWithoutCSRFToken = []string{"GET", "HEAD", "OPTIONS"}

type VerifyCSRFToken struct {
}

func NewVerifyCSRFToken() *VerifyCSRFToken {
	return &VerifyCSRFToken{}
}

func (m *VerifyCSRFToken) VerifyCSRFToken(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sessionCSRFToken := sess.CSRFToken()
	req := ctx.Request
	requestCSRFToken := req.Header.Get("X-CSRF-TOKEN")

	// Skip CSRF token verification for some methods
	for _, method := range methodsWithoutCSRFToken {
		if req.Method == method {
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
