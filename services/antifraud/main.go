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

	logger.Info("Consumer started, listening to topic", zap.String("topic", shared.PaymentTopic))

	defer consumer.Close()

	for {
		msg, err := consumer.ReadMessage(-1)

		if err != nil {
			logger.Error("Consumer error", zap.Error(err))
			continue
		}

		logger.Info("Received message from topic", zap.String("topic", *msg.TopicPartition.Topic), zap.String("message", string(msg.Value)))

		if shared.GetKafkaHeader(*msg, "event_type") == "fraud-validation" {
			controller.ProcessMessage(msg)
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
