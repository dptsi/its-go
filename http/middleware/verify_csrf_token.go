package middleware

import (
	"fmt"
	"regexp"

	"bitbucket.org/dptsi/its-go/app/errors"
	"bitbucket.org/dptsi/its-go/contracts"
	"bitbucket.org/dptsi/its-go/http"
	"bitbucket.org/dptsi/its-go/web"
)

var errInvalidCSRFToken = errors.NewForbidden(errors.ForbiddenParam{
	Message: "invalid_csrf_token",
	Details: "Ambil token CSRF dari cookie dan masukkan ke header X-CSRF-TOKEN",
})

type VerifyCSRFToken struct {
	cfg            http.CSRFConfig
	sessionService contracts.SessionService
	pathException  [](*regexp.Regexp)
}

func NewVerifyCSRFToken(cfg http.CSRFConfig, sessionService contracts.SessionService) (*VerifyCSRFToken, error) {
	exceptions := make([](*regexp.Regexp), len(cfg.Except))
	for i, p := range cfg.Except {
		regex, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("new verify csrf token: error when compiling regex \"%s\": %w", p, err)
		}

		exceptions[i] = regex
	}

	return &VerifyCSRFToken{
		cfg:            cfg,
		sessionService: sessionService,
		pathException:  exceptions,
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

		for _, regex := range m.pathException {
			if !shouldCheckToken {
				break
			}
			if regex.MatchString(req.URL.Path) {
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
