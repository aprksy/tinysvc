.PHONY: help build run test clean migrate dev

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/tinysvc cmd/server/main.go

run: build ## Build and run the application
	./bin/tinysvc

dev: ## Run with auto-reload (requires air: go install github.com/cosmtrek/air@latest)
	air

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

coverage: test ## Show test coverage
	go tool cover -html=coverage.out

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

deps: ## Download dependencies
	go mod download
	go mod tidy

docker-build: ## Build Docker image
	docker build -t tinysvc:latest .

.DEFAULT_GOAL := help