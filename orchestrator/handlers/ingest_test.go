package handlers

import (
	"bytes"
	"database/sql"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIngestHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	s, err := setupServer()
	if err != nil {
		t.Fatalf("Could not set up server due to error: %s", err.Error())
	}

	handler := createIngestHandler(&s)

	t.Run("Enqueue scan", func(t *testing.T) {

		req, err := http.NewRequest("POST", "/api/v1/ingest", bytes.NewBufferString("{\"Domain\":\"praetorian.com\"}"))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
	})
}

func setupServer() (Server, error) {

	db, err := getDbConnection()
	if err != nil {
		return Server{}, err
	}

	consumer, err := setupKafkaConsumer()
	if err != nil {
		return Server{}, err
	}

	producer, err := setupKafkaProducer()
	if err != nil {
		return Server{}, err
	}

	return Server{
		SqlClient: SqlClient{DB: db},
		Consumer:  consumer,
		Producer:  producer,
	}, nil
}

func createIngestHandler(s *Server) http.HandlerFunc {
	handler := http.HandlerFunc(s.IngestHandler)
	return handler
}

func setupKafkaConsumer() (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"test.mock.num.brokers": "3",
		"group.id":              "orchestrator",
		"auto.offset.reset":     "smallest"})
	if err != nil {
		return nil, err
	}
	err = consumer.SubscribeTopics([]string{"domainEvent", "subdomainEvent"}, nil)

	return consumer, err
}

func setupKafkaProducer() (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{"test.mock.num.brokers": "3"}) // Default port is 9092
}

func getDbConnection() (*sql.DB, error) {
	return sql.Open("mysql", "root:root@tcp(127.0.0.1:3307)/scanner")
}
