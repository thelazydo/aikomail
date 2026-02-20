.PHONY: help setup run build test fmt up down logs

.DEFAULT_GOAL := help

help:
	@echo "Available commands:"
	@echo "  make setup   - Initialize the project (.env, compose.env and go modules)"
	@echo "  make run     - Run the API locally"
	@echo "  make build   - Build a statically linked binary"
	@echo "  make test    - Run all tests with race detector"
	@echo "  make fmt     - Format all Go code"
	@echo "  make up      - Spin up docker containers"
	@echo "  make down    - Tear down docker containers"
	@echo "  make logs    - View docker logs"
	@echo "  make restart - Restart docker containers"

setup:
	@if [ ! -f .env ]; then cp .env.schema .env; echo "✅ Created .env from .env.schema"; else echo "⚡ .env already exists, skipping copy."; fi
	@if [ ! -f compose.env ]; then cp .env.schema compose.env; echo "✅ Created compose.env from .env.schema"; else echo "⚡ compose.env already exists, skipping copy."; fi
	go mod download
	go mod tidy

run: setup
	go run ./cmd/api

build:
	CGO_ENABLED=0 go build -o bin/api ./cmd/api
	@echo "✅ Build complete: ./bin/api"

test:
	go test -v -race ./...

fmt:
	go fmt ./...

up:
	docker compose --env-file compose.env up -d --build

down:
	docker compose --env-file compose.env down

logs:
	docker compose --env-file compose.env logs -f

restart:
	make down && make up
