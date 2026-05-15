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
