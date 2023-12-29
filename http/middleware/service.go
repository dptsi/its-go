package middleware

import (
	"fmt"

	"bitbucket.org/dptsi/go-framework/app"
	"bitbucket.org/dptsi/go-framework/contracts"
)

type Config struct {
	Groups map[string]([]string)
}

type Service struct {
	app *app.Application
	cfg Config
}

func NewService(app *app.Application, cfg Config) *Service {
	return &Service{
		app: app,
		cfg: cfg,
	}
}

func (s *Service) Get(name string) (contracts.Middleware, error) {
	return app.Make[contracts.Middleware](s.app, fmt.Sprintf("http.middleware.handler.%s", name))
}

func (s *Service) Global() ([]contracts.Middleware, error) {
	global := s.cfg.Groups["global"]

	return s.group(global)
}

func (s *Service) group(group []string) ([]contracts.Middleware, error) {
	middlewares := make([]contracts.Middleware, len(group))
	for i, name := range group {
		m, err := app.Make[contracts.Middleware](s.app, fmt.Sprintf("http.middleware.handler.%s", name))
		if err != nil {
			return nil, err
		}
		middlewares[i] = m
	}

	return middlewares, nil
}
