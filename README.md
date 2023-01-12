# product-golang-backend-interview

This interview station involves a small microservice architecture that scans root domains for subdomains and persists that information in a SQL database.

![Microservice Documentation](microservice_documentation.png "Microservice Documentation")

## Dependencies

In order to run the Kafka zookeeper, Kafka broker, and MySQL server, you'll need to install [docker-compose](https://docs.docker.com/compose/install/).

In order to compile the two microservices, you'll need [Go 1.17](https://go.dev/doc/install) or newer.

## Component Breakdown

There are four components in this project:
 - Orchestrator (`orchestrator/`)
   - API server that serves for requests to scan a domain and view information about previously scanned domains and their subdomains
 - Enumeration (`enumeration/`)
   - Ingests a domain, iterates over its subdomains, and publishes individual events for the subdomains
 - MySQL (`docker-compose.yaml`)
   - Consists of a `scanner` database accessible to the root user
 - Kafka (`docker-compose.yaml`)
    - Consists of a zookeeper and a single broker

### Usage
To run the Kafka and MySQL containers, spin them up with docker-compose:

```
docker-compose up -d
```

Don't forget to create the kafka topics:
```
./create-topics.sh
```

To run the orchestrator and enumeration services, `cd` into the respective directory and run 
```
./run.sh
```

Once both services are running, start a subdomain by calling the orchestrator's ingest endpoint:
```
curl -v -X POST "localhost:9000/api/v1/ingest" -H "Content-type: application/json" -d '{"Domain": "praetorian.com"}'
```

### API

These are the endpoints that the orchestrator supports:
```
POST /api/v1/ingest
  example body: { "Domain": "praetorian.com" }
  
  Initiates a scan of the provided domain.

GET  /api/v1/domain
  Returns all domains that have been scanned since startup

GET  /api/v1/subdomain?page=<uint32>&limit=<uint32>
  Returns up to <limit> starting on page <page>.

```
