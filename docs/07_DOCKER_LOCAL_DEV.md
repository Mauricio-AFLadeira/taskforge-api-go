# TaskForge MVP Plan

**Project:** TaskForge  
**Repository:** `taskforge-api-go`  
**Main backend language:** Go  
**Database:** PostgreSQL  
**Cache:** Redis  
**Local environment:** Docker Compose  
**CI/CD:** GitHub Actions  
**Architecture:** Modular monolith first, microservices later  

---


## Docker and Local Development Plan

Use Docker Compose for local infrastructure.

MVP services:

```text
postgres
redis
```

Optional later:

```text
api
```

At first, run the Go API directly on your machine with `go run`, and run PostgreSQL/Redis in Docker. This is easier for development.

---

## docker-compose.yml Plan

Services:

```text
postgres:
  image: postgres
  port: 5432
  database: taskforge
  user: taskforge
  password: taskforge

redis:
  image: redis
  port: 6379
```

---

## Environment Variables

`.env.example`:

```text
APP_ENV=development
APP_PORT=8080

DATABASE_URL=postgres://taskforge:taskforge@localhost:5432/taskforge?sslmode=disable

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=change-me-in-development
JWT_ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h
```

---

## Makefile Commands

Recommended commands:

```makefile
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
```

---

## First Local Success Criteria

Commands:

```bash
make up
make run
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

---

## Dockerfile Later

Add Dockerfile after the API foundation works.

Purpose:

- Build production image
- Let GitHub Actions test Docker builds
- Prepare for future AWS/ECS deployment

---

## Development Rules

- Do not install PostgreSQL locally if Docker already provides it.
- Do not commit `.env`.
- Commit `.env.example`.
- Keep local ports standard unless conflict happens.
- Use Docker Compose for repeatable setup.
