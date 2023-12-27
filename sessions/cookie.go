package sessions

import (
	"net/http"

	"bitbucket.org/dptsi/go-framework/web"
)

type SessionsConfig struct {
	Name           string
	CsrfCookieName string
	MaxAge         int
	Path           string
	Domain         string
	Secure         bool
}

type CookieUtil struct {
	cfg SessionsConfig
}

func NewCookieUtil(cfg SessionsConfig) *CookieUtil {
	return &CookieUtil{
		cfg: cfg,
	}
}

func (c *CookieUtil) Write(ctx *web.Context, sess *Data) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	// Set session cookie
	ctx.SetCookie(c.cfg.Name, sess.Id(), c.cfg.MaxAge, c.cfg.Path, c.cfg.Domain, c.cfg.Secure, true)
	if c.cfg.CsrfCookieName == "" {
		c.cfg.CsrfCookieName = "CSRF-TOKEN"
	}
	ctx.SetCookie(c.cfg.CsrfCookieName, sess.CSRFToken(), c.cfg.MaxAge, c.cfg.Path, c.cfg.Domain, c.cfg.Secure, false)
}
