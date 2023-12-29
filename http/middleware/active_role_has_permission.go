package middleware

import (
	"fmt"

	"bitbucket.org/dptsi/go-framework/app/errors"
	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/web"
)

type ActiveRoleHasPermission struct {
	service contracts.AuthService
}

func NewActiveRoleHasPermission(service contracts.AuthService) *ActiveRoleHasPermission {
	return &ActiveRoleHasPermission{
		service: service,
	}
}

type ActiveRoleHasPermissionParam struct {
	NeededPermission string
}

func (a *ActiveRoleHasPermission) Handle(param interface{}) web.HandlerFunc {
	return func(ctx *web.Context) {
		activeRoleHasPermissionParam, ok := param.(ActiveRoleHasPermissionParam)
		if !ok {
			ctx.Error(fmt.Errorf("active role in middleware: handle: invalid parameter type"))
			ctx.Abort()
			return
		}
		user, err := a.service.User(ctx)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}
		if user.HasPermission(activeRoleHasPermissionParam.NeededPermission) {
			ctx.Next()
			return
		}

		msg := fmt.Sprintf("current user active role (%s) doesn't have permission to access this resource", user.ActiveRoleName())
		// details := fmt.Sprintf("permission to access this resource is: %s", neededPermission)
		details := ""
		ctx.Error(errors.NewForbidden(errors.ForbiddenParam{
			Message: msg,
			Details: details,
		}))
		ctx.Abort()
	}
}
