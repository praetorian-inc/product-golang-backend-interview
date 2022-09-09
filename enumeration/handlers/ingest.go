package handlers

import (
	"encoding/json"
	"enumeration/dto"
	"fmt"
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
		domainDto := dto.DomainDto{
			Id:     ingestDto.Id,
			Root:   ingestDto.Domain,
			Domain:	subdomain,
			Status: "SCANNED",
			Owner:  "",
		}

		domainMarshal, err := marshalDomainDtoHelper(domainDto)
		if err != nil {
			return err
		}

		domainEvent := dto.KafkaMessage{
			Type:    "domainEvent",
			Payload: domainMarshal,
		}

		domainEventMessage, err := json.Marshal(domainEvent)
		if err != nil {
			return err
		}

		fmt.Printf("Producing message: %s\n", string(domainEventMessage))
		err = s.Produce("domainEvent", domainEventMessage)
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

func marshalDomainDtoHelper(domainDto dto.DomainDto) (map[string]interface{}, error) {
	rawJson, err := json.Marshal(domainDto)
	if err != nil {
		return map[string]interface{}{}, err
	}

	var domainMarshal map[string]interface{}
	if err := json.Unmarshal(rawJson, &domainMarshal); err != nil {
		return map[string]interface{}{}, err
	}

	return domainMarshal, nil
}
