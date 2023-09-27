.PHONY: build
build:
	go build --gcflags="all=-N -l" -o main ./cmd/main.go

.PHONY: run
run:
	go build --gcflags="all=-N -l" -o ./cmd/main ./cmd/main.go
	./cmd/main

.PHONY: docker recreate
docker recreate:
	docker compose -f deploy/docker-compose.local.yml up --build --force-recreate

.PHONY: docker recreate database
docker recreate:
	docker compose -p lingua-evo -f deploy/docker-compose.database.local.yml up --build --force-recreate

.PHONY: lint
lint:
	@go version
	@golangci-lint --version
	golangci-lint run ./...