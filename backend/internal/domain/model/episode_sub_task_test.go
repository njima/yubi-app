package model

import (
	"testing"
	"time"
)

func TestInitEpisodeSubTask(t *testing.T) {
	tests := []struct {
		name           string
		organizationID string
		episodeID      string
		subtaskID      string
		wantErr        bool
	}{
		{
			name:           "success with valid inputs",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			episodeID:      "550e8400-e29b-41d4-a716-446655440002",
			subtaskID:      "550e8400-e29b-41d4-a716-446655440003",
			wantErr:        false,
		},
		{
			name:           "error when organization_id is empty",
			organizationID: "",
			episodeID:      "550e8400-e29b-41d4-a716-446655440002",
			subtaskID:      "550e8400-e29b-41d4-a716-446655440003",
			wantErr:        true,
		},
		{
			name:           "error when episode_id is empty",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			episodeID:      "",
			subtaskID:      "550e8400-e29b-41d4-a716-446655440003",
			wantErr:        true,
		},
		{
			name:           "error when subtask_id is empty",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			episodeID:      "550e8400-e29b-41d4-a716-446655440002",
			subtaskID:      "",
			wantErr:        true,
		},
		{
			name:           "error when all fields are empty",
			organizationID: "",
			episodeID:      "",
			subtaskID:      "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitEpisodeSubTask(tt.organizationID, tt.episodeID, tt.subtaskID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("InitEpisodeSubTask() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("InitEpisodeSubTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.IDNatural == "" {
				t.Errorf("InitEpisodeSubTask() IDNatural is empty")
			}
			if got.OrganizationID != tt.organizationID {
				t.Errorf("InitEpisodeSubTask() OrganizationID = %v, want %v", got.OrganizationID, tt.organizationID)
			}
			if got.EpisodeID != tt.episodeID {
				t.Errorf("InitEpisodeSubTask() EpisodeID = %v, want %v", got.EpisodeID, tt.episodeID)
			}
			if got.SubTaskID != tt.subtaskID {
				t.Errorf("InitEpisodeSubTask() SubTaskID = %v, want %v", got.SubTaskID, tt.subtaskID)
			}
			if got.CollectionStatus != SubTaskCollectionStatusReady {
				t.Errorf("InitEpisodeSubTask() CollectionStatus = %v, want %v", got.CollectionStatus, SubTaskCollectionStatusReady)
			}
			if got.CreatedAt.IsZero() {
				t.Errorf("InitEpisodeSubTask() CreatedAt is zero")
			}
		})
	}
}

func TestNewEpisodeSubTask(t *testing.T) {
	now := time.Now()
	updatedAt := now

	tests := []struct {
		name             string
		id               int64
		idNatural        string
		organizationID   string
		episodeID        string
		subtaskID        string
		collectionStatus SubTaskCollectionStatus
		createdAt        time.Time
		updatedAt        *time.Time
	}{
		{
			name:             "create with all fields",
			id:               1,
			idNatural:        "550e8400-e29b-41d4-a716-446655440000",
			organizationID:   "550e8400-e29b-41d4-a716-446655440001",
			episodeID:        "550e8400-e29b-41d4-a716-446655440002",
			subtaskID:        "550e8400-e29b-41d4-a716-446655440003",
			collectionStatus: SubTaskCollectionStatusCompleted,
			createdAt:        now,
			updatedAt:        &updatedAt,
		},
		{
			name:             "create with nil updated_at",
			id:               2,
			idNatural:        "550e8400-e29b-41d4-a716-446655440004",
			organizationID:   "550e8400-e29b-41d4-a716-446655440005",
			episodeID:        "550e8400-e29b-41d4-a716-446655440006",
			subtaskID:        "550e8400-e29b-41d4-a716-446655440007",
			collectionStatus: SubTaskCollectionStatusReady,
			createdAt:        now,
			updatedAt:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEpisodeSubTask(
				tt.id,
				tt.idNatural,
				tt.organizationID,
				tt.episodeID,
				tt.subtaskID,
				tt.collectionStatus,
				tt.createdAt,
				tt.updatedAt,
			)

			if got.ID != tt.id {
				t.Errorf("NewEpisodeSubTask() ID = %v, want %v", got.ID, tt.id)
			}
			if got.IDNatural != tt.idNatural {
				t.Errorf("NewEpisodeSubTask() IDNatural = %v, want %v", got.IDNatural, tt.idNatural)
			}
			if got.OrganizationID != tt.organizationID {
				t.Errorf("NewEpisodeSubTask() OrganizationID = %v, want %v", got.OrganizationID, tt.organizationID)
			}
			if got.EpisodeID != tt.episodeID {
				t.Errorf("NewEpisodeSubTask() EpisodeID = %v, want %v", got.EpisodeID, tt.episodeID)
			}
			if got.SubTaskID != tt.subtaskID {
				t.Errorf("NewEpisodeSubTask() SubTaskID = %v, want %v", got.SubTaskID, tt.subtaskID)
			}
			if got.CollectionStatus != tt.collectionStatus {
				t.Errorf("NewEpisodeSubTask() CollectionStatus = %v, want %v", got.CollectionStatus, tt.collectionStatus)
			}
			if got.CreatedAt != tt.createdAt {
				t.Errorf("NewEpisodeSubTask() CreatedAt = %v, want %v", got.CreatedAt, tt.createdAt)
			}
		})
	}
}

