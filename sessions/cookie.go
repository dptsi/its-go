package sessions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddSessionCookieToResponseAttributes struct {
	Name           string
	CsrfCookieName string
	MaxAge         int
	Path           string
	Domain         string
	Secure         bool
}

func AddSessionCookieToResponse(ctx *gin.Context, attr AddSessionCookieToResponseAttributes, sess *Data) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	// Set session cookie
	ctx.SetCookie(attr.Name, sess.Id(), attr.MaxAge, attr.Path, attr.Domain, attr.Secure, true)
	if attr.CsrfCookieName == "" {
		attr.CsrfCookieName = "CSRF-TOKEN"
	}
	ctx.SetCookie(attr.CsrfCookieName, sess.CSRFToken(), attr.MaxAge, attr.Path, attr.Domain, attr.Secure, false)
}
