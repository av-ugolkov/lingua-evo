.PHONY: tool
tool: ## Run tools (vet, gofmt, goimports, tidy, etc.)
	@go version
	go mod tidy
	gofmt -w .
	goimports -w .
