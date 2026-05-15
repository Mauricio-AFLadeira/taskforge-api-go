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


## 1. MVP Goal

Build a backend-first task management SaaS API that looks like a simplified Jira/Trello backend.

The MVP must support:

- User registration and login
- JWT access tokens
- Refresh tokens
- Logout
- Organizations/workspaces
- Organization members
- Role-based access control
- Projects
- Tasks
- Task assignment
- Task status updates
- Task comments
- Basic task search and filters
- PostgreSQL persistence
- Redis rate limiting
- Redis permission cache
- Docker local development
- Automated tests
- GitHub Actions CI

The MVP should be strong enough to show in a backend/cloud portfolio, but small enough to actually finish.

---

## 2. What Is Not in the MVP

Avoid these until the backend MVP is finished:

- Frontend
- AWS deployment
- Java notification service
- File uploads
- Task labels
- WebSockets
- Kubernetes
- Terraform
- Email notifications
- Payment system

These are excellent post-MVP features, but building them too early will slow the project down.

---

## 3. Main System Description

TaskForge is a SaaS API where users can create organizations, invite/manage members, create projects, create tasks, assign tasks to members, update task status, and discuss tasks with comments.

Example flow:

1. User registers.
2. User logs in.
3. User creates an organization.
4. User becomes organization owner.
5. Owner adds members.
6. Owner/admin creates projects.
7. Members create and update tasks.
8. Members comment on tasks.
9. API protects access with JWT and organization permissions.
10. Redis helps with rate limiting and permission caching.

---

## 4. Recommended Development Order

Build in this order:

1. Repository setup
2. Docker Compose with PostgreSQL and Redis
3. Go HTTP server
4. Health endpoint
5. Config loading
6. PostgreSQL connection
7. Redis connection
8. Database migrations
9. User table and model
10. Register endpoint
11. Login endpoint
12. JWT middleware
13. Refresh token system
14. Logout endpoint
15. Organizations
16. Organization members
17. Permission checks
18. Projects
19. Tasks
20. Task comments
21. Task filters/search
22. Redis rate limiting
23. Redis permission cache
24. Tests
25. GitHub Actions
26. Documentation

After **step 1** (repository setup through **running foundation**: Compose, minimal server, `/health`, DB + Redis wired), prefer **Cursor Multitask** in the batches documented under **Milestone Plan → Cursor Multitask execution structure** (`B1`–`B9`). Keep **foundation** single-threaded to avoid conflicting edits to `main`, `server`, and compose files.
---

## 5. Final MVP Definition

The MVP is complete when you can say:

> TaskForge is a Go backend SaaS API for task management. It supports user authentication with JWT and refresh tokens, organizations, role-based access control, projects, tasks, task assignment, comments, task filtering, PostgreSQL persistence, Redis rate limiting/cache, Docker local development, tests, and GitHub Actions CI.

---

## 6. Final MVP Checklist

- [ ] API runs locally
- [ ] PostgreSQL runs with Docker
- [ ] Redis runs with Docker
- [ ] Migrations work
- [ ] Users can register
- [ ] Users can login
- [ ] Users can refresh tokens
- [ ] Users can logout
- [ ] Protected routes require JWT
- [ ] Users can create organizations
- [ ] Organizations have members
- [ ] Members have roles
- [ ] Users can create projects
- [ ] Users can create tasks
- [ ] Users can assign tasks
- [ ] Users can update task status
- [ ] Users can comment on tasks
- [ ] Users can filter/search tasks
- [ ] Redis rate limiting works
- [ ] Redis permission cache works
- [ ] Tests exist
- [ ] GitHub Actions works
- [ ] README is complete


---

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


---

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

