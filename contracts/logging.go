package contracts

type LoggingDriver interface {
	Debug(message string, context map[string]interface{})
	Info(message string, context map[string]interface{})
	Notice(message string, context map[string]interface{})
	Warning(message string, context map[string]interface{})
	Error(message string, context map[string]interface{})
	Critical(message string, context map[string]interface{})
	Alert(message string, context map[string]interface{})
	Emergency(message string, context map[string]interface{})
}

type LoggingService interface {
	Debug(message string, context map[string]interface{}) error
	Info(message string, context map[string]interface{}) error
	Notice(message string, context map[string]interface{}) error
	Warning(message string, context map[string]interface{}) error
	Error(message string, context map[string]interface{}) error
	Critical(message string, context map[string]interface{}) error
	Alert(message string, context map[string]interface{}) error
	Emergency(message string, context map[string]interface{}) error
	RegisterDriver(name string, construct LoggingDriverConstructor) error
}

type LoggingDriverConstructor func(config map[string]interface{}) LoggingDriver
