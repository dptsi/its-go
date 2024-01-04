package module

import (
	"bitbucket.org/dptsi/its-go/contracts"
)

type Module struct {
	name string
	contracts.Application
}

type Provider[T any] func(mod contracts.Module) (T, error)

func (m *Module) Name() string {
	return m.name
}

func (m *Module) App() contracts.Application {
	return m.Application
}
