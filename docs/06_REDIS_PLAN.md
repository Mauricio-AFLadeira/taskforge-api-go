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
