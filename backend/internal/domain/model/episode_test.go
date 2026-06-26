package model

import (
	"testing"
	"time"
)

func TestInitEpisode(t *testing.T) {
	tests := []struct {
		name           string
		organizationID string
		taskID         string
		locationID     string
		robotID        string
		userID         string
		wantErr        bool
	}{
		{
			name:           "success with valid inputs",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			taskID:         "550e8400-e29b-41d4-a716-446655440002",
			locationID:     "550e8400-e29b-41d4-a716-446655440003",
			robotID:        "550e8400-e29b-41d4-a716-446655440004",
			userID:         "550e8400-e29b-41d4-a716-446655440005",
			wantErr:        false,
		},
		{
			name:           "error when organization_id is empty",
			organizationID: "",
			taskID:         "550e8400-e29b-41d4-a716-446655440002",
			locationID:     "550e8400-e29b-41d4-a716-446655440003",
			robotID:        "550e8400-e29b-41d4-a716-446655440004",
			userID:         "550e8400-e29b-41d4-a716-446655440005",
			wantErr:        true,
		},
		{
			name:           "error when task_id is empty",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			taskID:         "",
			locationID:     "550e8400-e29b-41d4-a716-446655440003",
			robotID:        "550e8400-e29b-41d4-a716-446655440004",
			userID:         "550e8400-e29b-41d4-a716-446655440005",
			wantErr:        true,
		},
		{
			name:           "error when location_id is empty",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			taskID:         "550e8400-e29b-41d4-a716-446655440002",
			locationID:     "",
			robotID:        "550e8400-e29b-41d4-a716-446655440004",
			userID:         "550e8400-e29b-41d4-a716-446655440005",
			wantErr:        true,
		},
		{
			name:           "error when robot_id is empty",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			taskID:         "550e8400-e29b-41d4-a716-446655440002",
			locationID:     "550e8400-e29b-41d4-a716-446655440003",
			robotID:        "",
			userID:         "550e8400-e29b-41d4-a716-446655440005",
			wantErr:        true,
		},
		{
			name:           "error when user_id is empty",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			taskID:         "550e8400-e29b-41d4-a716-446655440002",
			locationID:     "550e8400-e29b-41d4-a716-446655440003",
			robotID:        "550e8400-e29b-41d4-a716-446655440004",
			userID:         "",
			wantErr:        true,
		},
		{
			name:           "error when multiple fields are empty",
			organizationID: "",
			taskID:         "",
			locationID:     "",
			robotID:        "",
			userID:         "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitEpisode(tt.organizationID, tt.taskID, tt.locationID, tt.robotID, tt.userID, nil)

			if tt.wantErr {
				if err == nil {
					t.Errorf("InitEpisode() error = nil, wantErr %v", tt.wantErr)
				}
				if got.IDNatural != "" {
					t.Errorf("InitEpisode() IDNatural = %v, want empty", got.IDNatural)
				}
				return
			}

			if err != nil {
				t.Errorf("InitEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.IDNatural == "" {
				t.Errorf("InitEpisode() IDNatural is empty")
			}
			if got.OrganizationID != tt.organizationID {
				t.Errorf("InitEpisode() OrganizationID = %v, want %v", got.OrganizationID, tt.organizationID)
			}
			if got.TaskID != tt.taskID {
				t.Errorf("InitEpisode() TaskID = %v, want %v", got.TaskID, tt.taskID)
			}
			if got.LocationID != tt.locationID {
				t.Errorf("InitEpisode() LocationID = %v, want %v", got.LocationID, tt.locationID)
			}
			if got.RobotID != tt.robotID {
				t.Errorf("InitEpisode() RobotID = %v, want %v", got.RobotID, tt.robotID)
			}
			if got.UserID != tt.userID {
				t.Errorf("InitEpisode() UserID = %v, want %v", got.UserID, tt.userID)
			}
			if got.Status != EpisodeStatusReady {
				t.Errorf("InitEpisode() Status = %v, want %v", got.Status, EpisodeStatusReady)
			}
			if got.CreatedAt.IsZero() {
				t.Errorf("InitEpisode() CreatedAt is zero")
			}
		})
	}
}

