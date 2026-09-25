package app

import (
	"context"
	"fmt"

	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/web"
	"github.com/samber/do"
	typetostring "github.com/samber/go-type-to-string"
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

// inferServiceName uses type inference to determine the service name
// based on the generic type parameter T. This is used internally
// to automatically generate service names from types.
//
// original: https://github.com/samber/do/blob/master/service.go
func inferServiceName[T any]() string {
	return typetostring.GetType[T]()
}

type Provider[T any] func(application contracts.Application) (T, error)

func bind[T any](app contracts.Application, name string, provider Provider[T]) {
	do.ProvideNamed[T](app.Injector(), name, func(i *do.Injector) (T, error) {
		return provider(app)
	})
}

func autoBind[T any](app contracts.Application, provider Provider[T]) {
	name := inferServiceName[T]()
	do.ProvideNamed[T](app.Injector(), name, func(i *do.Injector) (T, error) {
		return provider(app)
	})
}

// binding service in application
//
// example:
//
//	// bind using explicit name
//	app.Bind(application, "auth.service", auth.NewDefaultAuthService)
//	app.Bind[contracts.AuthService](application, "auth.service", auth.NewDefaultAuthService)
//
//	// or bind using implicit name
//	app.Bind[contracts.AuthService](application, auth.NewDefaultAuthService)
func Bind[T any](app contracts.Application, name any, provider ...Provider[T]) {
	switch _name := name.(type) {
	case string:
		var _provider Provider[T] = provider[0]
		bind(app, _name, _provider)
	default:
		var _provider Provider[T] = name.(func(application contracts.Application) (T, error))
		autoBind(app, _provider)
	}
}

func mustMake[T any](app contracts.Application, name string) T {
	instance, err := do.InvokeNamed[T](app.Injector(), name)
	if err != nil {
		panic(fmt.Errorf("error when creating object %s: %w", name, err))
	}
	return instance
}

func autoMustMake[T any](app contracts.Application) T {
	name := inferServiceName[T]()
	instance, err := do.InvokeNamed[T](app.Injector(), name)
	if err != nil {
		panic(fmt.Errorf("error when creating object %s: %w", name, err))
	}
	return instance
}

// get service instance and panic if failed
//
// example:
//
//	// get using explicit name
//	app.MustMake[contracts.AuthService](application, "auth.service")
//	// or get using implicit name
//	app.MustMake[contracts.AuthService](application)
func MustMake[T any](app contracts.Application, args ...any) T {
	var name string

	if len(args) > 0 {
		name = args[0].(string)
		return mustMake[T](app, name)
	} else {
		return autoMustMake[T](app)
	}
}

func make[T any](app contracts.Application, name string) (T, error) {
	return do.InvokeNamed[T](app.Injector(), name)
}

func autoMake[T any](app contracts.Application) (T, error) {
	name := inferServiceName[T]()
	return do.InvokeNamed[T](app.Injector(), name)
}

// get service instance
//
// example:
//
//	// get using explicit name
//	app.Make[contracts.AuthService](application, "auth.service")
//	// or get using implicit name
//	app.Make[contracts.AuthService](application)
func Make[T any](app contracts.Application, args ...any) (T, error) {
	var name string

	if len(args) > 0 {
		name = args[0].(string)
		return make[T](app, name)
	} else {
		return autoMake[T](app)
	}
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

func (app *Application) Shutdown() error {
	return app.i.Shutdown()
}

func (app *Application) Services() contracts.ApplicationServices {
	return contracts.ApplicationServices{
		Auth:        MustMake[contracts.AuthService](app, "auth.service"),
		ActivityLog: MustMake[contracts.ActivityLogService](app, "activity_log.service"),
		Cache:       MustMake[contracts.CacheService](app, "cache.service"),
		Crypt:       MustMake[contracts.CryptService](app, "crypt.service"),
		Database:    MustMake[contracts.DatabaseService](app, "database.service"),
		Event:       MustMake[contracts.EventService](app, "event.service"),
		Logging:     MustMake[contracts.LoggingService](app, "logging.service"),
		Middleware:  MustMake[contracts.MiddlewareService](app, "http.middleware.service"),
		Module:      MustMake[contracts.ModuleService](app, "module.service"),
		Redis:       MustMake[contracts.RedisService](app, "redis.service"),
		Session:     MustMake[contracts.SessionService](app, "sessions.service"),
		WebEngine:   MustMake[*web.Engine](app, "web.engine"),
	}
}
