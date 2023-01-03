package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"net/http"
	"orchestrator/handlers"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	port := flag.Int("listenPort", 9000, "Port on which to serve HTTP requests")
	rand.Seed(time.Now().UnixNano())

	// Setup Kafka consumer
	consumer, err := setupKafkaConsumer()
	if err != nil {
		fmt.Printf("Error while setting up Kafka consumer: %s\n", err.Error())
		return
	}

	// Setup Kafka producer
	producer, err := setupKafkaProducer()
	if err != nil {
		fmt.Printf("Error while setting up Kafka producer: %s\n", err.Error())
		return
	}
	defer producer.Close()

	// Setup MySQL connection
	db, err := getDbConnection()
	if err != nil {
		fmt.Printf("Error while getting MySQL connection: %s\n", err.Error())
		return
	}
	defer db.Close()

	s := handlers.Server{
		Producer:  producer,
		Consumer:  consumer,
		SqlClient: handlers.SqlClient{DB: db},
	}

	// Consume kafka events asynchronously
	go s.PollKafkaEvents()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/ingest", s.IngestHandler)
	mux.HandleFunc("/api/v1/domain", s.GetDomainsHandler)
	mux.HandleFunc("/api/v1/subdomain", s.GetSubdomainsHandler)

	fmt.Println("Listening on localhost:" + strconv.Itoa(*port))
	err = http.ListenAndServe("localhost:"+strconv.Itoa(*port), mux)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func setupKafkaConsumer() (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "orchestrator",
		"auto.offset.reset": "smallest"})
	if err != nil {
		return nil, err
	}
	err = consumer.SubscribeTopics([]string{"domainEvent", "subdomainEvent"}, nil)

	return consumer, err
}

func setupKafkaProducer() (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"}) // Default port is 9092
}

func getDbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3307)/scanner")
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
