package middleware

import (
	"errors"
	"fmt"

	"bitbucket.org/dptsi/base-go-libraries/contracts"
	"bitbucket.org/dptsi/base-go-libraries/sessions"
)

type StartSession struct {
	sessionsConfig sessions.SessionsConfig
	storage        contracts.SessionStorage
	cookieUtil     sessions.CookieUtil
}

func (m *StartSession) Execute(ctx contracts.WebFrameworkContext) {
	if m.storage == nil {
		err := errors.New("session storage not configured")
		ctx.Error(fmt.Errorf("start session middleware: %w", err))
		ctx.Abort()
	}

	// Initialize session data
	var data *sessions.Data
	sessionId, err := ctx.Cookie(m.sessionsConfig.Name)

	if err == nil {
		// Get session data from storage
		sessInterface, err := m.storage.Get(ctx, sessionId)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}
		sessionData, ok := sessInterface.(*sessions.Data)
		if !ok {
			ctx.Error(fmt.Errorf("start session middleware: session data is not of type session.Data"))
			ctx.Abort()
			return
		}
		if sessionData != nil {
			data = sessionData
		}
	}
	if data == nil {
		data = sessions.NewEmptyData(int64(m.sessionsConfig.MaxAge))
		if err := m.storage.Save(ctx, data); err != nil {
			ctx.Error(fmt.Errorf("start session middleware: %w", err))
			ctx.Abort()
		}
	}
	ctx.Set("session", data)
	m.cookieUtil.AddSessionCookieToResponse(ctx, data)
	ctx.Next()
}
