package shared

import "github.com/confluentinc/confluent-kafka-go/kafka"

func GetKafkaHeader(msg kafka.Message, key string) string {
	for _, h := range msg.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}
