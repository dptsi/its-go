package middleware

import (
	"fmt"
	"strings"

	"bitbucket.org/dptsi/its-go/app/errors"
	"bitbucket.org/dptsi/its-go/contracts"
	"bitbucket.org/dptsi/its-go/web"
)

type UserHasRole struct {
	service contracts.AuthService
}

func NewUserHasRole(service contracts.AuthService) *UserHasRole {
	return &UserHasRole{
		service: service,
	}
}

type UserHasRoleParam struct {
	Roles []string
}

func (a *UserHasRole) Handle(param interface{}) web.HandlerFunc {
	return func(ctx *web.Context) {
		userHasRoleParam, ok := param.(UserHasRoleParam)
		if !ok {
			ctx.Error(fmt.Errorf("user has role middleware: handle: invalid parameter type"))
			ctx.Abort()
			return
		}
		user, err := a.service.User(ctx)

		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}
		userRoles := make([]string, len(user.Roles()))
		for i, role := range user.Roles() {
			userRoles[i] = role.Id
		}

		for _, neededRole := range userHasRoleParam.Roles {
			for _, userRole := range userRoles {
				if neededRole == userRole {
					ctx.Next()
					return
				}
			}
		}

		ctx.Error(errors.NewForbidden(errors.ForbiddenParam{
			Message: "user doesn't have permission to access this resource",
			Details: fmt.Sprintf(
				"allowed roles: %s, user roles: %s",
				strings.Join(userHasRoleParam.Roles, ","),
				strings.Join(userRoles, ","),
			),
		}))
		ctx.Abort()
	}
}
