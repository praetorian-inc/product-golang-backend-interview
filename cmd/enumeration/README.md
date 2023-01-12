# Enumeration Service

## Usage

```
./run.sh
```

To install packages, execute go mod tidy.

## File structure

```
.
├── README.md
├── dto - Basic definitions of structs sent over Kafka
│   ├── ingest.go
│   ├── kafka_message.go
│   ├── root_domain.go
│   └── subdomain.go
├── go.mod
├── go.sum
├── handlers
│   ├── domain_resolver.go - helper function for getting the status of a domain (either resolving or not resolving)
│   ├── domain_resolver_test.go
│   ├── enumerate.go - business logic for enumerating the subdomains of a root domain
│   ├── enumerate_test.go
│   ├── ingest.go - top-level handler for ingesting domain events
│   ├── kafka_consumer.go - definition of the kafka consumer, multiplexes events based on "type" field
│   ├── kafka_helpers.go - helper function for publishing kafka messages
│   └── server.go - definition of static objects in use by the handlers
├── main.go
└── run.sh
```

## Interaction with other components

The service listens for 1 Kafka events, with a handler defined in:

- `ingestDomain` - `handlers/ingest.go` - `IngestHandler`

This service also produces events with type `subdomainEvent` on the `subdomainEvent` topic.