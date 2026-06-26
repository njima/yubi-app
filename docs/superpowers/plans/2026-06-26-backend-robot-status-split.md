# Backend Robot Status Split Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Split persistent robot operation status from Redis-derived online/offline display state without changing OpenAPI or generated code.

**Architecture:** Keep `model.RobotStatus` as the API-compatible status enum, but make domain helpers distinguish persistent operation states from connection-only states. Usecases should resolve online/offline only for returned copies and pass heartbeat liveness explicitly when starting teleoperation.

**Tech Stack:** Go backend, Clean Architecture layers, Bun persistence, Docker-based `make be-*` commands.

---

## File Structure

- Modify `backend/internal/domain/model/robot.go`: add status policy helpers, make heartbeat resolution return a status instead of mutating the entity, and make teleoperation start require heartbeat liveness.
- Modify `backend/internal/domain/model/robot_test.go`: update/add RED tests for non-mutating display status and teleoperation readiness.
- Modify `backend/internal/usecase/robot.go`: assign resolved display status only after persistence operations.
- Modify `backend/internal/usecase/episode.go`: stop mutating robot status to online before `StartTeleoperation`; pass heartbeat liveness explicitly.
- Modify `backend/internal/infra/persistence/query.go`: keep online/offline filters scoped to operation-ready robots and legacy online rows.
- Modify `backend/internal/infra/persistence/robot.go`: mirror the same filter behavior for robot type listing.

### Task 1: Domain Status Policy

**Files:**
- Modify: `backend/internal/domain/model/robot_test.go`
- Modify: `backend/internal/domain/model/robot.go`

- [ ] **Step 1: Write failing tests**

Add or replace tests near `TestRobot_StartTeleoperation`:

```go
func TestRobot_ResolvedStatusDoesNotMutateOperationStatus(t *testing.T) {
	tests := []struct {
		name           string
		initialStatus  RobotStatus
		heartbeatAlive bool
		wantResolved   RobotStatus
		wantStored     RobotStatus
	}{
		{name: "ready with heartbeat appears online", initialStatus: RobotStatusReady, heartbeatAlive: true, wantResolved: RobotStatusOnline, wantStored: RobotStatusReady},
		{name: "ready without heartbeat appears offline", initialStatus: RobotStatusReady, heartbeatAlive: false, wantResolved: RobotStatusOffline, wantStored: RobotStatusReady},
		{name: "busy ignores heartbeat", initialStatus: RobotStatusBusy, heartbeatAlive: true, wantResolved: RobotStatusBusy, wantStored: RobotStatusBusy},
		{name: "faulted ignores heartbeat", initialStatus: RobotStatusFaulted, heartbeatAlive: true, wantResolved: RobotStatusFaulted, wantStored: RobotStatusFaulted},
		{name: "maintenance ignores heartbeat", initialStatus: RobotStatusMaintenance, heartbeatAlive: true, wantResolved: RobotStatusMaintenance, wantStored: RobotStatusMaintenance},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRobotWithStatus(tt.initialStatus)

			got := r.ResolvedStatus(tt.heartbeatAlive)

			if got != tt.wantResolved {
				t.Fatalf("Robot.ResolvedStatus() = %v, want %v", got, tt.wantResolved)
			}
			if r.Status != tt.wantStored {
				t.Fatalf("Robot.Status mutated to %v, want %v", r.Status, tt.wantStored)
			}
		})
	}
}

func TestRobot_StartTeleoperationRequiresReadyAndHeartbeat(t *testing.T) {
	episodeID := "550e8400-e29b-41d4-a716-446655440010"
	userID := "550e8400-e29b-41d4-a716-446655440011"

	tests := []struct {
		name           string
		initialStatus  RobotStatus
		heartbeatAlive bool
		wantErr        bool
		wantStatus     RobotStatus
	}{
		{name: "ready with heartbeat starts", initialStatus: RobotStatusReady, heartbeatAlive: true, wantStatus: RobotStatusBusy},
		{name: "ready without heartbeat fails", initialStatus: RobotStatusReady, heartbeatAlive: false, wantErr: true, wantStatus: RobotStatusReady},
		{name: "busy with heartbeat fails", initialStatus: RobotStatusBusy, heartbeatAlive: true, wantErr: true, wantStatus: RobotStatusBusy},
		{name: "legacy online status fails as non-persistent operation state", initialStatus: RobotStatusOnline, heartbeatAlive: true, wantErr: true, wantStatus: RobotStatusOnline},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRobotWithStatus(tt.initialStatus)

			err := r.StartTeleoperation(episodeID, userID, tt.heartbeatAlive)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("Robot.StartTeleoperation() error = nil, want error")
				}
				if r.Status != tt.wantStatus {
					t.Fatalf("Robot.Status = %v, want %v", r.Status, tt.wantStatus)
				}
				return
			}
			if err != nil {
				t.Fatalf("Robot.StartTeleoperation() error = %v", err)
			}
			if r.Status != tt.wantStatus {
				t.Fatalf("Robot.Status = %v, want %v", r.Status, tt.wantStatus)
			}
			if r.ActiveEpisodeID == nil || *r.ActiveEpisodeID != episodeID {
				t.Fatalf("Robot.ActiveEpisodeID = %v, want %v", r.ActiveEpisodeID, episodeID)
			}
			if r.ActiveUserID == nil || *r.ActiveUserID != userID {
				t.Fatalf("Robot.ActiveUserID = %v, want %v", r.ActiveUserID, userID)
			}
		})
	}
}
```

- [ ] **Step 2: Verify RED**

Run:

```bash
make be-test TEST_ARGS=./internal/domain/model
```

Expected: FAIL because `ResolvedStatus` currently returns no value and `StartTeleoperation` does not accept heartbeat liveness.

