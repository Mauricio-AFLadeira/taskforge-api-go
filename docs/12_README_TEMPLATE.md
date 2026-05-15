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