func TestNewEpisode(t *testing.T) {
	now := time.Now()
	startedAt := now.Add(-time.Hour)
	finishedAt := now
	updatedAt := now
	errorDetails := "some error"

	recordedByID := "550e8400-e29b-41d4-a716-446655440099"

	tests := []struct {
		name          string
		id            int64
		idNatural     string
		taskID        string
		taskVersionID string
		locationID    string
		robotID       string
		userID        string
		errorDetails  *string
		startedAt     *time.Time
		finishedAt    *time.Time
		status        EpisodeStatus
		createdAt     time.Time
		updatedAt     *time.Time
		recordedByID  *string
	}{
		{
			name:          "create with all fields",
			id:            1,
			idNatural:     "550e8400-e29b-41d4-a716-446655440000",
			taskID:        "550e8400-e29b-41d4-a716-446655440001",
			taskVersionID: "550e8400-e29b-41d4-a716-446655440002",
			locationID:    "550e8400-e29b-41d4-a716-446655440003",
			robotID:       "550e8400-e29b-41d4-a716-446655440004",
			userID:        "550e8400-e29b-41d4-a716-446655440005",
			errorDetails:  &errorDetails,
			startedAt:     &startedAt,
			finishedAt:    &finishedAt,
			status:        EpisodeStatusCompleted,
			createdAt:     now,
			updatedAt:     &updatedAt,
			recordedByID:  &recordedByID,
		},
		{
			name:          "create with nil optional fields",
			id:            2,
			idNatural:     "550e8400-e29b-41d4-a716-446655440006",
			taskID:        "550e8400-e29b-41d4-a716-446655440007",
			taskVersionID: "550e8400-e29b-41d4-a716-446655440008",
			locationID:    "550e8400-e29b-41d4-a716-446655440009",
			robotID:       "550e8400-e29b-41d4-a716-446655440010",
			userID:        "550e8400-e29b-41d4-a716-446655440011",
			errorDetails:  nil,
			startedAt:     nil,
			finishedAt:    nil,
			status:        EpisodeStatusReady,
			createdAt:     now,
			updatedAt:     nil,
			recordedByID:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEpisode(
				tt.id,
				tt.idNatural,
				tt.taskID,
				tt.taskVersionID,
				tt.locationID,
				tt.robotID,
				tt.userID,
				tt.errorDetails,
				tt.startedAt,
				tt.finishedAt,
				tt.status,
				tt.createdAt,
				tt.updatedAt,
				tt.recordedByID,
			)

			if got.ID != tt.id {
				t.Errorf("NewEpisode() ID = %v, want %v", got.ID, tt.id)
			}
			if got.IDNatural != tt.idNatural {
				t.Errorf("NewEpisode() IDNatural = %v, want %v", got.IDNatural, tt.idNatural)
			}
			if got.TaskID != tt.taskID {
				t.Errorf("NewEpisode() TaskID = %v, want %v", got.TaskID, tt.taskID)
			}
			if got.TaskVersionID != tt.taskVersionID {
				t.Errorf("NewEpisode() TaskVersionID = %v, want %v", got.TaskVersionID, tt.taskVersionID)
			}
			if got.LocationID != tt.locationID {
				t.Errorf("NewEpisode() LocationID = %v, want %v", got.LocationID, tt.locationID)
			}
			if got.RobotID != tt.robotID {
				t.Errorf("NewEpisode() RobotID = %v, want %v", got.RobotID, tt.robotID)
			}
			if got.UserID != tt.userID {
				t.Errorf("NewEpisode() UserID = %v, want %v", got.UserID, tt.userID)
			}
			if got.Status != tt.status {
				t.Errorf("NewEpisode() Status = %v, want %v", got.Status, tt.status)
			}
			if got.CreatedAt != tt.createdAt {
				t.Errorf("NewEpisode() CreatedAt = %v, want %v", got.CreatedAt, tt.createdAt)
			}
			if got.RecordedByID != tt.recordedByID {
				if tt.recordedByID == nil {
					t.Errorf("NewEpisode() RecordedByID = %v, want nil", got.RecordedByID)
				} else if got.RecordedByID == nil {
					t.Errorf("NewEpisode() RecordedByID is nil, want %v", *tt.recordedByID)
				} else if *got.RecordedByID != *tt.recordedByID {
					t.Errorf("NewEpisode() RecordedByID = %v, want %v", *got.RecordedByID, *tt.recordedByID)
				}
			}
		})
	}
}

