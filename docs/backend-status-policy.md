# Backend Status Policy

This document defines which backend statuses are stored as domain state and which ones are derived for API display. Keep new status logic in `backend/internal/domain/model`; usecases should call domain helpers instead of duplicating transition rules.

## Robot

Robot status has two meanings:

- **Persistent operation status**: `Ready`, `Busy`, `Faulted`, `Maintenance`
- **Connection-only display status**: `Online`, `Offline`

`Ready` is the stored default. `Online` and `Offline` are derived from Redis heartbeat state when the stored operation status is `Ready`. They should not be used as manual operation updates. If legacy data stores `Online`, display resolution still applies the heartbeat rule.

Teleoperation transitions are:

```text
Ready + heartbeat alive -> Busy
Busy -> Ready
```

`Faulted` and `Maintenance` are manual operation states. `Busy` is controlled by teleoperation lifecycle actions and must not be overwritten by manual status updates.

## Episode

Episode status is lifecycle-owned:

- `Ready`: created but not recording
- `Recording`: active recording
- `Completed`: terminal and successful
- `Cancel`: terminal but not successful

Use `Start`, `Finish`, and `Cancel` for lifecycle transitions. Direct status updates should be no-ops only; changing lifecycle state directly bypasses domain rules.

## Episode Subtask

`EpisodeSubTask.CollectionStatus` tracks collection workflow state:

- `Ready` and `InProgress` are open states.
- `Completed`, `Skipped`, and `Cancelled` are terminal states.
- `Completed` and `Skipped` are workflow-resolved states.
- Only `Completed` is a successful completion.

Use `StartProgress`, `Complete`, `Skip`, and `Cancel` for transitions.

## Episode Subtask Execution

Execution status tracks individual execution attempts:

- `Ready` and `Started` are open states.
- `Finished` and `Cancelled` are terminal and workflow-resolved states.
- Only `Finished` is a successful completion.

Use execution lifecycle helpers instead of setting status fields directly.

## Task

Task progress status is derived from duration:

- `Planning`: actual duration is zero.
- `Doing`: actual duration is greater than zero and less than target duration.
- `Completed`: actual duration is greater than or equal to target duration.

`Canceled` is a separate explicit state and should not be inferred from duration.
