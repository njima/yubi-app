# Frontend Refactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Improve frontend maintainability by tightening `app / features / shared / lib/api` boundaries, unifying repeated API patterns, and splitting high-churn feature code without changing product behavior.

**Architecture:** Keep the existing feature-based Next.js structure. Move concrete feature dependencies out of `shared`, add small API helpers before migrating hooks, preserve public import paths during server wrapper splits, and update docs only after code reflects the new rules.

**Tech Stack:** Next.js 16 App Router, TypeScript 5, TanStack Query, Zod, generated Zodios/OpenAPI client, Tailwind CSS, React Hook Form.

---

## File Map

- `frontend/src/shared/components/layout-renderer.tsx`: keep generic layout rendering only.
- `frontend/src/shared/lib/layout-registry.tsx`: keep component registry and context types.
- `frontend/src/features/robots/lib/register-teleop-components.tsx`: register robot-specific camera and teleoperation components.
- `frontend/src/lib/api/client-fetch.ts`: create browser-side JSON fetch and schema parse helpers.
- `frontend/src/lib/api/query-string.ts`: create shared query string builder.
- `frontend/src/lib/api/backend-client/`: split server-side backend wrappers by endpoint group.
- `frontend/src/lib/api/backend-client.ts`: preserve existing export path by re-exporting split modules.
- `frontend/src/features/**/hooks/*.ts`: migrate repeated fetch/query parsing in batches.
- `frontend/src/shared/lib/status-constants.ts`, `frontend/src/shared/lib/status-utils.ts`, `frontend/src/shared/hooks/use-status-labels.ts`: centralize frontend status display policy.
- `frontend/src/features/robots/components/*.tsx`, `frontend/src/features/tasks/components/*.tsx`, `frontend/src/features/episodes/components/*.tsx`: split hotspots only when scoped.
- `docs/frontend-architecture.md`, `docs/ja/frontend-architecture.md`: update final architecture guidance.

## PR 1: Remove Shared-To-Feature Layout Dependency

**Files:**
- Modify: `frontend/src/shared/components/layout-renderer.tsx`
- Modify: `frontend/src/shared/lib/layout-registry.tsx`
- Modify: `frontend/src/features/robots/lib/register-teleop-components.tsx`

- [ ] **Step 1: Add a generic layout component registration path for camera items**

Extend `LayoutContext` only if needed, and make `layout-renderer.tsx` look up `"camera"` through `getLayoutComponent()` instead of importing `CameraView`.

- [ ] **Step 2: Move camera resolution/rendering into robots feature**

Register `"camera"` in `registerTeleopComponents()` and keep robot-specific `CameraView` props there.

- [ ] **Step 3: Verify dependency direction**

Run:

```bash
rg "@/features" frontend/src/shared frontend/src/lib
```

Expected: no concrete feature imports from `shared` or `lib` except allowed generated type references if still necessary.

- [ ] **Step 4: Run checks**

Run:

```bash
make fe-lint
make fe-ci
git diff --check
```

- [ ] **Step 5: Commit and create PR**

Use:

```bash
git commit -m "refactor(frontend): decouple shared layout renderer"
```

## PR 2: Add Client API Helpers And Migrate Representative Hooks

**Files:**
- Create: `frontend/src/lib/api/query-string.ts`
- Create: `frontend/src/lib/api/client-fetch.ts`
- Modify: `frontend/src/features/tasks/hooks/use-tasks-query.ts`
- Modify: `frontend/src/features/robots/hooks/use-robots-query.ts`
- Modify: `frontend/src/features/episodes/hooks/use-episodes-query.ts`

- [ ] **Step 1: Add query string helper**

Create `buildQueryString(params)` supporting strings, numbers, booleans, arrays, `null`, and `undefined`.

- [ ] **Step 2: Add client JSON helper**

Create `fetchJson(path)` and `fetchAndParse(path, schema, errorLabel)` for `/web/api/*` calls.

- [ ] **Step 3: Migrate list/detail query hooks for tasks, robots, and episodes**

Replace manual `URLSearchParams`, `fetch`, `response.ok`, and schema parsing with the helpers. Preserve existing normalization logic such as robot `robot_config: null -> undefined`.

- [ ] **Step 4: Run checks**

Run:

```bash
make fe-lint
make fe-ci
git diff --check
```

- [ ] **Step 5: Commit and create PR**

Use:

```bash
git commit -m "refactor(frontend): share client api fetch helpers"
```

## PR 3: Split Server Backend API Wrapper

**Files:**
- Create directory: `frontend/src/lib/api/backend-client/`
- Create: `frontend/src/lib/api/backend-client/core.ts`
- Create: `frontend/src/lib/api/backend-client/types.ts`
- Create endpoint-group files under `frontend/src/lib/api/backend-client/`
- Modify: `frontend/src/lib/api/backend-client.ts`
- Modify imports only if circular exports require it.

