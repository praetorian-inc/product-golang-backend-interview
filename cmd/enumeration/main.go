package main

import (
	"fmt"

	"github.com/praetorian-inc/product-golang-backend-interview/internal/enumeration"
	"github.com/praetorian-inc/product-golang-backend-interview/internal/enumeration/client"
	"github.com/praetorian-inc/product-golang-backend-interview/internal/kafka"
)

func main() {
	consumer, err := kafka.SetupKafkaConsumer("enumeration", []string{"ingest"})
	if err != nil {
		fmt.Printf("Error while setting up Kafka consumer: %s\n", err.Error())
		return
	}

	producer, err := kafka.SetupKafkaProducer()
	if err != nil {
		fmt.Printf("Error while setting up Kafka producer: %s\n", err.Error())
		return
	}
	defer producer.Close()

	// Enter main loop and start listening for messages
	client.PollAndUpdate(consumer, producer, enumeration.Subdomains)
}
