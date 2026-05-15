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


## CI/CD Plan

For MVP, CI means automated validation, not deployment.

Use GitHub Actions.

---

## CI Goals

On every pull request and push to main:

- Run formatting check
- Run `go vet`
- Run tests
- Start PostgreSQL
- Start Redis
- Run migrations
- Build Docker image

---

## Workflow File

Location:

```text
.github/workflows/ci.yml
```

---

## Jobs

### lint

Checks:

```bash
gofmt
go vet ./...
```

### test

Checks:

```bash
go test ./...
```

Needs:

- PostgreSQL service
- Redis service
- Test environment variables

### docker-build

Checks:

```bash
docker build -t taskforge-api-go .
```

---

## CI Triggers

```yaml
on:
  pull_request:
  push:
    branches:
      - main
```

---

## CI Environment Variables

```text
APP_ENV=test
APP_PORT=8080
DATABASE_URL=postgres://taskforge:taskforge@localhost:5432/taskforge_test?sslmode=disable
REDIS_ADDR=localhost:6379
JWT_SECRET=test-secret
JWT_ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h
```

---

## Later CD Plan

After MVP, add deployment:

1. Build Docker image
2. Push image to AWS ECR
3. Deploy to ECS Fargate
4. Run migrations safely
5. Monitor with CloudWatch

Do not add deployment before the local backend and tests are stable.
