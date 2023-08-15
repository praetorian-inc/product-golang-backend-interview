package client

import (
	"fmt"

	dto "github.com/praetorian-inc/product-golang-backend-interview/internal"
	"github.com/praetorian-inc/product-golang-backend-interview/internal/orchestrator/mysql"
)

func domainEventHandler(sqlClient mysql.SqlClient, m dto.KafkaMessage) error {
	domainDto, ok := m.Payload.(dto.RootDomainDto)
	if !ok {
		return fmt.Errorf("domain event was wrong type: %T", m.Payload)
	}

	fmt.Printf("Received RootDomainDto: %v\n", domainDto)

	err := sqlClient.SaveDomain(domainDto)
	if err != nil {
		return err
	}

	fmt.Println("Inserted domain")

	return nil
}

func subdomainEventHandler(sqlClient mysql.SqlClient, m dto.KafkaMessage) error {
	subdomainDto, ok := m.Payload.(dto.SubdomainDto)
	if !ok {
		return fmt.Errorf("subdomain event was wrong type: %T", m.Payload)
	}

	fmt.Printf("Received SubdomainDto: %v\n", subdomainDto)

	err := sqlClient.SaveSubdomain(subdomainDto)
	if err != nil {
		return err
	}

	fmt.Println("Inserted subdomain")

	return nil
}
