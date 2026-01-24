package logger_test

import (
	"context"
	"github.com/TommyLearning/go-rest-api-project/internal/logger"
	"log/slog"
	"os"
	"testing"
)

func Test_CtxWithLogger(t *testing.T) {
	testCases := []struct {
		name   string
		ctx    context.Context
		logger *slog.Logger
		exists bool
	}{
		{
			name: "return context without logger",
			ctx:  context.Background(),
		},
		{
			name:   "return ctx as it is",
			ctx:    context.WithValue(context.Background(), logger.CtxKey{}, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))),
			exists: true,
		},
		{
			name:   "inject logger",
			ctx:    context.Background(),
			logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})),
			exists: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := logger.CtxWithLogger(tc.ctx, tc.logger)
			_, ok := ctx.Value(logger.CtxKey{}).(*slog.Logger)

			if tc.exists != ok {
				t.Errorf("expected :%v got %v", tc.exists, ok)
			}

		})
	}
}

func Test_FromContext(t *testing.T) {
	testCases := []struct {
		name     string
		ctx      context.Context
		expected bool
	}{
		{
			name: "logger exists",
			ctx:  context.WithValue(context.Background(), logger.CtxKey{}, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))),
		},
		{
			name: "new logger returned",
			ctx:  context.Background(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := logger.FromContext(tc.ctx)

			if tc.expected && logger == nil {
				t.Errorf("expected %v ,got %v", tc.expected, logger)

			}
		})
	}
}
