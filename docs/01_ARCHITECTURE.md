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


## Architecture Strategy

Use a **modular monolith** in Go for the MVP.

This means the project is deployed as one Go API, but the code is separated by business modules:

- Auth
- Users
- Organizations
- Projects
- Tasks
- Comments
- Shared utilities

This is better than starting with microservices because your Go level is still growing. It keeps the project realistic and easier to finish.

Later, parts of the system can become independent services, for example:

- Java notification service
- Go worker service
- Frontend app
- Infrastructure repository

---

## MVP Architecture

```text
taskforge-api-go
  |
  |-- HTTP REST API
  |-- Auth middleware
  |-- Permission middleware
  |-- PostgreSQL
  |-- Redis
```

Local development:

```text
Client/Postman
    |
    v
Go API :8080
    |
    |-- PostgreSQL :5432
    |-- Redis :6379
```

---

## Future Architecture

After the MVP:

```text
React Frontend
    |
    v
Go API
    |
    |-- PostgreSQL / AWS RDS
    |-- Redis / AWS ElastiCache
    |-- S3 for files
    |-- SQS for async events
    |-- Java Notification Service
```

---

## Suggested Go Libraries

Use simple and professional libraries.

HTTP router options:

- `chi` — recommended, simple and clean
- `gin` — popular and beginner friendly
- `echo` — also good

Recommended for this project:

```text
chi
```

Database options:

- `pgx` — recommended for PostgreSQL
- `sqlc` — excellent later if you want generated type-safe queries
- `gorm` — easier but hides SQL details

Recommended for MVP:

```text
pgxpool + handwritten SQL repositories
```

Other libraries:

- `go-redis` for Redis
- `golang-jwt/jwt` for JWT
- `bcrypt` for password hashing
- `golang-migrate` for migrations
- `testify` for tests
- `zerolog` or Go standard `slog` for logging

Recommended logging:

```text
slog from the Go standard library
```

---

## Internal Layering

Each module should follow a simple pattern:

```text
handler -> service -> repository -> database
```

Example:

```text
auth handler
  receives HTTP request

auth service
  validates business rules

user repository
  reads/writes PostgreSQL
```

Keep business rules out of handlers.

---

## Error Handling Pattern

Use shared application errors:

```text
ErrNotFound
ErrUnauthorized
ErrForbidden
ErrConflict
ErrValidation
ErrInternal
```

Convert them to HTTP responses in one place.

Example:

```text
ErrValidation -> 400
ErrUnauthorized -> 401
ErrForbidden -> 403
ErrNotFound -> 404
ErrConflict -> 409
ErrInternal -> 500
```
