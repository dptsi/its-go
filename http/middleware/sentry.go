package middleware

import (
	"fmt"
	"os"

	"github.com/dptsi/its-go/app/errors"
	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/web"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

type SentryGin struct {
	service     contracts.SentryService
	authService contracts.AuthService
}

func NewSentryGin(service contracts.SentryService, authService contracts.AuthService) (*SentryGin, error) {
	return &SentryGin{service, authService}, nil
}

/**
 * will return NOOP middleware if sentry is disabled,
 * otherwise returns sentrygin middleware.
 */
func (s *SentryGin) Handle(interface{}) web.HandlerFunc {
	service := s.service
	auth := s.authService

	if !service.IsEnabled() {
		fmt.Printf("sentry SDK is disabled\n")
		return func(ctx *web.Context) {} // NOOP middleware
	}

	if err := sentry.Init(sentry.ClientOptions{
		Environment:      os.Getenv("APP_ENV"),
		Debug:            service.IsDebug(),
		Dsn:              service.GetDsn(),
		EnableTracing:    service.IsTracingEnabled(),
		TracesSampleRate: service.GetTracesSampleRate(),
	}); err != nil {
		panic(errors.Errorf("sentry SDK initialization failed: %w", err))
	}

	middleware := sentrygin.New(sentrygin.Options{
		Repanic:         service.MustRepanicGin(),
		WaitForDelivery: service.MustWaitForDeliveryGin(),
		Timeout:         service.GetTimeoutGin(),
	})

	extraScopeHandler := func(ctx *gin.Context) {
		user, _ := auth.User(ctx)

		sentry.ConfigureScope(func(scope *sentry.Scope) {
			if user != nil {
				scope.SetUser(sentry.User{
					ID:    user.Id(),
					Name:  user.Name(),
					Email: user.Email(),
				})
			}
		})
		middleware(ctx)
	}

	return extraScopeHandler
}