func newValidEpisodeSubTask() EpisodeSubTask {
	return EpisodeSubTask{
		IDNatural:        "550e8400-e29b-41d4-a716-446655440000",
		OrganizationID:   "550e8400-e29b-41d4-a716-446655440001",
		EpisodeID:        "550e8400-e29b-41d4-a716-446655440002",
		SubTaskID:        "550e8400-e29b-41d4-a716-446655440003",
		CollectionStatus: SubTaskCollectionStatusReady,
	}
}

func newEpisodeSubTaskWithStatus(status SubTaskCollectionStatus) EpisodeSubTask {
	est := newValidEpisodeSubTask()
	est.CollectionStatus = status
	return est
}

func TestEpisodeSubTask_StartProgress(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus SubTaskCollectionStatus
		wantErr       bool
		wantStatus    SubTaskCollectionStatus
	}{
		{
			name:          "success: Ready → InProgress",
			initialStatus: SubTaskCollectionStatusReady,
			wantErr:       false,
			wantStatus:    SubTaskCollectionStatusInProgress,
		},
		{
			name:          "idempotent: InProgress → InProgress (no error)",
			initialStatus: SubTaskCollectionStatusInProgress,
			wantErr:       false,
			wantStatus:    SubTaskCollectionStatusInProgress,
		},
		{
			name:          "error: Completed → cannot start progress",
			initialStatus: SubTaskCollectionStatusCompleted,
			wantErr:       true,
		},
		{
			name:          "error: Cancelled → cannot start progress",
			initialStatus: SubTaskCollectionStatusCancelled,
			wantErr:       true,
		},
		{
			name:          "error: Skipped → cannot start progress",
			initialStatus: SubTaskCollectionStatusSkipped,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			est := newEpisodeSubTaskWithStatus(tt.initialStatus)
			err := est.StartProgress()

			if tt.wantErr {
				if err == nil {
					t.Errorf("EpisodeSubTask.StartProgress() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("EpisodeSubTask.StartProgress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if est.CollectionStatus != tt.wantStatus {
				t.Errorf("EpisodeSubTask.StartProgress() CollectionStatus = %v, want %v", est.CollectionStatus, tt.wantStatus)
			}
		})
	}
}

func TestEpisodeSubTask_Complete(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus SubTaskCollectionStatus
		wantErr       bool
		wantStatus    SubTaskCollectionStatus
	}{
		{
			name:          "success: Ready → Completed",
			initialStatus: SubTaskCollectionStatusReady,
			wantErr:       false,
			wantStatus:    SubTaskCollectionStatusCompleted,
		},
		{
			name:          "success: InProgress → Completed",
			initialStatus: SubTaskCollectionStatusInProgress,
			wantErr:       false,
			wantStatus:    SubTaskCollectionStatusCompleted,
		},
		{
			name:          "idempotent: Completed → Completed (no error)",
			initialStatus: SubTaskCollectionStatusCompleted,
			wantErr:       false,
			wantStatus:    SubTaskCollectionStatusCompleted,
		},
		{
			name:          "error: Cancelled → cannot complete",
			initialStatus: SubTaskCollectionStatusCancelled,
			wantErr:       true,
		},
		{
			name:          "error: Skipped → cannot complete",
			initialStatus: SubTaskCollectionStatusSkipped,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			est := newEpisodeSubTaskWithStatus(tt.initialStatus)
			err := est.Complete()

			if tt.wantErr {
				if err == nil {
					t.Errorf("EpisodeSubTask.Complete() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("EpisodeSubTask.Complete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if est.CollectionStatus != tt.wantStatus {
				t.Errorf("EpisodeSubTask.Complete() CollectionStatus = %v, want %v", est.CollectionStatus, tt.wantStatus)
			}
		})
	}
}

func TestEpisodeSubTask_Skip(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus SubTaskCollectionStatus
		wantErr       bool
		wantStatus    SubTaskCollectionStatus
	}{
		{
			name:          "success: Ready → Skipped",
			initialStatus: SubTaskCollectionStatusReady,
			wantErr:       false,
			wantStatus:    SubTaskCollectionStatusSkipped,
		},
		{
			name:          "success: InProgress → Skipped",
			initialStatus: SubTaskCollectionStatusInProgress,
			wantErr:       false,
			wantStatus:    SubTaskCollectionStatusSkipped,
		},
		{
			name:          "idempotent: Skipped → Skipped (no error)",
			initialStatus: SubTaskCollectionStatusSkipped,
			wantErr:       false,
			wantStatus:    SubTaskCollectionStatusSkipped,
		},
		{
			name:          "error: Cancelled → cannot skip",
			initialStatus: SubTaskCollectionStatusCancelled,
			wantErr:       true,
		},
		{
			name:          "error: Completed → cannot skip",
			initialStatus: SubTaskCollectionStatusCompleted,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			est := newEpisodeSubTaskWithStatus(tt.initialStatus)
			err := est.Skip()

			if tt.wantErr {
				if err == nil {
					t.Errorf("EpisodeSubTask.Skip() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("EpisodeSubTask.Skip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if est.CollectionStatus != tt.wantStatus {
				t.Errorf("EpisodeSubTask.Skip() CollectionStatus = %v, want %v", est.CollectionStatus, tt.wantStatus)
			}
		})
	}
}

func TestEpisodeSubTask_Cancel(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus SubTaskCollectionStatus
		wantErr       bool
		wantStatus    SubTaskCollectionStatus
	}{
		{
			name:          "success: Ready → Cancelled",
			initialStatus: SubTaskCollectionStatusReady,
			wantErr:       false,
			wantStatus:    SubTaskCollectionStatusCancelled,
		},
		{
			name:          "success: InProgress → Cancelled",
			initialStatus: SubTaskCollectionStatusInProgress,
			wantErr:       false,
			wantStatus:    SubTaskCollectionStatusCancelled,
		},
		{
			name:          "success: Skipped → Cancelled",
			initialStatus: SubTaskCollectionStatusSkipped,
			wantErr:       false,
			wantStatus:    SubTaskCollectionStatusCancelled,
		},
		{
			name:          "idempotent: Cancelled → Cancelled (no error)",
			initialStatus: SubTaskCollectionStatusCancelled,
			wantErr:       false,
			wantStatus:    SubTaskCollectionStatusCancelled,
		},
		{
			name:          "error: Completed → cannot cancel",
			initialStatus: SubTaskCollectionStatusCompleted,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			est := newEpisodeSubTaskWithStatus(tt.initialStatus)
			err := est.Cancel()

			if tt.wantErr {
				if err == nil {
					t.Errorf("EpisodeSubTask.Cancel() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("EpisodeSubTask.Cancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if est.CollectionStatus != tt.wantStatus {
				t.Errorf("EpisodeSubTask.Cancel() CollectionStatus = %v, want %v", est.CollectionStatus, tt.wantStatus)
			}
		})
	}
}

func TestSubTaskCollectionStatusPolicy(t *testing.T) {
	tests := []struct {
		name                   string
		status                 SubTaskCollectionStatus
		wantTerminal           bool
		wantWorkflowResolved   bool
		wantSuccessfulComplete bool
	}{
		{name: "ready is open", status: SubTaskCollectionStatusReady},
		{name: "in progress is open", status: SubTaskCollectionStatusInProgress},
		{
			name:                   "completed is terminal, resolved, and successful",
			status:                 SubTaskCollectionStatusCompleted,
			wantTerminal:           true,
			wantWorkflowResolved:   true,
			wantSuccessfulComplete: true,
		},
		{
			name:                 "skipped is terminal and workflow resolved",
			status:               SubTaskCollectionStatusSkipped,
			wantTerminal:         true,
			wantWorkflowResolved: true,
		},
		{
			name:         "cancelled is terminal but unresolved",
			status:       SubTaskCollectionStatusCancelled,
			wantTerminal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsTerminal(); got != tt.wantTerminal {
				t.Errorf("SubTaskCollectionStatus.IsTerminal() = %v, want %v", got, tt.wantTerminal)
			}
			if got := tt.status.IsWorkflowResolved(); got != tt.wantWorkflowResolved {
				t.Errorf("SubTaskCollectionStatus.IsWorkflowResolved() = %v, want %v", got, tt.wantWorkflowResolved)
			}
			if got := tt.status.IsSuccessfulCompletion(); got != tt.wantSuccessfulComplete {
				t.Errorf("SubTaskCollectionStatus.IsSuccessfulCompletion() = %v, want %v", got, tt.wantSuccessfulComplete)
			}

			est := newEpisodeSubTaskWithStatus(tt.status)
			if got := est.IsTerminal(); got != tt.wantTerminal {
				t.Errorf("EpisodeSubTask.IsTerminal() = %v, want %v", got, tt.wantTerminal)
			}
			if got := est.IsWorkflowResolved(); got != tt.wantWorkflowResolved {
				t.Errorf("EpisodeSubTask.IsWorkflowResolved() = %v, want %v", got, tt.wantWorkflowResolved)
			}
			if got := est.IsSuccessfulCompletion(); got != tt.wantSuccessfulComplete {
				t.Errorf("EpisodeSubTask.IsSuccessfulCompletion() = %v, want %v", got, tt.wantSuccessfulComplete)
			}
		})
	}
}
