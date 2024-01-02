package middleware

import (
	"fmt"

	"bitbucket.org/dptsi/go-framework/app"
	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/web"
)

type Config struct {
	Groups map[string]([]string)
}

type Service struct {
	app contracts.Application
	cfg Config
}

func NewService(app contracts.Application, cfg Config) *Service {
	return &Service{
		app: app,
		cfg: cfg,
	}
}

func (s *Service) Use(name string, params interface{}) web.HandlerFunc {
	m, err := app.Make[contracts.Middleware](s.app, fmt.Sprintf("http.middleware.handler.%s", name))
	if err != nil {
		panic(fmt.Errorf("middleware %s not found", name))
	}
	return m.Handle(params)
}

func (s *Service) Global() web.HandlerFunc {
	global := s.cfg.Groups["global"]

	return s.group(global)
}

func (s *Service) group(group []string) web.HandlerFunc {
	return func(ctx *web.Context) {
		for _, name := range group {
			m := s.Use(name, nil)
			m(ctx)
		}
	}
}
