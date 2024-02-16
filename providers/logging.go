package providers

import (
	"log"

	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/logging/driver"
)

func registerLogger(application contracts.Application) error {
	service := application.Services().Logging
	service.RegisterDriver("go", func(config map[string]interface{}) contracts.LoggingDriver {
		return driver.NewGoLogger(log.Default())
	})

	return nil
}
