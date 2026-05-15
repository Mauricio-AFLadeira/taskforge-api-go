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
