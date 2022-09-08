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

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "orchestrator",
		"auto.offset.reset": "smallest"})
	if err != nil {
		panic(err)
	}
	err = consumer.SubscribeTopics([]string{"domainEvent", "subdomainEvent"}, nil)
	if err != nil {
		panic(err)
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/scanner")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.Exec("DROP TABLE IF EXISTS `root_domain`;")
	db.Exec("CREATE TABLE `root_domain` (`id` int(8) unsigned NOT NULL, root varchar(32) NOT NULL, status varchar(32), owner varchar(32), PRIMARY KEY (`id`));")

	db.Exec("DROP TABLE IF EXISTS `subdomain`;")
	db.Exec("CREATE TABLE `subdomain` (`id` int(8) unsigned NOT NULL, root varchar(32) NOT NULL, source varchar(32), PRIMARY KEY (`id`));")

	s := handlers.Server{
		Producer: p,
	}

	// Consume kafka events asynchronously
	go s.PollKafkaEvents(consumer)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/ingest", s.IngestHandler)
	mux.HandleFunc("/api/v1/domain", s.IngestHandler)
	mux.HandleFunc("/api/v1/subdomain", s.IngestHandler)

	fmt.Println("Listening on localhost:" + strconv.Itoa(*port))
	err = http.ListenAndServe("localhost:"+strconv.Itoa(*port), mux)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
