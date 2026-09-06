package utils

import (
	"concurrency-simulator/services/shared"

	"go.uber.org/zap"
)

func NewRequestLogger() *zap.Logger {
	var loggerInstance = shared.NewLogger("account-svc")

	logger := loggerInstance.With(
		zap.String("topic", shared.AccountTopic),
	)

	return logger
}
