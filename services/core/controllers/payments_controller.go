package controllers

import (
	"concurrency-simulator/services/core/utils"
	"concurrency-simulator/services/shared"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

func PaymentsController(w http.ResponseWriter, r *http.Request, kafkaProducer *kafka.Producer) {
	logger := utils.NewRequestLogger(r)

	topic := shared.PaymentTopic

	bodyMap, err := shared.ParseBody(r)

	if err != nil {
		logger.Error("Failed to parse request body",
			zap.Error(err),
			zap.String("content_type", r.Header.Get("Content-Type")),
		)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := validatePaymentBody(bodyMap); err != nil {
		logger.Warn("Invalid request body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accountId := bodyMap["account_id"]
	amount := bodyMap["amount"]
	installments := bodyMap["installments"]

	amountFloat, ok := amount.(float64)
	if !ok {
		logger.Error("Error: amount is not a float64")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalMessage := map[string]interface{}{
		"account_id":   accountId,
		"amount":       amountFloat,
		"installments": installments,
	}

	jsonMessage, err := json.Marshal(originalMessage)
	if err != nil {
		logger.Error("Failed to marshal payment message",
			zap.Error(err),
			zap.Int("message_size", len(jsonMessage)),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: shared.PartitionAlias["starting"],
		},
		Value: jsonMessage,
	}

	deliveryChan := make(chan kafka.Event)
	kafkaProducer.Produce(message, deliveryChan)

	go func() {
		for e := range deliveryChan {
			msg := e.(*kafka.Message)
			if msg.TopicPartition.Error != nil {
				logger.Error("Kafka delivery failed",
					zap.Error(msg.TopicPartition.Error),
					zap.String("topic", topic),
					zap.Int32("partition", msg.TopicPartition.Partition),
				)
			} else {
				logger.Info("Kafka delivery confirmed",
					zap.String("topic", topic),
					zap.Int32("partition", msg.TopicPartition.Partition),
					zap.Int64("offset", int64(msg.TopicPartition.Offset)),
				)
			}
		}
	}()

	pendingMessages := kafkaProducer.Flush(1000)
	if pendingMessages > 0 {
		logger.Warn("Messages still pending after flush",
			zap.Int("pending_count", pendingMessages),
		)
	}

	close(deliveryChan)
	logger.Info("Payment processing completed",
		zap.String("status", "success"),
	)

	w.Write([]byte("sent to kafka"))
}

func validatePaymentBody(bodyMap map[string]interface{}) error {
	if bodyMap["amount"] == nil || bodyMap["installments"] == nil || bodyMap["account_id"] == nil {
		return errors.New("amount, installments, account_id are required")
	}

	return nil
}
