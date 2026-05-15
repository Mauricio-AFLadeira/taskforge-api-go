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
