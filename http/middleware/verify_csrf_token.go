package middleware

import (
	"net/http"

	"bitbucket.org/dptsi/go-framework/app/errors"
	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/web"
	"github.com/gin-gonic/gin"
)

var errInvalidCSRFToken = errors.NewForbidden(errors.ForbiddenParam{
	Message: "invalid_csrf_token",
	Details: "Ambil token CSRF dari cookie dan masukkan ke header X-CSRF-TOKEN",
})
var methodsWithoutCSRFToken = []string{"GET", "HEAD", "OPTIONS"}

type VerifyCSRFToken struct {
	sessionService contracts.SessionService
}

func NewVerifyCSRFToken(sessionService contracts.SessionService) *VerifyCSRFToken {
	return &VerifyCSRFToken{sessionService}
}

func (m *VerifyCSRFToken) Execute(ctx *web.Context) {
	req := ctx.Request
	requestCSRFToken := req.Header.Get("X-CSRF-TOKEN")

	// Skip CSRF token verification for some methods
	for _, method := range methodsWithoutCSRFToken {
		if req.Method == method {
			ctx.Next()
			return
		}
	}

	if requestCSRFToken == "" {
		ctx.Error(errInvalidCSRFToken)
		ctx.Abort()
		return
	}

	isMatch, err := m.sessionService.IsTokenMatch(ctx, requestCSRFToken)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	if !isMatch {
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
