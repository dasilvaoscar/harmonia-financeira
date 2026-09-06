package entrypoints

import (
	"net/http"

	"concurrency-simulator/services/core/controllers"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func AccountsEntrypoint(kafkaProducer *kafka.Producer) {
	http.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		controllers.AccountsController(w, r, kafkaProducer)
	})
}
