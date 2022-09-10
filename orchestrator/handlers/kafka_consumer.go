package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"orchestrator/dto"
	"os"
)

func (s *Server) PollKafkaEvents() {
	run := true
	fmt.Printf("Consuming kafka events from port 9092\n")
	for run == true {
		ev := s.Consumer.Poll(0)
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

			switch message.Type {
			case "domainEvent":
				err := s.DomainEventHandler(message)
				if err != nil {
					fmt.Printf("Could not ingest domain due to error %s\n", err)
					continue
				}
			case "subdomainEvent":
				err := s.SubdomainEventHandler(message)
				if err != nil {
					fmt.Printf("Could not ingest subdomain due to error %s\n", err)
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

	fmt.Printf("Kafka Consumer exited.")
}
