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

func (s *Server) DomainEventHandler(m dto.KafkaMessage) error {
	domainDto, err := unmarshalDomainDtoHelper(m.Payload)
	if err != nil {
		return err
	}

	fmt.Printf("Received DomainDto: %v\n", domainDto)

	err = s.SqlClient.SaveDomain(domainDto)
	if err != nil {
		return err
	}

	fmt.Println("Inserted domain")

	return nil
}

func (s *Server) SubdomainEventHandler(m dto.KafkaMessage) error {
	subdomainDto, err := unmarshalSubdomainDtoHelper(m.Payload)
	if err != nil {
		return err
	}

	fmt.Printf("Received SubdomainDto: %v\n", subdomainDto)

	err = s.SqlClient.SaveSubdomain(subdomainDto)
	if err != nil {
		return err
	}

	fmt.Println("Inserted subdomain")

	return nil
}

func unmarshalDomainDtoHelper(raw map[string]interface{}) (dto.DomainDto, error) {
	rawJson, err := json.Marshal(raw)
	if err != nil {
		return dto.DomainDto{}, err
	}

	// Convert json string to struct
	var domainDto dto.DomainDto
	if err := json.Unmarshal(rawJson, &domainDto); err != nil {
		return dto.DomainDto{}, err
	}

	return domainDto, nil
}

func unmarshalSubdomainDtoHelper(raw map[string]interface{}) (dto.SubdomainDto, error) {
	rawJson, err := json.Marshal(raw)
	if err != nil {
		return dto.SubdomainDto{}, err
	}

	// Convert json string to struct
	var subdomainDto dto.SubdomainDto
	if err := json.Unmarshal(rawJson, &subdomainDto); err != nil {
		return dto.SubdomainDto{}, err
	}

	return subdomainDto, nil
}
