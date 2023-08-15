package kafka

import (
	"encoding/json"
	"fmt"
	"os"

	dto "github.com/praetorian-inc/product-golang-backend-interview/internal"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Consumer wraps a kafka consumer.
type Consumer struct{ *kafka.Consumer }

// Ensure Consumer satisfies the required interface.
var _ dto.Poller = (*Consumer)(nil)

// SetupKafkaConsumer returns a consumer subscribed to the defined topics.
func SetupKafkaConsumer(groupId string, topics []string) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          groupId,
		"auto.offset.reset": "smallest"})
	if err != nil {
		return nil, err
	}
	err = consumer.SubscribeTopics(topics, nil)

	return &Consumer{consumer}, err
}

// PollKafka polls kafka for the subscribed messages and calls the MessageHandlerFn on each.
func (c *Consumer) PollKafka(mh dto.MessageHandlerFn) {
	run := true
	fmt.Printf("Consuming kafka events from port 9092\n")
	for run {
		// 1 to lower CPU
		ev := c.Poll(1)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("Received message: %s\n", string(e.Value))

			// Unmarshal JSON into struct
			var message dto.KafkaMessage
			err := json.Unmarshal(e.Value, &message)
			if err != nil {
				fmt.Printf("Could not unmarshal message due to error: %s\n", err)
				continue
			}

			err = mh(message)
			if err != nil {
				fmt.Printf("Failed to handle message: %s\n", err)
				continue
			}

		case kafka.Error:
			_, err := fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			if err != nil {
				return
			}
			run = false
		}
	}

	fmt.Printf("Kafka Consumer exited.")
}
