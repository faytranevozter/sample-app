# Golang API Boilerplate Template

[![Go Unit Test](https://github.com/faytranevozter/go-template/actions/workflows/unit-test.yml/badge.svg)](https://github.com/faytranevozter/go-template/actions/workflows/unit-test.yml)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> A production-ready boilerplate for building RESTful APIs with Clean Architecture in Go

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Available Commands](#available-commands)
- [Docker](#docker)
- [Usage Guide](#usage-guide)
- [Authentication](#authentication)
- [Database Support](#database-support)
- [Environment Variables](#environment-variables)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## 🎯 Overview

This is a comprehensive boilerplate template for building scalable RESTful APIs in [Golang](https://golang.org/). It implements [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) principles, providing a solid foundation for your next microservice or API project.

The architecture is inspired by [go-clean-arch](https://github.com/bxcodec/go-clean-arch) with extensive customizations and modern best practices. This template helps you get started quickly with well-organized code structure and industry-standard patterns.

## ✨ Features

- 🏗️ **Clean Architecture** — Separation of concerns with clear boundaries between layers
- 🔐 **JWT Authentication** — Built-in JWT token generation and validation
- 🗄️ **Multiple Database Support** — MySQL, PostgreSQL, SQLite, SQL Server, and MongoDB
- 💾 **Redis Caching** — Integrated caching layer with middleware
- ☁️ **AWS S3 Integration** — File storage and upload support
- 🚀 **RESTful API** — Built with Gin framework for high performance
- 📝 **Structured Logging** — JSON logging with Logrus and log rotation via Lumberjack
- 🔄 **CORS Support** — Configurable cross-origin resource sharing
- 🛡️ **Middleware Suite** — Auth, cache, CORS, logger, and recovery
- 🐳 **Docker & Docker Compose** — Full containerization with service orchestration

## 📋 Tech Stack

| Category | Technologies |
|---|---|
| **Core** | Go 1.25+, [Gin](https://github.com/gin-gonic/gin), [Logrus](https://github.com/sirupsen/logrus) |
| **SQL Databases** | MySQL, PostgreSQL, SQLite, SQL Server via [GORM](https://gorm.io/) |
| **NoSQL** | [MongoDB](https://www.mongodb.com/) 6.0+ |
| **Caching** | [Redis](https://redis.io/) 7+ |
| **Auth** | [JWT v5](https://github.com/golang-jwt/jwt) |
| **Storage** | [AWS S3](https://aws.amazon.com/s3/) (SDK v2) |
| **DevOps** | Docker, Docker Compose, GitHub Actions |

## 🏛️ Architecture

This project follows Clean Architecture principles. Dependencies point **inward** — outer layers depend on inner layers, never the reverse.

```
┌──────────────────────────────────────────┐
│              Delivery Layer              │
│         (REST handlers, middleware)       │
├──────────────────────────────────────────┤
│             Use Case Layer               │
│          (Services, business logic)       │
├──────────────────────────────────────────┤
│            Repository Layer              │
│     (Database, cache, storage access)     │
├──────────────────────────────────────────┤
│              Domain Layer                │
│      (Entities, interfaces, models)       │
└──────────────────────────────────────────┘
```

**Flow:** HTTP Request → Middleware → Handler → Service → Repository → Database

## 📁 Project Structure

```
.
├── app/                    # Application entry point
│   └── main.go
├── domain/                 # Business entities and interfaces
│   ├── article.go
│   ├── user.go
│   ├── error.go
│   └── model/             # Data models and DTOs
│       ├── auth/          # JWT claims
│       ├── gorm/          # GORM base models
│       ├── mongo/         # MongoDB base models
│       ├── request/       # Request DTOs
│       ├── response/      # Response DTOs
│       ├── storage/       # Storage DTOs
│       └── filter.go      # Query filters
├── helpers/               # Utility functions and helpers
│   ├── connection/        # Database connection helpers
│   ├── jsonwebtoken/      # JWT utilities
│   ├── common.go          # Common helpers
│   └── pagination.go      # Pagination utilities
├── internal/              # Private application code
│   ├── repository/        # Data access layer
│   │   ├── gorm/         # SQL repositories (GORM)
│   │   ├── mongo/        # MongoDB repositories
│   │   ├── redis/        # Cache repositories
│   │   └── s3/           # S3 storage repositories
│   ├── rest/             # HTTP handlers
│   │   └── middleware/   # HTTP middlewares
│   └── worker/           # Background workers
├── user/                 # User module (example)
│   ├── service.go        # Business logic
│   ├── storage.go        # Storage interface
│   └── cache.go          # Cache decorator
├── docker-compose.yml    # Docker Compose orchestration
├── Dockerfile            # Multi-stage Docker build
├── Makefile             # Build automation
├── .env.example         # Environment variable template
└── go.mod               # Go module dependencies
```

## 🚀 Getting Started

### Prerequisites

- **Go** 1.25 or higher
- **Database** — MySQL, PostgreSQL, MongoDB, or SQLite
- **Redis** (optional) — for caching layer
- **Docker** (optional) — for containerized deployment

### Installation

1. **Clone this repository**
   ```bash
   git clone https://github.com/faytranevozter/go-template.git
   cd go-template
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```

3. **Configure your `.env` file**

   Edit the `.env` file with your database credentials, Redis configuration, and other settings.

4. **Run the application**
   ```bash
   make run
   ```

   Or using Go directly:
   ```bash
   go run app/main.go
   ```

The application will download all dependencies automatically and start the server on port `5050`.

### Quick Start with Docker Compose

The fastest way to get the full stack running:

```bash
cp .env.example .env
docker compose up -d
```

This starts the application along with MySQL and Redis. The API will be available at `http://localhost:5050`.

## 🛠️ Available Commands

```bash
make run              # Run the application
make test             # Run all tests with coverage
make testcoverage     # Run tests and generate coverage report
make generate-mocks   # Generate mock files for testing
make build            # Build optimized binary (output: build-app)
```

## 🐳 Docker

### Docker Compose (Recommended)

Docker Compose sets up the complete environment with the app, MySQL, and Redis:

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f app

# Stop all services
docker compose down

# Stop and remove volumes (reset data)
docker compose down -v
```

**Services included:**

| Service | Image | Port |
|---|---|---|
| `app` | Built from `Dockerfile` | `5050` |
| `mysql` | `mysql:8.0` | `3306` |
| `redis` | `redis:7-alpine` | `6379` |

You can customize ports and credentials via environment variables in your `.env` file:

```env
# Docker Compose overrides
MYSQL_ROOT_PASSWORD=rootpassword
MYSQL_DATABASE=app_db
MYSQL_USER=user
MYSQL_PASSWORD=password
MYSQL_PORT=3306
REDIS_PORT=6379
```

### Standalone Docker

Build and run the application container individually:

```bash
# Build the image
docker build -t build-app .

# Run the container
docker run -p 5050:5050 --env-file .env -d build-app
```

## 📖 Usage Guide

### Adding a New Module

1. **Define domain entity** in `domain/`
2. **Create storage interface** in your module directory
3. **Implement repository** in `internal/repository/`
4. **Add business logic** in service layer
5. **Create REST handlers** in `internal/rest/`
6. **Register routes** in `app/main.go`

### Example module structure (User):
```
user/
├── service.go        # Business logic and use cases
├── service_test.go   # Service tests
├── storage.go        # Storage interface definition
└── cache.go          # Caching decorator layer
```

## 🔒 Authentication

JWT authentication is built-in. Use the `auth` middleware to protect your routes:

```go
router.Use(middleware.Auth())
```

Token generation utilities are available in `helpers/jsonwebtoken/`.

## 🗄️ Database Support

This boilerplate supports multiple databases simultaneously:

- **SQL databases** via GORM (MySQL, PostgreSQL, SQLite, SQL Server)
- **NoSQL** via MongoDB driver
- **Caching** via Redis

Connection helpers are in `helpers/connection/`. Switch the active database driver by uncommenting the appropriate line in `app/main.go`.

## ⚙️ Environment Variables

| Variable | Description | Default |
|---|---|---|
| `PORT` | Server port | `5050` |
| `GO_ENV` | Environment (`production` / `development`) | `production` |
| `LOG_TO_STDOUT` | Log to standard output | `true` |
| `LOG_TO_FILE` | Log to file | `false` |
| `LOG_MAX_SIZE` | Max log file size (MB) | `50` |
| `DB_URL` | Database connection URL | — |
| `USE_REDIS` | Enable Redis caching | `false` |
| `REDIS_URL` | Redis connection URL | — |
| `REDIS_TTL` | Cache TTL | `60s` |
| `REDIS_KEY_PREFIX` | Redis key prefix | `app:` |
| `JWT_MEMBER_SECRET_KEY` | JWT signing secret | — |
| `JWT_MEMBER_TTL` | JWT token TTL (minutes) | `60` |
| `S3_ENDPOINT` | S3 endpoint URL | — |
| `S3_REGION` | S3 region | `auto` |
| `S3_ACCESS_KEY` | S3 access key | — |
| `S3_SECRET_KEY` | S3 secret key | — |
| `S3_BUCKET_NAME` | S3 bucket name | — |

See [`.env.example`](.env.example) for the full template.

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Please read the [Contributing Guide](CONTRIBUTING.md) for details on our development workflow, code style, and how to submit pull requests.

- 🐛 [Report a Bug](../../issues/new?template=bug_report.md)
- 💡 [Request a Feature](../../issues/new?template=feature_request.md)

## 📝 License

This project is licensed under the [MIT License](LICENSE).

Copyright (c) 2025 Fahrur Rifai

## 📧 Contact

You can reach the maintainer at [mfahrurrifai@gmail.com](mailto:mfahrurrifai@gmail.com).

---

**Happy coding! 🚀**
