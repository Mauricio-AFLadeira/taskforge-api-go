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
