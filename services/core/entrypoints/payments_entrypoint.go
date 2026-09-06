package entrypoints

import (
	"net/http"

	"concurrency-simulator/services/core/controllers"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func PaymentsEntrypoints(kafkaProducer *kafka.Producer) {
	http.HandleFunc("/payments", func(w http.ResponseWriter, r *http.Request) {
		controllers.PaymentsController(w, r, kafkaProducer)
	})
}
