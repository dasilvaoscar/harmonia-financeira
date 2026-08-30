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

	assingPartitions(consumer, logger)

	for {
		msg, err := consumer.ReadMessage(-1)

		if err != nil {
			logger.Error("Consumer error", zap.Error(err))
			continue
		}

		logger.Info("Mensagem lida", zap.String("Value", string(msg.Value)))
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

func assingPartitions(consumer *kafka.Consumer, logger *zap.Logger) {
	topic := shared.PaymentTopic
	err := consumer.Assign([]kafka.TopicPartition{
		{
			Topic:     &topic,
			Partition: shared.PartitionAlias["starting"],
		},
	})

	if err != nil {
		logger.Error("Failed to assign partitions", zap.Error(err))
		panic(err)
	}
}