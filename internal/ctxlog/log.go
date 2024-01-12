package ctxlog

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const emptyCtxLoggerMsg = "failed to get logger from context"

var ctxLoggerKey = &struct{}{}

type CtxLogger struct {
	Logger zerolog.Context
	mux    sync.Mutex
}

func GetLogger(ctx context.Context) zerolog.Context {
	logger, err := getCtxLogger(ctx)
	if err != nil {
		log.Warn().Err(err).Msg(emptyCtxLoggerMsg)
		return log.With()
	}

	return logger.Logger
}

func getCtxLogger(ctx context.Context) (*CtxLogger, error) {
	logger, ok := ctx.Value(ctxLoggerKey).(*CtxLogger)
	if !ok || logger == nil {
		return nil, fmt.Errorf("invalid logger")
	}

	return logger, nil
}

func WithLogger(ctx context.Context, opts ...func(l *CtxLogger)) context.Context {
	logger := log.With()
	ctxLogger := &CtxLogger{Logger: logger}
	for _, opt := range opts {
		opt(ctxLogger)
	}

	return context.WithValue(ctx, ctxLoggerKey, ctxLogger)
}

func WithAny(ctx context.Context, k string, v any) {
	logger, err := getCtxLogger(ctx)
	if err != nil {
		log.Warn().Err(err).Msg(emptyCtxLoggerMsg)
		return
	}
	logger.mux.Lock()
	defer logger.mux.Unlock()

	l := logger.Logger.Interface(k, v)
	logger.Logger = l
}

func WithStr(ctx context.Context, k, v string) {
	logger, err := getCtxLogger(ctx)
	if err != nil {
		log.Warn().Err(err).Msg(emptyCtxLoggerMsg)
		return
	}
	logger.mux.Lock()
	defer logger.mux.Unlock()

	l := logger.Logger.Str(k, v)
	logger.Logger = l
}

func WithInt(ctx context.Context, k string, v int) {
	logger, err := getCtxLogger(ctx)
	if err != nil {
		log.Warn().Err(err).Msg(emptyCtxLoggerMsg)
		return
	}
	logger.mux.Lock()
	defer logger.mux.Unlock()

	l := logger.Logger.Int(k, v)
	logger.Logger = l
}

func WithInt32(ctx context.Context, k string, v int32) {
	logger, err := getCtxLogger(ctx)
	if err != nil {
		log.Warn().Err(err).Msg(emptyCtxLoggerMsg)
		return
	}
	logger.mux.Lock()
	defer logger.mux.Unlock()

	l := logger.Logger.Int32(k, v)
	logger.Logger = l
}

func WithInt64(ctx context.Context, k string, v int64) {
	logger, err := getCtxLogger(ctx)
	if err != nil {
		log.Warn().Err(err).Msg(emptyCtxLoggerMsg)
		return
	}
	logger.mux.Lock()
	defer logger.mux.Unlock()

	l := logger.Logger.Int64(k, v)
	logger.Logger = l
}

func WithTime(ctx context.Context, k string, v time.Time) {
	logger, err := getCtxLogger(ctx)
	if err != nil {
		log.Warn().Err(err).Msg(emptyCtxLoggerMsg)
		return
	}
	logger.mux.Lock()
	defer logger.mux.Unlock()

	l := logger.Logger.Time(k, v)
	logger.Logger = l
}
