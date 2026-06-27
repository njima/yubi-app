# Frontend Architecture

## 概要

frontend は Next.js 16 application で、robot fleets、tasks、teleoperation episodes を管理します。OpenAPI specification から生成された type-safe API clients を使用します。

## 技術スタック

| Category | Technology |
|----------|------------|
| Framework | Next.js 16 (App Router) |
| Language | TypeScript 5 |
| UI Components | Radix UI, Tailwind CSS 4 |
| State Management | TanStack Query (React Query) |
| Forms | React Hook Form + Zod |
| API Client | Zodios (OpenAPI から自動生成) |
| URL State | nuqs |
| Notifications | Sonner |

## ディレクトリ構成

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router
│   │   ├── (app)/              # main application routes
│   │   │   ├── dashboard/      # Dashboard page
│   │   │   ├── episodes/       # Episodes management
│   │   │   ├── robots/         # Robots management
│   │   │   ├── tasks/          # Tasks management
│   │   │   ├── users/          # Users management
│   │   │   └── profile/        # User profile
│   │   └── api/                # API route handlers (backend への proxy)
│   │
│   ├── components/             # domain-neutral React components
│   │   ├── layout/             # app shell and reusable layout components
│   │   └── ui/                 # shadcn-style UI primitives
│   │
│   ├── features/               # feature modules (domain-driven)
│   │   ├── api-keys/           # API key management
│   │   ├── dashboard/          # Dashboard composition
│   │   ├── episodes/           # Episode CRUD, hooks, components
│   │   ├── locations/          # Location management
│   │   ├── organizations/      # Organization management
│   │   ├── reporting/          # reporting and export workflows
│   │   ├── robots/             # Robot CRUD, status, camera viewer
│   │   ├── tasks/              # Task management
│   │   └── users/              # User CRUD, roles
│   │
│   ├── lib/                    # app-wide non-UI libraries
│   │   ├── api/                # API client configuration
│   │   │   ├── client.ts       # interceptors 付き Zodios client
│   │   │   ├── client-fetch.ts # browser-side fetch/schema helpers
│   │   │   ├── query-string.ts # query string builder
│   │   │   ├── backend-client.ts # server-side re-export facade
│   │   │   ├── backend-client/ # endpoint group ごとの server-side backend wrappers
│   │   │   ├── config.ts       # API URL configuration
│   │   │   └── generated/      # auto-generated Zodios client
│   │   ├── auth/               # session management
│   │   ├── hooks/              # app-wide hooks
│   │   ├── i18n/               # i18n setup, language storage, locales
│   │   ├── providers/          # React providers (QueryProvider)
│   │   └── status/             # app-wide status constants and display metadata
│
├── public/                     # static assets
├── next.config.ts              # Next.js configuration
├── Dockerfile                  # production image
├── Dockerfile.dev              # development image
└── package.json
```

## API Client Architecture

```
Browser
    │
    ▼
React Component (Client Component)
    │
    ▼
TanStack Query Hook (e.g., useEpisodesQuery)
    │
    ▼
