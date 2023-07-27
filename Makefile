.PHONY: build
build:
	go build -o main ./cmd/main.go

.PHONY: run
run:
	go build -o ./cmd/main ./cmd/main.go
	./cmd/main

.PHONY: docker recreate
docker recreate:
	docker compose -f docker-compose.local.yml up --build --force-recreate

.PHONY: lint
lint:
	@go version
	@golangci-lint --version
	golangci-lint run ./...