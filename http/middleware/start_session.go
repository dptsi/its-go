package middleware

import (
	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/web"
)

type StartSession struct {
	service contracts.SessionService
}

func NewStartSession(service contracts.SessionService) *StartSession {
	return &StartSession{
		service: service,
	}
}

func (m *StartSession) Handle(interface{}) web.HandlerFunc {
	return func(ctx *web.Context) {
		if err := m.service.Start(ctx); err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
