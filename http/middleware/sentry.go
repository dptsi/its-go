package middleware

import (
	"fmt"

	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/web"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
)

type SentryGin struct {
	service contracts.SentryService
}

func NewSentryGin(service contracts.SentryService) (*SentryGin, error) {
	return &SentryGin{service}, nil
}

/**
 * will return NOOP middleware if sentry is disabled,
 * otherwise returns sentrygin middleware.
 */
func (s *SentryGin) Handle(interface{}) web.HandlerFunc {
	service := s.service

	if !service.IsEnabled() {
		fmt.Printf("sentry SDK is disabled\n")
		return func(ctx *web.Context) {} // NOOP middleware
	}

	if err := sentry.Init(sentry.ClientOptions{
		Debug:            service.IsDebug(),
		Dsn:              service.GetDsn(),
		EnableTracing:    service.IsTracingEnabled(),
		TracesSampleRate: service.GetTracesSampleRate(),
	}); err != nil {
		panic(fmt.Errorf("sentry SDK initialization failed: %w", err))
	}

	middleware := sentrygin.New(sentrygin.Options{
		Repanic:         service.MustRepanicGin(),
		WaitForDelivery: service.MustWaitForDeliveryGin(),
		Timeout:         service.GetTimeoutGin(),
	})

	return middleware
}
