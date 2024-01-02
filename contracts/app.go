package contracts

import (
	"context"

	"bitbucket.org/dptsi/go-framework/web"
	"github.com/samber/do"
)

type ApplicationServices struct {
	Auth       AuthService
	Database   DatabaseService
	Event      EventService
	Middleware MiddlewareService
	Module     ModuleService
	Session    SessionService
	WebEngine  *web.Engine
}

type Application interface {
	Context() context.Context
	Config() map[string]interface{}
	ListProvidedServices() []string
	Injector() *do.Injector
	Services() ApplicationServices
}
