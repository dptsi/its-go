package middleware

import (
	"github.com/dptsi/its-go/app/errors"
	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/web"
)

type Auth struct {
	service contracts.AuthService
}

func NewAuth(service contracts.AuthService) *Auth {
	return &Auth{service}
}

func (a *Auth) Handle(interface{}) web.HandlerFunc {
	return func(ctx *web.Context) {
		user, err := a.service.User(ctx)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}
		if user == nil {
			ctx.Error(errors.NewUnauthorized(errors.UnauthorizedParam{}))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