Next.js API Route (/web/api/*)
    │
    ▼
backend-client.ts (X-User-ID header を追加)
    │
    ▼
Backend API (http://backend:8000)
```

### Server Components と Client Components

- **Server Components / Route Handlers** (`backend-client.ts`): `X-User-ID` header 付きで backend から直接 data を取得します。facade は `lib/api/backend-client/` 配下の endpoint-specific modules を re-export します。
- **Client Components**: `/web/api/*` routes を呼ぶ TanStack Query hooks を使います。繰り返しになる fetch、schema parsing、query string building には `client-fetch.ts` と `query-string.ts` を優先して使います。
- **SSE Streams**: real-time updates のため、`sse-proxy.ts` 経由で `X-User-ID` header を付けて proxy します。

### 依存関係のルール

dependency direction は次の形に保ちます。

```text
app -> features -> lib
app -> components
app/api -> lib/api
features -> lib
features -> components
lib/api -> generated, auth
```

`lib` は app-wide non-UI code を置く場所で、concrete な `app`、`features`、`components` modules を import しません。feature-specific helpers は `features/*/lib` に置きます。feature 間の import は、target feature の小さな public API 経由にします。

`app/**/page.tsx` と `app/**/layout.tsx` は薄い route entrypoint として扱います。feature-owned page composition は `features/*/components`、navigation などの app shell components は `components/layout` に置きます。`app/**/_components` は、feature や app-wide library area に自然に置けない小さな route glue に限ります。

### Component Placement

- `components/ui`: `Button`、`Dialog`、`Table`、`DropdownMenu` などの shadcn-style primitives。
- `components/layout`: `TopNav`、navigation items、user menu composition などの app shell components。app shell components は current user menu などの feature state を compose できます。
- `features/*/components`: list pages、export menus、teleoperation screens を含む feature-owned UI。feature modules は page と 1:1 でなくてもよく、`reporting` のような capability modules も許容します。
- `app/**`: route entrypoints、route groups、layouts、API routes。route-local components は少数かつ小さく保ちます。
- `lib/*`: API clients、auth、hooks、i18n、providers、status metadata、utilities などの app-wide non-UI code。React UI primitives はここに置きません。

import boundary rules は frontend container 内の `npm run lint:boundaries`、または `make fe-ci` で確認します。

### API Code Generation

API client は OpenAPI spec から自動生成します。

```bash
make fe-generate-api
```

この command は `src/lib/api/generated/api.ts` を生成し、次の内容を含みます。

- 全 request/response types の Zod schemas
- typed endpoints を持つ Zodios API client

**generated files を直接編集しないでください。** `openapi/openapi.yaml` を変更して再生成してください。

## Feature Module Pattern

`src/features/` 配下の各 feature は自己完結した構成にします。

```
features/episodes/
├── components/          # feature-specific React components
│   ├── episode-columns.tsx
│   ├── episode-detail.tsx
│   └── create-episode-dialog.tsx
├── hooks/               # data fetching and mutation hooks
│   ├── use-episodes-query.ts
│   └── use-create-episode-mutation.ts
├── schemas/             # Zod validation schemas (forms 用)
└── index.ts             # public API (re-exports)
```

feature files は責務を絞ります。大きくなった feature-specific renderer は `teleop-layout-components.tsx` のような local file に分け、registration や page wiring とは分離します。form ownership は、section の props contract が安定している場合を除き parent form に残します。

Robot teleoperation layout rendering は、robot config、robot status streams、episode/task context に依存するため feature-owned とします。registry、renderer、layout config types は app-wide `lib` や generic app layout directory ではなく `features/robots` 配下に置きます。

## Status Display Policy

status values は `lib/status/constants.ts` に置きます。badge color、label key、terminal state、successful completion などの display metadata は `lib/status/display.ts` に置きます。feature badge components は feature-specific のままでよいですが、status map を重複させず app-wide status metadata を使います。

## 認証

frontend は cookie で active user を識別し、なければ `DEFAULT_USER_ID` 環境変数へ fallback します。

- `session.ts` の `getUserId()` はまず `active_user_id` cookie を読み、なければ `process.env.DEFAULT_USER_ID` を返します
- `switch-user.ts` は user switch 時に cookie を設定する Server Action です
- `backend-client.ts` は `getUserId()` を使い、全 backend requests に `X-User-ID` header を付与します
- `SessionProvider` が app を wrap し、components に session context を提供します

## Robot Camera Configuration

MJPEG stream settings は robot ごとの `robot_config` field に設定します。

```json
{
  "host": "192.168.1.101",
  "port": 9090,
  "cameras": [
    { "namespace": "camera_0", "name": "Front Camera" },
    { "namespace": "camera_1", "name": "Top Camera" }
  ]
}
```

## 環境変数

`frontend/.env` に設定します (`.env.sample` から copy)。

| Variable | Default | 説明 |
|----------|---------|------|
| `NEXT_PUBLIC_API_BASE_URL` | `/web/api` | client-side API calls の base URL |
| `NEXT_PUBLIC_API_TIMEOUT` | `30000` | request timeout (ms) |
| `BACKEND_API_URL` | `http://localhost:8000` | server-side requests 用 backend URL |
| `DEFAULT_USER_ID` | (required) | X-User-ID header に使う user UUID |

## 主な規約

- base path は `/web` です (`next.config.ts` で設定)
- `next-themes` による dark mode をサポートします
- bookmark 可能な filters のため、URL state は `nuqs` で管理します
- i18n は `react-i18next` を使います (en/ja)
