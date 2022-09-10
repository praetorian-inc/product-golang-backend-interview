# Enumeration Service

## Usage

```
./run.sh
```

If you get errors, try running `go mod tidy`.

## File structure

```
enumeration/
  dto/
   - Basic definitions of structs sent over Kafka
  handlers
   - domain_resolver.go: helper function for getting the status of a domain (either resolving or not resolving)
   - enumerate.go: business logic for enumerating the subdomains of a root domain
   - ingest.go: top-level handler for ingesting domain events
   - kafka_consumer.go: definition of the kafka consumer, multiplexes events based on "type" field
   - kafka_helpers.go: helper function for publishing kafka messages
   - server.go: definition of static objects in use by the handlers
```

## Interaction with other components

The enumeration service interfaces with Kafka only. The service consumes messages with the type `domainEvent` on the `ingest` topic, and it produces events with type `subdomainEvent` on the `subdomainEvent topic`.