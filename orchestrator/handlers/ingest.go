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

	ingestMarshal, err := marshalIngestDtoHelper(ingestDto)
	if err != nil {
		w.WriteHeader(500) // INTERNAL_SERVER_ERROR
		return
	}

	ingestMessage := dto.KafkaMessage{
		Type:    "ingestDomain",
		Payload: ingestMarshal,
	}

	messageBytes, err := json.Marshal(ingestMessage)
	if err != nil {
		w.WriteHeader(500) // INTERNAL_SERVER_ERROR
		return
	}

	fmt.Printf("Producing message %s\n", string(messageBytes))

	err = s.Produce(topic, messageBytes)
	if err != nil {
		w.WriteHeader(500) // INTERNAL_SERVER_ERROR
		return
	}
}

func marshalIngestDtoHelper(ingestDto dto.IngestDto) (map[string]interface{}, error) {
	ingestJson, err := json.Marshal(ingestDto)
	if err != nil {
		return map[string]interface{}{}, err
	}

	var ingestMarshal map[string]interface{}
	if err := json.Unmarshal(ingestJson, &ingestMarshal); err != nil {
		return map[string]interface{}{}, err
	}

	return ingestMarshal, nil
}