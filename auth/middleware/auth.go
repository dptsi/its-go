package middleware

import (
	"encoding/json"

	internalContract "bitbucket.org/dptsi/base-go-libraries/auth/internal/contracts"
	"bitbucket.org/dptsi/base-go-libraries/auth/internal/utils"
	"bitbucket.org/dptsi/base-go-libraries/contracts"
	"bitbucket.org/dptsi/base-go-libraries/sessions"
)

func Auth() contracts.HandlerFunc {
	return func(ctx contracts.WebFrameworkContext) {
		sess := sessions.Default(ctx)
		userIf, ok := sess.Get("user")
		if !ok {
			ctx.Error(unauthorizedError)
			ctx.Abort()
			return
		}
		userJson, ok := userIf.(string)
		if !ok {
			ctx.Error(unauthorizedError)
			ctx.Abort()
			return
		}
		var userData internalContract.UserSessionData
		err := json.Unmarshal([]byte(userJson), &userData)
		if err != nil {
			ctx.Error(unauthorizedError)
			ctx.Abort()
			return
		}

		u := contracts.NewUser(userData.Id)
		u.SetEmail(userData.Email)
		u.SetName(userData.Name)
		u.SetPreferredUsername(userData.PreferredUsername)
		u.SetPicture(userData.Picture)
		for _, role := range userData.Roles {
			u.AddRole(role.Id, role.Name, role.Permissions, role.IsDefault)
		}
		u.SetActiveRole(userData.ActiveRole)

		ctx.Set(utils.UserKey, u)
		ctx.Next()
	}
}
