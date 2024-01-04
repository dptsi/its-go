package templates

const ModuleEntrypointTemplate = `package {{.Name}}

import (
	"bitbucket.org/dptsi/its-go/contracts"
	"{{.BaseMod}}/modules/{{.Name}}/internal/app/providers"
	"{{.BaseMod}}/modules/{{.Name}}/internal/presentation/routes"
)

func SetupModule(mod contracts.Module) {
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
	g.POST("/hello", middlewareService.Use("auth", nil), func(ctx *web.Context) {
		ctx.JSON(http.StatusOK, web.H{
			"code":    0,
			"message": "hello from module {{.Name}}",
			"data":    nil,
		})
	})
}
`
