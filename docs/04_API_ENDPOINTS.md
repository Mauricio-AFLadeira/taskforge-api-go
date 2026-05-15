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
