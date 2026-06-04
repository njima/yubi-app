# Robot API Guide

This document describes how robots communicate with the Yubi App platform to execute teleoperation episodes.

## Overview

Robots interact with the system through the **Robot Device API** (`/api/robot/*`). Two authentication methods are supported.

```
Robot
  │  X-API-Key header (recommended)
  │  — or —
  │  X-User-ID + X-Robot-ID headers (fallback)
  ▼
Backend API (port 8000)
  ├─ Auth middleware: validates API key or user/robot headers
  ├─ Robot Device endpoints: /api/robot/*
  ▼
PostgreSQL
```

## Authentication

### Method 1: API Key (Recommended)

Use an API key issued from the Web UI. This is the recommended approach for production deployments.

| Header | Required | Description |
|--------|----------|-------------|
| `X-API-Key` | Yes | API key issued from the admin UI |

The backend hashes the key, looks it up in the database, and validates that it is active (not revoked, not expired). The user and robot bound to the key are used to set the request context.

```bash
curl -H "X-API-Key: <your-api-key>" http://localhost:8000/api/robot/me
```

**How to issue an API key**

1. Log in as an Admin user
2. Go to **API Keys** in the navigation
3. Click **"Create API Key"**
4. Select the robot to bind the key to
5. Copy the raw key (shown only once)

### Method 2: X-User-ID + X-Robot-ID Headers (Fallback)

If no `X-API-Key` header is present, the backend falls back to header-based authentication.

| Header | Required | Description |
|--------|----------|-------------|
| `X-User-ID` | Yes | UUID of the user operating the robot |
| `X-Robot-ID` | Yes | UUID of the robot |

The backend validates that the user and robot both exist in the database and belong to the same organization.

```bash
curl -H "X-User-ID: <user-uuid>" -H "X-Robot-ID: <robot-uuid>" \
  http://localhost:8000/api/robot/me
```

### Error Responses

```json
{"error": "Invalid API key"}
{"error": "X-API-Key or X-User-ID header is required"}
{"error": "User not found"}
{"error": "Robot not found"}
{"error": "User and robot do not belong to the same organization"}
```

### Setup

All examples in this guide use API key authentication.

```bash
export BASE_URL="http://localhost:8000"
export API_KEY="<your-api-key>"
```

> **Tip**: Issue an API key from the Web UI (Admin → API Keys → Create). Alternatively, use header-based auth with UUIDs from `backend/internal/database/seeder/initial_data.sql`.

## Prerequisites

Before a robot can execute episodes, the following must be set up via the Web UI.

1. **Register the robot** — Create a robot entry at `/web/robots`
2. **Create a task with subtasks** — Define the work to be done at `/web/tasks`
3. **Create episodes** — Assign episodes (task + robot + location) at `/web/episodes`

See the [User Guide](user-guide.md) for step-by-step instructions.

## Episode Execution Flow

This is the typical sequence a robot follows to execute an episode.

```
[Preparation - Web UI]
1. Create episodes (status: Ready)

[Execution - Robot API]
2. List available episodes        GET  /api/robot/episodes
3. Get episode details            GET  /api/robot/episodes/{id}
4. Start episode                  POST /api/robot/episodes/{id}/start
5. For each subtask:
   a. Create execution            POST /api/robot/episodes/{id}/subtasks/{subtaskId}/executions
   b. Start execution             POST /api/robot/episodes/{id}/subtasks/{subtaskId}/executions/{execId}/start
   c. (Robot performs the task and collects data)
   d. Finish execution            POST /api/robot/episodes/{id}/subtasks/{subtaskId}/executions/{execId}/finish
   e. Complete subtask            POST /api/robot/episodes/{id}/subtasks/{subtaskId}/complete
6. Finish episode                 POST /api/robot/episodes/{id}/finish

[Review - Web UI]
7. View results at /web/episodes/{id}
```

## API Reference

### Get Robot Info

Verify that the robot's credentials are working and see its current configuration.

```bash
curl -H "X-API-Key: ${API_KEY}" \
  "${BASE_URL}/api/robot/me"
```

**Response (200)** — Returns a `Robot` object.

```json
{
  "id": "c2f8e62b-ea23-4a50-8660-d707e4d5c2bc",
  "name": "Sample Robot",
  "robot_type": "yubi",
  "status": 0,
  "organization_id": "7bfbe942-5fd6-4525-ac13-0356147c202b",
  "organization_name": "Sample Organization",
  "location_id": "91154897-df4b-4b39-8c4c-b48daf4a3b37",
  "location_name": "Sample Location",
  "leader_status": null,
  "battery_level": 85,
  "last_heartbeat_at": "2026-01-01T10:00:00Z",
  "offline_reason": null,
  "robot_config": {
    "host": "localhost",
    "port": 9090,
    "cameras": [{"namespace": "camera_0", "name": "Front Camera"}]
  },
  "active_episode_id": null,
  "active_user_id": null,
  "active_operator": null
  // Other fields: consecutive_fault_days, leader_consecutive_fault_days,
  // leader_fault_started_at — see openapi.yaml Robot schema for full definition
}
```

