package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/praetorian-inc/product-golang-backend-interview/internal/kafka"
	"github.com/praetorian-inc/product-golang-backend-interview/internal/orchestrator/mysql"
)

func TestIngestHandler(t *testing.T) {
	db, err := mysql.GetDbConnection("root:root@tcp(127.0.0.1:3317)/scanner")
	if err != nil {
		t.Fatal(err)
	}

	producer, err := kafka.SetupKafkaProducer()
	if err != nil {
		t.Fatal(err)
	}

	handler := ingestHandler(mysql.SqlClient{DB: db}, producer)

	t.Run("Enqueue scan", func(t *testing.T) {
		rr := httptest.NewRecorder()

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
