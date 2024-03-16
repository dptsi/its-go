package middleware

import (
	"github.com/dptsi/its-go/app/errors"
	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/http"
	"github.com/dptsi/its-go/web"
)

var errInvalidCSRFToken = errors.NewForbidden(errors.ForbiddenParam{
	Message: "invalid_csrf_token",
	Details: "Ambil token CSRF dari cookie dan masukkan ke header X-CSRF-TOKEN",
})

type VerifyCSRFToken struct {
	cfg            http.CSRFConfig
	sessionService contracts.SessionService
	pathException  []string
}

func NewVerifyCSRFToken(cfg http.CSRFConfig, sessionService contracts.SessionService) (*VerifyCSRFToken, error) {
	return &VerifyCSRFToken{
		cfg:            cfg,
		sessionService: sessionService,
		pathException:  cfg.Except,
	}, nil
}

func (m *VerifyCSRFToken) Handle(interface{}) web.HandlerFunc {
	return func(ctx *web.Context) {
		req := ctx.Request
		shouldCheckToken := false
		for _, method := range m.cfg.Methods {
			if req.Method != method {
				continue
			}
			shouldCheckToken = true
			break
		}

		for _, e := range m.pathException {
			if !shouldCheckToken {
				break
			}
			if ctx.FullPath() == e {
				shouldCheckToken = false
			}
		}

		if !shouldCheckToken {
			ctx.Next()
			return
		}

		requestCSRFToken := req.Header.Get("X-CSRF-TOKEN")
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
}
