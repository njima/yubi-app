# Backend Architecture

## Overview

The backend is a Go-based REST API using Clean Architecture principles. It manages robots, tasks, episodes, and user data for the teleoperation platform.

## Directory Structure

```
backend/
├── cmd/
│   ├── server/                 # Main API server
│   ├── create-db-schema/       # Schema generation tool
│   ├── aggregate-episode-stats/ # Stats aggregation batch
│   └── write-robot-status-metrics/ # Metrics batch
├── internal/
│   ├── domain/                 # Domain models and business rules
│   │   ├── authz/              # Authorization policy
│   │   └── model/              # Domain entities and status policies
│   ├── gen/                    # Generated code (OpenAPI)
│   ├── infra/                  # External adapters
│   │   ├── cache/              # Redis-backed adapters
│   │   ├── database/           # Bun connection, entities, schema, migrations
│   │   ├── persistence/        # Repository implementations
│   │   └── storage/            # S3-backed adapters
│   ├── interfaces/             # HTTP controllers, middleware, SSE handlers
│   ├── platform/               # Runtime config and logging
│   ├── repository/             # Repository interfaces
│   ├── shared/                 # Cross-cutting helpers (errors, request context)
│   └── usecase/                # Application business logic
├── openapi.yaml                # oapi-codegen configuration
├── atlas.hcl                   # Atlas migration configuration
├── Dockerfile                  # Production image
└── Dockerfile.dev              # Development image (with Air hot reload)
```

## Clean Architecture Layers

```
HTTP Request
    │
    ▼
interfaces/ (controllers, middleware)
    │
    ▼
usecase/ (business logic)
    │
    ▼
repository/ (interfaces)  ←──  infra/persistence, infra/cache, infra/storage
    │
    ▼
infra/database (entities, ORM)
```

| Layer | Description |
|-------|-------------|
| **domain** | Core business entities and validation rules |
| **usecase** | Application-specific business logic, orchestrates domain objects |
| **repository** | Data access interface definitions (ports) |
| **infra** | Concrete implementations of repository interfaces (adapters) |
| **interfaces** | HTTP controllers, middleware, SSE handlers |
| **platform** | Runtime configuration and logging |
| **shared** | Cross-cutting helpers that do not own business rules |

### Key Principles

- **Dependency direction**: Outer layers depend on inner layers, never the reverse
- **Domain models** use `Init*()` constructors for creation with validation, `New*()` for reconstruction from DB
- **Status policy** belongs in domain models; usecases should call lifecycle helpers such as `Start`, `Finish`, and `Cancel` instead of setting lifecycle state directly
- **Usecase files** may be split by workflow (for example, episode lifecycle vs. recording/stat lookups) while sharing the same usecase type
- **HTTP stream presenters** live under `interfaces/http/handler` and keep OpenAPI response construction out of stream control flow
- **Dual ID system**: `ID` (int64, internal PK) + `IDNatural` (UUID, exposed via API)
- **Organization scoping**: `OrgScoped` entity hook automatically filters queries by `organization_id` from request context

## Authentication & Authorization

### Authentication (Middleware)

The auth middleware supports two methods for Robot API paths (`/api/robot/*`).

1. **API Key** (recommended): `X-API-Key` header → hash lookup in Postgres → sets context (userID, robotID, orgID, role)
2. **Header fallback**: `X-User-ID` + `X-Robot-ID` headers → DB lookup → validates same organization

For non-robot paths (`/api/*`), `X-User-ID` header is required.

API key authentication is preferred for production robot deployments. The header fallback is available for development and testing.

### Authorization (RBAC)

The `authz` middleware checks user role against operation permissions defined per OpenAPI operationID:

| Role | Access |
|------|--------|
| Admin | Full access |
| DataEngineer | Location read-only; full Task/Robot/Episode access |
| Manager | Same as DataEngineer |
| Operator | Location/Task/Robot read-only; Episode create/update |
| Viewer | Read-only all resources |

## Database

### Tech Stack

- **PostgreSQL 17.5** — primary database
- **Bun ORM** — SQL-first ORM with struct tags
- **Atlas** — migration management

### Migration Workflow

1. Edit entity in `internal/infra/database/entity/`
2. `make be-schema-gen` — generates `schema.up.sql` from entities
3. `make be-migrate-diff NAME=description` — generates migration SQL
4. Review SQL in `internal/infra/database/migrate/`
5. `make migrate` — apply migration

### Hand-edited Migrations

When a migration needs production-safe directives (e.g., `CREATE INDEX CONCURRENTLY`), follow these steps.

1. Define the index in the entity (so `schema-gen` includes it)
2. Run `make be-migrate-diff NAME=name` to generate baseline
3. Hand-edit the migration (add `-- atlas:txmode none`, etc.)
4. Run `atlas migrate hash` to recompute `atlas.sum`

## Batch Commands

### Episode Stats Aggregation

`cmd/aggregate-episode-stats/main.go` aggregates episode data into hourly, daily, and monthly stats tables.

```bash
cd backend

# Regular aggregation (previous period)
go run ./cmd/aggregate-episode-stats/main.go
go run ./cmd/aggregate-episode-stats/main.go --period daily
go run ./cmd/aggregate-episode-stats/main.go --period monthly

# Backfill historical data [from, to)
go run ./cmd/aggregate-episode-stats/main.go \
  --period hourly --backfill \
  --from 2026-03-01T00:00:00Z --to 2026-03-02T00:00:00Z
```

| Flag | Default | Description |
|------|---------|-------------|
| `--period` | `hourly` | `hourly`, `daily`, `monthly` |
| `--backfill` | `false` | Enable historical backfill |
| `--from` | - | Backfill start |
| `--to` | - | Backfill end (exclusive) |

DB connection is configured via flags or environment variables (`DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`).

From the repo root, `make be-aggregate` and `make be-aggregate-backfill` wrap these
commands and run them inside the backend container.

### Robot Uptime Metrics

`cmd/write-robot-status-metrics/main.go` is a long-running worker that flushes
Redis-buffered robot uptime deltas into the `robot_uptime_hourly` table on a timer.
The dashboard's uptime figures read from that table.

```bash
cd backend
go run ./cmd/write-robot-status-metrics/main.go   # or: make be-uptime-writer
```

## API

- OpenAPI spec: [`openapi/openapi.yaml`](../openapi/openapi.yaml)
- Code generation: `make be-generate-api`
- Health check: `GET /health-check`
