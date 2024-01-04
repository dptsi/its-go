package templates

const ModuleEntrypointTemplate = `package {{.Name}}

import (
	"bitbucket.org/dptsi/its-go/contracts"
	"{{.BaseMod}}/modules/{{.Name}}/internal/app/providers"
	"{{.BaseMod}}/modules/{{.Name}}/internal/presentation/routes"
)

func SetupModule(mod contracts.Module) {
	providers.ExtendAuth(mod)
	providers.RegisterEvents(mod)
	providers.RegisterMiddlewares(mod)
	providers.RegisterDependencies(mod)
	routes.RegisterRoutes(mod)
}
`

const ModuleDeps = `package providers

import (
	"bitbucket.org/dptsi/its-go/contracts"
)

func RegisterDependencies(mod contracts.Module) {
	// Libraries

	// Queries

	// Repositories

	// Controllers
}
`

const ModuleRoutes = `package routes

import (
	"net/http"

	"bitbucket.org/dptsi/its-go/contracts"
	"bitbucket.org/dptsi/its-go/web"
)

func RegisterRoutes(mod contracts.Module) {
	engine := mod.App().Services().WebEngine
	middlewareService := mod.App().Services().Middleware

	// Routing
	g := engine.Group("/{{.Name}}")

	// Controllers

	// Routes
	g.GET("/hello", middlewareService.Use("auth", nil), func(ctx *web.Context) {
		ctx.JSON(http.StatusOK, web.H{
			"code":    0,
			"message": "hello from module {{.Name}}",
			"data":    nil,
		})
	})
}
`

const ModuleAuth = `package providers

import (
	"bitbucket.org/dptsi/its-go/contracts"
)

func ExtendAuth(mod contracts.Module) {
	// service := mod.App().Services().Auth
	// service.RegisterGuard("example", func(application contracts.Application) (contracts.AuthGuard, error) {
	// 	return nil, nil
	// })
}
`

const ModuleEvent = `package providers

import (
	"bitbucket.org/dptsi/its-go/contracts"
)

type Listener struct {
	eventName            string
	listenersConstructor []func(application contracts.Application) (contracts.EventListener, error)
}

var listen []Listener = []Listener{}

func RegisterEvents(mod contracts.Module) {
	service := mod.App().Services().Event
	for _, l := range listen {
		service.Register(l.eventName, l.listenersConstructor)
	}
}
`

const ModuleMiddleware = `package providers

import "bitbucket.org/dptsi/its-go/contracts"

type Middleware struct {
	name        string
	constructor contracts.MiddlewareConstructor
}

var middlewares []Middleware = []Middleware{}

func RegisterMiddlewares(mod contracts.Module) {
	service := mod.App().Services().Middleware
	for _, m := range middlewares {
		service.Register(m.name, m.constructor)
	}
}
`
