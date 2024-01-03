package module

import (
	"fmt"

	"bitbucket.org/dptsi/go-framework/app"
	"bitbucket.org/dptsi/go-framework/contracts"
)

type DependencyScope = string

const (
	DependencyScopeModule = "module"
	DependencyScopeGlobal = "global"
)

type Module struct {
	name string
	contracts.Application
}

type Provider[T any] func(mod contracts.Module) (T, error)

func Bind[T any](mod contracts.Module, name string, provider Provider[T]) {
	key := getServiceKey(mod.Name(), name)
	app.Bind[T](mod.App(), key, func(a contracts.Application) (T, error) {
		return provider(mod)
	})
}

func MustMake[T any](mod contracts.Module, name string, scope DependencyScope) T {
	key := name
	if scope == DependencyScopeModule {
		key = getServiceKey(mod.Name(), name)
	}
	return app.MustMake[T](mod.App(), key)
}

func Make[T any](mod contracts.Module, name string, scope DependencyScope) (T, error) {
	key := name
	if scope == DependencyScopeModule {
		key = getServiceKey(mod.Name(), name)
	}
	return app.Make[T](mod.App(), key)
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) App() contracts.Application {
	return m.Application
}

func getServiceKey(moduleName, key string) string {
	return fmt.Sprintf("modules.%s.dependencies.%s", moduleName, key)
}
