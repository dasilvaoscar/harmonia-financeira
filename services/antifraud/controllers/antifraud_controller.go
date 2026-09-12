package controllers

import (
	"encoding/json"
	"fmt"

	"concurrency-simulator/services/antifraud/internal/core/services"
	"concurrency-simulator/services/shared/topic_messages"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type AntifraudController struct {
	logger  *zap.Logger
	service *services.AntifraudAnalisysService
}

func (pc *AntifraudController) ProcessMessage(msg *kafka.Message) {
	var paymentMsg topic_messages.Payment
	if err := json.Unmarshal(msg.Value, &paymentMsg); err != nil {
		pc.logger.Error("Error unmarshalling message", zap.Error(err))
		return
	}

	payment := pc.service.Execute(paymentMsg)

	message := fmt.Sprintf("ANTIFRAUD_RESULT: Email=%s, Amount=%.2f, Installments=%d, Status=%s",
		payment.Email, payment.Amount, payment.Installments, *payment.Status)

	pc.logger.Info(message)
}
