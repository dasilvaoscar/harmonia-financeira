package controllers

import "github.com/confluentinc/confluent-kafka-go/kafka"

func NewAccountsController(kafkaProducer *kafka.Producer) *AccountsController {
	return &AccountsController{
		producer: kafkaProducer,
	}
}

func NewPaymentsController(kafkaProducer *kafka.Producer) *PaymentsController {
	return &PaymentsController{
		producer: kafkaProducer,
	}
}
