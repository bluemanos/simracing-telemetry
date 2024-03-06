.PHONY: format gqlgen init lint lint-fix test test-go-reset-db-local dump_testsetup_db

format:
	golines --base-formatter="goimports" -w -m 120 --ignored-dirs="sdk vendor generated" --ignore-generated .
	gofumpt -w .

lint: ## Run Go linter locally
	golangci-lint --version
	golangci-lint run ./...

lint-docker: ## Run Go linter in Docker
	docker-compose exec publicbox golangci-lint run ./...

format_lint:
	make format
	make lint

lint-fix: ## Run Go linter locally & autofix issues
	golangci-lint run --fix ./...

test: ## Run Go tests locally
	go test ./...

add-vendor: ## Add vendor folder
	go mod tidy
	go mod verify
	go mod vendor

init:
	go mod download
	go install github.com/segmentio/golines@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
