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
| `make migrate-install` | Install [`migrate`](https://github.com/golang-migrate/migrate) CLI into `GOBIN` |
| `make migrate-up` | Apply SQL in **`migrations/`** (needs `DATABASE_URL`, e.g. from `.env`) |
| `make migrate-down` | Roll back **one** migration (`down 1`) |
| `make migrate-force` | Repair version table (`VERSION=N make migrate-force`) |

After `make up`, run **`make migrate-install`** once, then **`make migrate-up`** (needs `DATABASE_URL`, e.g. from `.env`). The Makefile invokes **`migrate` from `$(go env GOPATH)/bin`** when it is not on your shell `PATH`.

## API

All responses are JSON. Authenticated routes expect `Authorization: Bearer <access token>`.

| Method | Route | Auth | Purpose |
|--------|-------|:----:|---------|
| `POST` | `/auth/register` | — | Create an account |
| `POST` | `/auth/login` | — | Exchange credentials for an access + refresh token pair |
| `POST` | `/auth/refresh` | — | Mint a new access token from a refresh token |
| `POST` | `/auth/logout` | — | Invalidate the current refresh token |
| `GET` | `/me` | ✓ | Current authenticated user |
| `GET` `POST` | `/organizations` | ✓ | List your organizations / create one |
| `GET` `PATCH` `DELETE` | `/organizations/{id}` | ✓ | Read, rename, or delete an organization |
| `GET` `POST` | `/organizations/{id}/members` | ✓ | List members / add a member |
| `DELETE` | `/organizations/{id}/members/{userID}` | ✓ | Remove a member |
| `PATCH` | `/organizations/{id}/members/{userID}/role` | ✓ | Change a member's role |
| `GET` | `/health` | — | Liveness probe |

Membership is role-based: authorization is enforced in the service layer, so a
caller's role within an organization gates every mutating operation on it.

### Implemented vs. planned

`auth`, `organizations`, and `users` are implemented and tested. `projects`,
`tasks`, and `comments` are declared as modules with their boundaries staked out
(`doc.go`) but are not built yet — the domain model is deliberately settled
before the surface area grows.

## Why a modular monolith

The code is organized as one deployable binary split into modules by business
capability (`internal/auth`, `internal/organizations`, …) rather than by
technical layer. Each module owns its own handlers, service, and repository, and
communicates outward through exported functions instead of shared table access.

The point is to keep the seams a service split would need *without* paying
distributed-systems costs the project has no use for — no network hops between
modules, no cross-service transactions, one process to run and debug. Because
the boundaries are already explicit, extracting a module into its own service
later is mechanical rather than archaeological.

## Design decisions

**Standard library routing.** Routes are wired on `net/http.ServeMux` with no
third-party router. Path parameters are parsed explicitly (see
`organizationsSubtree`), which costs a little verbosity and removes a dependency
from the request path.

**Handlers stay thin.** Handlers decode, delegate, and encode. Business rules and
authorization live in each module's service, so they're testable without an HTTP
server — the tests in `internal/auth` and `internal/organizations` exercise
services directly.

**Domain errors, not HTTP errors.** Services return domain-level errors that a
single `writeErr` translation layer maps to status codes. Modules never import
HTTP status semantics into their business logic.

**`server.go` is the sole mux owner.** Route registration is isolated in
`routes.go` so wiring stays in one reviewable place as modules are added.

**Explicit SQL over an ORM.** Queries go through `pgx/v5` in hand-written
repositories, keeping the mapping between Go types and columns visible and the
query plan predictable.

**Migrations are versioned SQL** under `migrations/`, applied by
`golang-migrate` — schema changes are reviewable artifacts, not ORM side effects.

**Stateless access tokens, opaque refresh tokens.** Access tokens are
short-lived JWTs, so the common request path validates a signature instead of
hitting the database. Refresh tokens are opaque random values persisted only as
hashes (`NewRefreshOpaque`), which keeps them revocable and means a database
leak does not hand an attacker usable tokens.

## Docs

Planning and milestones live under [`docs/`](docs/).