### List Episodes

Fetch episodes assigned to this robot. Only episodes with `Ready` status are returned — these are waiting for the robot to pick up.

```bash
curl -H "X-API-Key: ${API_KEY}" \
  "${BASE_URL}/api/robot/episodes"
```

**Response (200)** — Returns an array of `Episode` objects.

```json
[
  {
    "id": "358637c3-bf03-4408-845f-f9b189ec767f",
    "user_id": "69fad3df-d73f-45e1-9fb4-df52bd4857b0",
    "robot_id": "c2f8e62b-ea23-4a50-8660-d707e4d5c2bc",
    "location_id": "91154897-df4b-4b39-8c4c-b48daf4a3b37",
    "task_id": "6013935a-ab9c-4bd8-b59d-49958f516d47",
    "task_name": "Sample Task",
    "task_version_id": "5437b101-6d9d-495f-a4e8-45420eb10d99",
    "status": 0,
    "created_at": "2026-01-01T09:00:00Z"
    // Other fields: task_description, task_version_display_name, started_at,
    // ended_at, recorded_by, parameter_values, subtasks, average_grade,
    // grade_count — see openapi.yaml Episode schema for full definition
  }
]
```

`status` values: `0` = Ready, `1` = Recording, `2` = Cancel, `3` = Completed

### Get Episode Details

Fetch full episode information including subtasks. Use this to know what subtasks the robot needs to execute.

```bash
curl -H "X-API-Key: ${API_KEY}" \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}"
```

**Response (200)** — Returns an `Episode` object with subtasks.

```json
{
  "id": "358637c3-bf03-4408-845f-f9b189ec767f",
  "user_id": "69fad3df-d73f-45e1-9fb4-df52bd4857b0",
  "robot_id": "c2f8e62b-ea23-4a50-8660-d707e4d5c2bc",
  "location_id": "91154897-df4b-4b39-8c4c-b48daf4a3b37",
  "task_id": "6013935a-ab9c-4bd8-b59d-49958f516d47",
  "task_name": "Sample Task",
  "task_version_id": "5437b101-6d9d-495f-a4e8-45420eb10d99",
  "status": 0,
  "created_at": "2026-01-01T09:00:00Z",
  "subtasks": [
    {
      "id": "b1c352f7-54c6-4d74-ba6d-aa1352dfaee0",
      "subtask_id": "7065d47f-8de7-4b6d-af34-4aa924dfa98e",
      "name": "Sample SubTask 1",
      "order_index": 0,
      "status": 0,
      "executions": []
    },
    {
      "id": "1cfb5edb-ce00-4c6c-8c58-69cdbd06ad5e",
      "subtask_id": "1d8744a2-cfe0-4cfd-88ea-f38fbc50c640",
      "name": "Sample SubTask 2",
      "order_index": 1,
      "status": 0,
      "executions": []
    }
  ]
  // Other fields: task_description, task_version_display_name, started_at,
  // ended_at, parameter_values, average_grade, grade_count
  // — see openapi.yaml Episode schema for full definition
}
```

The subtask `id` (not `subtask_id`) is what you use in subsequent API calls for subtask operations.

### Start Episode

Tells the system that the robot is beginning to execute this episode. The episode transitions from `Ready` to `Recording`, and the robot status changes from `Online` to `Busy`.

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"occurred_at": "2026-01-01T10:00:00Z"}' \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/start"
```

| Field | Required | Format | Description |
|-------|----------|--------|-------------|
| `occurred_at` | Yes | RFC 3339 | When the episode actually started on the robot |

**Response (200)** — Episode started successfully.

**Error (400)** — Episode is not in `Ready` status or validation failed.

### Create Execution

Before performing a subtask, the robot creates an execution record. A subtask can have multiple executions (retries).

```bash
SUBTASK_ID="b1c352f7-54c6-4d74-ba6d-aa1352dfaee0"  # subtask id from Get Episode Details

curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/subtasks/${SUBTASK_ID}/executions"
```

**Response (201)**

```json
{
  "execution_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

Save the `execution_id` — you need it for the start/finish/cancel calls below.

### Start Execution

Marks the beginning of an execution attempt.

```bash
EXECUTION_ID="a1b2c3d4-e5f6-7890-abcd-ef1234567890"

curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"occurred_at": "2026-01-01T10:05:00Z"}' \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/subtasks/${SUBTASK_ID}/executions/${EXECUTION_ID}/start"
```

| Field | Required | Format | Description |
|-------|----------|--------|-------------|
| `occurred_at` | Yes | RFC 3339 | When the execution actually started |

### Finish Execution

Marks a successful completion of this execution attempt.

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"occurred_at": "2026-01-01T10:10:00Z"}' \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/subtasks/${SUBTASK_ID}/executions/${EXECUTION_ID}/finish"
```

### Cancel Execution

Cancels an in-progress execution (e.g., the robot encountered an obstacle). You can create a new execution to retry.

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/subtasks/${SUBTASK_ID}/executions/${EXECUTION_ID}/cancel"
```

### Complete Subtask