func TestEpisode_validate(t *testing.T) {
	tests := []struct {
		name    string
		episode Episode
		wantErr bool
	}{
		{
			name: "valid with all required fields",
			episode: Episode{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				TaskID:         "550e8400-e29b-41d4-a716-446655440002",
				LocationID:     "550e8400-e29b-41d4-a716-446655440003",
				RobotID:        "550e8400-e29b-41d4-a716-446655440004",
				UserID:         "550e8400-e29b-41d4-a716-446655440005",
			},
			wantErr: false,
		},
		{
			name: "error when id_natural is empty",
			episode: Episode{
				IDNatural:      "",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				TaskID:         "550e8400-e29b-41d4-a716-446655440002",
				LocationID:     "550e8400-e29b-41d4-a716-446655440003",
				RobotID:        "550e8400-e29b-41d4-a716-446655440004",
				UserID:         "550e8400-e29b-41d4-a716-446655440005",
			},
			wantErr: true,
		},
		{
			name: "error when organization_id is empty",
			episode: Episode{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "",
				TaskID:         "550e8400-e29b-41d4-a716-446655440002",
				LocationID:     "550e8400-e29b-41d4-a716-446655440003",
				RobotID:        "550e8400-e29b-41d4-a716-446655440004",
				UserID:         "550e8400-e29b-41d4-a716-446655440005",
			},
			wantErr: true,
		},
		{
			name: "error when task_id is empty",
			episode: Episode{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				TaskID:         "",
				LocationID:     "550e8400-e29b-41d4-a716-446655440003",
				RobotID:        "550e8400-e29b-41d4-a716-446655440004",
				UserID:         "550e8400-e29b-41d4-a716-446655440005",
			},
			wantErr: true,
		},
		{
			name: "error when location_id is empty",
			episode: Episode{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				TaskID:         "550e8400-e29b-41d4-a716-446655440002",
				LocationID:     "",
				RobotID:        "550e8400-e29b-41d4-a716-446655440004",
				UserID:         "550e8400-e29b-41d4-a716-446655440005",
			},
			wantErr: true,
		},
		{
			name: "error when robot_id is empty",
			episode: Episode{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				TaskID:         "550e8400-e29b-41d4-a716-446655440002",
				LocationID:     "550e8400-e29b-41d4-a716-446655440003",
				RobotID:        "",
				UserID:         "550e8400-e29b-41d4-a716-446655440005",
			},
			wantErr: true,
		},
		{
			name: "error when user_id is empty",
			episode: Episode{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				TaskID:         "550e8400-e29b-41d4-a716-446655440002",
				LocationID:     "550e8400-e29b-41d4-a716-446655440003",
				RobotID:        "550e8400-e29b-41d4-a716-446655440004",
				UserID:         "",
			},
			wantErr: true,
		},
		{
			name: "error when multiple fields are empty",
			episode: Episode{
				IDNatural:      "",
				OrganizationID: "",
				TaskID:         "",
				LocationID:     "",
				RobotID:        "",
				UserID:         "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.episode.validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Episode.validate() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Episode.validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func newValidEpisode() Episode {
	return Episode{
		IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
		OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
		TaskID:         "550e8400-e29b-41d4-a716-446655440002",
		LocationID:     "550e8400-e29b-41d4-a716-446655440003",
		RobotID:        "550e8400-e29b-41d4-a716-446655440004",
		UserID:         "550e8400-e29b-41d4-a716-446655440005",
	}
}

func TestEpisode_SetTaskVersionID(t *testing.T) {
	tests := []struct {
		name          string
		taskVersionID string
		wantErr       bool
	}{
		{
			name:          "success with valid task_version_id",
			taskVersionID: "550e8400-e29b-41d4-a716-446655440006",
			wantErr:       false,
		},
		{
			name:          "success with empty task_version_id",
			taskVersionID: "",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newValidEpisode()
			err := e.SetTaskVersionID(tt.taskVersionID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Episode.SetTaskVersionID() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Episode.SetTaskVersionID() error = %v, wantErr %v", err, tt.wantErr)
				}
				if e.TaskVersionID != tt.taskVersionID {
					t.Errorf("Episode.TaskVersionID = %v, want %v", e.TaskVersionID, tt.taskVersionID)
				}
			}
		})
	}
}

func TestEpisode_SetStartedAt(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		startedAt time.Time
		wantErr   bool
	}{
		{
			name:      "success with valid started_at",
			startedAt: now,
			wantErr:   false,
		},
		{
			name:      "success with past time",
			startedAt: now.Add(-24 * time.Hour),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newValidEpisode()
			err := e.SetStartedAt(tt.startedAt)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Episode.SetStartedAt() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Episode.SetStartedAt() error = %v, wantErr %v", err, tt.wantErr)
				}
				if e.StartedAt == nil {
					t.Errorf("Episode.StartedAt is nil")
				} else if *e.StartedAt != tt.startedAt {
					t.Errorf("Episode.StartedAt = %v, want %v", *e.StartedAt, tt.startedAt)
				}
			}
		})
	}
}

