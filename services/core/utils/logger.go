package utils

import (
	"concurrency-simulator/services/shared"
	"net/http"

	"go.uber.org/zap"
)

func NewRequestLogger(r *http.Request) *zap.Logger {
	var loggerInstance = shared.NewLogger("core-svc")

	logger := loggerInstance.With(
		zap.String("endpoint", r.URL.Path),
		zap.String("method", r.Method),
		zap.String("user_agent", r.UserAgent()),
	)

	return logger
}
