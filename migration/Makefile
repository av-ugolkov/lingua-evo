deps:
	go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: tool
tool: ## Run tools (vet, gofmt, goimports, tidy, etc.)
	@go version
	go mod tidy
	gofmt -w .
	goimports -w .
