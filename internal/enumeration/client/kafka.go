package client

import (
	"fmt"

	dto "github.com/praetorian-inc/product-golang-backend-interview/internal"
)

func handleMessage(producer dto.Producer, enumerateFn dto.EnumerateFn) func(message dto.KafkaMessage) error {
	return func(message dto.KafkaMessage) error {
		switch message.Type {
		case "ingestDomain":
			err := ingestHandler(enumerateFn, producer, message)
			if err != nil {
				return fmt.Errorf("could not ingest domain due to error %s", err)
			}
		default:
			return fmt.Errorf("unexpected message type: %s", message.Type)
		}

		return nil
	}
}

// PollAndUpdate runs the enumeration handler for when kafka messages are received.
func PollAndUpdate(consumer dto.Poller, producer dto.Producer, enumerateFn dto.EnumerateFn) {
	consumer.PollKafka(handleMessage(producer, enumerateFn))
}
