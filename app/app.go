package app

import (
	"context"
	"fmt"

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

type Provider[T any] func(*Application) (T, error)

func Bind[T any](app *Application, name string, provider Provider[T]) {
	do.ProvideNamed[T](app.i, name, func(i *do.Injector) (T, error) {
		return provider(app)
	})
}

func MustMake[T any](app *Application, name string) T {
	instance, err := do.InvokeNamed[T](app.i, name)
	if err != nil {
		panic(fmt.Errorf("error when creating object %s: %w", name, err))
	}
	return instance
}

func Make[T any](app *Application, name string) (T, error) {
	return do.InvokeNamed[T](app.i, name)
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
