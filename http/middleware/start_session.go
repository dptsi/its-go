package middleware

import (
	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/web"
)

type StartSession struct {
	service contracts.SessionService
}

func NewStartSession(service contracts.SessionService) *StartSession {
	return &StartSession{
		service: service,
	}
}

func (m *StartSession) Execute(ctx *web.Context) {
	if err := m.service.Start(ctx); err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.Next()
}
