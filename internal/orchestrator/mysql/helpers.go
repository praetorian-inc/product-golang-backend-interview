package mysql

import (
	"database/sql"
	"fmt"

	dto "github.com/praetorian-inc/product-golang-backend-interview/internal"

	_ "github.com/go-sql-driver/mysql"
)

// GetDbConnection initializes and returns a mysql sql.DB.
func GetDbConnection(conn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `root_domain` (`id` int(8) unsigned NOT NULL, root varchar(32) NOT NULL, status varchar(32), owner varchar(32), PRIMARY KEY (`id`));")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `subdomain` (`id` int(8) unsigned NOT NULL, root varchar(32) NOT NULL, source varchar(256), PRIMARY KEY (`id`));")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return db, nil
}

// SqlClient wraps the DB and provides some helper functions for interacting with the data.
type SqlClient struct {
	DB *sql.DB
}

// GetSubdomains returns subdomains from the database.
func (s SqlClient) GetSubdomains(page uint, limit uint) ([]dto.SubdomainDto, error) {
	fmt.Printf("query: %d %d", (page-1)*limit, limit)
	rows, err := s.DB.Query("SELECT * FROM subdomain LIMIT ?, ?;", (page-1)*limit, limit)
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

// SaveSubdomain inserts the subdomain into the database.
func (s SqlClient) SaveSubdomain(subdomain dto.SubdomainDto) error {
	_, err := s.DB.Exec("INSERT INTO subdomain (id, root, source) VALUES(?, ?, ?);",
		subdomain.Id, subdomain.Root, subdomain.Source)
	return err
}

// GetAllDomains returns all domains from the database
func (s SqlClient) GetAllDomains() ([]dto.RootDomainDto, error) {
	query := "SELECT * FROM root_domain;"
	rows, err := s.DB.Query(query)
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

// SaveDomain writes a domain to the database.
func (s SqlClient) SaveDomain(domain dto.RootDomainDto) error {
	_, err := s.DB.Exec("INSERT INTO root_domain (id, root, status, owner) VALUES(?, ?, ?, ?);",
		domain.Id, domain.Root, domain.Status, domain.Owner)
	return err
}
