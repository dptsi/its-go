package module

import (
	"log"

	"github.com/dptsi/its-go/contracts"
)

type Service struct {
	application contracts.Application
	modules     map[string]bool
}

func NewService(application contracts.Application) *Service {
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
