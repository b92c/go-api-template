package zaplogger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	l *zap.SugaredLogger
}

func New(env string) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	cfg.Encoding = "json"
	if env == "dev" || env == "local" {
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	log, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{l: log.Sugar()}, nil
}

func (z *Logger) Sync() { _ = z.l.Sync() }

func (z *Logger) Debug(msg string, args ...any) { z.l.Debugw(msg, toFields(args...)...) }
func (z *Logger) Info(msg string, args ...any)  { z.l.Infow(msg, toFields(args...)...) }
func (z *Logger) Warn(msg string, args ...any)  { z.l.Warnw(msg, toFields(args...)...) }
func (z *Logger) Error(msg string, args ...any) { z.l.Errorw(msg, toFields(args...)...) }

func toFields(args ...any) []any {
	// aceita pares chave,valor; se inválido, anexa como extra
	if len(args)%2 != 0 {
		return []any{"extra", fmt.Sprint(args...)}
	}
	return args
}

// Helper para integração simples no main
func FromEnv() (*Logger, error) {
	return New(os.Getenv("APP_ENV"))
}
