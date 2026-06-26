# Backend Status Policy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Centralize episode and episode sub-task status policy in domain models and prevent generic update flows from bypassing lifecycle transitions.

**Architecture:** Domain models expose lifecycle and status classification helpers. Usecases call those helpers instead of direct enum comparisons where behavior is policy-driven. OpenAPI and generated files remain unchanged.

**Tech Stack:** Go backend, colocated `*_test.go`, Docker-based `make be-*` commands.

---

## Task 1: Episode Lifecycle Guard

**Files:**
- Modify: `backend/internal/domain/model/episode.go`
- Modify: `backend/internal/domain/model/episode_test.go`
- Modify: `backend/internal/usecase/episode.go`
- Modify: `backend/internal/usecase/episode_test.go`

- [ ] Add RED tests proving direct status changes through generic update are rejected except no-op same-status updates.
- [ ] Add domain helpers `IsTerminal`, `IsSuccessfulCompletion`, and `CanSetStatusFromUpdate`.
- [ ] Update `usecase/episode.go` so `Update` rejects lifecycle status changes and preserves `Start`, `Finish`, `Cancel` as transition entry points.
- [ ] Run focused domain/usecase tests.

## Task 2: Episode Sub-Task Status Helpers

**Files:**
- Modify: `backend/internal/domain/model/episode_sub_task.go`
- Modify: `backend/internal/domain/model/episode_sub_task_test.go`
- Modify: `backend/internal/domain/model/episode_sub_task_execution.go`
- Modify: `backend/internal/domain/model/episode_sub_task_execution_test.go`

- [ ] Add RED tests for terminal, workflow-resolved, and successful-completion classification.
- [ ] Add helper methods on status value types and model structs.
- [ ] Keep existing transition behavior unchanged.
- [ ] Run focused domain tests.

## Task 3: Replace Policy Comparisons In Usecases

**Files:**
- Modify usecase files found by `rg "EpisodeStatusCompleted|EpisodeStatusCancel|SubTaskCollectionStatus|ExecutionStatus" backend/internal/usecase`.

- [ ] Replace direct comparisons only where the comparison represents terminal/completion policy.
- [ ] Leave persistence filters and API enum mapping unchanged.
- [ ] Run focused usecase tests.

## Task 4: Verification And Commit

- [ ] Run `make be-fmt`.
- [ ] Run `make be-test`.
- [ ] Run `make be-lint`.
- [ ] Run `git diff --check`.
- [ ] Commit with `refactor(backend): centralize episode status policy`.

## Self-Review

- Covers the design spec items for episode lifecycle guard and sub-task completion policy.
- Excludes robot status split, which is already complete.
- Excludes HTTP/SSE handler and query/usecase structure cleanup, which will be later PRs.
