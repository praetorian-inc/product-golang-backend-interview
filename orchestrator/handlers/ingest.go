package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"orchestrator/dto"
)

func (s *Server) IngestHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(405) // METHOD_NOT_ALLOWED
		return
	}

	topic := "ingest"

	// Decode JSON
	var payload json.RawMessage
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&payload)
	if err != nil {
		w.WriteHeader(400) // BAD_REQUEST
		return
	}

	// Unmarshal JSON into struct
	var ingestDto dto.IngestDto
	err = json.Unmarshal(payload, &ingestDto)
	if err != nil {
		w.WriteHeader(400) // BAD_REQUEST
		return
	}

	// If the ingestDto does not have an ID, we should generate a random one.
	if ingestDto.Id == 0 {
		ingestDto.Id = rand.Uint32()
	}

	// Save the domain so we can track scanning progress
	err = s.SqlClient.SaveDomain(dto.RootDomainDto{
		Id:   ingestDto.Id,
		Root: ingestDto.Domain,
	})
	if err != nil {
		fmt.Printf("Error while saving domain: %s\n", err.Error())
		w.WriteHeader(500) // INTERNAL_SERVER_ERROR
		return
	}

	// Marshall into a send-able message
	ingestMarshal, err := marshalPayloadHelper(ingestDto)
	if err != nil {
		fmt.Printf("Error while marshalling payload: %s\n", err.Error())
		w.WriteHeader(500) // INTERNAL_SERVER_ERROR
		return
	}
	ingestMessage := dto.KafkaMessage{
		Type:    "ingestDomain",
		Payload: ingestMarshal,
	}

	messageBytes, err := json.Marshal(ingestMessage)
	if err != nil {
		fmt.Printf("Error while marshalling injest message: %s\n", err.Error())
		w.WriteHeader(500) // INTERNAL_SERVER_ERROR
		return
	}

	fmt.Printf("Producing message %s\n", string(messageBytes))

	err = s.Produce(topic, messageBytes)
	if err != nil {
		fmt.Printf("Error while producing ingest message: %s\n", err.Error())
		w.WriteHeader(500) // INTERNAL_SERVER_ERROR
		return
	}
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
