package contracts

import (
	"context"

	"github.com/dptsi/its-go/web"
	"github.com/samber/do"
)

type ApplicationServices struct {
	Auth        AuthService
	ActivityLog ActivityLogService
	Crypt       CryptService
	Database    DatabaseService
	Event       EventService
	Logging     LoggingService
	Middleware  MiddlewareService
	Module      ModuleService
	Session     SessionService
	WebEngine   *web.Engine
}

type Application interface {
	Context() context.Context
	Config() map[string]interface{}
	ListProvidedServices() []string
	Injector() *do.Injector
	Services() ApplicationServices
}
