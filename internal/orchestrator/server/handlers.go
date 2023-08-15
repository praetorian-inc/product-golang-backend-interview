package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	dto "github.com/praetorian-inc/product-golang-backend-interview/internal"
	"github.com/praetorian-inc/product-golang-backend-interview/internal/orchestrator/mysql"
)

func ingestHandler(sqlClient mysql.SqlClient, producer dto.Producer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(405) // METHOD_NOT_ALLOWED
			return
		}

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
		err = sqlClient.SaveDomain(dto.RootDomainDto{
			Id:   ingestDto.Id,
			Root: ingestDto.Domain,
		})
		if err != nil {
			fmt.Printf("Error while saving domain: %s\n", err.Error())
			w.WriteHeader(500) // INTERNAL_SERVER_ERROR
			return
		}

		ingestMessage := dto.KafkaMessage{
			Type:    "ingestDomain",
			Payload: ingestDto,
		}

		messageBytes, err := json.Marshal(ingestMessage)
		if err != nil {
			fmt.Printf("Error while marshalling injest message: %s\n", err.Error())
			w.WriteHeader(500) // INTERNAL_SERVER_ERROR
			return
		}

		fmt.Printf("Producing message %s\n", string(messageBytes))

		err = producer.ProduceMsg("ingest", messageBytes)
		if err != nil {
			fmt.Printf("Error while producing ingest message: %s\n", err.Error())
			w.WriteHeader(500) // INTERNAL_SERVER_ERROR
			return
		}
	}
}

func getDomainsHandler(sqlClient mysql.SqlClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405) // METHOD_NOT_ALLOWED
			return
		}

		domains, err := sqlClient.GetAllDomains()
		if err != nil {
			fmt.Printf("Could not get all domains due to error: %s", err.Error())
			w.WriteHeader(500) // INTERNAL_SERVER_ERROR
			return
		}

		domainJson, err := json.Marshal(domains)
		if err != nil {
			fmt.Printf("Could not marshall domains due to error: %s", err.Error())
			w.WriteHeader(500) // INTERNAL_SERVER_ERROR
			return
		}

		w.WriteHeader(200) // OK
		w.Write(domainJson)
	}
}

func getSubdomainsHandler(sqlClient mysql.SqlClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405) // METHOD_NOT_ALLOWED
			return
		}

		queryVars := r.URL.Query()
		pageResults, pageOk := queryVars["page"]
		limitResults, limitOk := queryVars["limit"]

		if !pageOk || !limitOk {
			w.WriteHeader(400) // BAD_REQUEST
			log.Println("page and limit are required query parameters")
			return
		}
		page, err := strconv.Atoi(pageResults[0])
		if err != nil {
			w.WriteHeader(400) // BAD_REQUEST
			return
		}
		limit, err := strconv.Atoi(limitResults[0])
		if err != nil {
			w.WriteHeader(400) // BAD_REQUEST
			return
		}

		log.Printf("page: %d, limit: %d", uint(page), uint(limit))

		subdomains, err := sqlClient.GetSubdomains(uint(page), uint(limit))
		if err != nil {
			fmt.Printf("Could not get all subdomains due to error: %s", err.Error())
			w.WriteHeader(500) // INTERNAL_SERVER_ERROR
			return
		}

		subdomainJson, err := json.Marshal(subdomains)
		if err != nil {
			fmt.Printf("Could not marshall subdomains due to error: %s", err.Error())
			w.WriteHeader(500) // INTERNAL_SERVER_ERROR
			return
		}

		w.WriteHeader(200) // OK
		w.Write(subdomainJson)
	}
}
