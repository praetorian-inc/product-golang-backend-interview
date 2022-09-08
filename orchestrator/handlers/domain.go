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

func (sqlClient SqlClient) GetAllDomains() ([]dto.DomainDto, error) {
	query := fmt.Sprintf("SELECT * FROM root_domain;")
	rows, err := sqlClient.DB.Query(query)
	if err != nil {
		return []dto.DomainDto{}, err
	}

	defer rows.Close()

	var domains []dto.DomainDto
	for rows.Next() {
		var domain dto.DomainDto

		err := rows.Scan(&domain.Id, &domain.Root, &domain.Status, &domain.Owner)
		if err != nil {
			return []dto.DomainDto{}, err
		}

		domains = append(domains, domain)
	}

	return domains, nil
}

func (sqlClient SqlClient) SaveDomain(domain dto.DomainDto) error {
	query := fmt.Sprintf("INSERT INTO root_domain (id, root, status, owner) VALUES(%d, '%s', '%s', '%s');", domain.Id, domain.Root, domain.Status, domain.Owner)
	_, err := sqlClient.DB.Exec(query)
	return err
}
