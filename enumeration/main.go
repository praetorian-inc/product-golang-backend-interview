package main

import (
	"enumeration/handlers"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	consumer, err := setupKafkaConsumer()
	if err != nil {
		fmt.Printf("Error while setting up Kafka consumer: %s\n", err.Error())
		return
	}

	producer, err := setupKafkaProducer()
	if err != nil {
		fmt.Printf("Error while setting up Kafka producer: %s\n", err.Error())
		return
	}
	defer producer.Close()

	s := handlers.Server{
		Producer: producer,
		Consumer: consumer,
	}

	// Enter main loop and start listening for messages
	s.PollKafkaEvents()
}

func setupKafkaConsumer() (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "enumeration",
		"auto.offset.reset": "smallest"})
	if err != nil {
		return nil, err
	}
	err = consumer.SubscribeTopics([]string{"ingest"}, nil)

	return consumer, err
}

func setupKafkaProducer() (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"}) // Default port is 9092
}
