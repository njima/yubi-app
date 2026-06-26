package handler

import (
	"testing"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func TestRobotStatusStreamDetail(t *testing.T) {
	gate := &model.GateConditionStatus{
		GateLevel: 2,
		Groups: map[string]model.GateGroupStatus{
			"battery": {
				Level:   2,
				Settled: true,
				Conditions: []model.GateCondition{
					{Name: "battery_ok", Passed: true, Reason: "ok", Escalation: 0},
				},
			},
		},
	}

	got := robotStatusStreamDetail(usecase.RobotDeviceStatusDetail{
		Battery:        usecase.RobotBatteryStatus{Pct: 87},
		Connection:     usecase.RobotConnectionStatus{QualityPct: 92},
		UptimeSec:      12.6,
		GateConditions: gate,
	})

	if got.BatteryPct != 87 {
		t.Fatalf("BatteryPct = %v, want 87", got.BatteryPct)
	}
	if got.ConnectionPct != 92 {
		t.Fatalf("ConnectionPct = %v, want 92", got.ConnectionPct)
	}
	if got.UptimeSec != 13 {
		t.Fatalf("UptimeSec = %v, want 13", got.UptimeSec)
	}
	if got.GateConditions == nil {
		t.Fatal("GateConditions = nil, want value")
	}
	if got.GateConditions.GateLevel != 2 {
		t.Fatalf("GateLevel = %v, want 2", got.GateConditions.GateLevel)
	}
}

func TestBuildEpisodeStreamResponse(t *testing.T) {
	startedAt := time.Date(2026, 6, 27, 10, 0, 0, 0, time.UTC)
	finishedAt := startedAt.Add(5 * time.Minute)
	recordedBy := "operator-1"
	avg := 0.8
	params := map[string]string{"mode": "fast"}
	taskDescription := "Inspect area"
	displayName := "Inspection v2"
	ep := &model.Episode{
		IDNatural:       "ep-1",
		LocationID:      "loc-1",
		UserID:          "user-1",
		RobotID:         "robot-1",
		Status:          model.EpisodeStatusCompleted,
		TaskID:          "task-1",
		TaskVersionID:   "tv-1",
		StartedAt:       &startedAt,
		FinishedAt:      &finishedAt,
		RecordedByID:    &recordedBy,
		AverageGrade:    &avg,
		GradeCount:      3,
		ParameterValues: params,
		CreatedAt:       startedAt.Add(-time.Minute),
	}
	task := &model.Task{Name: "Inspection", Description: &taskDescription}
	subtasks := []openapi.EpisodeSubTask{{Id: "subtask-1"}}

	got := buildEpisodeStreamResponse(ep, subtasks, task, &displayName)

	if got.Id != ep.IDNatural {
		t.Fatalf("Id = %v, want %v", got.Id, ep.IDNatural)
	}
	if got.TaskName == nil || *got.TaskName != "Inspection" {
		t.Fatalf("TaskName = %v, want Inspection", got.TaskName)
	}
	if got.TaskDescription == nil || *got.TaskDescription != taskDescription {
		t.Fatalf("TaskDescription = %v, want %v", got.TaskDescription, taskDescription)
	}
	if got.TaskVersionDisplayName == nil || *got.TaskVersionDisplayName != displayName {
		t.Fatalf("TaskVersionDisplayName = %v, want %v", got.TaskVersionDisplayName, displayName)
	}
	if got.ParameterValues == nil || (*got.ParameterValues)["mode"] != "fast" {
		t.Fatalf("ParameterValues = %v, want mode=fast", got.ParameterValues)
	}
	if got.Subtasks == nil || len(*got.Subtasks) != 1 || (*got.Subtasks)[0].Id != "subtask-1" {
		t.Fatalf("Subtasks = %v, want subtask-1", got.Subtasks)
	}
}
