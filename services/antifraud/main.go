package main

import (
	"concurrency-simulator/services/antifraud/controllers"
	"concurrency-simulator/services/antifraud/utils"
	"concurrency-simulator/services/shared"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

func main() {
	logger := utils.NewRequestLogger()
	execution(logger)
}

func execution(logger *zap.Logger) {
	consumer := createConsumer(logger)

	controller := controllers.NewAntifraudController()
	
	subscribeToTopic(consumer, logger)

	logger.Error("Consumer started, listening to topic", zap.String("topic", shared.PaymentTopic))

	defer consumer.Close()

	for {
		msg, err := consumer.ReadMessage(-1)

		if err != nil {
			logger.Error("Consumer error", zap.Error(err))
			continue
		}

		controller.ProcessMessage(msg)
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
