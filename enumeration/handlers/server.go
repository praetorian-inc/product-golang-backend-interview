package handlers

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Server struct {
	Producer *kafka.Producer
}
