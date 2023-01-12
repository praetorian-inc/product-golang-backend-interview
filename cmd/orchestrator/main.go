package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/praetorian-inc/product-golang-backend-interview/internal/kafka"
	"github.com/praetorian-inc/product-golang-backend-interview/internal/orchestrator/client"
	"github.com/praetorian-inc/product-golang-backend-interview/internal/orchestrator/mysql"
	"github.com/praetorian-inc/product-golang-backend-interview/internal/orchestrator/server"
)

func main() {
	port := flag.Int("listenPort", 9000, "Port on which to serve HTTP requests")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// Setup Kafka consumer
	consumer, err := kafka.SetupKafkaConsumer("orchestrator", []string{"domainEvent", "subdomainEvent"})
	if err != nil {
		fmt.Printf("Error while setting up Kafka consumer: %s\n", err.Error())
		return
	}

	// Setup Kafka producer
	producer, err := kafka.SetupKafkaProducer()
	if err != nil {
		fmt.Printf("Error while setting up Kafka producer: %s\n", err.Error())
		return
	}
	defer producer.Close()

	// Setup MySQL connection
	db, err := mysql.GetDbConnection("root:root@tcp(127.0.0.1:3307)/scanner")
	if err != nil {
		fmt.Printf("Error while getting MySQL connection: %s\n", err.Error())
		return
	}
	defer db.Close()

	sqlClient := mysql.SqlClient{DB: db}

	// Consume kafka events asynchronously
	go client.PollAndUpdate(consumer, sqlClient)

	// Serve the API
	err = server.ListenAndServe(*port, producer, sqlClient)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
