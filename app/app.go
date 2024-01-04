package app

import (
	"context"
	"fmt"

	"bitbucket.org/dptsi/its-go/contracts"
	"bitbucket.org/dptsi/its-go/web"
	"github.com/samber/do"
)

type Application struct {
	ctx context.Context
	i   *do.Injector
	cfg map[string]interface{}
}

func NewApplication(ctx context.Context, i *do.Injector, cfg map[string]interface{}) *Application {
	return &Application{
		ctx: ctx,
		i:   i,
		cfg: cfg,
	}
}

type Provider[T any] func(application contracts.Application) (T, error)

func Bind[T any](app contracts.Application, name string, provider Provider[T]) {
	do.ProvideNamed[T](app.Injector(), name, func(i *do.Injector) (T, error) {
		return provider(app)
	})
}

func MustMake[T any](app contracts.Application, name string) T {
	instance, err := do.InvokeNamed[T](app.Injector(), name)
	if err != nil {
		panic(fmt.Errorf("error when creating object %s: %w", name, err))
	}
	return instance
}

func Make[T any](app contracts.Application, name string) (T, error) {
	return do.InvokeNamed[T](app.Injector(), name)
}

func (app *Application) Context() context.Context {
	return app.ctx
}

func (app *Application) Config() map[string]interface{} {
	return app.cfg
}

func (app *Application) ListProvidedServices() []string {
	return app.i.ListProvidedServices()
}

func (app *Application) Injector() *do.Injector {
	return app.i
}

func (app *Application) Services() contracts.ApplicationServices {
	return contracts.ApplicationServices{
		Auth:       MustMake[contracts.AuthService](app, "auth.service"),
		Crypt:      MustMake[contracts.CryptService](app, "crypt.service"),
		Database:   MustMake[contracts.DatabaseService](app, "database.service"),
		Event:      MustMake[contracts.EventService](app, "event.service"),
		Middleware: MustMake[contracts.MiddlewareService](app, "http.middleware.service"),
		Module:     MustMake[contracts.ModuleService](app, "module.service"),
		Session:    MustMake[contracts.SessionService](app, "sessions.service"),
		WebEngine:  MustMake[*web.Engine](app, "web.engine"),
	}
}
