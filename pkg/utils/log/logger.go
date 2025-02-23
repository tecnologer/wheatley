package log

import "io"

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	SetOutput(w io.Writer)
	SetLevel(level Level)
	Level() Level
}

var defaultLogger Logger = NewHandler() //nolint:gochecknoglobals

func Default() Logger {
	return defaultLogger
}

func SetLogger(logger Logger) {
	defaultLogger = logger
}

func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

func Debug(msg string) {
	defaultLogger.Debug(msg)
}

func Debugf(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

func Warn(msg string) {
	defaultLogger.Warn(msg)
}

func Warnf(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

func Info(msg string) {
	defaultLogger.Info(msg)
}

func Infof(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func Error(msg string) {
	defaultLogger.Error(msg)
}

func Errorf(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}