func TestEpisode_SetFinishedAt(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		finishedAt time.Time
		wantErr    bool
	}{
		{
			name:       "success with valid finished_at",
			finishedAt: now,
			wantErr:    false,
		},
		{
			name:       "success with future time",
			finishedAt: now.Add(24 * time.Hour),
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newValidEpisode()
			err := e.SetFinishedAt(tt.finishedAt)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Episode.SetFinishedAt() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Episode.SetFinishedAt() error = %v, wantErr %v", err, tt.wantErr)
				}
				if e.FinishedAt == nil {
					t.Errorf("Episode.FinishedAt is nil")
				} else if *e.FinishedAt != tt.finishedAt {
					t.Errorf("Episode.FinishedAt = %v, want %v", *e.FinishedAt, tt.finishedAt)
				}
			}
		})
	}
}

func TestEpisode_SetStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  EpisodeStatus
		wantErr bool
	}{
		{
			name:    "success with completed status",
			status:  EpisodeStatusCompleted,
			wantErr: false,
		},
		{
			name:    "success with ready status",
			status:  EpisodeStatusReady,
			wantErr: false,
		},
		{
			name:    "success with recording status",
			status:  EpisodeStatusRecording,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newValidEpisode()
			err := e.SetStatus(tt.status)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Episode.SetStatus() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Episode.SetStatus() error = %v, wantErr %v", err, tt.wantErr)
				}
				if e.Status != tt.status {
					t.Errorf("Episode.Status = %v, want %v", e.Status, tt.status)
				}
			}
		})
	}
}

