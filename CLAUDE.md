# CLAUDE.md — Claude Code Instructions

This file provides guidance to Claude (Anthropic's AI assistant) when working on this repository.

## Quick Reference

- **Language:** Go 1.25+
- **Module:** `app`
- **Build:** `go build ./...` or `make build`
- **Test:** `go test -v -cover ./...` or `make test`
- **Run:** `make run` or `go run app/main.go`
- **Lint:** use standard `go vet ./...`

## Project Architecture

This is a **Clean Architecture** Go REST API boilerplate using the Gin framework.

### Layer Overview

```
domain/          → Entities, interfaces, models (no external dependencies)
user/            → Business logic services (depends on domain)
internal/
  repository/    → Data access (GORM, MongoDB, Redis, S3)
  rest/          → HTTP handlers + middleware (Gin)
  worker/        → Background job processors
helpers/         → Shared utilities (connections, JWT, pagination)
app/             → Entry point wiring everything together
```

### Dependency Flow

```
REST Handler → Service → Repository
     ↓            ↓          ↓
         domain (shared interfaces & models)
```

## Code Style & Conventions

- **No external test frameworks** — use Go's standard `testing` package with manual mocks.
- **Interface-driven design** — domain interfaces defined in `domain/`, implementations elsewhere.
- **Error types** are centralized in `domain/error.go`.
- **Models** are organized under `domain/model/` by category (auth, gorm, mongo, request, response, storage).
- **Environment config** — all configuration via environment variables loaded from `.env` by `godotenv`.
- **Logging** — use `logrus` with JSON formatter. Never use `fmt.Println` for logging.
- **HTTP responses** — use the response model in `domain/model/response/response.go`.

## Adding New Features

When adding a new domain module (e.g., `article`):

1. Define the entity interface in `domain/article.go`
2. Create the module directory: `article/service.go`, `article/storage.go`
3. Implement the repository in `internal/repository/gorm/article.go` (or mongo)
4. Create REST handlers in `internal/rest/article.go`
5. Wire it up in `app/main.go`
6. Write tests alongside each file using `_test.go` suffix

## Common Pitfalls

- The `internal/` directory is enforced by the Go compiler — packages here cannot be imported by external modules.
- Database connections use URL format strings, not individual host/port/user params.
- Redis is optional — check `USE_REDIS` env var before initializing.
- The default database driver is MySQL via GORM. Switch by uncommenting alternative drivers in `app/main.go`.

## Testing Guidelines

- Tests live next to the code they test (e.g., `service.go` → `service_test.go`).
- Mock dependencies manually by implementing interfaces, not via code generation frameworks.
- Run the full test suite before submitting changes: `make test`
- Check test coverage: `make testcoverage`

## File Locations

| What | Where |
|---|---|
| Entry point | `app/main.go` |
| Domain interfaces | `domain/*.go` |
| Data models | `domain/model/` |
| DB connections | `helpers/connection/` |
| JWT utilities | `helpers/jsonwebtoken/` |
| Middleware | `internal/rest/middleware/` |
| REST handlers | `internal/rest/` |
| Repositories | `internal/repository/` |
| Example module | `user/` |
| Docker config | `Dockerfile`, `docker-compose.yml` |
| CI/CD | `.github/workflows/unit-test.yml` |
| Env template | `.env.example` |
