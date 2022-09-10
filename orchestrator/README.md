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
├── dto
│   ├── domain.go
│   ├── ingest.go
│   ├── kafka_message.go
│   └── subdomain.go
├── go.mod
├── go.sum
├── handlers
│   ├── domain.go
│   ├── ingest.go
│   ├── kafka_consumer.go
│   ├── kafka_helpers.go
│   ├── server.go
│   └── subdomain.go
├── main.go
└── run.sh
```

## Interaction with other components