# TaskForge API

Backend-first task management API (Go, PostgreSQL, Redis). Modular monolith; see [`docs/00_OVERVIEW.md`](docs/00_OVERVIEW.md).

## Prerequisites

- Go 1.21+
- Docker & Docker Compose

## Quick start

Copy environment file and adjust if needed:

```bash
cp .env.example .env
```

Start databases:

```bash
make up
```

Run the API (loads `.env` automatically when present):

```bash
make run
```

Health check:

```bash
curl -s http://localhost:8080/health
```

Expect: `{"status":"ok"}`

Graceful shutdown: `Ctrl+C`.

## Commands

| Make target   | Purpose                |
|---------------|------------------------|
| `make up`     | Postgres + Redis (detached) |
| `make down`   | Stop containers        |
| `make run`    | Run API locally        |
| `make test`   | Tests                  |

Migrations (`migrate-up` / `migrate-down`) require the [`golang-migrate/migrate`](https://github.com/golang-migrate/migrate) CLI installed; SQL migrations will land in **`migrations/`** in Milestone 2.

## Module path

`go.mod` uses `github.com/mauricio-reportei/taskforge-api-go`. Change it to match **your** canonical import path before open-sourcing.

## Docs

Planning and milestones live under [`docs/`](docs/).
