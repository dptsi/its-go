package auth

import (
	"encoding/json"
	"fmt"
	"strings"

	internalContract "bitbucket.org/dptsi/base-go-libraries/auth/internal/contracts"
	"bitbucket.org/dptsi/base-go-libraries/auth/internal/utils"
	"bitbucket.org/dptsi/base-go-libraries/contracts"
	"bitbucket.org/dptsi/base-go-libraries/sessions"
)

type Service struct {
	sessionStorage contracts.SessionStorage
}

func (s *Service) Login(ctx contracts.WebFrameworkContext, u *contracts.User) error {
	sess := sessions.Default(ctx)
	userData := internalContract.UserSessionData{
		Id:                strings.ToLower(u.Id()),
		ActiveRole:        u.ActiveRole(),
		Name:              u.Name(),
		PreferredUsername: u.PreferredUsername(),
		Email:             u.Email(),
		Picture:           u.Picture(),
		Roles:             u.Roles(),
	}
	userJson, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("login service failed: %w", err)
	}
	sess.Set("user", string(userJson))
	return s.sessionStorage.Save(ctx, sess)
}

func (s *Service) Logout(ctx contracts.WebFrameworkContext) error {
	sess := sessions.Default(ctx)
	sess.Delete("user.id")
	sess.Delete("user.active_role")
	sess.Delete("user.roles")
	return s.sessionStorage.Save(ctx, sess)
}

func (s *Service) User(ctx contracts.WebFrameworkContext) *contracts.User {
	uInterface, exist := ctx.Get(utils.UserKey)
	if !exist {
		panic("cannot get user info, forgot to add auth middleware?")
	}
	u, ok := uInterface.(*contracts.User)
	if !ok {
		panic("cannot get user info, forgot to add auth middleware?")
	}

	return u
}
