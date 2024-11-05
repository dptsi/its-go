package module

import (
	"fmt"

	"github.com/dptsi/its-go/contracts"
)

type Service struct {
	application contracts.Application
	modules     map[string]bool
	logger      contracts.LoggingService
}

func NewService(application contracts.Application, logger contracts.LoggingService) *Service {
	return &Service{application, make(map[string]bool), logger}
}

func (s *Service) Register(name string, entrypoint contracts.ModuleEntrypoint) {
	if _, exists := s.modules[name]; exists {
		s.logger.Warning(s.application.Context(), fmt.Sprintf("module with name %s already exists", name))
	}
	s.modules[name] = true

	module := &Module{
		name:        name,
		Application: s.application,
	}
	entrypoint(module)
}
