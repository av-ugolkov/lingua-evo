.PHONY: build
build:
	go build --gcflags="all=-N -l" -o ./cmd/main ./cmd/main.go

.PHONY: run
run:
	./cmd/main

.PHONY: backup
backup:
	./backup.sh

.PHONY: release
release:
	@echo ${shell ./deploy.sh release}

.PHONY: dev
dev:
	@echo ${shell ./deploy.sh dev}

.PHONY: run.docker.database
run.docker.database:
	docker compose -p lingua-evo-dev -f deploy/docker-compose.dev.yml up redis postgres migration --build --force-recreate

.PHONY: lint
lint:
	@go version
	@golangci-lint --version
	GOWORK=off golangci-lint run ./...

.PHONY: test
test: 
	go test ./... -count=1

.PHONY: count line
count line:
	find . -name '*.go' | xargs wc -l

.PHONY: pprof.cpu
pprof.cpu:
	go tool pprof -http=:8080 profile

.PHONY: pprof.heap
pprof.heap:
	go tool pprof -http=:8080 heap

.PHONY: pprof.trace
pprof.trace:
	go tool pprof -http=:8080 trace