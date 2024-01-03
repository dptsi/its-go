package providers

import (
	"bitbucket.org/dptsi/go-framework/app"
	"bitbucket.org/dptsi/go-framework/auth"
	"bitbucket.org/dptsi/go-framework/contracts"
)

func registerAuthGuard(application contracts.Application) error {
	service := application.Services().Auth
	service.RegisterGuard("sessions", func(application contracts.Application) (contracts.AuthGuard, error) {
		return auth.NewSessionGuard(app.MustMake[contracts.SessionService](application, "sessions.service")), nil
	})

	return nil
}