func TestEpisodeStatusPolicy(t *testing.T) {
	tests := []struct {
		name                 string
		status               EpisodeStatus
		wantTerminal         bool
		wantSuccessfulFinish bool
	}{
		{
			name:                 "ready is not terminal",
			status:               EpisodeStatusReady,
			wantTerminal:         false,
			wantSuccessfulFinish: false,
		},
		{
			name:                 "recording is not terminal",
			status:               EpisodeStatusRecording,
			wantTerminal:         false,
			wantSuccessfulFinish: false,
		},
		{
			name:                 "cancel is terminal but not successful",
			status:               EpisodeStatusCancel,
			wantTerminal:         true,
			wantSuccessfulFinish: false,
		},
		{
			name:                 "completed is terminal and successful",
			status:               EpisodeStatusCompleted,
			wantTerminal:         true,
			wantSuccessfulFinish: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsTerminal(); got != tt.wantTerminal {
				t.Errorf("EpisodeStatus.IsTerminal() = %v, want %v", got, tt.wantTerminal)
			}
			if got := tt.status.IsSuccessfulCompletion(); got != tt.wantSuccessfulFinish {
				t.Errorf("EpisodeStatus.IsSuccessfulCompletion() = %v, want %v", got, tt.wantSuccessfulFinish)
			}

			ep := newEpisodeWithStatus(tt.status)
			if got := ep.IsTerminal(); got != tt.wantTerminal {
				t.Errorf("Episode.IsTerminal() = %v, want %v", got, tt.wantTerminal)
			}
			if got := ep.IsSuccessfulCompletion(); got != tt.wantSuccessfulFinish {
				t.Errorf("Episode.IsSuccessfulCompletion() = %v, want %v", got, tt.wantSuccessfulFinish)
			}
		})
	}
}

func TestEpisode_CanApplyStatusUpdate(t *testing.T) {
	t.Run("allows no-op status update", func(t *testing.T) {
		ep := newEpisodeWithStatus(EpisodeStatusRecording)

		if err := ep.CanApplyStatusUpdate(EpisodeStatusRecording); err != nil {
			t.Fatalf("Episode.CanApplyStatusUpdate() error = %v, want nil", err)
		}
	})

	t.Run("rejects lifecycle status change", func(t *testing.T) {
		ep := newEpisodeWithStatus(EpisodeStatusReady)

		if err := ep.CanApplyStatusUpdate(EpisodeStatusRecording); err == nil {
			t.Fatal("Episode.CanApplyStatusUpdate() error = nil, want error")
		}
	})
}

func TestEpisode_SetErrorDetails(t *testing.T) {
	tests := []struct {
		name         string
		errorDetails string
		wantErr      bool
	}{
		{
			name:         "success with error details",
			errorDetails: "Connection timeout",
			wantErr:      false,
		},
		{
			name:         "success with empty error details",
			errorDetails: "",
			wantErr:      false,
		},
		{
			name:         "success with long error details",
			errorDetails: "This is a very long error message that contains detailed information about what went wrong during the episode execution.",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newValidEpisode()
			err := e.SetErrorDetails(tt.errorDetails)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Episode.SetErrorDetails() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Episode.SetErrorDetails() error = %v, wantErr %v", err, tt.wantErr)
				}
				if e.ErrorDetails == nil {
					t.Errorf("Episode.ErrorDetails is nil")
				} else if *e.ErrorDetails != tt.errorDetails {
					t.Errorf("Episode.ErrorDetails = %v, want %v", *e.ErrorDetails, tt.errorDetails)
				}
			}
		})
	}
}

func TestEpisode_SetRecordedByID(t *testing.T) {
	tests := []struct {
		name         string
		recordedByID string
		wantErr      bool
	}{
		{
			name:         "success with valid user id",
			recordedByID: "550e8400-e29b-41d4-a716-446655440099",
			wantErr:      false,
		},
		{
			name:         "success with empty user id",
			recordedByID: "",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newValidEpisode()
			err := e.SetRecordedByID(tt.recordedByID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Episode.SetRecordedByID() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Episode.SetRecordedByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				if e.RecordedByID == nil {
					t.Errorf("Episode.RecordedByID is nil")
				} else if *e.RecordedByID != tt.recordedByID {
					t.Errorf("Episode.RecordedByID = %v, want %v", *e.RecordedByID, tt.recordedByID)
				}
			}
		})
	}
}

func newEpisodeWithStatus(status EpisodeStatus) Episode {
	e := newValidEpisode()
	e.Status = status
	return e
}

