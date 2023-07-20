.PHONY: build
build:
	go build -o main .

.PHONY: postgres
postgres:

.PHONY: docker recreate
docker recreate:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./app/main ./app/main.go
	docker compose -f docker-compose.local.yml up --build --force-recreate
