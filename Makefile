.PHONY: help build run test migrate-up migrate-down clean docker-build docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-up    - Start all services with Docker Compose"
	@echo "  make docker-down  - Stop all services"
	@echo "  make clean        - Clean build artifacts"

build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test -v -race ./...

migrate-up:
	migrate -path migrations -database "postgres://localhost/resume_builder?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://localhost/resume_builder?sslmode=disable" down

clean:
	rm -rf bin/
	go clean

docker-build:
	docker build -t resume-builder .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
