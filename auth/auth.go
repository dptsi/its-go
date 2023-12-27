package auth

import (
	"fmt"

	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/models"
	"bitbucket.org/dptsi/go-framework/web"
)

const defaultGuard = "web"
const guardContextKey = "auth.guard"
const errorPrefix = "auth service"

type GuardsConfig struct {
	Driver contracts.AuthGuard
}

type Config struct {
	Guards map[string]GuardsConfig
}

type Service struct {
	cfg Config
}

func NewService(cfg Config) *Service {
	return &Service{
		cfg: cfg,
	}
}

func (a *Service) Login(ctx *web.Context, user *models.User) error {
	key, guard, err := a.getGuard(ctx)
	if err != nil {
		return fmt.Errorf("%s: login: %w", errorPrefix, err)
	}

	statefulGuard, ok := guard.(contracts.StatefulAuthGuard)
	if !ok {
		return fmt.Errorf("%s: guard \"%s\" doesn't support login because it is not a stateful guard", errorPrefix, key)
	}

	return statefulGuard.Login(ctx, user)
}

func (a *Service) Logout(ctx *web.Context) error {
	errorPrefix := fmt.Sprintf("%s: logout", errorPrefix)
	key, guard, err := a.getGuard(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", errorPrefix, err)
	}

	statefulGuard, ok := guard.(contracts.StatefulAuthGuard)
	if !ok {
		return fmt.Errorf("%s: guard \"%s\" doesn't support logout because it is not a stateful guard", errorPrefix, key)
	}

	return statefulGuard.Logout(ctx)
}

func (a *Service) User(ctx *web.Context) (*models.User, error) {
	errorPrefix := fmt.Sprintf("%s: user", errorPrefix)
	_, guard, err := a.getGuard(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorPrefix, err)
	}

	return guard.User(ctx), nil
}

func (a *Service) getGuard(ctx *web.Context) (key string, guard contracts.AuthGuard, err error) {
	key = ctx.GetString(guardContextKey)
	if key == "" {
		key = defaultGuard
	}
	cfg, exists := a.cfg.Guards[key]
	if !exists {
		return key, nil, fmt.Errorf("auth guard \"%s\" not found", key)
	}

	return key, cfg.Driver, nil
}
