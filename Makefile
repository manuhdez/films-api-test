.PHONY: dev test

setup: build start seed
	@echo "Application is setup and running 🚀"

format:
	@gofmt -w .

tidy:
	@go mod tidy

test:
	@go test ./internal/...

# Docker commands
build:
	@docker compose build

start:
	@docker compose up -d

stop:
	@docker compose down

# Database commands
goose:
	@docker compose ps server --format '{{.Name}}' | xargs -I % docker exec % goose $(cmd)

migration:
	@make goose cmd="create $(name) sql"

migrate:
	@make goose cmd="up"

rollback:
	@make goose cmd="down"

db-reset:
	@make goose cmd="reset"

db-status:
	@make goose cmd="status"

seed: db-reset migrate
	@docker compose ps server --format '{{.Name}}' | xargs -I % docker exec % go run cmd/seed/main.go
