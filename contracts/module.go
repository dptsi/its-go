package contracts

type Module interface {
	Name() string
	App() Application
}

type ModuleEntrypoint = func(mod Module)

type ModuleService interface {
	Register(name string, entrypoint ModuleEntrypoint)
}
