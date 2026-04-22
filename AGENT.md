# AGENT.md — AI Agent Instructions

This file provides instructions for AI coding agents working on this repository.

## Project Overview

This is a Go RESTful API boilerplate template built with **Clean Architecture** principles. It uses the **Gin** web framework, **GORM** ORM, **MongoDB** driver, **Redis** caching, and **AWS S3** for storage.

## Tech Stack

- **Language:** Go 1.25+
- **Module name:** `app`
- **Web framework:** Gin (`github.com/gin-gonic/gin`)
- **ORM:** GORM (`gorm.io/gorm`)
- **Logging:** Logrus (`github.com/sirupsen/logrus`)
- **Auth:** JWT v5 (`github.com/golang-jwt/jwt/v5`)
- **Cache:** Redis (`github.com/redis/go-redis/v9`)
- **Cloud storage:** AWS S3 SDK v2

## Directory Structure

```
app/                  → Application entry point (main.go)
domain/               → Business entities, interfaces, models, DTOs
helpers/              → Utility packages (DB connections, JWT, pagination)
internal/             → Private implementation code
  ├── repository/     → Data access layer (gorm, mongo, redis, s3)
  ├── rest/           → HTTP handlers and middleware
  └── worker/         → Background workers
user/                 → Example domain module (service, storage, cache)
```

## Build & Run Commands

```bash
make run              # Run the application (go run app/main.go)
make test             # Run all tests with verbose output and coverage
make testcoverage     # Run tests and generate coverage report
make generate-mocks   # Generate mock files using mockery
make build            # Build optimized binary (output: build-app)
```

## Testing

- Tests use Go's standard `testing` package with manual mock implementations.
- No external test frameworks (e.g., testify) are used.
- Run tests: `go test -v -cover ./...` or `make test`
- Test files follow the `_test.go` convention alongside source files.

## Architecture Pattern

This project follows **Clean Architecture** with these layers:

1. **Domain** (`domain/`) — Entities, interfaces, and value objects
2. **Use Case / Service** (`user/`, etc.) — Business logic
3. **Repository** (`internal/repository/`) — Data access implementations
4. **Delivery** (`internal/rest/`) — HTTP handlers and middleware

Dependencies point inward: delivery → use case → domain ← repository.

## Conventions

- Environment variables are loaded from `.env` via `godotenv`.
- Configuration is done through environment variables (see `.env.example`).
- The application runs on port `5050` by default.
- Database connection strings use URL format (e.g., `mysql://user:password@host:port/db`).
- Internal packages under `internal/` are not importable by external projects.
- Each domain module has its own directory with `service.go`, `storage.go`, and optionally `cache.go`.

## Docker

- **Dockerfile** uses multi-stage build: Go 1.25 builder → Alpine runtime.
- **docker-compose.yml** orchestrates the app with MySQL and Redis.
- The app container exposes port `5050`.

## Environment Variables

Key environment variables (see `.env.example` for full list):

| Variable | Description | Default |
|---|---|---|
| `PORT` | Server port | `5050` |
| `GO_ENV` | Environment mode | `production` |
| `DB_URL` | Database connection URL | — |
| `USE_REDIS` | Enable Redis caching | `false` |
| `REDIS_URL` | Redis connection URL | — |
| `JWT_MEMBER_SECRET_KEY` | JWT signing secret | — |
