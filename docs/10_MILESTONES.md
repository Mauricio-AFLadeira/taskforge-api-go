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

# Milestone 1 — Backend Foundation

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
