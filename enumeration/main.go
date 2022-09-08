package main

import (
	"encoding/json"
	"enumeration/dto"
	"enumeration/handlers"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "enumeration",
		"auto.offset.reset": "smallest"})
	if err != nil {
		panic(err)
	}
	err = consumer.SubscribeTopics([]string{"ingest"}, nil)
	if err != nil {
		panic(err)
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	s := handlers.Server{
		Producer: p,
	}

	run := true
	for run == true {
		ev := consumer.Poll(0)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("Received message: %s\n", string(e.Value))

			// Unmarshal JSON into struct
			var message dto.KafkaMessage
			err = json.Unmarshal(e.Value, &message)
			if err != nil {
				fmt.Printf("Could not unmarshal message due to error: %s\n", err)
				continue
			}

			switch message.Type {
			case "ingestDomain":
				err := s.IngestHandler(message)
				if err != nil {
					fmt.Printf("Could not ingest domain due to error %s\n", err)
					continue
				}
			default:
				fmt.Printf("Unexpected message type: %s\n", message.Type)
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
}
