# Backend Architecture

## 概要

backend は Go 製の REST API で、Clean Architecture の考え方に沿って構成されています。teleoperation platform の robots、tasks、episodes、users data を管理します。

## ディレクトリ構成

```
backend/
├── cmd/
│   ├── server/                 # main API server
│   ├── create-db-schema/       # schema generation tool
│   ├── aggregate-episode-stats/ # stats aggregation batch
│   └── write-robot-status-metrics/ # metrics batch
├── internal/
│   ├── apperror/               # application error definitions
│   ├── authz/                  # Role-based access control (RBAC)
│   ├── ccontext/               # request context helpers
│   ├── config/                 # configuration management
│   ├── database/               # schema, migrations, entities, seeders
│   ├── domain/                 # domain models and business rules
│   ├── event/                  # real-time updates 用 event bus (SSE)
│   ├── gateway/                # repository implementations
│   ├── gen/                    # generated code (OpenAPI)
│   ├── interfaces/             # HTTP controllers and middleware
│   ├── log/                    # logging utilities
│   ├── pagination/             # pagination helpers
│   ├── redis/                  # Redis client
│   ├── repository/             # repository interfaces
│   ├── s3/                     # S3 client
│   ├── stack/                  # stack trace utilities
│   └── usecase/                # application business logic
├── openapi.yaml                # oapi-codegen configuration
├── atlas.hcl                   # Atlas migration configuration
├── Dockerfile                  # production image
└── Dockerfile.dev              # development image (Air hot reload)
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
repository/ (interfaces)  ←──  gateway/ (implementations)
    │
    ▼
database/ (entities, ORM)
```

| Layer | 説明 |
|-------|------|
| **domain** | 中核となる business entities と validation rules |
| **usecase** | application-specific business logic。domain objects を orchestrate する |
| **repository** | data access interface definitions (ports) |
| **gateway** | repository interfaces の concrete implementations (adapters) |
| **interfaces** | HTTP controllers, middleware, SSE handlers |

### 基本方針

- **Dependency direction**: outer layers が inner layers に依存し、逆方向には依存しない
- **Domain models**: creation と validation には `Init*()` constructors、DB からの reconstruction には `New*()` を使う
- **Dual ID system**: `ID` (int64, internal PK) + `IDNatural` (UUID, API に公開)
- **Organization scoping**: `OrgScoped` entity hook が request context の `organization_id` で queries を自動 filter する

## 認証と認可

### 認証 (Middleware)

auth middleware は Robot API paths (`/api/robot/*`) に対して 2 つの方式をサポートします。

1. **API Key** (推奨): `X-API-Key` header → Postgres で hash lookup → context (userID, robotID, orgID, role) を設定
2. **Header fallback**: `X-User-ID` + `X-Robot-ID` headers → DB lookup → 同じ organization か検証

robot 以外の paths (`/api/*`) では `X-User-ID` header が必要です。

production robot deployments では API key authentication を推奨します。header fallback は development/testing 向けです。

### 認可 (RBAC)

`authz` middleware は OpenAPI operationID ごとに定義された operation permissions と user role を照合します。

| Role | Access |
|------|--------|
| Admin | Full access |
| DataEngineer | Location read-only; Task/Robot/Episode full access |
| Manager | DataEngineer と同等 |
| Operator | Location/Task/Robot read-only; Episode create/update |
| Viewer | すべて read-only |

## Database

### 技術スタック

- **PostgreSQL 17.5** — primary database
- **Bun ORM** — SQL-first ORM with struct tags
- **Atlas** — migration management

### Migration Workflow

1. `internal/database/entity/` の entity を編集する
2. `make be-schema-gen` — entities から `schema.up.sql` を生成
3. `make be-migrate-diff NAME=description` — migration SQL を生成
4. `internal/database/migrate/` の SQL を確認する
5. `make migrate` — migration を適用する

### 手動編集が必要な Migration

production-safe な directives が必要な場合 (例: `CREATE INDEX CONCURRENTLY`) は次の手順に従います。

1. entity に index を定義する (`schema-gen` に含めるため)
2. `make be-migrate-diff NAME=name` で baseline を生成する
3. migration を手動編集する (`-- atlas:txmode none` などを追加)
4. `atlas migrate hash` で `atlas.sum` を再計算する

## Batch Commands

### Episode Stats Aggregation

`cmd/aggregate-episode-stats/main.go` は episode data を hourly、daily、monthly stats tables に集計します。

```bash
cd backend

# 通常集計 (直前期間)
go run ./cmd/aggregate-episode-stats/main.go
go run ./cmd/aggregate-episode-stats/main.go --period daily
go run ./cmd/aggregate-episode-stats/main.go --period monthly

# historical data の backfill [from, to)
go run ./cmd/aggregate-episode-stats/main.go \
  --period hourly --backfill \
  --from 2026-03-01T00:00:00Z --to 2026-03-02T00:00:00Z
```

| Flag | Default | 説明 |
|------|---------|------|
| `--period` | `hourly` | `hourly`, `daily`, `monthly` |
| `--backfill` | `false` | historical backfill を有効化 |
| `--from` | - | backfill start |
| `--to` | - | backfill end (exclusive) |

DB connection は flags または環境変数 (`DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`) で設定します。

repo root からは `make be-aggregate` と `make be-aggregate-backfill` がこれらの commands を backend container 内で実行します。

### Robot Uptime Metrics

`cmd/write-robot-status-metrics/main.go` は long-running worker です。Redis-buffered robot uptime deltas を timer で `robot_uptime_hourly` table に flush します。dashboard の uptime 値はこの table を参照します。

```bash
cd backend
go run ./cmd/write-robot-status-metrics/main.go   # または: make be-uptime-writer
```

## API

- OpenAPI spec: [`openapi/openapi.yaml`](../../openapi/openapi.yaml)
- Code generation: `make be-generate-api`
- Health check: `GET /health-check`
