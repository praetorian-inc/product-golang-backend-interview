package handlers

import (
	"database/sql"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Server struct {
	Producer  *kafka.Producer
	SqlClient SqlClient
}

type SqlClient struct {
	DB *sql.DB
}
