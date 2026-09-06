package external

import (
	"errors"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewKafkaProducer() (*kafka.Producer, error) {
	server := os.Getenv("KAFKA_BROKER")

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  server,
		"enable.idempotence": "true",
		"acks":               "all",
		"linger.ms":          0,
		"batch.size":         1,
	})

	if err != nil {
		return nil, errors.New("Producer creation error")
	}

	return producer, nil
}
