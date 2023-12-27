package contracts

import (
	"bitbucket.org/dptsi/go-framework/models"
	"bitbucket.org/dptsi/go-framework/web"
)

type AuthService interface {
	Login(ctx *web.Context, user *models.User) error
	Logout(ctx *web.Context) error
	User(ctx *web.Context) (*models.User, error)
}

type AuthGuard interface {
	User(ctx *web.Context) *models.User
	SetUser(ctx *web.Context, user *models.User)
}

type StatefulAuthGuard interface {
	AuthGuard
	Login(ctx *web.Context, user *models.User) error
	Logout(ctx *web.Context) error
}
