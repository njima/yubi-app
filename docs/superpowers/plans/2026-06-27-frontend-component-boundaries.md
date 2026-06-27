# Frontend Component Boundaries Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move shared UI and app shell components into clearer frontend component boundaries.

**Architecture:** Keep `app/**` focused on route entrypoints. Put shadcn-style primitives under `components/ui`, app shell components under `components/layout`, and feature-owned composition under `features/*/components`.

**Tech Stack:** Next.js App Router, TypeScript, React, shadcn-style UI primitives.

---

### Task 1: Move Shared UI Primitives

**Files:**
- Move: `frontend/src/shared/ui/*` -> `frontend/src/components/ui/*`
- Modify: all imports from `@/shared/ui/*` to `@/components/ui/*`

- [x] Create `frontend/src/components/ui`.
- [x] Move every UI primitive from `frontend/src/shared/ui`.
- [x] Update all `@/shared/ui/*` imports.

### Task 2: Move Layout-Level Components

**Files:**
- Move: `frontend/src/shared/components/layout-renderer.tsx` -> `frontend/src/components/layout/layout-renderer.tsx`
- Move: `frontend/src/app/(app)/_components/{top-nav,nav-item,user-menu,language-switcher}.tsx` -> `frontend/src/components/layout/`
- Modify: imports in `frontend/src/app/(app)/layout.tsx` and teleoperation views.

- [x] Create `frontend/src/components/layout`.
- [x] Move app shell components and layout renderer.
- [x] Update imports to `@/components/layout/*`.

### Task 3: Move Feature-Owned Components

**Files:**
- Move: `frontend/src/app/(app)/_components/switch-user-dialog.tsx` -> `frontend/src/features/users/components/switch-user-dialog.tsx`
- Move: `frontend/src/app/(app)/episodes/_components/export-menu.tsx` -> `frontend/src/features/episodes/components/export-menu.tsx`
- Move: `frontend/src/shared/components/parameterized-name.tsx` -> `frontend/src/features/tasks/components/parameterized-name.tsx`
- Modify: feature public exports and imports.

- [x] Move user switching, episode export, and parameterized name components to owning features.
- [x] Export moved feature components through their feature `index.ts` files where external features use them.
- [x] Remove empty `_components` and `shared/components` directories.

### Task 4: Update Documentation and Verify

**Files:**
- Modify: `docs/frontend-architecture.md`
- Modify: `docs/ja/frontend-architecture.md`

- [x] Document `components/ui`, `components/layout`, and minimal `app/**/_components` usage.
- [x] Run `make fe-ci`.
- [ ] Commit with `refactor(frontend): clarify component boundaries`.
