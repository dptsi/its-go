package sentry

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dptsi/its-go/contracts"
)

type Service struct {
	application contracts.Application
	config      contracts.SentryConfig
}

func NewService(application contracts.Application) (*Service, error) {
	config, err := loadConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return &Service{
		application: application,
		config:      *config,
	}, nil
}

func (s *Service) IsEnabled() bool {
	return s.config.Enabled
}

func (s *Service) IsDebug() bool {
	return s.config.Debug
}

func (s *Service) IsTracingEnabled() bool {
	return s.config.EnableTracing
}

func (s *Service) IsDefaultPiiSent() bool {
	return s.config.SendDefaultPii
}

func (s *Service) MustRepanicGin() bool {
	return s.config.GinRepanic
}

func (s *Service) GetDsn() string {
	return s.config.Dsn
}

func (s *Service) GetTracesSampleRate() float64 {
	return s.config.TracesSampleRate
}

func (s *Service) GetTimeoutGin() time.Duration {
	return time.Duration(s.config.GinTimeout) * time.Second
}

func (s *Service) MustWaitForDeliveryGin() bool {
	return s.config.GinWaitForDelivery
}

func loadConfigFromEnv() (*contracts.SentryConfig, error) {
	var cfg contracts.SentryConfig

	enabled, err := strconv.ParseBool(os.Getenv("SENTRY_ENABLED"))
	if err != nil {
		enabled = false
	}
	// set enabled status before any error happens
	cfg.Enabled = enabled
	debug, err := strconv.ParseBool(os.Getenv("SENTRY_DEBUG"))
	if err != nil {
		debug = false
	}
	// set debug status before any error happens
	cfg.Debug = debug

	dsn := os.Getenv("SENTRY_DSN")
	if enabled && dsn == "" {
		return nil, fmt.Errorf("SENTRY_DSN cannot be empty")
	}

	tracingEnabled, err := strconv.ParseBool(os.Getenv("SENTRY_ENABLE_TRACING"))
	if err != nil {
		tracingEnabled = false
	}

	tracesSampleRate, err := strconv.ParseFloat(os.Getenv("SENTRY_TRACES_SAMPLE_RATE"), 64)
	if err != nil {
		tracesSampleRate = 0
	}

	sendDefaultPii, err := strconv.ParseBool(os.Getenv("SENTRY_SEND_DEFAULT_PII"))
	if err != nil {
		sendDefaultPii = true
	}

	sentryginRepanic, err := strconv.ParseBool(os.Getenv("SENTRY_GIN_REPANIC"))
	if err != nil {
		sentryginRepanic = true
	}

	sentryginWaitForDelivery, err := strconv.ParseBool(os.Getenv("SENTRY_GIN_WAIT_FOR_DELIVERY"))
	if err != nil {
		sentryginWaitForDelivery = false
	}

	sentryginTimeout, err := strconv.Atoi(os.Getenv("SENTRY_GIN_TIMEOUT"))
	if err != nil {
		sentryginTimeout = 5
	}

	cfg.EnableTracing = tracingEnabled
	cfg.TracesSampleRate = tracesSampleRate
	cfg.SendDefaultPii = sendDefaultPii

	cfg.GinRepanic = sentryginRepanic
	cfg.GinWaitForDelivery = sentryginWaitForDelivery
	cfg.GinTimeout = sentryginTimeout

	return &cfg, nil
}
