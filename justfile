# default command (help)
default: help

# help
@help:
    just --list

# Initialize the project (.env, compose.env and go modules)
setup:
    @if [ ! -f .env ]; then cp .env.schema .env; echo "✅ Created .env from .env.schema"; else echo "⚡ .env already exists, skipping copy."; fi
    @if [ ! -f compose.env ]; then cp .env.schema compose.env; echo "✅ Created compose.env from .env.schema"; else echo "⚡ comopose.env already exists, skipping copy."; fi
    go mod download
    go mod tidy

# Run the API locally
run: setup
    @just --dotenv-filename .env run_internal

[private]
run_internal:
    go run ./cmd/api

# Build a statiscally linked binary
build:
    CGO_ENABLED=0 go build -o bin/api ./cmd/api
    @echo "✅ Build complete: ./bin/api"

# Run all tests with race detector
test:
    go test -v -race ./...

# Format all Go code
fmt:
    go fmt ./...

# Spin up docker containers
up:
  docker compose --env-file compose.env up -d --build

# Tear down docker containers
down:
  docker compose --env-file compose.env down

# View docker logs
logs:
  docker compose --env-file compose.env logs -f

# Restart docker containers
restart:
  just down && just up
