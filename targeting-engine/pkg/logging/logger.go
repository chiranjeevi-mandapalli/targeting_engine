package logging

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New(serviceName string) *Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}

	logger := zerolog.New(output).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()

	return &Logger{logger}
}

func (l *Logger) Debug() *zerolog.Event {
	return l.Logger.Debug()
}

func (l *Logger) Info() *zerolog.Event {
	return l.Logger.Info()
}

func (l *Logger) Warn() *zerolog.Event {
	return l.Logger.Warn()
}

func (l *Logger) Error() *zerolog.Event {
	return l.Logger.Error()
}

func (l *Logger) Fatal() *zerolog.Event {
	return l.Logger.Fatal()
}

func (l *Logger) LogError(err error, msg string) {
	l.Error().Err(err).Msg(msg)
}

func (l *Logger) LogFatal(err error, msg string) {
	l.Fatal().Err(err).Msg(msg)
	os.Exit(1)
}

func (l *Logger) With() zerolog.Context {
	return l.Logger.With()
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{l.Logger.With().Interface(key, value).Logger()}
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	ctx := l.Logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &Logger{ctx.Logger()}
}
