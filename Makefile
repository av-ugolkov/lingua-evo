.PHONY: build
build:
	go build -o main ./app/main.go

.PHONY: run
run:
	go build -o main ./app/main.go
	./main

.PHONY: docker recreate
docker recreate:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./main ./app/main.go
	docker compose -f docker-compose.local.yml up --build --force-recreate
