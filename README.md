# Yubi App

[日本語版 README](README.ja.md)

A web platform for collecting, processing, and managing teleoperation data from robot fleets. Built with Go (backend) and Next.js (frontend).

## Project Structure

```
.
├── backend/          # REST API (Go, Gin, Bun ORM, Clean Architecture)
├── frontend/         # Web UI (Next.js 16, TypeScript, React 19)
├── openapi/          # OpenAPI specification (single source of truth)
├── compose.yaml      # Docker Compose for local development
└── Makefile          # Unified development commands
```

## Quick Start

### Prerequisites

- Docker & Docker Compose v2.20+
- Make

### Setup

```bash
# 1. Clone the repository
git clone https://github.com/airoa-org/yubi-app.git
cd yubi-app

# 2. Copy environment files
cp backend/.env.example backend/.env
cp frontend/.env.sample frontend/.env

# 3. Start all services
make up PLATFORM=arm64    # Apple Silicon
# or
make up                   # Intel / Linux (amd64)

# 4. Apply database migrations and seed data
make migrate
make seed

# 5. Open the application
open http://localhost:3000/web
```

### Service URLs

| Service | URL | Description |
|---------|-----|-------------|
| Frontend | http://localhost:3000/web | Web UI |
| Backend API | http://localhost:8000 | REST API |
| LocalStack | http://localhost:4566 | S3-compatible storage (dev) |

### Stop Services

```bash
make down       # Stop all services
make reset      # Stop and delete all data (volumes)
```

## Authentication

> **Note**: This OSS version uses a simplified, header-based authentication intended for local development and evaluation. It does not provide production-grade security. For production use, consider adding an authentication layer (e.g., OAuth2, API gateway).

- **Frontend**: Google sign-in is handled by Auth.js / NextAuth. The server-side backend client sends `X-User-ID` from the authenticated session, with `DEFAULT_USER_ID` available as a local development fallback.
- **Robot API**: Robots send `X-User-ID` and `X-Robot-ID` headers directly. No API key or token is required.
- **RBAC**: Role-based access control is enforced based on the user's active organization membership role.

A default Admin user and organization membership are created via `make seed`. The user ID is configured in `frontend/.env`.

See [Authentication and Workspace Setup](docs/authentication.md) for Google OAuth setup, the local development fallback, and dashboard 403 troubleshooting.

## Development

All development commands run inside Docker containers. Run `make up` first.

```bash
make help           # Show all available commands
```

### Backend

```bash
make be-test        # Run tests
make be-lint        # Run linter (staticcheck)
make be-fmt         # Format code
make be-tidy        # Tidy Go modules
make be-generate-api  # Regenerate code from OpenAPI spec
```

### Dashboard / Batch

The dashboard reads pre-aggregated stats tables. Run the aggregation batch to
populate them.

```bash
make be-aggregate                     # Aggregate the previous period (PERIOD=hourly|daily|monthly)
make be-aggregate PERIOD=monthly      # Aggregate the previous month
make be-aggregate-backfill PERIOD=monthly FROM=2025-11-01 TO=2026-06-01  # Backfill a range
make be-uptime-writer                 # Run the robot uptime metrics writer (long-running daemon)
```

`be-aggregate` covers the *previous* period only. To populate the dashboard's
default 6-month window, use `be-aggregate-backfill` over that range.

### Frontend

```bash
make fe-fmt         # Format with Prettier
make fe-lint        # Run ESLint
make fe-typecheck   # TypeScript type check
make fe-ci          # Run all CI checks (lint, format, typecheck, build)
make fe-generate-api  # Regenerate API client from OpenAPI spec
```

### Database

```bash
make migrate        # Apply pending migrations
make migrate-status # Show migration status
make seed           # Insert seed data
make reset          # Delete all data and volumes
```

### API Development Workflow

1. Edit `openapi/openapi.yaml`
2. `make be-generate-api` (regenerate Go server stubs)
3. `make fe-generate-api` (regenerate TypeScript client)
4. Implement backend handler
5. Connect frontend

### Database Migration Workflow

1. Edit entity in `backend/internal/database/entity/`
2. `make be-schema-gen` (regenerate schema.up.sql)
3. `make be-migrate-diff NAME=description` (generate migration SQL)
4. Review generated SQL in `backend/internal/database/migrate/`
5. `make migrate` (apply migration)

## Technology Stack

### Backend

- **Language**: Go 1.25+
- **Framework**: Gin (HTTP), Bun (ORM)
- **Database**: PostgreSQL 17.5
- **Cache**: Redis
- **Migration**: Atlas
- **Architecture**: Clean Architecture

### Frontend

- **Framework**: Next.js 16
- **Language**: TypeScript
- **UI**: React 19, Radix UI, Tailwind CSS
- **State**: TanStack Query
- **API Client**: Zodios (generated from OpenAPI)

## Configuration

### Host Port Overrides

Default ports can be changed via environment variables:

```bash
HOST_BACKEND_PORT=9000 HOST_DB_PORT=5433 make up
```

| Variable | Default | Description |
|----------|---------|-------------|
| `HOST_BACKEND_PORT` | 8000 | Backend API port |
| `HOST_FRONTEND_PORT` | 3000 | Frontend port |
| `HOST_DB_PORT` | 5432 | PostgreSQL port |
| `HOST_REDIS_PORT` | 6379 | Redis port |
| `HOST_LOCALSTACK_PORT` | 4566 | LocalStack (S3) port |
| `DOCKER_PLATFORM` | linux/amd64 | Docker platform |

## Documentation

New to the project? Start with the [User Guide](docs/user-guide.md) to understand the core concepts, then refer to the other documents as needed.

| Document | Description |
|----------|-------------|
| [User Guide](docs/user-guide.md) | Start here — core concepts, tutorial walkthrough, Web UI usage |
| [Robot API Guide](docs/robot-api-guide.md) | Robot authentication, episode execution flow, API examples |
| [Authentication and Workspace Setup](docs/authentication.md) | Local auth model, workspace membership setup, dashboard 403 troubleshooting |
| [Backend Architecture](docs/backend-architecture.md) | Clean Architecture layers, DB migration workflow, batch commands |
| [Frontend Architecture](docs/frontend-architecture.md) | Project structure, API client pattern, feature modules |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.
