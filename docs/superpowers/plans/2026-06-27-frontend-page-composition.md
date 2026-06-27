# Frontend Page Composition Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Keep Next.js route files thin by moving page-level composition into feature components and renaming report-specific feature code to `reporting`.

**Architecture:** `app/**/page.tsx` imports a feature-owned page component. Business page composition lives under `features/*/components`, while shared UI primitives remain in `components/ui`.

**Tech Stack:** Next.js App Router, TypeScript, React, TanStack Query.

---

### Task 1: Move List and Shell Pages Into Features

**Files:**
- Move page logic from `frontend/src/app/(app)/*/page.tsx` into `frontend/src/features/*/components/*-page.tsx`.
- Keep `frontend/src/app/(app)/*/page.tsx` as thin route entrypoints.

- [x] Move api key, dashboard, episode list, location list, profile, robot list, task list, and user list page composition.
- [x] Export moved page components through feature public APIs where useful.

### Task 2: Rename Reports Feature to Reporting

**Files:**
- Move: `frontend/src/features/reporting` -> `frontend/src/features/reporting`
- Modify imports from `@/features/reporting` to `@/features/reporting`.

- [x] Rename the feature directory.
- [x] Update imports and docs.

### Task 3: Verify and Commit

- [x] Run `git diff --check`.
- [x] Run `make fe-ci`.
- [ ] Commit with `refactor(frontend): move page composition into features`.
