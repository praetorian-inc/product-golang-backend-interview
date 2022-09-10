# Orchestrator Service


## Usage
```
./run.sh
```
To install packages, execute `go mod tidy`.

## File structure
```
.
├── README.md
├── dto - Basic definitions of API requests & Kafka messages
│   ├── ingest.go
│   ├── kafka_message.go
│   ├── root_domain.go
│   └── subdomain.go
├── go.mod
├── go.sum
├── handlers
│   ├── domain.go - Viewing and saving domains
│   ├── ingest.go - Initiating enumeration of new domains
│   ├── ingest_test.go
│   ├── kafka_consumer.go - Consume domainEvent and subdomainEvent Kafka messages
│   ├── kafka_helpers.go
│   ├── server.go
│   └── subdomain.go - Viewing and saving subdomains
├── main.go
└── run.sh
```

The Orchestrator service listens for 3 API requests, for which API handlers
are defined in their respective files in `handlers/`.
- `/api/v1/ingest` - `handlers/ingest.go` - `IngestHandler`
- `/api/v1/domain` - `handlers/domain.go` - `GetDomainsHandler`
- `/api/v1/subdomain` - `handlers/subdomain.go - GetSubdomainsHandler`

The service listens for 2 Kafka events, with handlers defined in:
- `domainEvent` - `handlers/domain.go` - `DomainEventHandler`
- `subdomainEvent` - `handers/subdomain.go` - `SubdomainEventHandler`