package middleware

import (
	"fmt"

	"bitbucket.org/dptsi/go-framework/app/errors"
	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/web"
)

type ActiveRole struct {
	service contracts.AuthService
}

func NewActiveRole(service contracts.AuthService) *ActiveRole {
	return &ActiveRole{
		service: service,
	}
}

func (a *ActiveRole) ActiveRoleIn(roles ...string) web.HandlerFunc {
	return func(ctx *web.Context) {
		user, err := a.service.User(ctx)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}
		for _, role := range roles {
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

func (a *ActiveRole) ActiveRoleHasPermission(neededPermission string) web.HandlerFunc {
	return func(ctx *web.Context) {
		user, err := a.service.User(ctx)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}
		if user.HasPermission(neededPermission) {
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
