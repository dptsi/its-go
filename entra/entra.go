package entra

import (
	"fmt"

	"bitbucket.org/dptsi/go-framework/models"
	"bitbucket.org/dptsi/go-framework/oidc"
	"bitbucket.org/dptsi/go-framework/web"
)

type entraIDClaim struct {
	ObjectId          string   `json:"oid"`
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	PreferredUsername string   `json:"preferred_username"`
	Roles             []string `json:"roles"`
}

func GetUserFromAuthorizationCode(ctx *web.Context, oidcClient *oidc.Client, code string, state string) (*models.User, error) {
	_, IDToken, err := oidcClient.ExchangeCodeForToken(ctx, code, state)
	if err != nil {
		return nil, fmt.Errorf("get user from entra id failed: %w", err)
	}

	var claims entraIDClaim
	if err := IDToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("get user from entra id failed: %w", err)
	}

	user := models.NewUser(claims.ObjectId)
	user.SetName(claims.Name)
	user.SetPreferredUsername(claims.PreferredUsername)
	user.SetEmail(claims.Email)
	for i, r := range claims.Roles {
		user.AddRole(r, r, make([]string, 0), i == 0)
	}

	return user, nil
}
