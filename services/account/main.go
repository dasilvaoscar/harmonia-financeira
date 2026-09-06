package main

import (
	"concurrency-simulator/services/account/controllers"
	"concurrency-simulator/services/account/utils"
	"concurrency-simulator/services/shared"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

func main() {
	logger := utils.NewRequestLogger()

	var wg sync.WaitGroup
	wg.Add(1)

	go execution(&wg, logger)

	wg.Wait()
}

func execution(wg *sync.WaitGroup, logger *zap.Logger) {
	defer wg.Done()

	consumer := createConsumer(logger)

	controller := controllers.NewAccountController()

	subscribeToTopic(consumer, logger)

	logger.Error("Consumer started, listening to topic", zap.String("topic", shared.AccountTopic))

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
	err := consumer.SubscribeTopics([]string{shared.AccountTopic}, nil)

	if err != nil {
		logger.Error("Failed to subscribe to topics", zap.Error(err))
		panic(err)
	}
}
