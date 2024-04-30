.PHONY: dev test

format:
	@gofmt -w .

tidy:
	@go mod tidy

test:
	@go test ./...

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
