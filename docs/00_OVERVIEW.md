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
