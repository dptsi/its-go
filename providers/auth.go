package providers

import (
	"bitbucket.org/dptsi/its-go/app"
	"bitbucket.org/dptsi/its-go/auth"
	"bitbucket.org/dptsi/its-go/contracts"
)

func registerAuthGuard(application contracts.Application) error {
	service := application.Services().Auth
	service.RegisterGuard("sessions", func(application contracts.Application) (contracts.AuthGuard, error) {
		return auth.NewSessionGuard(app.MustMake[contracts.SessionService](application, "sessions.service")), nil
	})

	return nil
}
