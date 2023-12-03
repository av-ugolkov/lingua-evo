.PHONY: build
build:
	go build --gcflags="all=-N -l" -o ./cmd/main ./cmd/main.go

.PHONY: run
run:
	./cmd/main

.PHONY: docker-create
docker-create:
	docker compose -p lingua-evo -f deploy/docker-compose.local.yml up --build --force-recreate

.PHONY: docker-create-database
docker-create-database:
	docker compose -p lingua-evo -f deploy/docker-compose.database.local.yml up --build --force-recreate

.PHONY: lint
lint:
	@go version
	@golangci-lint --version
	GOWORK=off golangci-lint run ./...

.PHONY: test
test: 
	go test ./... -count=1