After this layout exists and **Milestone 1 — Backend Foundation** is complete (Compose, server, `/health`, DB/Redis wired), drive the remaining MVP work through **parallel Cursor Multitask batches** documented under [Cursor Multitask execution structure](#cursor-multitask-execution-structure) (Phase B batches **B1–B9**).

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


---

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


## Database Plan

Use PostgreSQL.

Core MVP tables:

- users
- refresh_tokens
- organizations
- organization_members
- projects
- tasks
- task_comments

---

## Entity Relationships

```text
users 1---N refresh_tokens
users 1---N organizations as owner
users N---N organizations through organization_members
organizations 1---N projects
projects 1---N tasks
tasks 1---N task_comments
users 1---N task_comments
users 1---N tasks as assignee
users 1---N tasks as creator
```

---

## Table: users

Purpose: stores system users.

Columns:

```text
id UUID PRIMARY KEY
name TEXT NOT NULL
email TEXT NOT NULL UNIQUE
password_hash TEXT NOT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

Indexes:

```sql
CREATE UNIQUE INDEX idx_users_email ON users(email);
```

---

## Table: refresh_tokens

Purpose: stores hashed refresh tokens for secure session renewal.

Columns:

```text
id UUID PRIMARY KEY
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
token_hash TEXT NOT NULL UNIQUE
expires_at TIMESTAMPTZ NOT NULL
revoked_at TIMESTAMPTZ NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

Indexes:

```sql
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE UNIQUE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
```

---

## Table: organizations

Purpose: workspaces/companies/teams.

Columns:

```text
id UUID PRIMARY KEY
name TEXT NOT NULL
owner_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

Indexes:

```sql
CREATE INDEX idx_organizations_owner_id ON organizations(owner_id);
```

---

## Table: organization_members

Purpose: connects users to organizations with roles.

Columns:

```text
id UUID PRIMARY KEY
organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
role TEXT NOT NULL CHECK (role IN ('owner', 'admin', 'member'))
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

Constraints:

```sql
UNIQUE (organization_id, user_id)
```

Indexes:

```sql
CREATE INDEX idx_organization_members_user_id ON organization_members(user_id);
CREATE INDEX idx_organization_members_organization_id ON organization_members(organization_id);
```

---

## Table: projects

Purpose: projects inside organizations.

Columns:

```text
id UUID PRIMARY KEY
organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE
name TEXT NOT NULL
description TEXT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

Indexes:

```sql
CREATE INDEX idx_projects_organization_id ON projects(organization_id);
```

---

## Table: tasks

Purpose: tasks inside projects.

Columns:

```text
id UUID PRIMARY KEY
project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE
title TEXT NOT NULL
description TEXT NULL
status TEXT NOT NULL CHECK (status IN ('todo', 'in_progress', 'done', 'cancelled')) DEFAULT 'todo'
priority TEXT NOT NULL CHECK (priority IN ('low', 'medium', 'high', 'urgent')) DEFAULT 'medium'
assignee_id UUID NULL REFERENCES users(id) ON DELETE SET NULL
created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT
due_date TIMESTAMPTZ NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

Indexes:

```sql
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
```

---

## Table: task_comments

Purpose: comments on tasks.

Columns:

```text
id UUID PRIMARY KEY
task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
content TEXT NOT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

Indexes:

```sql
CREATE INDEX idx_task_comments_task_id ON task_comments(task_id);
CREATE INDEX idx_task_comments_user_id ON task_comments(user_id);
```

---

## Migration Strategy

Use `golang-migrate`.

Commands:

```bash
make migrate-up
make migrate-down
make migrate-force version=1
```

Rules:

- Every schema change must be a migration.
- Do not manually edit the database without a migration.
- Keep up/down migrations paired.
- Use UUIDs for public IDs.

---

## Suggested Migration Order

1. Create users
2. Create refresh_tokens
3. Create organizations and organization_members
4. Create projects
5. Create tasks
6. Create task_comments

---

## Seed Data Later

Optional later:

- Demo user
- Demo organization
- Demo project
- Demo tasks

Do not seed until the core API works.


---

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


## API Endpoint Plan

Base URL locally:

```text
http://localhost:8080
```

---

## Health

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

---

## Auth Endpoints

### Register

```http
POST /auth/register
```

Request:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "strong-password"
}
```

Response:

```json
{
  "user": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com"
  },
  "access_token": "jwt",
  "refresh_token": "random-token"
}
```

### Login

```http
POST /auth/login
```

Request:

```json
{
  "email": "john@example.com",
  "password": "strong-password"
}
```

### Refresh Token

```http
POST /auth/refresh
```

Request:

```json
{
  "refresh_token": "random-token"
}
```

### Logout

```http
POST /auth/logout
```

Request:

```json
{
  "refresh_token": "random-token"
}
```

### Current User

```http
GET /me
```

Requires JWT.

---

## Organization Endpoints

### Create Organization

```http
POST /organizations
```

Request:

```json
{
  "name": "Acme Inc"
}
```

### List My Organizations

```http
GET /organizations
```

### Get Organization

```http
GET /organizations/{organizationId}
```

### Update Organization

```http
PATCH /organizations/{organizationId}
```

Request:

```json
{
  "name": "New Name"
}
```

### Delete Organization

```http
DELETE /organizations/{organizationId}
```

---

## Organization Member Endpoints

### Add Member

```http
POST /organizations/{organizationId}/members
```

Request:

```json
{
  "email": "member@example.com",
  "role": "member"
}
```

### List Members

```http
GET /organizations/{organizationId}/members
```

### Change Member Role

```http
PATCH /organizations/{organizationId}/members/{userId}/role
```

Request:

```json
{
  "role": "admin"
}
```

### Remove Member

```http
DELETE /organizations/{organizationId}/members/{userId}
```

---

## Project Endpoints

### Create Project

```http
POST /organizations/{organizationId}/projects
```

Request:

```json
{
  "name": "Backend API",
  "description": "Build the TaskForge backend"
}
```

### List Projects

```http
GET /organizations/{organizationId}/projects
```

### Get Project

```http
GET /projects/{projectId}
```

### Update Project

```http
PATCH /projects/{projectId}
```

### Delete Project

```http
DELETE /projects/{projectId}
```

---

## Task Endpoints

### Create Task

```http
POST /projects/{projectId}/tasks
```

Request:

```json
{
  "title": "Create auth module",
  "description": "Implement register, login and JWT middleware",
  "priority": "high",
  "assignee_id": "uuid-or-null",
  "due_date": "2026-06-01T12:00:00Z"
}
```

### List Tasks

```http
GET /projects/{projectId}/tasks
```

Filters:

```http
GET /projects/{projectId}/tasks?status=todo
GET /projects/{projectId}/tasks?assignee_id=uuid
GET /projects/{projectId}/tasks?priority=high
GET /projects/{projectId}/tasks?search=auth
```

### Get Task

```http
GET /tasks/{taskId}
```

### Update Task

```http
PATCH /tasks/{taskId}
```

### Delete Task

```http
DELETE /tasks/{taskId}
```

### Change Task Status

```http
PATCH /tasks/{taskId}/status
```

Request:

```json
{
  "status": "in_progress"
}
```

### Assign Task

```http
PATCH /tasks/{taskId}/assignee
```

Request:

```json
{
  "assignee_id": "uuid"
}
```

---

## Comment Endpoints

### Add Comment

```http
POST /tasks/{taskId}/comments
```

Request:

```json
{
  "content": "I started working on this."
}
```

### List Comments

```http
GET /tasks/{taskId}/comments
```

### Update Comment

```http
PATCH /comments/{commentId}
```

### Delete Comment

```http
DELETE /comments/{commentId}
```

---

## Standard Error Response

Use one consistent format:

```json
{
  "error": {
    "code": "validation_error",
    "message": "Email is required"
  }
}
```

Common HTTP statuses:

```text
400 Bad Request
401 Unauthorized
403 Forbidden
404 Not Found
409 Conflict
422 Unprocessable Entity
429 Too Many Requests
500 Internal Server Error
```

---

## Pagination Later

For the first version, simple lists are okay.

Then add:

```http
GET /projects/{projectId}/tasks?page=1&page_size=20
```


---

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


## Authentication Plan

Use email/password authentication with JWT access tokens and refresh tokens.

---

## Auth Flow

### Register

1. Validate name, email, password.
2. Check if email already exists.
3. Hash password with bcrypt.
4. Create user.
5. Generate access token.
6. Generate refresh token.
7. Hash refresh token.
8. Store refresh token hash in PostgreSQL.
9. Return tokens and user info.

### Login

1. Validate email and password.
2. Find user by email.
3. Compare password with password hash.
4. Generate access token.
5. Generate refresh token.
6. Store refresh token hash.
7. Return tokens.

### Refresh

1. Receive refresh token.
2. Hash received token.
3. Find token hash in database.
4. Check token is not expired.
5. Check token is not revoked.
6. Generate new access token.
7. Optional: rotate refresh token.

For MVP, you can rotate refresh tokens or keep it simpler. Rotation is better.

### Logout

1. Receive refresh token.
2. Hash token.
3. Mark token as revoked.

---

## Token Durations

Recommended:

```text
Access token: 15 minutes
Refresh token: 7 days or 30 days
```

For MVP:

```text
Access token: 15 minutes
Refresh token: 7 days
```

---

## Password Rules

Minimum rules:

```text
At least 8 characters
Cannot be empty
```

Later, improve with:

```text
Must include uppercase
Must include lowercase
Must include number
Must include special character
```

---

## JWT Claims

Access token should include:

```json
{
  "sub": "user-id",
  "email": "user@example.com",
  "exp": 1234567890,
  "iat": 1234567890
}
```

Keep organization permissions out of JWT for MVP. Read permissions from PostgreSQL/Redis so role changes apply quickly.

---

# Role-Based Access Control

Roles:

```text
owner
admin
member
```

---

## Permission Matrix

| Action | Owner | Admin | Member |
|---|---:|---:|---:|
| View organization | Yes | Yes | Yes |
| Update organization | Yes | No | No |
| Delete organization | Yes | No | No |
| Add member | Yes | Yes | No |
| Remove member | Yes | Yes | No |
| Remove owner | No | No | No |
| Change member role | Yes | Yes | No |
| Change owner role | No | No | No |
| Create project | Yes | Yes | No |
| View project | Yes | Yes | Yes |
| Update project | Yes | Yes | No |
| Delete project | Yes | Yes | No |
| Create task | Yes | Yes | Yes |
| View task | Yes | Yes | Yes |
| Update task | Yes | Yes | Yes |
| Delete task | Yes | Yes | No |
| Assign task | Yes | Yes | Yes |
| Comment on task | Yes | Yes | Yes |
| Delete own comment | Yes | Yes | Yes |
| Delete any comment | Yes | Yes | No |

---

## Permission Helper

Create functions like:

```text
RequireOrganizationRole(userID, organizationID, allowedRoles)
CanAccessProject(userID, projectID)
CanAccessTask(userID, taskID)
```

Do not duplicate permission SQL in every handler.

---

## Permission Cache With Redis

Since permission checks happen often, cache this:

```text
org_member:{organization_id}:{user_id}
```

Value:

```json
{
  "role": "admin"
}
```

TTL:

```text
5 to 15 minutes
```

Invalidate when:

- Member is removed
- Role is changed
- Organization is deleted


---

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


## Redis Plan

Use Redis only for real backend features. Do not add Redis just to say the project uses Redis.

MVP Redis features:

1. Rate limiting
2. Permission cache

---

## Feature 1: Rate Limiting

Protect auth endpoints:

```http
POST /auth/login
POST /auth/register
```

Rules:

```text
Max 5 login attempts per minute per IP
Max 5 login attempts per minute per email
```

Keys:

```text
rate_limit:login:ip:{ip}
rate_limit:login:email:{email}
rate_limit:register:ip:{ip}
```

Behavior:

1. User sends request.
2. Middleware increments Redis key.
3. If key does not exist, set expiration.
4. If count exceeds limit, return 429.

Error response:

```json
{
  "error": {
    "code": "rate_limited",
    "message": "Too many requests. Try again later."
  }
}
```

---

## Feature 2: Permission Cache

Permission checks are common. Cache organization membership role.

Key:

```text
org_member:{organization_id}:{user_id}
```

Value:

```json
{
  "role": "member"
}
```

TTL:

```text
10 minutes
```

Flow:

1. Permission check starts.
2. Check Redis for member role.
3. If found, use cached role.
4. If not found, query PostgreSQL.
5. Store result in Redis.
6. Continue permission check.

Invalidate cache when:

- Member role changes
- Member is removed
- Organization is deleted

---

## Redis Connection Config

Environment variables:

```text
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

For Docker Compose, API uses:

```text
REDIS_ADDR=redis:6379
```

For local machine without API container, use:

```text
REDIS_ADDR=localhost:6379
```

---

## Later Redis Features

After MVP:

- Background job queue
- WebSocket pub/sub
- Notification fanout
- Temporary invitation tokens
- Cache task lists
- Cache project summaries

Do not add these before the MVP is stable.


---

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


---

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


## Testing Plan

Goal: write meaningful tests that prove the backend works.

Do not chase 100% coverage. Focus on business-critical behavior.

---

## Test Types

### Unit Tests

Test service logic without real PostgreSQL/Redis.

Examples:

- Password validation
- Token generation
- Permission rules
- Task status transitions

### Integration Tests

Test repositories with real PostgreSQL.

Examples:

- Create user
- Find user by email
- Create organization
- Add member
- Create project
- Create task

### Handler/API Tests

Test HTTP endpoints.

Examples:

- Register endpoint
- Login endpoint
- Protected route with JWT
- Create organization
- Create task

---

## Minimum Test Checklist

Auth:

- [ ] User can register
- [ ] User cannot register with duplicate email
- [ ] User cannot register with invalid email
- [ ] User can login
- [ ] Invalid password fails
- [ ] Refresh token returns new access token
- [ ] Revoked refresh token fails
- [ ] Protected route fails without token
- [ ] Protected route works with token

Organizations:

- [ ] User can create organization
- [ ] Creator becomes owner
- [ ] Owner can add member
- [ ] Member cannot add member
- [ ] Admin can add member
- [ ] Owner cannot be removed
- [ ] Non-member cannot access organization

Projects:

- [ ] Owner can create project
- [ ] Admin can create project
- [ ] Member cannot create project
- [ ] Member can view project
- [ ] Non-member cannot view project

Tasks:

- [ ] Member can create task
- [ ] Member can update task
- [ ] Owner/admin can delete task
- [ ] Member cannot delete task
- [ ] Task can be assigned to organization member
- [ ] Task cannot be assigned to non-member
- [ ] Filter by status works
- [ ] Search by title works

Comments:

- [ ] Member can comment
- [ ] User can update own comment
- [ ] User cannot update another user's comment
- [ ] Owner/admin can delete any comment

Redis:

- [ ] Rate limiter blocks after limit
- [ ] Permission role is cached
- [ ] Permission cache is invalidated after role change

---

## Test Commands

```bash
make test
```

CI command:

```bash
go test ./...
```

With race detector later:

```bash
go test -race ./...
```

---

## Testing Strategy by Phase

Do not wait until the end to test everything.

Phase order:

1. Auth tests after auth is implemented
2. Organization permission tests after RBAC
3. Project tests after project module
4. Task tests after task module
5. Comment tests after comments
6. Redis tests after Redis features

---

## Test Database

Use a separate test database:

```text
taskforge_test
```

Or create temporary schemas during tests.

For a first MVP, using Docker PostgreSQL and a test database is enough.


---

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


---

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


## Milestone Plan

Use this as your project board in GitHub Projects, Trello, Notion, or VS Code.

---

## Cursor Multitask execution structure

Use **one agent** until the repository skeleton and wiring exist. After that, use **Multitask** (parallel agents) only when work is on **different files/modules** and **merge order** is explicit.

### Phase A — Sequential only (do not multitask)

Complete **Milestone 1 — Backend Foundation** in a single thread:

- Repository, `go.mod`, folder layout, Docker Compose, config, HTTP server, `/health`, PostgreSQL + Redis connections, graceful shutdown.

**Why:** Everything touches `main`, `server`, `routes`, shared bootstrapping, and `docker-compose.yml`. Parallel edits here cause constant merge conflicts.

### Phase B — Multitask batches (after Phase A is green)

Run parallel agents **within** a batch. **Merge and verify** (`make up`, `make run`, `curl /health`, `make migrate-up` when migrations exist) before starting the next batch.

| Batch | Multitask tracks (parallel) | Merge note |
| ----- | ----------------------------- | ---------- |
| **B1 — Schema** | Track 1: `users` + `refresh_tokens` migrations and indexes. Track 2: `organizations` + `organization_members`. Track 3: `projects`. Track 4: `tasks`. Track 5: `task_comments`. Track 6: Makefile / migrate tooling if not already done. | One PR or sequential commits; all migrations must apply cleanly together. |
| **B2 — Auth vertical slice** | Track 1: `internal/users` (model, repo, DTOs). Track 2: `internal/auth` (service, JWT, refresh, password). Track 3: HTTP handlers + routes wiring for register/login/refresh/logout. Track 4: `internal/middleware/auth.go` + `/me`. | Integrate on `routes`/`server` in a final pass or dedicated “integration” sub-agent after tracks land. |
| **B3 — Auth tests & hardening** | Track 1: repository tests. Track 2: service/JWT tests. Track 3: HTTP/e2e-style API tests. | Depends on B2. |
| **B4 — Organizations RBAC** | Track 1: org repository + service. Track 2: handlers + DTOs. Track 3: `permissions.go` + permission helper. Track 4: permission tests. | Single owner for route registration to avoid duplicate route blocks. |
| **B5 — Projects** | Track 1: repository. Track 2: service + permissions hooks. Track 3: handlers + tests. | Depends on B4 for org membership checks. |
| **B6 — Tasks** | Track 1: repository + filters. Track 2: service + assignment/status rules. Track 3: handlers. Track 4: task tests. | Depends on B5 (project scope). |
| **B7 — Comments** | Track 1: repository + service. Track 2: handlers + comment permissions. Track 3: tests. | Depends on B6 (task IDs). |
| **B8 — Redis** | Track 1: rate limit middleware + login/register protection. Track 2: permission cache + invalidation. Track 3: Redis-focused tests. | Coordinate on middleware order in `server` / `routes`. |
| **B9 — CI & docs** | Track 1: `.github/workflows/ci.yml` + Dockerfile polish. Track 2: README runbook. Track 3: `docs/API.md` / `DATABASE.md` / architecture notes. | Low conflict if code is stable; run CI after merging doc-only + workflow changes. |

### Rules when using Multitask

1. **One writer per hot file** — e.g. only one agent edits `routes.go` or `main.go` per batch; others deliver packages and a small integration agent wires them.
2. **Batch boundary = integration point** — compile, run migrations, run tests before the next batch.
3. **Do not** parallelize two milestones that both change the same migration version or the same shared error/response types without agreeing on the types first (short sequential “contract” step if needed).

This maps **Milestones 2–9** to **B1–B9**: each milestone aligns with the batch of the same number for planning, but **B2–B9** are where Multitask is intended; **Milestone 1** stays Phase A.

---

# Milestone 1 — Backend Foundation

**Cursor Multitask:** do **not** use parallel agents here ([Phase A — Sequential only](#phase-a--sequential-only-do-not-multitask)); start Multitask at **Milestone 2** / batch **B1**.

Goal: first working API with infrastructure.

Tasks:

- [ ] Create repository `taskforge-api-go`
- [ ] Initialize Go module
- [ ] Create folder structure
- [ ] Add `.gitignore`
- [ ] Add `.env.example`
- [ ] Add `README.md`
- [ ] Add `Makefile`
- [ ] Add `docker-compose.yml`
- [ ] Add PostgreSQL service
- [ ] Add Redis service
- [ ] Add config loader
- [ ] Add HTTP server
- [ ] Add `/health` endpoint
- [ ] Connect to PostgreSQL
- [ ] Connect to Redis
- [ ] Add graceful shutdown

Done when:

```bash
make up
make run
curl http://localhost:8080/health
```

returns:

```json
{
  "status": "ok"
}
```

---

# Milestone 2 — Database Migrations

**Cursor Multitask:** batch **B1** ([Cursor Multitask execution structure](#cursor-multitask-execution-structure)).

Goal: schema is managed professionally.

Tasks:

- [ ] Add `golang-migrate`
- [ ] Create users migration
- [ ] Create refresh_tokens migration
- [ ] Create organizations migration
- [ ] Create organization_members migration
- [ ] Create projects migration
- [ ] Create tasks migration
- [ ] Create task_comments migration
- [ ] Add indexes
- [ ] Add constraints
- [ ] Add migration Makefile commands

Done when:

```bash
make migrate-up
```

creates all tables.

---

# Milestone 3 — Authentication

**Cursor Multitask:** batches **B2** (implementation) then **B3** (tests) ([Cursor Multitask execution structure](#cursor-multitask-execution-structure)).

Goal: users can create accounts and authenticate.

Tasks:

- [ ] Create user repository
- [ ] Create auth service
- [ ] Create password hashing helper
- [ ] Create JWT helper
- [ ] Register endpoint
- [ ] Login endpoint
- [ ] Refresh endpoint
- [ ] Logout endpoint
- [ ] Auth middleware
- [ ] Current user endpoint `/me`
- [ ] Auth tests

Done when:

- [ ] User can register
- [ ] User can login
- [ ] Protected endpoint works with token
- [ ] Protected endpoint fails without token

---

# Milestone 4 — Organizations and Members

**Cursor Multitask:** batch **B4** ([Cursor Multitask execution structure](#cursor-multitask-execution-structure)).

Goal: workspace system with RBAC.

Tasks:

- [ ] Create organization repository
- [ ] Create organization service
- [ ] Create organization handlers
- [ ] Create organization endpoint
- [ ] List organizations endpoint
- [ ] Get organization endpoint
- [ ] Update organization endpoint
- [ ] Delete organization endpoint
- [ ] Add member endpoint
- [ ] List members endpoint
- [ ] Change role endpoint
- [ ] Remove member endpoint
- [ ] Permission helper
- [ ] Permission tests

Done when:

- [ ] User can create organization
- [ ] Creator becomes owner
- [ ] Owner can add members
- [ ] Member cannot manage members

---

# Milestone 5 — Projects

**Cursor Multitask:** batch **B5** ([Cursor Multitask execution structure](#cursor-multitask-execution-structure)).

Goal: organizations can contain projects.

Tasks:

- [ ] Create project repository
- [ ] Create project service
- [ ] Create project handlers
- [ ] Create project endpoint
- [ ] List projects endpoint
- [ ] Get project endpoint
- [ ] Update project endpoint
- [ ] Delete project endpoint
- [ ] Project permission tests

Done when:

- [ ] Owner/admin can create projects
- [ ] Member can view projects
- [ ] Non-member cannot access projects

---

# Milestone 6 — Tasks

**Cursor Multitask:** batch **B6** ([Cursor Multitask execution structure](#cursor-multitask-execution-structure)).

Goal: project task management works.

Tasks:

- [ ] Create task repository
- [ ] Create task service
- [ ] Create task handlers
- [ ] Create task endpoint
- [ ] List tasks endpoint
- [ ] Get task endpoint
- [ ] Update task endpoint
- [ ] Delete task endpoint
- [ ] Change task status endpoint
- [ ] Assign task endpoint
- [ ] Add filters
- [ ] Add basic search
- [ ] Task tests

Done when:

- [ ] Members can create tasks
- [ ] Tasks can be assigned
- [ ] Task status can change
- [ ] Task filters work

---

# Milestone 7 — Comments

**Cursor Multitask:** batch **B7** ([Cursor Multitask execution structure](#cursor-multitask-execution-structure)).

Goal: users can discuss tasks.

Tasks:

- [ ] Create comment repository
- [ ] Create comment service
- [ ] Create comment handlers
- [ ] Add comment endpoint
- [ ] List comments endpoint
- [ ] Update comment endpoint
- [ ] Delete comment endpoint
- [ ] Comment permission tests

Done when:

- [ ] Members can comment on tasks
- [ ] Users can edit own comments
- [ ] Owner/admin can delete any comment

---

# Milestone 8 — Redis Features

**Cursor Multitask:** batch **B8** ([Cursor Multitask execution structure](#cursor-multitask-execution-structure)).

Goal: Redis is used in realistic backend features.

Tasks:

- [ ] Add rate limiter middleware
- [ ] Protect login endpoint
- [ ] Protect register endpoint
- [ ] Add permission cache
- [ ] Cache organization roles
- [ ] Invalidate cache on role changes
- [ ] Redis tests

Done when:

- [ ] Too many login attempts return 429
- [ ] Permission checks use Redis cache

---

# Milestone 9 — CI and Documentation

**Cursor Multitask:** batch **B9** ([Cursor Multitask execution structure](#cursor-multitask-execution-structure)).

Goal: project looks professional on GitHub.

Tasks:

- [ ] Add GitHub Actions workflow
- [ ] Run gofmt check
- [ ] Run go vet
- [ ] Run tests
- [ ] Build Docker image
- [ ] Improve README
- [ ] Add API documentation
- [ ] Add database documentation
- [ ] Add architecture documentation

Done when:

- [ ] Pull request runs CI successfully
- [ ] README explains how to run the project
- [ ] Docs are clear enough for a recruiter/developer


---

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


## Post-MVP Roadmap

Only start this after the MVP is finished.

---

# Phase 1 — Frontend

Repository:

```text
taskforge-web
```

Tech options:

```text
React + TypeScript
Next.js
Vite + React
```

Suggested:

```text
React + TypeScript + Vite
```

Features:

- Login/register pages
- Organization dashboard
- Project list
- Task board
- Task details
- Comments

---

# Phase 2 — Java Notification Service

Repository:

```text
taskforge-notifications-java
```

Tech:

```text
Java Spring Boot
```

Responsibilities:

- Receive notification events
- Send email notifications
- Store notification history
- Handle task assignment notifications
- Handle due date reminders

Possible communication:

- Start simple with HTTP calls from Go API
- Later use AWS SQS or RabbitMQ

---

# Phase 3 — AWS Deployment

Repository:

```text
taskforge-infra
```

AWS services:

- ECS Fargate
- ECR
- RDS PostgreSQL
- ElastiCache Redis
- S3
- Secrets Manager
- CloudWatch
- Application Load Balancer

Deploy order:

1. Dockerize API
2. Push image to ECR
3. Deploy API to ECS
4. Move PostgreSQL to RDS
5. Move Redis to ElastiCache
6. Add domain and HTTPS

---

# Phase 4 — File Uploads

Use S3.

Features:

- Attach files to tasks
- Store metadata in PostgreSQL
- Upload directly to S3 with presigned URLs

Tables:

```text
task_attachments
```

---

# Phase 5 — Task Labels

Features:

- Create labels per organization or project
- Attach labels to tasks
- Filter tasks by label

Tables:

```text
labels
task_labels
```

---

# Phase 6 — Observability

Add:

- Structured logs
- Request IDs
- Metrics
- Health checks
- Readiness checks
- Error tracking

---

# Phase 7 — Advanced Search

Improve from simple `ILIKE` to:

- PostgreSQL full-text search
- Ranking
- Multiple filters
- Pagination
- Sorting

---

# Phase 8 — Microservice Split

Only split services when the modular monolith is stable.

Possible split:

```text
taskforge-api-go
  auth/users/organizations/projects/tasks

taskforge-notifications-java
  notifications

taskforge-worker-go
  background jobs

taskforge-web
  frontend

taskforge-infra
  AWS/Terraform
```


---

# README Template

# TaskForge

TaskForge is a backend-first task management SaaS API built with Go, PostgreSQL, Redis, Docker, and GitHub Actions.

It is designed as a professional portfolio project to demonstrate backend engineering, authentication, authorization, database design, Redis usage, testing, Docker, and CI/CD.

---

## Tech Stack

- Go
- PostgreSQL
- Redis
- Docker Compose
- GitHub Actions
- JWT authentication
- Refresh tokens

---

## Features

- User registration and login
- JWT access tokens
- Refresh tokens
- Logout
- Organizations/workspaces
- Organization members
- Role-based access control
- Projects
- Tasks
- Task assignment
- Task status updates
- Task comments
- Basic task search and filters
- Redis rate limiting
- Redis permission cache
- Tests
- CI pipeline

---

## Architecture

```text
Client/Postman
    |
    v
Go API
    |
    |-- PostgreSQL
    |-- Redis
```

The project starts as a modular monolith. Each domain is separated into its own package.

---

## Local Setup

### Requirements

- Go
- Docker
- Docker Compose
- Git

### Start Infrastructure

```bash
make up
```

### Run API

```bash
make run
```

### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

---

## Environment Variables

Copy `.env.example` to `.env`.

```bash
cp .env.example .env
```

Variables:

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

## Database Migrations

Run migrations:

```bash
make migrate-up
```

Rollback migrations:

```bash
make migrate-down
```

---

## Tests

```bash
make test
```

---

## API Endpoints

### Auth

```http
POST /auth/register
POST /auth/login
POST /auth/refresh
POST /auth/logout
GET  /me
```

### Organizations

```http
POST   /organizations
GET    /organizations
GET    /organizations/{organizationId}
PATCH  /organizations/{organizationId}
DELETE /organizations/{organizationId}
```

### Members

```http
POST   /organizations/{organizationId}/members
GET    /organizations/{organizationId}/members
PATCH  /organizations/{organizationId}/members/{userId}/role
DELETE /organizations/{organizationId}/members/{userId}
```

### Projects

```http
POST   /organizations/{organizationId}/projects
GET    /organizations/{organizationId}/projects
GET    /projects/{projectId}
PATCH  /projects/{projectId}
DELETE /projects/{projectId}
```

### Tasks

```http
POST   /projects/{projectId}/tasks
GET    /projects/{projectId}/tasks
GET    /tasks/{taskId}
PATCH  /tasks/{taskId}
DELETE /tasks/{taskId}
PATCH  /tasks/{taskId}/status
PATCH  /tasks/{taskId}/assignee
```

### Comments

```http
POST   /tasks/{taskId}/comments
GET    /tasks/{taskId}/comments
PATCH  /comments/{commentId}
DELETE /comments/{commentId}
```

---

## Roadmap

- Frontend with React
- Java Spring Boot notification service
- AWS deployment
- S3 file uploads
- Task labels
- WebSockets
- Terraform


---

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


## VS Code Execution Checklist

Use this file while coding.

---

# Day 1 / Session 1 — Create Base Project

Commands:

```bash
mkdir taskforge-api-go
cd taskforge-api-go
go mod init github.com/YOUR_USERNAME/taskforge-api-go
mkdir -p cmd/api internal/{config,database,redis,server,middleware,auth,users,organizations,projects,tasks,comments,shared} migrations docs .github/workflows
```

Create files:

- [ ] `cmd/api/main.go`
- [ ] `.env.example`
- [ ] `.gitignore`
- [ ] `docker-compose.yml`
- [ ] `Makefile`
- [ ] `README.md`

---

# Session 2 — Docker Compose

Tasks:

- [ ] Add PostgreSQL service
- [ ] Add Redis service
- [ ] Run `make up`
- [ ] Confirm containers are running

Commands:

```bash
docker compose ps
```

---

# Session 3 — HTTP Server

Tasks:

- [ ] Create config loader
- [ ] Create server
- [ ] Add `/health`
- [ ] Run API

Commands:

```bash
make run
curl http://localhost:8080/health
```

---

# Session 4 — PostgreSQL and Redis Connections

Tasks:

- [ ] Add PostgreSQL connection with `pgxpool`
- [ ] Ping PostgreSQL on startup
- [ ] Add Redis connection with `go-redis`
- [ ] Ping Redis on startup
- [ ] Return health status

---

# Session 5 — Migrations

Tasks:

- [ ] Install migration tool
- [ ] Add migration commands to Makefile
- [ ] Create first migration
- [ ] Run migrations

---

# Session 6 — Auth Register

Tasks:

- [ ] Create users table
- [ ] Create user model
- [ ] Create user repository
- [ ] Create password hashing helper
- [ ] Create register handler
- [ ] Test with Postman/Insomnia

---

# Session 7 — Auth Login and JWT

Tasks:

- [ ] Create login handler
- [ ] Validate password
- [ ] Generate access token
- [ ] Create auth middleware
- [ ] Create `/me`

---

# Session 8 — Refresh Tokens

Tasks:

- [ ] Create refresh_tokens table
- [ ] Generate refresh token
- [ ] Store hash
- [ ] Refresh access token
- [ ] Logout/revoke refresh token

---

# Session 9 — Organizations

Tasks:

- [ ] Create organizations table
- [ ] Create organization_members table
- [ ] Create organization module
- [ ] Add create/list/get/update/delete endpoints
- [ ] Add owner on organization creation

---

# Session 10 — Members and RBAC

Tasks:

- [ ] Add member by email
- [ ] List members
- [ ] Change role
- [ ] Remove member
- [ ] Add permission helper
- [ ] Add permission middleware

---

# Session 11 — Projects

Tasks:

- [ ] Create projects table
- [ ] Create project module
- [ ] Add project CRUD
- [ ] Add permission checks

---

# Session 12 — Tasks

Tasks:

- [ ] Create tasks table
- [ ] Create task module
- [ ] Add task CRUD
- [ ] Add status update
- [ ] Add assignee update
- [ ] Add filters/search

---

# Session 13 — Comments

Tasks:

- [ ] Create task_comments table
- [ ] Create comments module
- [ ] Add comment CRUD
- [ ] Add ownership rules

---

# Session 14 — Redis Features

Tasks:

- [ ] Add auth rate limiter
- [ ] Add permission role cache
- [ ] Add cache invalidation

---

# Session 15 — Tests and CI

Tasks:

- [ ] Add auth tests
- [ ] Add permission tests
- [ ] Add task tests
- [ ] Add GitHub Actions
- [ ] Fix CI failures

---

# Session 16 — Documentation

Tasks:

- [ ] Improve README
- [ ] Add API docs
- [ ] Add database docs
- [ ] Add architecture docs
- [ ] Add screenshots or terminal examples
