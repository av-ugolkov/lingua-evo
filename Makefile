.PHONY: build
build:
	go build --gcflags="all=-N -l" -o main ./cmd/main.go

.PHONY: run
run:
	go build --gcflags="all=-N -l" -o ./cmd/main ./cmd/main.go
	./cmd/main

.PHONY: docker recreate
docker recreate:
	docker compose -f docker-compose.local.yml up --build --force-recreate

.PHONY: lint
lint:
	@go version
	@golangci-lint --version
	golangci-lint run ./...