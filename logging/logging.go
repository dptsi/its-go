package logging

import (
	"fmt"

	"github.com/dptsi/its-go/app"
	"github.com/dptsi/its-go/contracts"
)

// ChannelConfig represents a logging channel configuration
type ChannelConfig struct {
	Driver       string
	DriverConfig map[string]interface{}
}

type Config struct {
	// Default log channel
	Default string

	// Log channels
	Channels map[string]ChannelConfig
}

type Service struct {
	application contracts.Application
	config      Config
}

func NewService(application contracts.Application, config Config) *Service {
	return &Service{
		application: application,
		config:      config,
	}
}

func (s *Service) RegisterDriver(name string, construct contracts.LoggingDriverConstructor) error {
	channelConfig, ok := s.config.Channels[name]
	if !ok {
		return fmt.Errorf("log channel \"%s\" not found", name)
	}

	app.Bind[contracts.LoggingDriver](
		s.application,
		fmt.Sprintf("logging.drivers.%s", name),
		func(application contracts.Application) (contracts.LoggingDriver, error) {
			return construct(channelConfig.DriverConfig), nil
		},
	)
	return nil
}

func (s *Service) getDefaultDriver() (contracts.LoggingDriver, error) {
	channel, ok := s.config.Channels[s.config.Default]
	if !ok {
		return nil, fmt.Errorf("default log channel \"%s\" not found", s.config.Default)
	}

	return app.Make[contracts.LoggingDriver](s.application, fmt.Sprintf("logging.drivers.%s", channel.Driver))
}

func (s *Service) Debug(message string, context map[string]interface{}) error {
	driver, err := s.getDefaultDriver()
	if err != nil {
		return err
	}
	driver.Debug(message, context)
	return nil
}

func (s *Service) Info(message string, context map[string]interface{}) error {
	driver, err := s.getDefaultDriver()
	if err != nil {
		return err
	}
	driver.Info(message, context)
	return nil

}

func (s *Service) Notice(message string, context map[string]interface{}) error {
	driver, err := s.getDefaultDriver()
	if err != nil {
		return err
	}
	driver.Notice(message, context)
	return nil
}

func (s *Service) Warning(message string, context map[string]interface{}) error {
	driver, err := s.getDefaultDriver()
	if err != nil {
		return err
	}
	driver.Warning(message, context)
	return nil
}

func (s *Service) Error(message string, context map[string]interface{}) error {
	driver, err := s.getDefaultDriver()
	if err != nil {
		return err
	}
	driver.Error(message, context)
	return nil
}

func (s *Service) Critical(message string, context map[string]interface{}) error {
	driver, err := s.getDefaultDriver()
	if err != nil {
		return err
	}
	driver.Critical(message, context)
	return nil
}

func (s *Service) Alert(message string, context map[string]interface{}) error {
	driver, err := s.getDefaultDriver()
	if err != nil {
		return err
	}
	driver.Alert(message, context)
	return nil
}

func (s *Service) Emergency(message string, context map[string]interface{}) error {
	driver, err := s.getDefaultDriver()
	if err != nil {
		return err
	}
	driver.Emergency(message, context)
	return nil
}
