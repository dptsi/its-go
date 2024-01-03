package providers

import (
	"fmt"

	"bitbucket.org/dptsi/its-go/app"
	"bitbucket.org/dptsi/its-go/contracts"
	"bitbucket.org/dptsi/its-go/http"
	"bitbucket.org/dptsi/its-go/http/middleware"
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

	service.Register("active_role_has_permission", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewActiveRoleHasPermission(app.MustMake[contracts.AuthService](application, "auth.service")), nil
	})
	service.Register("active_role_in", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewActiveRoleIn(app.MustMake[contracts.AuthService](application, "auth.service")), nil
	})
	service.Register("auth", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewAuth(app.MustMake[contracts.AuthService](application, "auth.service")), nil
	})
	service.Register("cors", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewCors(corsConfig), nil
	})
	service.Register("start_session", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewStartSession(app.MustMake[contracts.SessionService](application, "sessions.service")), nil
	})
	service.Register("verify_csrf_token", func(application contracts.Application) (contracts.Middleware, error) {
		return middleware.NewVerifyCSRFToken(
			csrfConfig,
			app.MustMake[contracts.SessionService](application, "sessions.service"),
		)
	})

	return nil
}
