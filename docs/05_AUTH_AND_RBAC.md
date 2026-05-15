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
