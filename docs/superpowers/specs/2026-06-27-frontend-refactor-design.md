# Frontend Refactor Design

## Goal

Improve frontend maintainability without changing product behavior or replacing the current folder architecture. The existing `app / features / shared / lib/api` structure is modern enough for this repository, so the refactor should tighten boundaries, unify repeated API patterns, and split high-churn feature code incrementally.

## Current Findings

- `shared` currently imports concrete feature code in a few places, most notably layout rendering depending on `features/robots`.
- Feature hooks repeat `fetch`, `URLSearchParams`, `response.ok` checks, JSON parsing, and schema validation.
- Query key factories generally follow a useful pattern, but invalidation scope differs by feature.
- `frontend/src/lib/api/backend-client.ts` contains generated-schema type aliases, raw fetch helpers, and many endpoint-specific server wrappers in one large file.
- Large feature components are concentrated in `robots`, `tasks`, and `episodes`.
- Status display policy is spread across shared utilities and feature-specific badge components.
- `docs/frontend-architecture.md` explains the intended structure but does not yet document the stricter dependency rules.

## Target Boundaries

- `app`: route segments, pages, layouts, and Next route handlers only.
- `app/api`: BFF/proxy route handlers. They may call `lib/api` server helpers but should not hold business display logic.
- `features`: business feature UI, hooks, schemas, and feature-specific helpers.
- `shared`: reusable UI, hooks, and utilities that do not import feature modules.
- `lib/api`: API clients, request helpers, generated-code access, server-side backend wrappers, and response helpers.
- `lib/api/generated`: generated files only; never edit directly.

Allowed dependency direction:

```txt
app -> features -> shared
app/api -> lib/api
features -> lib/api, shared
shared -> lib/api only for generated/shared types when unavoidable
lib/api -> generated, auth
```

`shared -> features` should be removed. Cross-feature imports are allowed only when the imported module is intentionally exported as a small public feature API; otherwise shared UI or query primitives should be extracted.

## Refactor Approach

### 1. Boundary Cleanup

Move feature-specific layout rendering out of `shared/components/layout-renderer.tsx`. The shared renderer should render generic registered layout components and delegate camera rendering through the same registry mechanism. Robot-specific camera resolution and `CameraView` usage should live under `features/robots`.

### 2. Client API Hook Utilities

Introduce small client-side helpers under `frontend/src/lib/api/` for:

- building query strings from scalar and array values
- fetching JSON from `/web/api/*`
- parsing responses through Zod schemas

Then migrate representative hooks first, followed by related feature hooks in batches. The helper should not replace TanStack Query or generated schemas; it should only remove repeated boilerplate.

### 3. Server API Wrapper Split

Split `backend-client.ts` by responsibility while preserving exports:

- core fetch/error helpers
- generated-schema API types
- endpoint groups such as tasks, robots, episodes, users, locations, fleet, api keys

Keep the existing import path stable initially by re-exporting from `backend-client.ts`. This makes the change safe and reviewable.

### 4. Query Key and Invalidation Consistency

Normalize query key factories and mutation invalidation in high-change features:

- list queries invalidate list keys
- detail updates invalidate both detail and relevant lists
- aggregate/dashboard queries have explicit keys and invalidation call sites

Do this after client API helpers exist, so hook files can be simplified while being touched.

### 5. Status Display Policy

Centralize frontend status metadata where it is shared across features. Feature badges may remain in features, but labels, terminal/successful status groupings, and colors should be derived from shared status config. Robot connection display must distinguish persisted robot readiness from live Redis/SSE connection state.

### 6. Feature File Splits

Split only files that are already being touched or are clear hotspots:

- robots: fleet stats table, teleoperation registration, robot forms
- tasks: task detail page, task forms, import/export dialogs
- episodes: create/edit forms, export dialog, detail page

Prefer extracting focused presentational components, form sections, column helpers, and pure formatting utilities. Avoid moving whole domains at once.

### 7. Documentation

Update `docs/frontend-architecture.md` and `docs/ja/frontend-architecture.md` after the implementation reflects the new rules. Document dependency direction, generated-code rule, API helper responsibilities, and feature split guidance.

## Testing And Verification

Each PR should run at least:

- `make fe-lint`
- `make fe-fmt` or `make fe-ci` depending on touched files
- `git diff --check`

For structural changes that may affect runtime rendering, run `make fe-ci`. No generated files should change unless OpenAPI is intentionally updated and regeneration is part of the PR.

## PR Sequence

1. Remove `shared -> features` layout dependency.
2. Add client API helper utilities and migrate one or two representative features.
3. Split server-side backend API wrapper while preserving exports.
4. Normalize query keys/invalidation across robots, tasks, and episodes.
5. Centralize status display policy.
6. Split largest feature files in focused batches.
7. Update frontend architecture docs.
