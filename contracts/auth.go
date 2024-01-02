package contracts

import (
	"bitbucket.org/dptsi/go-framework/app"
	"bitbucket.org/dptsi/go-framework/models"
	"bitbucket.org/dptsi/go-framework/web"
)

type AuthService interface {
	Login(ctx *web.Context, user *models.User) error
	Logout(ctx *web.Context) error
	User(ctx *web.Context) (*models.User, error)
	RegisterGuard(name string, constructor AuthGuardConstructor) error
}

type AuthGuard interface {
	User(ctx *web.Context) *models.User
	SetUser(ctx *web.Context, user *models.User)
}

type AuthGuardConstructor = func(application *app.Application) (AuthGuard, error)

type StatefulAuthGuard interface {
	AuthGuard
	Login(ctx *web.Context, user *models.User) error
	Logout(ctx *web.Context) error
}
