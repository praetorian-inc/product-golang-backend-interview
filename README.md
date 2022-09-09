# product-golang-backend-interview

This interview station involves a small microservice architecture that scans root domains for subdomains and persists that information in a SQL database.

![Microservice Documentation](microservice_documentation.png "Microservice Documentation")

### To run:
```$ docker-compose up -d```

Don't forget to create the kafka topic:
```
docker exec broker \
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic ingest
```

```
docker exec broker \
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic domainEvent
```

```
docker exec broker \
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic subdomainEvent
```