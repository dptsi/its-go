package middleware

import (
	"fmt"

	"bitbucket.org/dptsi/go-framework/app/errors"
	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/web"
)

type ActiveRoleIn struct {
	service contracts.AuthService
}

func NewActiveRoleIn(service contracts.AuthService) *ActiveRoleIn {
	return &ActiveRoleIn{
		service: service,
	}
}

type ActiveRoleInParam struct {
	Roles []string
}

func (a *ActiveRoleIn) Handle(param interface{}) web.HandlerFunc {
	return func(ctx *web.Context) {
		activeRoleInParam, ok := param.(ActiveRoleInParam)
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
		for _, role := range activeRoleInParam.Roles {
			if role == user.ActiveRole() {
				ctx.Next()
				return
			}
		}

		msg := fmt.Sprintf("current user active role (%s) doesn't have permission to access this resource", user.ActiveRoleName())
		// details := fmt.Sprintf("allowed role to access this resource are: %s", strings.Join(roles, ", "))
		details := ""
		ctx.Error(errors.NewForbidden(errors.ForbiddenParam{
			Message: msg,
			Details: details,
		}))
		ctx.Abort()
	}
}
