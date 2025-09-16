package contracts

import "time"

type SentryConfig struct {
	Enabled bool `env:"SENTRY_ENABLED,default=false"`
	Debug   bool `env:"SENTRY_DEBUG,default=false"`

	// sentry config
	Dsn              string  `env:"SENTRY_DSN,required"`
	EnableTracing    bool    `env:"SENTRY_ENABLE_TRACING,default=false"`
	TracesSampleRate float64 `env:"SENTRY_TRACES_SAMPLE_RATE,default=0"`
	SendDefaultPii   bool    `env:"SENTRY_SEND_DEFAULT_PII,default=true"`

	// sentrygin specific config
	GinRepanic         bool `env:"SENTRY_GIN_REPANIC,default=true"`
	GinWaitForDelivery bool `env:"SENTRY_GIN_WAIT_FOR_DELIVERY, default=false"`
	GinTimeout         int  `env:"SENTRY_GIN_TIMEOUT,default=5"`
}

type SentryService interface {
	IsEnabled() bool
	IsDebug() bool
	IsTracingEnabled() bool
	IsDefaultPiiSent() bool
	MustRepanicGin() bool
	MustWaitForDeliveryGin() bool
	GetDsn() string
	GetTracesSampleRate() float64
	GetTimeoutGin() time.Duration
}
