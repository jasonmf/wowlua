package wowlua

import "log"

const (
	// LogLevelDebug will log debug messages
	LogLevelDebug = iota
	// LogLevelErrors will output errors-only
	LogLevelErrors
)

var (
	logger Logger = DefaultLogger(LogLevelErrors)
)

// A Logger has a debug and error logging function
type Logger interface {
	Debugf(tmpl string, v ...interface{})
	Errorf(tmpl string, v ...interface{})
}

// DefaultLogger is a wrapper around the log package that conditionally logs
// based on its log level
type DefaultLogger int

// Debugf outputs a debug message if logging is set to LogLevelDebug
func (level DefaultLogger) Debugf(tmpl string, v ...interface{}) {
	if level >= LogLevelDebug {
		return
	}
	log.Printf("DEBUG: "+tmpl, v...)
}

// Errorf outputs and error message
func (level DefaultLogger) Errorf(tmpl string, v ...interface{}) {
	if level >= LogLevelErrors {
		return
	}
	log.Printf("ERROR: "+tmpl, v...)
}
