package app

import "github.com/samber/do"

type Application struct {
	i   *do.Injector
	cfg map[string]interface{}
}

func NewApplication(i *do.Injector, cfg map[string]interface{}) *Application {
	return &Application{
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

func Make[T any](app *Application, name string) T {
	return do.MustInvokeNamed[T](app.i, name)
}

func (app *Application) Config() map[string]interface{} {
	return app.cfg
}

func (app *Application) ListProvidedServices() []string {
	return app.i.ListProvidedServices()
}
