package logging

import (
	"context"
	"log"

	"github.com/dptsi/its-go/contracts"
)

type GoLogger struct {
	logger *log.Logger
}

func NewGoLogger(logger *log.Logger) contracts.LoggingService {
	return &GoLogger{
		logger: logger,
	}
}

func (l *GoLogger) Debug(_ context.Context, message string) error {
	l.logger.Printf("DEBUG: %s\n", message)
	return nil
}

func (l *GoLogger) Info(_ context.Context, message string) error {
	l.logger.Printf("INFO: %s\n", message)
	return nil
}

func (l *GoLogger) Notice(_ context.Context, message string) error {
	l.logger.Printf("NOTICE: %s\n", message)
	return nil
}

func (l *GoLogger) Warning(_ context.Context, message string) error {
	l.logger.Printf("WARNING: %s\n", message)
	return nil
}

func (l *GoLogger) Error(_ context.Context, message string) error {
	l.logger.Printf("ERROR: %s\n", message)
	return nil
}

func (l *GoLogger) Critical(_ context.Context, message string) error {
	l.logger.Printf("CRITICAL: %s\n", message)
	return nil
}

func (l *GoLogger) Alert(_ context.Context, message string) error {
	l.logger.Printf("ALERT: %s\n", message)
	return nil
}

func (l *GoLogger) Emergency(_ context.Context, message string) error {
	l.logger.Printf("EMERGENCY: %s\n", message)
	return nil
}
