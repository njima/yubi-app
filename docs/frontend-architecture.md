# Frontend Architecture

## Overview

The frontend is a Next.js 16 application for managing robot fleets, tasks, and teleoperation episodes. It uses type-safe API clients generated from the OpenAPI specification.

## Tech Stack

| Category | Technology |
|----------|------------|
| Framework | Next.js 16 (App Router) |
| Language | TypeScript 5 |
| UI Components | Radix UI, Tailwind CSS 4 |
| State Management | TanStack Query (React Query) |
| Forms | React Hook Form + Zod |
| API Client | Zodios (auto-generated from OpenAPI) |
| URL State | nuqs |
| Notifications | Sonner |

## Directory Structure

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router
│   │   ├── (app)/              # Main application routes
│   │   │   ├── dashboard/      # Dashboard page
│   │   │   ├── episodes/       # Episodes management
│   │   │   ├── robots/         # Robots management
│   │   │   ├── tasks/          # Tasks management
│   │   │   ├── users/          # Users management
│   │   │   └── profile/        # User profile
│   │   └── api/                # API route handlers (proxy to backend)
│   │
│   ├── components/             # Domain-neutral React components
│   │   ├── layout/             # App shell and reusable layout components
│   │   └── ui/                 # shadcn-style UI primitives
│   │
│   ├── features/               # Feature modules (domain-driven)
│   │   ├── api-keys/           # API key management
│   │   ├── dashboard/          # Dashboard composition
│   │   ├── episodes/           # Episode CRUD, hooks, components
│   │   ├── locations/          # Location management
│   │   ├── organizations/      # Organization management
│   │   ├── reporting/          # Reporting and export workflows
│   │   ├── robots/             # Robot CRUD, status, camera viewer
│   │   ├── tasks/              # Task management
│   │   └── users/              # User CRUD, roles
│   │
│   ├── lib/                    # Core libraries
│   │   ├── api/                # API client configuration
│   │   │   ├── client.ts       # Zodios client with interceptors
│   │   │   ├── client-fetch.ts # Browser-side fetch/schema helpers
│   │   │   ├── query-string.ts # Query string builder
│   │   │   ├── backend-client.ts # Server-side re-export facade
│   │   │   ├── backend-client/ # Server-side backend wrappers by endpoint group
│   │   │   ├── config.ts       # API URL configuration
│   │   │   └── generated/      # Auto-generated Zodios client
│   │   └── auth/               # Session management
│   │       └── session.ts      # getUserId(), getUserSession()
│   │
│   └── shared/                 # Shared code across features
│       ├── hooks/              # Shared hooks (status labels, etc.)
│       ├── lib/                # Utilities (date, status constants)
│       └── providers/          # React providers (QueryProvider)
│
├── public/                     # Static assets
├── next.config.ts              # Next.js configuration
├── Dockerfile                  # Production image
├── Dockerfile.dev              # Development image
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
backend-client.ts (adds X-User-ID header)
    │
    ▼
Backend API (http://backend:8000)
```

### Server Components vs Client Components

- **Server Components / Route Handlers** (`backend-client.ts`): Fetch data directly from backend with `X-User-ID` header. The facade re-exports endpoint-specific modules under `lib/api/backend-client/`.
- **Client Components**: Use TanStack Query hooks that call `/web/api/*` routes. Prefer `client-fetch.ts` and `query-string.ts` for repeated fetch, schema parsing, and query string building.
- **SSE Streams**: Proxied through `sse-proxy.ts` with `X-User-ID` header for real-time updates.

### Dependency Rules

Keep dependency direction explicit:

```text
app -> features -> shared
app -> components
app/api -> lib/api
features -> lib/api, shared
features -> components
lib/api -> generated, auth
```

`shared` must not import concrete feature modules. If shared rendering needs feature-specific behavior, register it from the feature layer through a registry, as teleoperation layout components do. Cross-feature imports should go through a small public API from the target feature.

Use `app/**/page.tsx` and `app/**/layout.tsx` as thin route entrypoints. Prefer `features/*/components` for feature-owned page composition and `components/layout` for app shell components such as navigation. Avoid adding `app/**/_components` unless the component is tiny route glue that cannot reasonably belong to a feature or shared component area.

### Component Placement

- `components/ui`: shadcn-style primitives such as `Button`, `Dialog`, `Table`, and `DropdownMenu`.
- `components/layout`: domain-neutral layout and app shell components such as `TopNav` and `LayoutRenderer`.
- `features/*/components`: feature-owned UI, including page-level composition such as list pages and export menus. Feature modules do not have to map 1:1 to pages; capability modules such as `reporting` are allowed.
- `app/**`: route entrypoints, route groups, layouts, and API routes. Keep route-local components rare and small.
- `shared/*`: cross-cutting hooks, providers, and utilities. Do not place React UI primitives here.

### API Code Generation

The API client is auto-generated from the OpenAPI spec:

```bash
make fe-generate-api
```

This generates `src/lib/api/generated/api.ts` containing the following.

- Zod schemas for all request/response types
- Zodios API client with typed endpoints

**Never edit generated files directly.** Modify `openapi/openapi.yaml` and regenerate.

## Feature Module Pattern

Each feature in `src/features/` is self-contained:

```
features/episodes/
├── components/          # Feature-specific React components
│   ├── episode-columns.tsx
│   ├── episode-detail.tsx
│   └── create-episode-dialog.tsx
├── hooks/               # Data fetching and mutation hooks
│   ├── use-episodes-query.ts
│   └── use-create-episode-mutation.ts
├── schemas/             # Zod validation schemas (for forms)
└── index.ts             # Public API (re-exports)
```

Feature files should stay focused. Large feature-specific renderers can be split into local files such as `teleop-layout-components.tsx` while keeping registration or page wiring separate. Keep form ownership in the parent form unless a section has a stable prop contract.

## Status Display Policy

Status values live in `shared/lib/status-constants.ts`. Display metadata such as badge color, label key, terminal state, and successful completion lives in `shared/lib/status-display.ts`. Feature badge components can remain feature-specific, but they should consume shared display metadata instead of duplicating status maps.

## Authentication

The frontend identifies the active user via a cookie, falling back to the `DEFAULT_USER_ID` environment variable.

- `session.ts` provides `getUserId()` which reads the `active_user_id` cookie first, then falls back to `process.env.DEFAULT_USER_ID`
- `switch-user.ts` is a Server Action that sets the cookie when a user switches accounts
- `backend-client.ts` attaches `X-User-ID` header to all backend requests using `getUserId()`
- `SessionProvider` wraps the app and provides session context to components

## Robot Camera Configuration

MJPEG stream settings are configured per-robot via the `robot_config` field:

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

## Environment Variables

Configure in `frontend/.env` (copy from `.env.sample`).

| Variable | Default | Description |
|----------|---------|-------------|
| `NEXT_PUBLIC_API_BASE_URL` | `/web/api` | Base URL for client-side API calls |
| `NEXT_PUBLIC_API_TIMEOUT` | `30000` | Request timeout (ms) |
| `BACKEND_API_URL` | `http://localhost:8000` | Backend URL for server-side requests |
| `DEFAULT_USER_ID` | (required) | User UUID for X-User-ID header |

## Key Conventions

- Base path is `/web` (configured in `next.config.ts`)
- Dark mode supported via `next-themes`
- URL state managed with `nuqs` for bookmarkable filters
- i18n via `react-i18next` (en/ja)
