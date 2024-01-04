package middleware

import (
	"fmt"

	"github.com/dptsi/its-go/app"
	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/web"
)

type Config struct {
	Groups map[string]([]string)
}

type Service struct {
	app         contracts.Application
	cfg         Config
	middlewares map[string]bool
}

func NewService(app contracts.Application, cfg Config) *Service {
	return &Service{
		app:         app,
		cfg:         cfg,
		middlewares: make(map[string]bool),
	}
}

func (s *Service) Use(name string, params interface{}) web.HandlerFunc {
	m, err := app.Make[contracts.Middleware](s.app, s.getServiceKey(name))
	if err != nil {
		panic(fmt.Errorf("middleware %s not found", name))
	}
	return m.Handle(params)
}

func (s *Service) Global() []web.HandlerFunc {
	global := s.cfg.Groups["global"]

	return s.group(global)
}

func (s *Service) Register(name string, constructor contracts.MiddlewareConstructor) error {
	if _, exists := s.middlewares[name]; exists {
		return fmt.Errorf("middleware service: register: middleware \"%s\" already exist", name)
	}
	s.middlewares[name] = true
	app.Bind[contracts.Middleware](s.app, s.getServiceKey(name), constructor)
	return nil
}

func (s *Service) group(group []string) []web.HandlerFunc {
	middlewares := make([]web.HandlerFunc, len(group))
	for i, name := range group {
		middlewares[i] = s.Use(name, nil)
	}
	return middlewares
}

func (s *Service) getServiceKey(name string) string {
	return fmt.Sprintf("http.middleware.handler.%s", name)
}
