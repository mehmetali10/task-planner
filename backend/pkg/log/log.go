package log

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	SetLogLevel(level LogLevel)
	GetLogLevel() LogLevel
	Trace(msg string, args ...any)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Panic(msg string, args ...any)
	Fatal(msg string, args ...any)
}

type LogLevel string

const (
	TraceLevel LogLevel = "Trace"
	DebugLevel LogLevel = "Debug"
	InfoLevel  LogLevel = "Info"
	WarnLevel  LogLevel = "Warn"
	ErrorLevel LogLevel = "Error"
	PanicLevel LogLevel = "Panic"
	FatalLevel LogLevel = "Fatal"
)

type customLogrus struct {
	lvl       LogLevel
	component string
	logger    *logrus.Entry
}

// NewLogger creates a new Logrus logger with a component name from config.
func NewLogger(component string, level string) Logger {
	llevel := getLogLevel(level)
	lg := &customLogrus{
		lvl:       llevel,
		component: component,
	}

	baseLogger := logrus.New()
	baseLogger.SetLevel(logrus.InfoLevel)

	fmt.Printf("component: %v || level: %v || logLevel: %v\n", component, level, llevel)
	baseLogger.SetFormatter(&logrus.JSONFormatter{})
	lg.logger = baseLogger.WithField("component", component)
	lg.SetLogLevel(llevel)

	return lg
}

// SetLogLevel sets the log level dynamically.
func (l *customLogrus) SetLogLevel(level LogLevel) {
	var logrusLevel logrus.Level
	switch level {
	case TraceLevel:
		logrusLevel = logrus.TraceLevel
	case DebugLevel:
		logrusLevel = logrus.DebugLevel
	case InfoLevel:
		logrusLevel = logrus.InfoLevel
	case WarnLevel:
		logrusLevel = logrus.WarnLevel
	case ErrorLevel:
		logrusLevel = logrus.ErrorLevel
	case FatalLevel:
		logrusLevel = logrus.FatalLevel
	case PanicLevel:
		logrusLevel = logrus.PanicLevel
	default:
		logrusLevel = logrus.InfoLevel
	}

	logrus.SetLevel(logrusLevel)
	l.logger.Logger.SetLevel(logrusLevel)
	l.lvl = level
}

// GetLogLevel returns the current log level.
func (l *customLogrus) GetLogLevel() LogLevel {
	return l.lvl
}

// Debug logs a debug message.
func (l *customLogrus) Debug(msg string, args ...any) {
	l.logger.Debugf(msg, args...)
}

// Error logs an error message.
func (l *customLogrus) Error(msg string, args ...any) {
	l.logger.Errorf(msg, args...)
}

// Fatal logs a fatal message and exits.
func (l *customLogrus) Fatal(msg string, args ...any) {
	l.logger.Fatalf(msg, args...)
}

// Info logs an info message.
func (l *customLogrus) Info(msg string, args ...any) {
	l.logger.Infof(msg, args...)
}

// Panic logs a panic message and panics.
func (l *customLogrus) Panic(msg string, args ...any) {
	l.logger.Panicf(msg, args...)
}

// Trace logs a trace message.
func (l *customLogrus) Trace(msg string, args ...any) {
	l.logger.Tracef(msg, args...)
}

// Warn logs a warning message.
func (l *customLogrus) Warn(msg string, args ...any) {
	l.logger.Warnf(msg, args...)
}

func getLogLevel(lvl string) LogLevel {
	switch strings.ToLower(lvl) {
	case "trace":
		return TraceLevel
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}
