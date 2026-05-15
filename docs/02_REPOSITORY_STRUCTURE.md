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


## Repository Structure

Create this repository first:

```text
taskforge-api-go/
```

Recommended structure:

```text
taskforge-api-go/
  cmd/
    api/
      main.go

  internal/
    config/
      config.go

    database/
      postgres.go

    redis/
      redis.go

    server/
      server.go
      routes.go

    middleware/
      auth.go
      rate_limit.go
      request_logger.go
      recoverer.go

    auth/
      handler.go
      service.go
      tokens.go
      password.go
      dto.go

    users/
      model.go
      repository.go
      service.go
      handler.go
      dto.go

    organizations/
      model.go
      repository.go
      service.go
      handler.go
      permissions.go
      dto.go

    projects/
      model.go
      repository.go
      service.go
      handler.go
      dto.go

    tasks/
      model.go
      repository.go
      service.go
      handler.go
      filters.go
      dto.go

    comments/
      model.go
      repository.go
      service.go
      handler.go
      dto.go

    shared/
      errors.go
      response.go
      validator.go
      pagination.go
      context.go

  migrations/
    000001_create_users.up.sql
    000001_create_users.down.sql
    000002_create_refresh_tokens.up.sql
    000002_create_refresh_tokens.down.sql
    000003_create_organizations.up.sql
    000003_create_organizations.down.sql
    000004_create_projects.up.sql
    000004_create_projects.down.sql
    000005_create_tasks.up.sql
    000005_create_tasks.down.sql
    000006_create_task_comments.up.sql
    000006_create_task_comments.down.sql

  docs/
    MVP_PLAN.md
    API.md
    DATABASE.md
    ARCHITECTURE.md

  .github/
    workflows/
      ci.yml

  .env.example
  .gitignore
  Dockerfile
  docker-compose.yml
  Makefile
  README.md
  go.mod
  go.sum
```

---

## Why This Structure

`cmd/api` contains the application entry point.

`internal` keeps application code private to this repository.

Each feature has its own module folder.

`shared` contains reusable helpers, but do not put business logic there.

`migrations` stores SQL database changes.

`docs` stores project documentation.

`.github/workflows` stores CI pipelines.

---

## Module Responsibilities

### auth

Responsible for:

- Register
- Login
- JWT generation
- Refresh token generation
- Logout
- Password hashing

### users

Responsible for:

- User persistence
- User lookup
- Current user profile

### organizations

Responsible for:

- Organizations/workspaces
- Members
- Roles
- Permission checks

### projects

Responsible for:

- Project CRUD
- Project visibility

### tasks

Responsible for:

- Task CRUD
- Task status
- Task assignment
- Task filters/search

### comments

Responsible for:

- Task comments
- Comment ownership rules

---

## Branch Strategy

Use simple branches:

```text
main
feature/setup-base
feature/auth
feature/organizations
feature/projects
feature/tasks
feature/comments
feature/redis
feature/tests
feature/ci
```

Merge with pull requests, even if you are working alone. This makes GitHub look more professional.
