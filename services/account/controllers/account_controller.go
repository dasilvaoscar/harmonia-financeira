package controllers

import (
	"encoding/json"
	"fmt"

	"concurrency-simulator/services/account/contracts"
	"concurrency-simulator/services/account/internal/core/services"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type AccountController struct {
	logger  *zap.Logger
	service *services.AccountService
}

func (pc *AccountController) ProcessMessage(msg *kafka.Message) {
	pc.logger.Info("Received message from topic", zap.String("topic", *msg.TopicPartition.Topic), zap.String("message", string(msg.Value)))

	var accountMessage contracts.AccountContract
	if err := json.Unmarshal(msg.Value, &accountMessage); err != nil {
		return
	}

	result := pc.service.Execute(accountMessage)

	message := fmt.Sprintf("ACCOUNT_RESULT: %v", result)

	pc.logger.Info(message)
}