func TestEpisode_Start(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		initialStatus EpisodeStatus
		occurredAt    time.Time
		wantErr       bool
		wantStatus    EpisodeStatus
	}{
		{
			name:          "success: Ready → Recording",
			initialStatus: EpisodeStatusReady,
			occurredAt:    now,
			wantErr:       false,
			wantStatus:    EpisodeStatusRecording,
		},
		{
			name:          "error: Recording → cannot start",
			initialStatus: EpisodeStatusRecording,
			occurredAt:    now,
			wantErr:       true,
		},
		{
			name:          "error: Completed → cannot start",
			initialStatus: EpisodeStatusCompleted,
			occurredAt:    now,
			wantErr:       true,
		},
		{
			name:          "error: Cancel → cannot start",
			initialStatus: EpisodeStatusCancel,
			occurredAt:    now,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEpisodeWithStatus(tt.initialStatus)
			err := e.Start(tt.occurredAt)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Episode.Start() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Episode.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if e.Status != tt.wantStatus {
				t.Errorf("Episode.Start() Status = %v, want %v", e.Status, tt.wantStatus)
			}
			if e.StartedAt == nil {
				t.Errorf("Episode.Start() StartedAt is nil")
			} else if !e.StartedAt.Equal(tt.occurredAt) {
				t.Errorf("Episode.Start() StartedAt = %v, want %v", *e.StartedAt, tt.occurredAt)
			}
		})
	}
}

func TestEpisode_Finish(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		initialStatus EpisodeStatus
		occurredAt    time.Time
		wantErr       bool
		wantStatus    EpisodeStatus
	}{
		{
			name:          "success: Recording → Completed",
			initialStatus: EpisodeStatusRecording,
			occurredAt:    now,
			wantErr:       false,
			wantStatus:    EpisodeStatusCompleted,
		},
		{
			name:          "error: Ready → cannot finish",
			initialStatus: EpisodeStatusReady,
			occurredAt:    now,
			wantErr:       true,
		},
		{
			name:          "error: Completed → cannot finish",
			initialStatus: EpisodeStatusCompleted,
			occurredAt:    now,
			wantErr:       true,
		},
		{
			name:          "error: Cancel → cannot finish",
			initialStatus: EpisodeStatusCancel,
			occurredAt:    now,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEpisodeWithStatus(tt.initialStatus)
			err := e.Finish(tt.occurredAt)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Episode.Finish() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Episode.Finish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if e.Status != tt.wantStatus {
				t.Errorf("Episode.Finish() Status = %v, want %v", e.Status, tt.wantStatus)
			}
			if e.FinishedAt == nil {
				t.Errorf("Episode.Finish() FinishedAt is nil")
			} else if !e.FinishedAt.Equal(tt.occurredAt) {
				t.Errorf("Episode.Finish() FinishedAt = %v, want %v", *e.FinishedAt, tt.occurredAt)
			}
		})
	}
}

func TestEpisode_Cancel(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus EpisodeStatus
		wantErr       bool
		wantStatus    EpisodeStatus
	}{
		{
			name:          "success: Recording → Cancel",
			initialStatus: EpisodeStatusRecording,
			wantErr:       false,
			wantStatus:    EpisodeStatusCancel,
		},
		{
			name:          "error: Ready → cannot cancel",
			initialStatus: EpisodeStatusReady,
			wantErr:       true,
		},
		{
			name:          "error: Completed → cannot cancel",
			initialStatus: EpisodeStatusCompleted,
			wantErr:       true,
		},
		{
			name:          "error: Cancel → cannot cancel again",
			initialStatus: EpisodeStatusCancel,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEpisodeWithStatus(tt.initialStatus)
			err := e.Cancel()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Episode.Cancel() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Episode.Cancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if e.Status != tt.wantStatus {
				t.Errorf("Episode.Cancel() Status = %v, want %v", e.Status, tt.wantStatus)
			}
		})
	}
}
