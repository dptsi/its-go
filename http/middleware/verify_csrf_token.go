package middleware

import (
	"net/http"

	"bitbucket.org/dptsi/go-framework/app/errors"
	"bitbucket.org/dptsi/go-framework/sessions"
	"bitbucket.org/dptsi/go-framework/web"
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

func (m *VerifyCSRFToken) Execute(ctx *web.Context) {
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

// CSRF cookie godoc
// @Summary		Rute dummy untuk set CSRF-TOKEN cookie
// @Router		/csrf-cookie [get]
// @Tags		CSRF Protection
// @Produce		json
// @Success		200 {object} responses.GeneralResponse{code=int,message=string} "Cookie berhasil diset"
// @Header      default {string} Set-Cookie "CSRF-TOKEN=00000000-0000-0000-0000-000000000000; Path=/"
func CSRFCookieRoute(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    nil,
	})
}
