# product-golang-backend-interview

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