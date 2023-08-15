package client

import (
	"encoding/json"
	"fmt"
	"math/rand"

	dto "github.com/praetorian-inc/product-golang-backend-interview/internal"
)

func ingestHandler(enumerateFn dto.EnumerateFn, producer dto.Producer, m dto.KafkaMessage) error {
	ingestDto, ok := m.Payload.(dto.IngestDto)
	if !ok {
		return fmt.Errorf("ingest was wrong type: %T", m.Payload)
	}

	for _, subdomain := range enumerateFn(ingestDto.Domain) {

		subdomainDto := dto.SubdomainDto{
			Id:     rand.Uint32(),
			Root:   ingestDto.Domain,
			Source: subdomain,
		}

		subdomainEvent := dto.KafkaMessage{
			Type:    "subdomainEvent",
			Payload: subdomainDto,
		}

		subdomainEventMessage, err := json.Marshal(subdomainEvent)
		if err != nil {
			return err
		}

		fmt.Printf("Producing message: %s\n", string(subdomainEventMessage))
		err = producer.ProduceMsg("subdomainEvent", subdomainEventMessage)
		if err != nil {
			return err
		}
	}

	return nil
}
