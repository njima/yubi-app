# Backend Status Management Refactor Design

## Context

Backend status handling is currently inconsistent across robots, episodes, and episode sub-tasks. Robot online/offline is derived from Redis heartbeat data, but the same enum is also used as a persistent robot status. Episode status can be changed directly through update flows even though start, finish, and cancel already express domain transitions. Sub-task and execution completion rules are scattered as raw enum comparisons.

This refactor keeps OpenAPI and generated code unchanged. It only clarifies domain behavior and moves status policy into domain models/usecases.

## Goals

- Separate persistent robot operation status from Redis-derived connection state at the domain/usecase boundary.
- Prevent unsafe direct episode status mutation from generic update flows.
- Centralize terminal/completion policy for episode sub-tasks and sub-task executions.
- Preserve current API response shapes and avoid editing generated files.
- Add focused backend tests before production changes.

## Non-Goals

- No OpenAPI schema change in this pass.
- No generated code edits.
- No database migration unless tests reveal an unavoidable persistence compatibility issue.
- No frontend behavior changes.

## Design

### Robot Status

Robot persistence should represent operation state only: ready, busy, faulted, and maintenance. Online/offline should be treated as connection state derived from Redis heartbeat data when building responses or checking teleoperation readiness.

The existing API-facing `RobotStatus` values can remain for compatibility, but domain helpers should make the split explicit:

- persistent operation states can be saved to the database;
- connection-only states are display/filter concepts;
- resolving heartbeat state should not mutate the persistent robot entity.

Teleoperation start should require both persistent `Ready` and live heartbeat. A ready robot without heartbeat should appear offline and must not start teleoperation.

### Episode Status

Episode lifecycle changes should go through domain transitions: `Start`, `Finish`, and `Cancel`. Generic episode update should not directly move the episode to another lifecycle status. A no-op status update may remain valid only when the requested status equals the current status.

This keeps lifecycle side effects centralized, including timestamps and robot status coordination.

### Sub-Task Status

Episode sub-task and execution status policy should be expressed with domain helpers instead of repeated raw comparisons.

Definitions:

- sub-task terminal states: completed, skipped, cancelled;
- sub-task workflow-resolved states: completed, skipped;
- sub-task successful completion: completed;
- execution terminal states: finished, cancelled;
- execution successful completion: finished.

Usecases should call these helpers when deciding whether a task is done, repeatable, cancellable, or included in completion metrics.

## Implementation Units

1. Robot operation/connection split.
   - Add RED tests for non-mutating heartbeat resolution and teleoperation readiness.
   - Update domain/usecase/persistence filters to avoid saving online/offline as persistent state.

2. Episode lifecycle guard.
   - Add RED tests for rejecting direct status changes through update.
   - Route lifecycle changes through existing start/finish/cancel flows.

3. Sub-task completion policy.
   - Add RED tests for terminal/resolved/successful helpers.
   - Replace raw status checks in usecases with the new helpers where behavior is equivalent.

## Verification

For each implementation PR:

- run focused Go tests for changed packages;
- run `make be-fmt`;
- run `make be-test`;
- run `make be-lint`;
- run `git diff --check`.

## Risks

Existing database rows may contain legacy robot `Online` or `Offline` values. The implementation should handle these defensively at read/filter boundaries while ensuring new writes use operation states only.