- [ ] **Step 1: Move `BackendApiError`, `fetchBackendRaw`, and generic JSON fetch into `core.ts`**

Keep behavior identical, including 401 cookie clearing and redirect.

- [ ] **Step 2: Move generated schema aliases into `types.ts`**

Export the same public type names as the current `backend-client.ts`.

- [ ] **Step 3: Move endpoint wrappers by domain**

Split tasks, robots, episodes, users, locations/sites/organizations, fleet/reports, api keys, and subtasks into separate files.

- [ ] **Step 4: Preserve public import path**

Make `frontend/src/lib/api/backend-client.ts` re-export from the split modules so `app/api` routes keep working.

- [ ] **Step 5: Run checks**

Run:

```bash
make fe-lint
make fe-ci
git diff --check
```

- [ ] **Step 6: Commit and create PR**

Use:

```bash
git commit -m "refactor(frontend): split backend api wrappers"
```

## PR 4: Normalize Query Keys And Invalidation

**Files:**
- Modify high-change hooks under:
  - `frontend/src/features/robots/hooks/`
  - `frontend/src/features/tasks/hooks/`
  - `frontend/src/features/episodes/hooks/`

- [ ] **Step 1: Audit query key factories**

Ensure each major feature has `all`, `lists`, `list`, `details`, and `detail` when applicable.

- [ ] **Step 2: Normalize mutation invalidation**

List-affecting creates/deletes/imports invalidate list keys and aggregate keys. Detail updates invalidate the changed detail and relevant lists.

- [ ] **Step 3: Run checks**

Run:

```bash
make fe-lint
make fe-ci
git diff --check
```

- [ ] **Step 4: Commit and create PR**

Use:

```bash
git commit -m "refactor(frontend): normalize query invalidation"
```

## PR 5: Centralize Status Display Policy

**Files:**
- Modify: `frontend/src/shared/lib/status-constants.ts`
- Modify: `frontend/src/shared/lib/status-utils.ts`
- Modify: `frontend/src/shared/hooks/use-status-labels.ts`
- Modify status badge components under robots, tasks, and episodes.

- [ ] **Step 1: Define shared status display metadata**

Centralize labels, color variants, terminal/success groupings, and robot connection display distinction.

- [ ] **Step 2: Update badges to consume shared metadata**

Keep feature-specific badge components, but remove duplicated status mapping.

- [ ] **Step 3: Run checks**

Run:

```bash
make fe-lint
make fe-ci
git diff --check
```

- [ ] **Step 4: Commit and create PR**

Use:

```bash
git commit -m "refactor(frontend): centralize status display policy"
```

## PR 6: Split Feature Hotspots

**Files:**
- Modify selected hotspot files after current diffs reveal the smallest safe splits:
  - `frontend/src/features/robots/components/fleet-stats-table.tsx`
  - `frontend/src/features/robots/lib/register-teleop-components.tsx`
  - `frontend/src/features/tasks/components/detail/task-detail-page.tsx`
  - `frontend/src/features/tasks/components/create-task-form.tsx`
  - `frontend/src/features/episodes/components/create-episode-form.tsx`

- [ ] **Step 1: Extract pure helpers and small presentational sections**

Only extract code with clear names and no behavior changes.

- [ ] **Step 2: Keep form state in parent forms**

Do not move `react-hook-form` ownership unless a section has a stable prop contract.

- [ ] **Step 3: Run checks**

Run:

```bash
make fe-lint
make fe-ci
git diff --check
```

- [ ] **Step 4: Commit and create PR**

Use:

```bash
git commit -m "refactor(frontend): split feature hotspots"
```

## PR 7: Update Frontend Architecture Docs

**Files:**
- Modify: `docs/frontend-architecture.md`
- Modify: `docs/ja/frontend-architecture.md`

- [ ] **Step 1: Document final dependency rules**

Add explicit allowed dependency direction and generated-file rules.

- [ ] **Step 2: Document API helper and feature split conventions**

Explain when to use `lib/api` helpers, when to keep code in features, and when to extract to shared.

- [ ] **Step 3: Run checks**

Run:

```bash
make fe-fmt
git diff --check
```

- [ ] **Step 4: Commit and create PR**

Use:

```bash
git commit -m "docs(frontend): document refactor boundaries"
```

## Self-Review

- The plan covers the spec's dependency boundary, API helper, server wrapper split, query invalidation, status policy, feature split, and documentation requirements.
- No generated files are implementation targets.
- Each PR can be reviewed and merged independently.
- Behavior changes are out of scope unless they fix an existing inconsistency discovered by tests or lint.
