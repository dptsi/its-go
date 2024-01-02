package contracts

import "bitbucket.org/dptsi/go-framework/app"

type Module interface {
	Name() string
	App() *app.Application
}

type ModuleEntrypoint = func(mod Module)

type ModuleService interface {
	Register(name string, entrypoint ModuleEntrypoint)
}
