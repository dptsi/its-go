package providers

import (
	"fmt"

	"github.com/dptsi/its-go/app"
	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/http"
	"github.com/dptsi/its-go/http/middleware"
)

func registerMiddlewares(application contracts.Application) error {
	config := application.Config()
	corsConfig, ok := config["cors"].(http.CorsConfig)
	if !ok {
		return fmt.Errorf("cors config is not available")
	}
	csrfConfig, ok := config["csrf"].(http.CSRFConfig)
	if !ok {
		return fmt.Errorf("csrf config is not available")
	}
	service := application.Services().Middleware

	authService := app.MustMake[contracts.AuthService](application, "auth.service")
	sessionService := app.MustMake[contracts.SessionService](application, "sessions.service")
	sentryService := app.MustMake[contracts.SentryService](application, "sentry.service")
	cacheService := app.MustMake[contracts.CacheService](application, "cache.service")

	service.Register("user_has_permission", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewUserHasPermission(authService), nil
	})
	service.Register("user_has_role", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewUserHasRole(authService), nil
	})
	service.Register("auth", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewAuth(authService), nil
	})
	service.Register("cors", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewCors(corsConfig), nil
	})
	service.Register("start_session", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewStartSession(sessionService), nil
	})
	service.Register("verify_csrf_token", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewVerifyCSRFToken(csrfConfig, sessionService)
	})

	/**
	 * register sentry middleware regardless of whether sentry is enabled or not.
	 * in some cases, sentry may be disabled on env, but still registered as
	 * global middleware by other codebase that use or inherit this repository.
	 * if env states that sentry is disabled, but it is still registered as global middleware,
	 * the code will throw an error due to missing provider for the sentry middleware.
	 */
	service.Register("sentry", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewSentryGin(sentryService, authService)
	})
	service.Register("cache", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewCache(cacheService), nil
	})

	return nil
}
