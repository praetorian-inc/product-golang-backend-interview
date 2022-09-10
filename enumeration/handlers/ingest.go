package handlers

import (
	"encoding/json"
	"enumeration/dto"
	"fmt"
	"math/rand"
)

func (s *Server) IngestHandler(m dto.KafkaMessage) error {
	ingestDto, err := unmarshalIngestDtoHelper(m.Payload)
	if err != nil {
		return err
	}

	subdomains, err := enumerateSubdomains(ingestDto.Domain)
	if err != nil {
		return err
	}

	for _, subdomain := range subdomains {

		subdomainDto := dto.SubdomainDto{
			Id:     rand.Uint32(),
			Root:   ingestDto.Domain,
			Source: subdomain,
		}

		subdomainMarshal, err := marshalPayloadHelper(subdomainDto)
		if err != nil {
			return err
		}

		subdomainEvent := dto.KafkaMessage{
			Type:    "subdomainEvent",
			Payload: subdomainMarshal,
		}

		subdomainEventMessage, err := json.Marshal(subdomainEvent)
		if err != nil {
			return err
		}

		fmt.Printf("Producing message: %s\n", string(subdomainEventMessage))
		err = s.Produce("subdomainEvent", subdomainEventMessage)
		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalIngestDtoHelper(raw map[string]interface{}) (dto.IngestDto, error) {
	rawJson, err := json.Marshal(raw)
	if err != nil {
		return dto.IngestDto{}, err
	}

	// Convert json string to struct
	var ingestDto dto.IngestDto
	if err := json.Unmarshal(rawJson, &ingestDto); err != nil {
		return dto.IngestDto{}, err
	}

	return ingestDto, nil
}

func marshalPayloadHelper(payload interface{}) (map[string]interface{}, error) {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return map[string]interface{}{}, err
	}

	var payloadMarshal map[string]interface{}
	if err := json.Unmarshal(payloadJson, &payloadMarshal); err != nil {
		return map[string]interface{}{}, err
	}

	return payloadMarshal, nil
}
