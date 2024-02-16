package driver

import (
	"log"
	"strings"
	"text/template"
)

type GoLogger struct {
	logger *log.Logger
}

func NewGoLogger(logger *log.Logger) *GoLogger {
	return &GoLogger{
		logger: logger,
	}
}

func (l *GoLogger) formatMessage(message string, context map[string]interface{}) string {
	tmpl, err := template.New("message").Parse(message)
	if err != nil {
		return message
	}

	var result strings.Builder
	err = tmpl.Execute(&result, context)
	if err != nil {
		return message
	}

	return result.String()
}

func (l *GoLogger) Debug(message string, context map[string]interface{}) {
	l.logger.Printf("DEBUG: %s", l.formatMessage(message, context))
}

func (l *GoLogger) Info(message string, context map[string]interface{}) {
	l.logger.Printf("INFO: %s", l.formatMessage(message, context))
}

func (l *GoLogger) Notice(message string, context map[string]interface{}) {
	l.logger.Printf("NOTICE: %s", l.formatMessage(message, context))
}

func (l *GoLogger) Warning(message string, context map[string]interface{}) {
	l.logger.Printf("WARNING: %s", l.formatMessage(message, context))
}

func (l *GoLogger) Error(message string, context map[string]interface{}) {
	l.logger.Printf("ERROR: %s", l.formatMessage(message, context))
}

func (l *GoLogger) Critical(message string, context map[string]interface{}) {
	l.logger.Printf("CRITICAL: %s", l.formatMessage(message, context))
}

func (l *GoLogger) Alert(message string, context map[string]interface{}) {
	l.logger.Printf("ALERT: %s", l.formatMessage(message, context))
}

func (l *GoLogger) Emergency(message string, context map[string]interface{}) {
	l.logger.Printf("EMERGENCY: %s", l.formatMessage(message, context))
}
