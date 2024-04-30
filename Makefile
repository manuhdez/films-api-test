.PHONY: dev

format:
	@gofmt -w .

tidy:
	@go mod tidy

# Docker commands
build:
	@docker compose build

start:
	@docker compose up -d

stop:
	@docker compose down
