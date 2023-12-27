package auth

import (
	"encoding/json"
	"fmt"
	"strings"

	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/models"
	"bitbucket.org/dptsi/go-framework/sessions"
	"bitbucket.org/dptsi/go-framework/web"
)

type UserSessionData struct {
	Id                string        `json:"id"`
	Name              string        `json:"name"`
	PreferredUsername string        `json:"preferred_username"`
	Email             string        `json:"email"`
	Picture           string        `json:"picture"`
	ActiveRole        string        `json:"active_role"`
	Roles             []models.Role `json:"roles"`
}

const userContextKey = "auth.user"

type SessionGuard struct {
	storage      contracts.SessionStorage
	cookieWriter contracts.SessionCookieWriter
}

func NewSessionGuard(storage contracts.SessionStorage, cookieWriter contracts.SessionCookieWriter) *SessionGuard {
	return &SessionGuard{
		storage:      storage,
		cookieWriter: cookieWriter,
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

	sess := sessions.Default(ctx)
	userIf, ok := sess.Get("user")
	if !ok {
		return nil
	}
	userJson, ok := userIf.(string)
	if !ok {
		return nil
	}
	var userData UserSessionData
	err := json.Unmarshal([]byte(userJson), &userData)
	if err != nil {
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
	user.SetActiveRole(userData.ActiveRole)
	g.SetUser(ctx, user)

	return user
}

func (g *SessionGuard) SetUser(ctx *web.Context, user *models.User) {
	ctx.Set(userContextKey, user)
}

func (g *SessionGuard) Login(ctx *web.Context, user *models.User) error {
	sess, err := g.updateSession(ctx, user)
	if err != nil {
		return fmt.Errorf("session guard: login: %w", err)
	}

	g.cookieWriter.Write(ctx, sess)

	return nil
}

func (g *SessionGuard) Logout(ctx *web.Context) error {
	sess, err := g.updateSession(ctx, nil)
	if err != nil {
		return fmt.Errorf("session guard: logout: %w", err)
	}

	g.cookieWriter.Write(ctx, sess)

	return nil
}

func (g *SessionGuard) updateSession(ctx *web.Context, user *models.User) (updated *sessions.Data, err error) {
	data := sessions.Default(ctx)

	if user != nil {
		userSessionData := UserSessionData{
			Id:                strings.ToLower(user.Id()),
			ActiveRole:        user.ActiveRole(),
			Name:              user.Name(),
			PreferredUsername: user.PreferredUsername(),
			Email:             user.Email(),
			Picture:           user.Picture(),
			Roles:             user.Roles(),
		}
		userJson, err := json.Marshal(userSessionData)
		if err != nil {
			return data, fmt.Errorf("session guard: update session: %w", err)
		}
		data.Set("user", string(userJson))
	} else {
		data.Delete("user")
	}

	return data, g.storage.Save(ctx, data)
}
