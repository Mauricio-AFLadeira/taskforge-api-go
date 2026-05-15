.PHONY: up down logs run test fmt vet migrate-install migrate-up migrate-down migrate-force

GOBIN := $(shell go env GOPATH)/bin
MIGRATE_CMD := $(shell command -v migrate 2>/dev/null)
ifeq ($(strip $(MIGRATE_CMD)),)
MIGRATE_CMD := $(GOBIN)/migrate
endif

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

migrate-install:
	go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1
	@echo "Installed: $(GOBIN)/migrate (add GOBIN to PATH or use make migrate-up, which calls this path automatically)"

migrate-up:
	@test -x "$(MIGRATE_CMD)" || (echo "migrate not found at $(MIGRATE_CMD). Run: make migrate-install" >&2; exit 127)
	"$(MIGRATE_CMD)" -path migrations -database "$${DATABASE_URL}" up

migrate-down:
	@test -x "$(MIGRATE_CMD)" || (echo "migrate not found at $(MIGRATE_CMD). Run: make migrate-install" >&2; exit 127)
	"$(MIGRATE_CMD)" -path migrations -database "$${DATABASE_URL}" down 1

migrate-force:
	test -n "$${VERSION}" || (echo "Usage: VERSION=42 make migrate-force" >&2; exit 1)
	@test -x "$(MIGRATE_CMD)" || (echo "migrate not found at $(MIGRATE_CMD). Run: make migrate-install" >&2; exit 127)
	"$(MIGRATE_CMD)" -path migrations -database "$${DATABASE_URL}" force $${VERSION}
