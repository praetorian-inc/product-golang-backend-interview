test:
	docker-compose -f docker-compose-test.yml up -d && \
	go test -cover -count 1 ./...; \
	docker-compose -f docker-compose-test.yml down

prepare:
	docker-compose up -d && bash ./create-topics.sh

tidy:
	docker-compose down

run-enumerate:
	go run ./cmd/enumeration/main.go

run-orchestrate:
	go run ./cmd/orchestrator/main.go -listenPort 9000

.PHONY: test prepare tidy run-enumerate run-orchestrate
