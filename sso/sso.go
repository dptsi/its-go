package sso

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dptsi/its-go/models"
	"github.com/dptsi/its-go/oidc"
	"github.com/dptsi/its-go/web"
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

type Role struct {
	RoleId    string       `json:"role_id"`
	RoleName  string       `json:"role_name"`
	IsDefault stringAsBool `json:"is_default"`
}

type Group struct {
	GroupId   string `json:"group_id"`
	GroupName string `json:"group_name"`
}

type userInfoRaw struct {
	Sub               string       `json:"sub"`
	Name              string       `json:"name"`
	Email             string       `json:"email"`
	EmailVerified     stringAsBool `json:"email_verified"`
	Phone             string       `json:"phone"`
	PhoneVerified     stringAsBool `json:"phone_verified"`
	Picture           string       `json:"picture"`
	PreferredUsername string       `json:"preferred_username"`
	Group             []Group      `json:"group"`
	Roles             []Role       `json:"role"`
	Resource          interface{}  `json:"resource"`
	OriginalUserId    *string      `json:"original_user_id"`
	RegId             *string      `json:"reg_id"`
}

type Roles map[string]Role
type RolePermissions map[string][]string
type GroupRoleMapping map[string]string
type Config struct {
	// Jika true, maka role dan permission user akan diambil dari SSO
	IsRoleFromSso bool

	// Jika true, maka group akan dimapping langsung ke role
	IsGroupMappedDirectlyToRole bool

	// Mekanisme mapping group ke role
	// Hanya digunakan ketika isGroupMappedDirectlyToRole = false

	// Referensi role yang tersedia
	Roles

	RolePermissions // Mapping role ke permission

	// Mapping group ke role ID
	GroupRoleMapping
}

type Sso struct {
	client *oidc.Client
	cfg    Config
}

func NewSso(appCfg map[string]interface{}, client *oidc.Client) *Sso {
	cfg, ok := appCfg["sso"].(Config)
	if !ok {
		return &Sso{
			client: client,
			cfg: Config{
				IsRoleFromSso:               true,
				IsGroupMappedDirectlyToRole: true,

				Roles:            nil,
				RolePermissions:  nil,
				GroupRoleMapping: nil,
			},
		}
	}

	return &Sso{
		client: client,
		cfg:    cfg,
	}
}

func (s *Sso) GetUserFromAuthorizationCode(ctx *web.Context, code string, state string) (*models.User, error) {
	token, _, err := s.client.ExchangeCodeForToken(ctx, code, state)
	if err != nil {
		return nil, fmt.Errorf("get user from myits sso failed: %w", err)
	}
	userInfo, err := s.userInfo(ctx, s.client, oauth2.StaticTokenSource(token))
	if err != nil {
		return nil, fmt.Errorf("get user from myits sso failed: %w", err)
	}

	user := models.NewUser(userInfo.Sub)
	user.SetName(userInfo.Name)
	user.SetPreferredUsername(userInfo.PreferredUsername)
	user.SetEmail(userInfo.Email)
	user.SetPhone(userInfo.Phone)
	user.SetPicture(userInfo.Picture)
	user.SetImpersonatorId(userInfo.OriginalUserId)

	if s.cfg.IsRoleFromSso {
		for _, r := range userInfo.Roles {
			permissions, ok := s.cfg.RolePermissions[r.RoleId]
			if !ok {
				permissions = make([]string, 0)
			}

			user.AddRole(r.RoleId, r.RoleName, permissions, bool(r.IsDefault), nil)
		}
	}

	for _, g := range userInfo.Group {
		if s.cfg.IsGroupMappedDirectlyToRole {
			user.AddRole(g.GroupId, g.GroupName, nil, false, nil)
			continue
		}

		roleId, ok := s.cfg.GroupRoleMapping[g.GroupId]
		if !ok {
			continue
		}
		role, ok := s.cfg.Roles[roleId]
		if !ok {
			continue
		}
		permissions, ok := s.cfg.RolePermissions[roleId]
		if !ok {
			permissions = make([]string, 0)
		}

		user.AddRole(role.RoleId, role.RoleName, permissions, bool(role.IsDefault), nil)
	}

	return user, nil
}

func (s *Sso) userInfo(ctx context.Context, oidcClient *oidc.Client, tokenSource oauth2.TokenSource) (*userInfoRaw, error) {
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
