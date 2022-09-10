package handlers

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"orchestrator/dto"
)

func (s *Server) GetDomainsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(405) // METHOD_NOT_ALLOWED
		return
	}

	domains, err := s.SqlClient.GetAllDomains()
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
	return

}

func (s *Server) DomainEventHandler(m dto.KafkaMessage) error {
	domainDto, err := unmarshalDomainDtoHelper(m.Payload)
	if err != nil {
		return err
	}

	fmt.Printf("Received RootDomainDto: %v\n", domainDto)

	err = s.SqlClient.SaveDomain(domainDto)
	if err != nil {
		return err
	}

	fmt.Println("Inserted domain")

	return nil
}

func unmarshalDomainDtoHelper(raw map[string]interface{}) (dto.RootDomainDto, error) {
	rawJson, err := json.Marshal(raw)
	if err != nil {
		return dto.RootDomainDto{}, err
	}

	// Convert json string to struct
	var domainDto dto.RootDomainDto
	if err := json.Unmarshal(rawJson, &domainDto); err != nil {
		return dto.RootDomainDto{}, err
	}

	return domainDto, nil
}

func (sqlClient SqlClient) GetAllDomains() ([]dto.RootDomainDto, error) {
	query := fmt.Sprintf("SELECT * FROM root_domain;")
	rows, err := sqlClient.DB.Query(query)
	if err != nil {
		return []dto.RootDomainDto{}, err
	}

	defer rows.Close()

	var domains []dto.RootDomainDto
	for rows.Next() {
		var domain dto.RootDomainDto

		err := rows.Scan(&domain.Id, &domain.Root, &domain.Status, &domain.Owner)
		if err != nil {
			return []dto.RootDomainDto{}, err
		}

		domains = append(domains, domain)
	}

	return domains, nil
}

func (sqlClient SqlClient) SaveDomain(domain dto.RootDomainDto) error {
	query := fmt.Sprintf("INSERT INTO root_domain (id, root, status, owner) VALUES(%d, '%s', '%s', '%s');", domain.Id, domain.Root, domain.Status, domain.Owner)
	_, err := sqlClient.DB.Exec(query)
	return err
}