Marks a subtask as completed after at least one execution has been finished. Call this after finishing the final execution for this subtask.

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/subtasks/${SUBTASK_ID}/complete"
```

### Skip Subtask

Skips a subtask entirely without executing it (e.g., not applicable in the current environment).

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/subtasks/${SUBTASK_ID}/skip"
```

### Finish Episode

Marks the episode as completed. The robot status returns from `Busy` to `Online`.

Call this after all subtasks have been completed or skipped.

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"occurred_at": "2026-01-01T11:00:00Z"}' \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/finish"
```

| Field | Required | Format | Description |
|-------|----------|--------|-------------|
| `occurred_at` | Yes | RFC 3339 | When the episode actually finished |

### Cancel Episode

Cancels an in-progress episode. Use this when the robot cannot complete the task (e.g., hardware failure, safety stop). The robot status returns from `Busy` to `Online`.

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/cancel"
```

### Repeat Last Episode

Creates a new episode with the same task, robot, and location as the last completed one. Useful for running the same task repeatedly without going through the Web UI.

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  "${BASE_URL}/api/robot/episodes/repeat-last"
```

**Response (201)** — The newly created episode object (status: Ready). You still need to call Start Episode to begin execution.

**Error (404)** — No previous completed episode found for this robot.

### Update Robot Status

Reports the robot's current status as a heartbeat. The backend uses this to determine whether the robot is online. Send this periodically (e.g., every 10-30 seconds).

The request body is a `RobotStatusUpdateRequest` containing `robot_type`, `reported_at`, and a `status` object with battery, connection, and metrics.

```bash
curl -X PUT \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "robot_type": "yubi",
    "reported_at": "2026-01-01T10:00:00Z",
    "status": {
      "battery": {"pct": 85, "charging": false},
      "connection": {"quality_pct": 100},
      "uptime_sec": 3600.0,
      "metrics": []
    }
  }' \
  "${BASE_URL}/api/robot/status"
```

| Field | Required | Description |
|-------|----------|-------------|
| `robot_type` | Yes | Robot type identifier (`yubi` or `yubi-portable`) |
| `reported_at` | Yes | Timestamp of the status report (RFC 3339) |
| `status` | Yes | `RobotStatusDetail` object (see below) |

**RobotStatusDetail fields**

| Field | Required | Description |
|-------|----------|-------------|
| `battery` | Yes | `{"pct": int, "charging": bool}` — battery percentage |
| `connection` | Yes | `{"quality_pct": int}` — connection quality (0-100) |
| `uptime_sec` | Yes | Robot uptime in seconds |
| `metrics` | No | Array of custom robot metrics |
| `gate_conditions` | No | Gate condition status |

**RobotStatus values** (as shown in the Robot object)

| Value | Status |
|-------|--------|
| 0 | Online |
| 1 | Busy |
| 2 | Offline |
| 3 | Faulted |
| 4 | Maintenance |
| 5 | Ready |

> The robot's status is determined by the backend based on heartbeat data. The `status` field in the request is the detailed status report, not the enum value above.

## Common Error Responses

All endpoints may return these errors.

| Status | Meaning | Example |
|--------|---------|---------|
| 401 | Authentication failed | Missing header, user/robot not found, org mismatch |
| 403 | Insufficient permissions | User role doesn't allow this operation |
| 404 | Resource not found | Episode or subtask ID doesn't exist |
| 409 | Conflict | Episode is not in the expected status for this operation |
| 500 | Internal server error | Unexpected backend failure |

Error response format.

```json
{
  "error": "Human-readable error message"
}
```

## Status Transitions

### Episode Status

```
Ready (created via Web UI)
    │ start
    ▼
Recording (robot executing)
    ├─ finish ──→ Completed
    └─ cancel ──→ Cancelled
```

### Subtask Status

```
Ready
    ├─ create execution + start ──→ InProgress
    │                                  ├─ complete ──→ Completed
    │                                  └─ cancel ───→ Cancelled
    └─ skip ──→ Skipped
```

### Robot Status

```
Offline (initial)
    │ heartbeat received
    ▼
Online (ready)
    │ episode start
    ▼
Busy (executing)
    │ episode finish/cancel
    ▼
Online
    │ heartbeat timeout
    ▼
Offline
```

## Real-time Updates (SSE)

The platform provides Server-Sent Events for real-time monitoring. These are primarily used by the Web UI but can also be consumed by robot-side tools.

| Endpoint | Description |
|----------|-------------|
| `GET /api/robots/{robotId}/status/stream` | Single robot status updates |
| `GET /api/robots/status/stream?robotIds=id1,id2` | Multiple robot status updates |
| `GET /api/episodes/{episodeId}/stream` | Episode progress updates |
| `GET /api/robots/{robotId}/teleop/stream` | Combined teleop stream (status + episode + task) |

These endpoints require `X-User-ID` header (not `X-Robot-ID`).

> These SSE endpoints are registered directly in the server code and are not included in the OpenAPI spec.

## Full OpenAPI Specification

For the complete API reference including all request/response schemas, see [`openapi/openapi.yaml`](../openapi/openapi.yaml).
