package application

import (
	"log/slog"
	"os"
)

const (
	SysVersion = "sys_version"
	SysBuildRv = "sys_buildrv"
	SysBuildAt = "sys_buildat"
	SysRuntime = "sys_runtime"
	InnerError = "inner_error"
	UserId     = "user_id"
	UserName   = "user_name"
	QueryId    = "query_type"
	ChatType   = "chat_type"
)

type Logger struct {
	log *slog.Logger
}

func NewLogger() *Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil)).With(
		SysVersion, Version,
		SysBuildRv, BuildRv,
		SysBuildAt, BuildAt,
		SysRuntime, GoVersion,
	)
	slog.SetDefault(logger)
	logger.Info("Starting dickbot (dickobrazz) ...")
	return &Logger{log: logger}
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{log: l.log.With(args...)}
}

func (l *Logger) D(msg string, args ...any) {
	l.log.Debug(msg, args...)
}

func (l *Logger) I(msg string, args ...any) {
	l.log.Info(msg, args...)
}

func (l *Logger) E(msg string, args ...any) {
	l.log.Error(msg, args...)
}

func (l *Logger) F(msg string, args ...any) {
	l.log.Error(msg, args...)
	panic(msg)
}
