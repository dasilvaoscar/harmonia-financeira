package entrypoints

import (
	"net/http"

	"concurrency-simulator/services/core/controllers"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func PaymentsEntrypoints(kafkaProducer *kafka.Producer) {
	ctrl := controllers.NewPaymentsController(kafkaProducer)
	http.HandleFunc("/payments", func(w http.ResponseWriter, r *http.Request) {
		ctrl.Execute(w, r)
	})
}
