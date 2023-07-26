.PHONY: build
build:
	go build -o main ./app/main.go

.PHONY: run
run:
	go build -o ./app/main ./app/main.go
	./app/main

.PHONY: docker recreate
docker recreate:
	docker compose -f docker-compose.local.yml up --build --force-recreate

.PHONY: lint
lint:
	@go version
	@golangci-lint --version
	golangci-lint run ./...