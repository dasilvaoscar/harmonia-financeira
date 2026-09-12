package main

import (
	"concurrency-simulator/services/account/utils"
	"concurrency-simulator/services/shared"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

func main() {
	logger := utils.NewRequestLogger()
	consumer := createConsumer(logger)

	subscribeToTopic(consumer, logger)

	for {
		msg, err := consumer.ReadMessage(-1)

		if err != nil {
			logger.Error("Consumer error", zap.Error(err))
			continue
		}

		logger.Info("Received message from topic", zap.String("topic", *msg.TopicPartition.Topic), zap.String("message", string(msg.Value)))

		if shared.GetKafkaHeader(*msg, "event_type") == "execute-transaction" {

		} else {
			logger.Info("Message skiped")
		}
	}
}

func createConsumer(logger *zap.Logger) *kafka.Consumer {
	consumer, err := kafka.NewConsumer(utils.GetKafkaConfig())

	if err != nil {
		logger.Error("Consumer creation error", zap.Error(err))
		panic(err)
	}

	return consumer
}

func subscribeToTopic(consumer *kafka.Consumer, logger *zap.Logger) {
	err := consumer.SubscribeTopics([]string{shared.PaymentTopic}, nil)

	if err != nil {
		logger.Error("Failed to subscribe to topics", zap.Error(err))
		panic(err)
	}
}
