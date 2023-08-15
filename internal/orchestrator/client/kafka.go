package client

import (
	"fmt"

	dto "github.com/praetorian-inc/product-golang-backend-interview/internal"
	"github.com/praetorian-inc/product-golang-backend-interview/internal/orchestrator/mysql"
)

func handleMessage(sqlClient mysql.SqlClient) func(message dto.KafkaMessage) error {
	return func(message dto.KafkaMessage) error {
		switch message.Type {
		case "domainEvent":
			err := domainEventHandler(sqlClient, message)
			if err != nil {
				return fmt.Errorf("could not ingest domain due to error %s", err)
			}
		case "subdomainEvent":
			err := subdomainEventHandler(sqlClient, message)
			if err != nil {
				return fmt.Errorf("could not ingest subdomain due to error %s", err)
			}
		default:
			return fmt.Errorf("unexpected message type: %s", message.Type)
		}

		return nil
	}
}

// PollAndUpdate runs the orchestration handler for when kafka messages are received.
func PollAndUpdate(consumer dto.Poller, sqlClient mysql.SqlClient) {
	consumer.PollKafka(handleMessage(sqlClient))
}
