package handlers

import (
	"database/sql"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Server struct {
	Producer *kafka.Producer
	DB       *sql.DB
}
