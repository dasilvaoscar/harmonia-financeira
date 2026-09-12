package entrypoints

import (
	"net/http"

	"concurrency-simulator/services/core/controllers"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func AccountsEntrypoint(kafkaProducer *kafka.Producer) {
	ctrl := controllers.NewAccountsController(kafkaProducer)
	http.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		ctrl.Execute(w, r)
	})
}
