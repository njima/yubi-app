# Yubi App

[English README](README.md)

ロボット fleet の遠隔操作データを収集・処理・管理するための Web プラットフォームです。バックエンドは Go、フロントエンドは Next.js で構成されています。

## プロジェクト構成

```
.
├── backend/          # REST API (Go, Gin, Bun ORM, Clean Architecture)
├── frontend/         # Web UI (Next.js 16, TypeScript, React 19)
├── openapi/          # OpenAPI 仕様 (single source of truth)
├── compose.yaml      # ローカル開発用 Docker Compose
└── Makefile          # 開発コマンドの入口
```

## クイックスタート

### 前提条件

- Docker & Docker Compose v2.20+
- Make

### セットアップ

```bash
# 1. リポジトリを clone
git clone https://github.com/airoa-org/yubi-app.git
cd yubi-app

# 2. 環境変数ファイルをコピー
cp backend/.env.example backend/.env
cp frontend/.env.sample frontend/.env

# 3. サービスを起動
make up PLATFORM=arm64    # Apple Silicon
# または
make up                   # Intel / Linux (amd64)

# 4. DB migration と seed data を投入
make migrate
make seed

# 5. アプリを開く
open http://localhost:3000/web
```

### サービス URL

| Service | URL | 説明 |
|---------|-----|------|
| Frontend | http://localhost:3000/web | Web UI |
| Backend API | http://localhost:8000 | REST API |
| LocalStack | http://localhost:4566 | S3 互換 storage (開発用) |

### 停止

```bash
make down       # 全サービスを停止
make reset      # 全サービスを停止し、DB などの volume を削除
```

## 認証

> **Note**: OSS 版では、ローカル開発・評価向けの簡易的な header-based authentication を使用します。本番レベルの security は提供しません。本番利用では OAuth2 や API gateway などの認証 layer 追加を検討してください。

- **Frontend**: server-side backend client が `X-User-ID` header を自動送信します。初期 user は `DEFAULT_USER_ID` 環境変数で指定します。右上の user menu から account を切り替えられます。
- **Robot API**: robots は `X-User-ID` と `X-Robot-ID` headers を直接送信します。API key や token は不要です。
- **RBAC**: database 上の user role に基づいて Role-based access control を適用します。

`make seed` で default Admin user が作成されます。user ID は `frontend/.env` で設定します。

## 開発

開発 command は Docker containers 内で実行されます。先に `make up` を実行してください。

```bash
make help           # 利用可能な command を表示
```

### Backend

```bash
make be-test          # test を実行
make be-lint          # staticcheck による lint
make be-fmt           # Go code を format
make be-tidy          # Go modules を整理
make be-generate-api  # OpenAPI spec から Go server code を再生成
```

### Dashboard / Batch

dashboard は事前集計済みの stats tables を参照します。集計 data を投入するには aggregation batch を実行してください。

```bash
make be-aggregate                     # 直前期間を集計 (PERIOD=hourly|daily|monthly)
make be-aggregate PERIOD=monthly      # 直前月を集計
make be-aggregate-backfill PERIOD=monthly FROM=2025-11-01 TO=2026-06-01  # 範囲を backfill
make be-uptime-writer                 # robot uptime metrics writer を起動 (long-running daemon)
```

`be-aggregate` は *直前の期間* だけを対象にします。dashboard の default 6 か月 window を埋めるには、その範囲に対して `be-aggregate-backfill` を実行してください。

### Frontend

```bash
make fe-fmt           # Prettier で format
make fe-lint          # ESLint を実行
make fe-typecheck     # TypeScript type check
make fe-ci            # CI 相当の checks (lint, format, typecheck, build)
make fe-generate-api  # OpenAPI spec から API client を再生成
```

### Database

```bash
make migrate        # 未適用 migrations を適用
make migrate-status # migration status を表示
make seed           # seed data を投入
make reset          # 全 data と volumes を削除
```

### API 開発フロー

1. `openapi/openapi.yaml` を編集する
2. `make be-generate-api` を実行する (Go server stubs を再生成)
3. `make fe-generate-api` を実行する (TypeScript client を再生成)
4. backend handler を実装する
5. frontend と接続する

### Database Migration フロー

1. `backend/internal/database/entity/` の entity を編集する
2. `make be-schema-gen` を実行する (`schema.up.sql` を再生成)
3. `make be-migrate-diff NAME=description` を実行する (migration SQL を生成)
4. `backend/internal/database/migrate/` の生成 SQL を確認する
5. `make migrate` を実行する

## 技術スタック

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
- **API Client**: Zodios (OpenAPI から生成)

## 設定

### Host Port の変更

default ports は環境変数で変更できます。

```bash
HOST_BACKEND_PORT=9000 HOST_DB_PORT=5433 make up
```

| Variable | Default | 説明 |
|----------|---------|------|
| `HOST_BACKEND_PORT` | 8000 | Backend API port |
| `HOST_FRONTEND_PORT` | 3000 | Frontend port |
| `HOST_DB_PORT` | 5432 | PostgreSQL port |
| `HOST_REDIS_PORT` | 6379 | Redis port |
| `HOST_LOCALSTACK_PORT` | 4566 | LocalStack (S3) port |
| `DOCKER_PLATFORM` | linux/amd64 | Docker platform |

## ドキュメント

初めて触る場合は、まず [ユーザーガイド](docs/ja/user-guide.md) で基本概念を確認し、必要に応じて他のドキュメントを参照してください。

| Document | 説明 |
|----------|------|
| [ユーザーガイド](docs/ja/user-guide.md) | 基本概念、tutorial、Web UI の使い方 |
| [Robot API ガイド](docs/ja/robot-api-guide.md) | robot authentication、episode execution flow、API examples |
| [Backend Architecture](docs/ja/backend-architecture.md) | Clean Architecture layers、DB migration workflow、batch commands |
| [Frontend Architecture](docs/ja/frontend-architecture.md) | project structure、API client pattern、feature modules |

## コントリビューション

guidelines は [CONTRIBUTING.md](CONTRIBUTING.md) を参照してください。

## ライセンス

この project は Apache License 2.0 のもとで提供されています。詳細は [LICENSE](LICENSE) を参照してください。
