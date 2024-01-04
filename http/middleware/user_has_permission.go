package middleware

import (
	"fmt"

	"bitbucket.org/dptsi/its-go/app/errors"
	"bitbucket.org/dptsi/its-go/contracts"
	"bitbucket.org/dptsi/its-go/web"
)

type UserHasPermission struct {
	service contracts.AuthService
}

func NewUserHasPermission(service contracts.AuthService) *UserHasPermission {
	return &UserHasPermission{
		service: service,
	}
}

type UserHasPermissionParam struct {
	NeededPermission string
}

func (a *UserHasPermission) Handle(param interface{}) web.HandlerFunc {
	return func(ctx *web.Context) {
		userHasPermissionParam, ok := param.(UserHasPermissionParam)
		if !ok {
			ctx.Error(fmt.Errorf("user has permission middleware: invalid parameter type"))
			ctx.Abort()
			return
		}
		user, err := a.service.User(ctx)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}
		if user.HasPermission(userHasPermissionParam.NeededPermission) {
			ctx.Next()
			return
		}

		ctx.Error(errors.NewForbidden(errors.ForbiddenParam{
			Message: "user doesn't have permission to access this resource",
			Details: fmt.Sprintf("needed permission: %s", userHasPermissionParam.NeededPermission),
		}))
		ctx.Abort()
	}
}
