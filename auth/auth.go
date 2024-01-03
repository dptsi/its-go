package auth

import (
	"fmt"

	"bitbucket.org/dptsi/go-framework/app"
	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/models"
	"bitbucket.org/dptsi/go-framework/web"
)

const defaultGuard = "sessions"
const guardContextKey = "auth.guard"
const errorPrefix = "auth service"

type Service struct {
	application contracts.Application
	guards      map[string]bool
}

func NewService(application contracts.Application) *Service {
	return &Service{
		application: application,
		guards:      make(map[string]bool),
	}
}

func (s *Service) Login(ctx *web.Context, user *models.User) error {
	key, guard, err := s.getGuard(ctx)
	if err != nil {
		return fmt.Errorf("%s: login: %w", errorPrefix, err)
	}

	statefulGuard, ok := guard.(contracts.StatefulAuthGuard)
	if !ok {
		return fmt.Errorf("%s: guard \"%s\" doesn't support login because it is not a stateful guard", errorPrefix, key)
	}

	return statefulGuard.Login(ctx, user)
}

func (s *Service) Logout(ctx *web.Context) error {
	errorPrefix := fmt.Sprintf("%s: logout", errorPrefix)
	key, guard, err := s.getGuard(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", errorPrefix, err)
	}

	statefulGuard, ok := guard.(contracts.StatefulAuthGuard)
	if !ok {
		return fmt.Errorf("%s: guard \"%s\" doesn't support logout because it is not a stateful guard", errorPrefix, key)
	}

	return statefulGuard.Logout(ctx)
}

func (s *Service) User(ctx *web.Context) (*models.User, error) {
	errorPrefix := fmt.Sprintf("%s: user", errorPrefix)
	_, guard, err := s.getGuard(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorPrefix, err)
	}

	return guard.User(ctx), nil
}

func (s *Service) getGuard(ctx *web.Context) (key string, guard contracts.AuthGuard, err error) {
	key = ctx.GetString(guardContextKey)
	if key == "" {
		key = defaultGuard
	}
	guard, err = app.Make[contracts.AuthGuard](s.application, s.getGuardKey(key))
	if err != nil {
		return key, nil, fmt.Errorf("auth guard \"%s\" not found", key)
	}

	return key, guard, nil
}

func (s *Service) RegisterGuard(name string, constructor contracts.AuthGuardConstructor) error {
	if _, exists := s.guards[name]; exists {
		return fmt.Errorf("auth service: register guard: guard \"%s\" already exist", name)
	}
	s.guards[name] = true

	app.Bind[contracts.AuthGuard](
		s.application,
		s.getGuardKey(name),
		constructor,
	)
	return nil
}

func (s *Service) getGuardKey(name string) string {
	return fmt.Sprintf("auth.guard.%s", name)
}
