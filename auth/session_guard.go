package auth

import (
	"encoding/json"
	"fmt"
	"strings"

	"bitbucket.org/dptsi/its-go/contracts"
	"bitbucket.org/dptsi/its-go/models"
	"bitbucket.org/dptsi/its-go/web"
)

type UserSessionData struct {
	Id                string        `json:"id"`
	Name              string        `json:"name"`
	PreferredUsername string        `json:"preferred_username"`
	Email             string        `json:"email"`
	Picture           string        `json:"picture"`
	Roles             []models.Role `json:"roles"`
}

const userContextKey = "auth.user"

type SessionGuard struct {
	service contracts.SessionService
}

func NewSessionGuard(service contracts.SessionService) *SessionGuard {
	return &SessionGuard{
		service: service,
	}
}

func (g *SessionGuard) userFromContext(ctx *web.Context) *models.User {
	uInterface, exist := ctx.Get(userContextKey)
	if !exist {
		return nil
	}
	user, ok := uInterface.(*models.User)
	if !ok {
		return nil
	}

	return user
}

func (g *SessionGuard) User(ctx *web.Context) *models.User {
	user := g.userFromContext(ctx)

	if user != nil {
		return user
	}

	userIf, err := g.service.Get(ctx, "user")
	if err != nil {
		return nil
	}

	userJson, ok := userIf.(string)
	if !ok {
		return nil
	}
	var userData UserSessionData
	if err := json.Unmarshal([]byte(userJson), &userData); err != nil {
		return nil
	}

	user = models.NewUser(userData.Id)
	user.SetEmail(userData.Email)
	user.SetName(userData.Name)
	user.SetPreferredUsername(userData.PreferredUsername)
	user.SetPicture(userData.Picture)
	for _, role := range userData.Roles {
		user.AddRole(role.Id, role.Name, role.Permissions, role.IsDefault)
	}
	g.SetUser(ctx, user)

	return user
}

func (g *SessionGuard) SetUser(ctx *web.Context, user *models.User) {
	ctx.Set(userContextKey, user)
}

func (g *SessionGuard) Login(ctx *web.Context, user *models.User) error {
	if err := g.updateSession(ctx, user); err != nil {
		return fmt.Errorf("session guard: login: %w", err)
	}

	return nil
}

func (g *SessionGuard) Logout(ctx *web.Context) error {
	if err := g.updateSession(ctx, nil); err != nil {
		return fmt.Errorf("session guard: logout: %w", err)
	}

	return nil
}

func (g *SessionGuard) updateSession(ctx *web.Context, user *models.User) error {
	if user != nil {
		userSessionData := UserSessionData{
			Id:                strings.ToLower(user.Id()),
			Name:              user.Name(),
			PreferredUsername: user.PreferredUsername(),
			Email:             user.Email(),
			Picture:           user.Picture(),
			Roles:             user.Roles(),
		}
		userJson, err := json.Marshal(userSessionData)
		if err != nil {
			return fmt.Errorf("session guard: update session: %w", err)
		}
		return g.service.Put(ctx, "user", string(userJson))
	}

	return g.service.Put(ctx, "user", nil)
}
