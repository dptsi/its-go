package middleware

import (
	"fmt"

	"bitbucket.org/dptsi/base-go-libraries/app/errors"
	"bitbucket.org/dptsi/base-go-libraries/auth"
	"bitbucket.org/dptsi/base-go-libraries/contracts"
)

type ActiveRole struct {
	service auth.Service
}

func (m *ActiveRole) ActiveRoleIn(roles ...string) contracts.HandlerFunc {
	return func(ctx contracts.WebFrameworkContext) {
		u := m.service.User(ctx)
		for _, role := range roles {
			if role == u.ActiveRole() {
				ctx.Next()
				return
			}
		}

		msg := fmt.Sprintf("current user active role (%s) doesn't have permission to access this resource", u.ActiveRoleName())
		// details := fmt.Sprintf("allowed role to access this resource are: %s", strings.Join(roles, ", "))
		details := ""
		ctx.Error(errors.NewForbidden(errors.ForbiddenParam{
			Message: msg,
			Details: details,
		}))
		ctx.Abort()
	}
}

func (m *ActiveRole) ActiveRoleHasPermission(neededPermission string) contracts.HandlerFunc {
	return func(ctx contracts.WebFrameworkContext) {
		u := m.service.User(ctx)
		if u.HasPermission(neededPermission) {
			ctx.Next()
			return
		}

		msg := fmt.Sprintf("current user active role (%s) doesn't have permission to access this resource", u.ActiveRoleName())
		// details := fmt.Sprintf("permission to access this resource is: %s", neededPermission)
		details := ""
		ctx.Error(errors.NewForbidden(errors.ForbiddenParam{
			Message: msg,
			Details: details,
		}))
		ctx.Abort()
	}
}
