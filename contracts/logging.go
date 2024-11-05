package contracts

import "context"

type LoggingService interface {
	Debug(ctx context.Context, message string) error
	Info(ctx context.Context, message string) error
	Notice(ctx context.Context, message string) error
	Warning(ctx context.Context, message string) error
	Error(ctx context.Context, message string) error
	Critical(ctx context.Context, message string) error
	Alert(ctx context.Context, message string) error
	Emergency(ctx context.Context, message string) error
}
