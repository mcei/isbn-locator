install-lint-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run

test:
	go test -v ./...

start:
	docker compose -f ./docker-compose.yml -p isbn-locator up -d

stop:
	docker compose down