- [ ] **Step 3: Implement minimal domain changes**

Change `robot.go` helpers to:

```go
func (s RobotStatus) IsPersistentOperationStatus() bool {
	switch s {
	case RobotStatusReady, RobotStatusBusy, RobotStatusFaulted, RobotStatusMaintenance:
		return true
	default:
		return false
	}
}

func (s RobotStatus) IsConnectionOnlyStatus() bool {
	return s == RobotStatusOnline || s == RobotStatusOffline
}

func (r *Robot) CanStartTeleoperation(heartbeatAlive bool) error {
	if r.Status != RobotStatusReady {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "robot operation status must be Ready to start teleoperation, current status: %d", r.Status),
		)
	}
	if !heartbeatAlive {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "robot must be online to start teleoperation"),
		)
	}
	return nil
}

func (r *Robot) StartTeleoperation(episodeID, userID string, heartbeatAlive bool) error {
	if err := r.CanStartTeleoperation(heartbeatAlive); err != nil {
		return err
	}
	r.Status = RobotStatusBusy
	r.ActiveEpisodeID = &episodeID
	r.ActiveUserID = &userID
	return nil
}

func (r Robot) ResolvedStatus(heartbeatAlive bool) RobotStatus {
	if r.Status == RobotStatusReady || r.Status == RobotStatusOnline {
		if heartbeatAlive {
			return RobotStatusOnline
		}
		return RobotStatusOffline
	}
	return r.Status
}
```

- [ ] **Step 4: Verify GREEN for domain**

Run:

```bash
make be-test TEST_ARGS=./internal/domain/model
```

Expected: PASS for domain tests after updating old call sites in tests.

### Task 2: Usecase Display Resolution

**Files:**
- Modify: `backend/internal/usecase/robot.go`
- Modify: `backend/internal/usecase/episode.go`

- [ ] **Step 1: Write failing usecase-facing test if an existing fake is available**

Search for robot usecase tests:

```bash
rg -n "NewRobot\\(|ResolvedStatus|StartTeleoperation" backend/internal/usecase backend/internal -g '*_test.go'
```

If a robot usecase test harness exists, add a test asserting that a ready robot returned from `GetByID` appears online with Redis heartbeat while the repository update path never saves online/offline. If no harness exists, rely on the domain RED test and compile failure from changed method signatures.

- [ ] **Step 2: Update robot usecase display mapping**

In `backend/internal/usecase/robot.go`, replace each call:

```go
rob.ResolvedStatus(status != nil)
```

with:

```go
rob.Status = rob.ResolvedStatus(status != nil)
```

For `Update`, apply this only to the returned `urob` after repository update, not to the entity before saving.

- [ ] **Step 3: Update episode start**

In `backend/internal/usecase/episode.go`, replace:

```go
robot.ResolvedStatus(robotStatus != nil)
...
if err := robot.StartTeleoperation(input.EpisodeID, activeUserID); err != nil {
	return err
}
```

with:

```go
heartbeatAlive := robotStatus != nil

if err := robot.StartTeleoperation(input.EpisodeID, activeUserID, heartbeatAlive); err != nil {
	return err
}
```

- [ ] **Step 4: Verify package compile/tests**

Run:

```bash
make be-test TEST_ARGS=./internal/usecase
```

Expected: PASS or reveal tests that need signature updates only.

### Task 3: Persistence Filter Compatibility

**Files:**
- Modify: `backend/internal/infra/persistence/query.go`
- Modify: `backend/internal/infra/persistence/robot.go`

- [ ] **Step 1: Keep filters defensive for legacy data**

Keep online/offline list filters limited to robots that can be connection-resolved. Use this list in both persistence filter locations:

```go
[]model.RobotStatus{
	model.RobotStatusReady,
	model.RobotStatusOnline,
}
```

Do not include `RobotStatusOffline` because new writes must not create it and offline legacy rows should not be considered ready for teleoperation.

- [ ] **Step 2: Verify no new direct writes of online/offline**

Run:

```bash
rg -n "SetStatus\\(model\\.RobotStatus(Online|Offline)|Status:\\s*model\\.RobotStatus(Online|Offline)|RobotStatusOnline|RobotStatusOffline" backend/internal --glob '!gen/**'
```

Expected: generated/controller filter references may remain; no production update path should save online/offline.

### Task 4: Full Verification And Commit

**Files:**
- Modify: all files above

- [ ] **Step 1: Format**

Run:

```bash
make be-fmt
```

Expected: command exits 0.

- [ ] **Step 2: Backend tests**

Run:

```bash
make be-test
```

Expected: command exits 0.

- [ ] **Step 3: Backend lint**

Run:

```bash
make be-lint
```

Expected: command exits 0.

- [ ] **Step 4: Diff hygiene**

Run:

```bash
git diff --check
git status --short
```

Expected: no whitespace errors; only intended files changed plus existing unrelated untracked spec files.

- [ ] **Step 5: Commit**

Stage only intended files:

```bash
git add backend/internal/domain/model/robot.go backend/internal/domain/model/robot_test.go backend/internal/usecase/robot.go backend/internal/usecase/episode.go backend/internal/infra/persistence/query.go backend/internal/infra/persistence/robot.go docs/superpowers/plans/2026-06-26-backend-robot-status-split.md
git commit -m "refactor(backend): split robot connection status from operation state"
```

## Self-Review

- Spec coverage: robot operation/connection split is covered. Episode and sub-task changes remain in the design document for later PRs.
- Placeholder scan: no `TBD`, `TODO`, or unspecified implementation steps are used.
- Type consistency: `ResolvedStatus` returns `RobotStatus`; `StartTeleoperation` and `CanStartTeleoperation` accept `heartbeatAlive bool`.
