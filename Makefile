.PHONY: up down logs run test fmt vet migrate-up migrate-down

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

run:
	go run ./cmd/api

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

migrate-up:
	migrate -path migrations -database "$${DATABASE_URL}" up

migrate-down:
	migrate -path migrations -database "$${DATABASE_URL}" down
