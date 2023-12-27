package myitssso

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/models"
	"bitbucket.org/dptsi/go-framework/oidc"
	"golang.org/x/oauth2"
)

type stringAsBool bool

func (sb *stringAsBool) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "1", `"1"`:
		*sb = true
	case "0", `"0"`:
		*sb = false
	default:
		return errors.New("invalid value for boolean")
	}
	return nil
}

type role struct {
	RoleId    string       `json:"role_id"`
	RoleName  string       `json:"role_name"`
	IsDefault stringAsBool `json:"is_default"`
}

type resource struct {
	Path string `json:"path"`
}

type userInfoRaw struct {
	Sub               string       `json:"sub"`
	Name              string       `json:"name"`
	Email             string       `json:"email"`
	EmailVerified     stringAsBool `json:"email_verified"`
	Picture           string       `json:"picture"`
	PreferredUsername string       `json:"preferred_username"`
	Roles             []role       `json:"role"`
	Resource          interface{}  `json:"resource"`
}

func GetUserFromAuthorizationCode(ctx context.Context, oidcClient *oidc.Client, sess contracts.SessionData, code string, state string) (*models.User, error) {
	token, _, err := oidcClient.ExchangeCodeForToken(ctx, sess, code, state)
	if err != nil {
		return nil, fmt.Errorf("get user from myits sso failed: %w", err)
	}
	// fmt.Println("token", token.AccessToken)
	userInfo, err := userInfo(ctx, oidcClient, oauth2.StaticTokenSource(token))
	if err != nil {
		return nil, fmt.Errorf("get user from myits sso failed: %w", err)
	}

	user := models.NewUser(userInfo.Sub)
	user.SetName(userInfo.Name)
	user.SetPreferredUsername(userInfo.PreferredUsername)
	user.SetEmail(userInfo.Email)
	user.SetPicture(userInfo.Picture)
	for _, r := range userInfo.Roles {
		permissions := make([]string, 0)
		userInfoResourceInterface, ok := userInfo.Resource.(map[string]interface{})
		var userInfoResource map[string][]resource
		// Convert to JSON first before parsing
		if ok {
			tmp, _ := json.Marshal(userInfoResourceInterface)
			json.Unmarshal(tmp, &userInfoResource)
		}

		resources, ok := userInfoResource[r.RoleName]
		if ok {
			permissions = make([]string, len(resources))
			for i, resource := range resources {
				permissions[i] = resource.Path
			}
		}

		user.AddRole(r.RoleId, r.RoleName, permissions, bool(r.IsDefault))
	}

	return user, nil
}

func userInfo(ctx context.Context, oidcClient *oidc.Client, tokenSource oauth2.TokenSource) (*userInfoRaw, error) {
	userInfoURL := oidcClient.UserInfoEndpoint()
	if userInfoURL == "" {
		return nil, errors.New("oidc: user info endpoint is not supported by this provider")
	}

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("oidc: create GET request: %w", err)
	}

	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("oidc: get access token: %w", err)
	}
	token.SetAuthHeader(req)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("oidc: user info request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, body)
	}

	var userInfo userInfoRaw
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("oidc: failed to decode userinfo: %v", err)
	}

	return &userInfo, nil
}
