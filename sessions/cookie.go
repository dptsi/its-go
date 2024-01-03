package sessions

import (
	"net/http"
	"net/url"

	"bitbucket.org/dptsi/its-go/contracts"
	"bitbucket.org/dptsi/its-go/web"
)

type CookieUtil struct {
	cfg CookieConfig
}

func NewCookieUtil(cfg CookieConfig) *CookieUtil {
	return &CookieUtil{
		cfg: cfg,
	}
}

func (c *CookieUtil) Write(ctx *web.Context, data contracts.SessionData) {
	path := c.cfg.Path
	if path == "" {
		path = "/"
	}
	// Set session cookie
	http.SetCookie(
		ctx.Writer,
		&http.Cookie{
			Name:     c.cfg.Name,
			Value:    url.QueryEscape(data.Id()),
			Path:     path,
			Domain:   c.cfg.Domain,
			Expires:  data.ExpiredAt(),
			SameSite: http.SameSiteLaxMode,
			Secure:   c.cfg.Secure,
			HttpOnly: true,
		},
	)
	if c.cfg.CsrfCookieName == "" {
		c.cfg.CsrfCookieName = "CSRF-TOKEN"
	}
	http.SetCookie(
		ctx.Writer,
		&http.Cookie{
			Name:     c.cfg.CsrfCookieName,
			Value:    url.QueryEscape(data.CSRFToken()),
			Path:     path,
			Domain:   c.cfg.Domain,
			Expires:  data.ExpiredAt(),
			SameSite: http.SameSiteLaxMode,
			Secure:   c.cfg.Secure,
			HttpOnly: false,
		},
	)
}
