package module

import (
	"log"

	"bitbucket.org/dptsi/go-framework/app"
	"bitbucket.org/dptsi/go-framework/contracts"
)

type Service struct {
	application *app.Application
	modules     map[string]bool
}

func NewService(application *app.Application) *Service {
	return &Service{application, make(map[string]bool)}
}

func (s *Service) Register(name string, entrypoint contracts.ModuleEntrypoint) {
	if _, exists := s.modules[name]; exists {
		log.Fatalf("module with name %s already exists", name)
	}
	s.modules[name] = true

	module := &Module{
		name:        name,
		Application: s.application,
	}
	entrypoint(module)
}
