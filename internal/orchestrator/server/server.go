package server

import (
	"fmt"
	"net/http"
	"strconv"

	dto "github.com/praetorian-inc/product-golang-backend-interview/internal"
	"github.com/praetorian-inc/product-golang-backend-interview/internal/orchestrator/mysql"
)

// ListenAndServe registers the HTTP API handlsers and serves the app on the specified port.
func ListenAndServe(port int, producer dto.Producer, sqlClient mysql.SqlClient) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/ingest", ingestHandler(sqlClient, producer))
	mux.HandleFunc("/api/v1/domain", getDomainsHandler(sqlClient))
	mux.HandleFunc("/api/v1/subdomain", getSubdomainsHandler(sqlClient))

	fmt.Println("Listening on localhost:" + strconv.Itoa(port))
	return http.ListenAndServe("localhost:"+strconv.Itoa(port), mux)
}
