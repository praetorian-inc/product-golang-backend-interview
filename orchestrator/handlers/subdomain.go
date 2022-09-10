package handlers

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"orchestrator/dto"
	"strconv"
)

func (s *Server) GetSubdomainsHandler(w http.ResponseWriter, r *http.Request) {
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

	subdomains, err := s.SqlClient.GetSubdomains(uint(page), uint(limit))
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
	return

}

func (sqlClient SqlClient) GetSubdomains(page uint, limit uint) ([]dto.SubdomainDto, error) {
	query := fmt.Sprintf("SELECT * FROM subdomain LIMIT %d, %d;", (page-1)*limit, limit)
	fmt.Printf("query: %s", query)
	rows, err := sqlClient.DB.Query(query)
	if err != nil {
		return []dto.SubdomainDto{}, err
	}

	defer rows.Close()

	var subdomains []dto.SubdomainDto
	for rows.Next() {
		var subdomain dto.SubdomainDto

		err := rows.Scan(&subdomain.Id, &subdomain.Root, &subdomain.Source)
		if err != nil {
			return []dto.SubdomainDto{}, err
		}

		subdomains = append(subdomains, subdomain)
	}

	return subdomains, nil
}

func (sqlClient SqlClient) SaveSubdomain(subdomain dto.SubdomainDto) error {
	query := fmt.Sprintf("INSERT INTO subdomain (id, root, source) VALUES(%d, '%s', '%s');", subdomain.Id, subdomain.Root, subdomain.Source)
	_, err := sqlClient.DB.Exec(query)
	return err
}
