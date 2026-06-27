package model

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestInitRobot(t *testing.T) {
	robotType := "TestModel"
	cameraConfig := json.RawMessage(`{"resolution": "1080p"}`)

	tests := []struct {
		name           string
		organizationID string
		locationID     string
		robotName      string
		robotType      *string
		cameraConfig   *json.RawMessage
		wantErr        bool
	}{
		{
			name:           "success with valid inputs",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			locationID:     "550e8400-e29b-41d4-a716-446655440002",
			robotName:      "Test Robot",
			robotType:      &robotType,
			cameraConfig:   &cameraConfig,
			wantErr:        false,
		},
		{
			name:           "success with nil robotType and cameraConfig",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			locationID:     "550e8400-e29b-41d4-a716-446655440002",
			robotName:      "Test Robot",
			robotType:      nil,
			cameraConfig:   nil,
			wantErr:        false,
		},
		{
			name:           "error when organization_id is empty",
			organizationID: "",
			locationID:     "550e8400-e29b-41d4-a716-446655440002",
			robotName:      "Test Robot",
			robotType:      &robotType,
			cameraConfig:   &cameraConfig,
			wantErr:        true,
		},
		{
			name:           "error when location_id is empty",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			locationID:     "",
			robotName:      "Test Robot",
			robotType:      &robotType,
			cameraConfig:   &cameraConfig,
			wantErr:        true,
		},
		{
			name:           "error when name is empty",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			locationID:     "550e8400-e29b-41d4-a716-446655440002",
			robotName:      "",
			robotType:      &robotType,
			cameraConfig:   &cameraConfig,
			wantErr:        true,
		},
		{
			name:           "error when name exceeds 100 characters",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			locationID:     "550e8400-e29b-41d4-a716-446655440002",
			robotName:      strings.Repeat("a", 101),
			robotType:      &robotType,
			cameraConfig:   &cameraConfig,
			wantErr:        true,
		},
		{
			name:           "success when name is exactly 100 characters",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			locationID:     "550e8400-e29b-41d4-a716-446655440002",
			robotName:      strings.Repeat("a", 100),
			robotType:      &robotType,
			cameraConfig:   &cameraConfig,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitRobot(tt.organizationID, tt.locationID, tt.robotName, tt.robotType, tt.cameraConfig)

			if tt.wantErr {
				if err == nil {
					t.Errorf("InitRobot() error = nil, wantErr %v", tt.wantErr)
				}
				if got.IDNatural != "" {
					t.Errorf("InitRobot() IDNatural = %v, want empty", got.IDNatural)
				}
				return
			}

			if err != nil {
				t.Errorf("InitRobot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.IDNatural == "" {
				t.Errorf("InitRobot() IDNatural is empty")
			}
			if got.OrganizationID != tt.organizationID {
				t.Errorf("InitRobot() OrganizationID = %v, want %v", got.OrganizationID, tt.organizationID)
			}
			if got.LocationID != tt.locationID {
				t.Errorf("InitRobot() LocationID = %v, want %v", got.LocationID, tt.locationID)
			}
			if got.Name != tt.robotName {
				t.Errorf("InitRobot() Name = %v, want %v", got.Name, tt.robotName)
			}
			if got.Status != RobotStatusReady {
				t.Errorf("InitRobot() Status = %v, want %v", got.Status, RobotStatusReady)
			}
			if got.CreatedAt.IsZero() {
				t.Errorf("InitRobot() CreatedAt is zero")
			}
		})
	}
}

func TestNewRobot(t *testing.T) {
	createdAt := time.Date(2026, 6, 26, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	robotType := "TestModel"
	leaderStatus := LeaderStatusFaulted
	leaderFaultStartedAt := createdAt.Add(30 * time.Minute)
	faultStartedAt := createdAt.Add(45 * time.Minute)
	lastHeartbeatAt := createdAt.Add(2 * time.Hour)
	offlineReason := "network"
	robotConfig := json.RawMessage(`{"mode":"test"}`)
	activeEpisodeID := "episode-1"
	activeUserID := "user-1"

	got := NewRobot(
		1,
		"550e8400-e29b-41d4-a716-446655440000",
		"550e8400-e29b-41d4-a716-446655440001",
		"Org",
		"site-1",
		"Site",
		"550e8400-e29b-41d4-a716-446655440002",
		"Location",
		"Test Robot",
		&robotType,
		RobotStatusBusy,
		&leaderStatus,
		&leaderFaultStartedAt,
		&faultStartedAt,
		&lastHeartbeatAt,
		&offlineReason,
		&robotConfig,
		&activeEpisodeID,
		&activeUserID,
		createdAt,
		&updatedAt,
	)

	if got.ID != 1 {
		t.Errorf("NewRobot() ID = %v, want 1", got.ID)
	}
	if got.IDNatural != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("NewRobot() IDNatural = %v", got.IDNatural)
	}
	if got.OrganizationName != "Org" || got.SiteName != "Site" || got.LocationName != "Location" {
		t.Errorf("NewRobot() relation names = org:%q site:%q location:%q", got.OrganizationName, got.SiteName, got.LocationName)
	}
	if got.RobotType == nil || *got.RobotType != robotType {
		t.Errorf("NewRobot() RobotType = %v, want %q", got.RobotType, robotType)
	}
	if got.Status != RobotStatusBusy {
		t.Errorf("NewRobot() Status = %v, want %v", got.Status, RobotStatusBusy)
	}
	if got.LeaderStatus == nil || *got.LeaderStatus != leaderStatus {
		t.Errorf("NewRobot() LeaderStatus = %v, want %v", got.LeaderStatus, leaderStatus)
	}
	if got.LeaderFaultStartedAt == nil || !got.LeaderFaultStartedAt.Equal(leaderFaultStartedAt) {
		t.Errorf("NewRobot() LeaderFaultStartedAt = %v, want %v", got.LeaderFaultStartedAt, leaderFaultStartedAt)
	}
	if got.FaultStartedAt == nil || !got.FaultStartedAt.Equal(faultStartedAt) {
		t.Errorf("NewRobot() FaultStartedAt = %v, want %v", got.FaultStartedAt, faultStartedAt)
	}
	if got.LastHeartbeatAt == nil || !got.LastHeartbeatAt.Equal(lastHeartbeatAt) {
		t.Errorf("NewRobot() LastHeartbeatAt = %v, want %v", got.LastHeartbeatAt, lastHeartbeatAt)
	}
	if got.OfflineReason == nil || *got.OfflineReason != offlineReason {
		t.Errorf("NewRobot() OfflineReason = %v, want %q", got.OfflineReason, offlineReason)
	}
	if got.RobotConfig == nil || string(*got.RobotConfig) != string(robotConfig) {
		t.Errorf("NewRobot() RobotConfig = %v, want %s", got.RobotConfig, robotConfig)
	}
	if got.ActiveEpisodeID == nil || *got.ActiveEpisodeID != activeEpisodeID {
		t.Errorf("NewRobot() ActiveEpisodeID = %v, want %q", got.ActiveEpisodeID, activeEpisodeID)
	}
	if got.ActiveUserID == nil || *got.ActiveUserID != activeUserID {
		t.Errorf("NewRobot() ActiveUserID = %v, want %q", got.ActiveUserID, activeUserID)
	}
	if !got.CreatedAt.Equal(createdAt) {
		t.Errorf("NewRobot() CreatedAt = %v, want %v", got.CreatedAt, createdAt)
	}
	if got.UpdatedAt == nil || !got.UpdatedAt.Equal(updatedAt) {
		t.Errorf("NewRobot() UpdatedAt = %v, want %v", got.UpdatedAt, updatedAt)
	}
}

func TestRobot_validate(t *testing.T) {
	tests := []struct {
		name    string
		robot   Robot
		wantErr bool
	}{
		{
			name: "valid with all required fields",
			robot: Robot{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				LocationID:     "550e8400-e29b-41d4-a716-446655440002",
				Name:           "Test Robot",
			},
			wantErr: false,
		},
		{
			name: "error when id_natural is empty",
			robot: Robot{
				IDNatural:      "",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				LocationID:     "550e8400-e29b-41d4-a716-446655440002",
				Name:           "Test Robot",
			},
			wantErr: true,
		},
		{
			name: "error when organization_id is empty",
			robot: Robot{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "",
				LocationID:     "550e8400-e29b-41d4-a716-446655440002",
				Name:           "Test Robot",
			},
			wantErr: true,
		},
		{
			name: "error when location_id is empty",
			robot: Robot{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				LocationID:     "",
				Name:           "Test Robot",
			},
			wantErr: true,
		},
		{
			name: "error when name is empty",
			robot: Robot{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				LocationID:     "550e8400-e29b-41d4-a716-446655440002",
				Name:           "",
			},
			wantErr: true,
		},
		{
			name: "error when name exceeds 100 characters",
			robot: Robot{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				LocationID:     "550e8400-e29b-41d4-a716-446655440002",
				Name:           strings.Repeat("a", 101),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.robot.validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Robot.validate() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Robot.validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func newValidRobot() Robot {
	return Robot{
		IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
		OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
		LocationID:     "550e8400-e29b-41d4-a716-446655440002",
		Name:           "Test Robot",
	}
}

func TestRobot_SetName(t *testing.T) {
	tests := []struct {
		name      string
		robotName string
		wantErr   bool
	}{
		{
			name:      "success with valid name",
			robotName: "New Robot",
			wantErr:   false,
		},
		{
			name:      "success with name at max length",
			robotName: strings.Repeat("b", 100),
			wantErr:   false,
		},
		{
			name:      "error when name is empty",
			robotName: "",
			wantErr:   true,
		},
		{
			name:      "error when name exceeds 100 characters",
			robotName: strings.Repeat("c", 101),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newValidRobot()
			err := r.SetName(tt.robotName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Robot.SetName() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Robot.SetName() error = %v, wantErr %v", err, tt.wantErr)
				}
				if r.Name != tt.robotName {
					t.Errorf("Robot.Name = %v, want %v", r.Name, tt.robotName)
				}
			}
		})
	}
}

func TestRobot_SetRobotType(t *testing.T) {
	tests := []struct {
		name      string
		robotType string
		wantErr   bool
	}{
		{
			name:      "success with valid robot type",
			robotType: "toyota-hsr",
			wantErr:   false,
		},
		{
			name:      "success with empty robot type",
			robotType: "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newValidRobot()
			err := r.SetRobotType(tt.robotType)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Robot.SetRobotType() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Robot.SetRobotType() error = %v, wantErr %v", err, tt.wantErr)
				}
				if r.RobotType == nil {
					t.Errorf("Robot.RobotType is nil")
				} else if *r.RobotType != tt.robotType {
					t.Errorf("Robot.RobotType = %v, want %v", *r.RobotType, tt.robotType)
				}
			}
		})
	}
}

func TestRobot_SetStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  RobotStatus
		wantErr bool
	}{
		{
			name:    "success with offline status",
			status:  RobotStatusOffline,
			wantErr: false,
		},
		{
			name:    "success with ready status",
			status:  RobotStatusReady,
			wantErr: false,
		},
		{
			name:    "success with busy status",
			status:  RobotStatusBusy,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newValidRobot()
			err := r.SetStatus(tt.status)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Robot.SetStatus() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Robot.SetStatus() error = %v, wantErr %v", err, tt.wantErr)
				}
				if r.Status != tt.status {
					t.Errorf("Robot.Status = %v, want %v", r.Status, tt.status)
				}
			}
		})
	}
}

func TestRobot_SetLastHeartbeatAt(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name            string
		lastHeartbeatAt time.Time
		wantErr         bool
	}{
		{
			name:            "success with valid time",
			lastHeartbeatAt: now,
			wantErr:         false,
		},
		{
			name:            "success with past time",
			lastHeartbeatAt: now.Add(-24 * time.Hour),
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newValidRobot()
			err := r.SetLastHeartbeatAt(tt.lastHeartbeatAt)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Robot.SetLastHeartbeatAt() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Robot.SetLastHeartbeatAt() error = %v, wantErr %v", err, tt.wantErr)
				}
				if r.LastHeartbeatAt == nil {
					t.Errorf("Robot.LastHeartbeatAt is nil")
				} else if *r.LastHeartbeatAt != tt.lastHeartbeatAt {
					t.Errorf("Robot.LastHeartbeatAt = %v, want %v", *r.LastHeartbeatAt, tt.lastHeartbeatAt)
				}
			}
		})
	}
}

func TestRobot_SetOfflineReason(t *testing.T) {
	tests := []struct {
		name          string
		offlineReason string
		wantErr       bool
	}{
		{
			name:          "success with valid reason",
			offlineReason: "Maintenance",
			wantErr:       false,
		},
		{
			name:          "success with empty reason",
			offlineReason: "",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newValidRobot()
			err := r.SetOfflineReason(tt.offlineReason)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Robot.SetOfflineReason() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Robot.SetOfflineReason() error = %v, wantErr %v", err, tt.wantErr)
				}
				if r.OfflineReason == nil {
					t.Errorf("Robot.OfflineReason is nil")
				} else if *r.OfflineReason != tt.offlineReason {
					t.Errorf("Robot.OfflineReason = %v, want %v", *r.OfflineReason, tt.offlineReason)
				}
			}
		})
	}
}

func TestRobot_SetRobotConfig(t *testing.T) {
	tests := []struct {
		name        string
		robotConfig json.RawMessage
		wantErr     bool
	}{
		{
			name:        "success with valid config",
			robotConfig: json.RawMessage(`{"resolution": "1080p"}`),
			wantErr:     false,
		},
		{
			name:        "success with empty config",
			robotConfig: json.RawMessage(`{}`),
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newValidRobot()
			err := r.SetRobotConfig(tt.robotConfig)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Robot.SetRobotConfig() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Robot.SetRobotConfig() error = %v, wantErr %v", err, tt.wantErr)
				}
				if r.RobotConfig == nil {
					t.Errorf("Robot.RobotConfig is nil")
				}
			}
		})
	}
}

func newRobotWithStatus(status RobotStatus) Robot {
	r := newValidRobot()
	r.Status = status
	return r
}

func TestRobotStatusPolicy(t *testing.T) {
	tests := []struct {
		name                    string
		status                  RobotStatus
		wantPersistentOperation bool
		wantConnectionOnly      bool
	}{
		{
			name:               "online is connection-only display state",
			status:             RobotStatusOnline,
			wantConnectionOnly: true,
		},
		{
			name:                    "busy is persistent operation state",
			status:                  RobotStatusBusy,
			wantPersistentOperation: true,
		},
		{
			name:               "offline is connection-only display state",
			status:             RobotStatusOffline,
			wantConnectionOnly: true,
		},
		{
			name:                    "faulted is persistent operation state",
			status:                  RobotStatusFaulted,
			wantPersistentOperation: true,
		},
		{
			name:                    "maintenance is persistent operation state",
			status:                  RobotStatusMaintenance,
			wantPersistentOperation: true,
		},
		{
			name:                    "ready is persistent operation state",
			status:                  RobotStatusReady,
			wantPersistentOperation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsPersistentOperationStatus(); got != tt.wantPersistentOperation {
				t.Errorf("RobotStatus.IsPersistentOperationStatus() = %v, want %v", got, tt.wantPersistentOperation)
			}
			if got := tt.status.IsConnectionOnlyStatus(); got != tt.wantConnectionOnly {
				t.Errorf("RobotStatus.IsConnectionOnlyStatus() = %v, want %v", got, tt.wantConnectionOnly)
			}
		})
	}
}

func TestRobot_ResolvedStatusDoesNotMutateOperationStatus(t *testing.T) {
	tests := []struct {
		name           string
		initialStatus  RobotStatus
		heartbeatAlive bool
		wantResolved   RobotStatus
		wantStored     RobotStatus
	}{
		{
			name:           "ready with heartbeat appears online",
			initialStatus:  RobotStatusReady,
			heartbeatAlive: true,
			wantResolved:   RobotStatusOnline,
			wantStored:     RobotStatusReady,
		},
		{
			name:           "ready without heartbeat appears offline",
			initialStatus:  RobotStatusReady,
			heartbeatAlive: false,
			wantResolved:   RobotStatusOffline,
			wantStored:     RobotStatusReady,
		},
		{
			name:           "busy ignores heartbeat",
			initialStatus:  RobotStatusBusy,
			heartbeatAlive: true,
			wantResolved:   RobotStatusBusy,
			wantStored:     RobotStatusBusy,
		},
		{
			name:           "faulted ignores heartbeat",
			initialStatus:  RobotStatusFaulted,
			heartbeatAlive: true,
			wantResolved:   RobotStatusFaulted,
			wantStored:     RobotStatusFaulted,
		},
		{
			name:           "maintenance ignores heartbeat",
			initialStatus:  RobotStatusMaintenance,
			heartbeatAlive: true,
			wantResolved:   RobotStatusMaintenance,
			wantStored:     RobotStatusMaintenance,
		},
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
		{
			name:           "ready with heartbeat starts",
			initialStatus:  RobotStatusReady,
			heartbeatAlive: true,
			wantErr:        false,
			wantStatus:     RobotStatusBusy,
		},
		{
			name:           "ready without heartbeat fails",
			initialStatus:  RobotStatusReady,
			heartbeatAlive: false,
			wantErr:        true,
			wantStatus:     RobotStatusReady,
		},
		{
			name:           "busy with heartbeat fails",
			initialStatus:  RobotStatusBusy,
			heartbeatAlive: true,
			wantErr:        true,
			wantStatus:     RobotStatusBusy,
		},
		{
			name:           "legacy online status fails as non-persistent operation state",
			initialStatus:  RobotStatusOnline,
			heartbeatAlive: true,
			wantErr:        true,
			wantStatus:     RobotStatusOnline,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRobotWithStatus(tt.initialStatus)
			err := r.StartTeleoperation(episodeID, userID, tt.heartbeatAlive)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Robot.StartTeleoperation() error = nil, wantErr %v", tt.wantErr)
				}
				if r.Status != tt.wantStatus {
					t.Errorf("Robot.StartTeleoperation() Status = %v, want %v", r.Status, tt.wantStatus)
				}
				return
			}
			if err != nil {
				t.Errorf("Robot.StartTeleoperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if r.Status != tt.wantStatus {
				t.Errorf("Robot.StartTeleoperation() Status = %v, want %v", r.Status, tt.wantStatus)
			}
			if r.ActiveEpisodeID == nil || *r.ActiveEpisodeID != episodeID {
				t.Errorf("Robot.StartTeleoperation() ActiveEpisodeID = %v, want %v", r.ActiveEpisodeID, episodeID)
			}
			if r.ActiveUserID == nil || *r.ActiveUserID != userID {
				t.Errorf("Robot.StartTeleoperation() ActiveUserID = %v, want %v", r.ActiveUserID, userID)
			}
		})
	}
}

func TestRobot_EndTeleoperation(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus RobotStatus
		wantErr       bool
		wantStatus    RobotStatus
	}{
		{
			name:          "success: Busy → Ready",
			initialStatus: RobotStatusBusy,
			wantErr:       false,
			wantStatus:    RobotStatusReady,
		},
		{
			name:          "error: Ready → cannot end teleoperation",
			initialStatus: RobotStatusReady,
			wantErr:       true,
		},
		{
			name:          "error: Offline → cannot end teleoperation",
			initialStatus: RobotStatusOffline,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRobotWithStatus(tt.initialStatus)
			episodeID := "550e8400-e29b-41d4-a716-446655440010"
			userID := "550e8400-e29b-41d4-a716-446655440011"
			if tt.initialStatus == RobotStatusBusy {
				r.ActiveEpisodeID = &episodeID
				r.ActiveUserID = &userID
			}
			err := r.EndTeleoperation()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Robot.EndTeleoperation() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Robot.EndTeleoperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if r.Status != tt.wantStatus {
				t.Errorf("Robot.EndTeleoperation() Status = %v, want %v", r.Status, tt.wantStatus)
			}
			if r.ActiveEpisodeID != nil {
				t.Errorf("Robot.EndTeleoperation() ActiveEpisodeID = %v, want nil", *r.ActiveEpisodeID)
			}
			if r.ActiveUserID != nil {
				t.Errorf("Robot.EndTeleoperation() ActiveUserID = %v, want nil", *r.ActiveUserID)
			}
		})
	}
}

func TestRobot_ConsecutiveFaultDays(t *testing.T) {
	tests := []struct {
		name      string
		status    RobotStatus
		startedAt *time.Time
		wantNil   bool
		wantDays  int
	}{
		{
			name:      "returns nil when status is not faulted",
			status:    RobotStatusReady,
			startedAt: nil,
			wantNil:   true,
		},
		{
			name:      "returns nil when faulted but start time is nil",
			status:    RobotStatusFaulted,
			startedAt: nil,
			wantNil:   true,
		},
		{
			name:      "returns floored elapsed days when faulted",
			status:    RobotStatusFaulted,
			startedAt: ptrTime(time.Now().Add(-49 * time.Hour)),
			wantNil:   false,
			wantDays:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newValidRobot()
			r.Status = tt.status
			r.FaultStartedAt = tt.startedAt

			got := r.ConsecutiveFaultDays()

			if tt.wantNil {
				if got != nil {
					t.Errorf("Robot.ConsecutiveFaultDays() = %v, want nil", *got)
				}
				return
			}

			if got == nil {
				t.Fatalf("Robot.ConsecutiveFaultDays() = nil, want %d", tt.wantDays)
			}
			if *got != tt.wantDays {
				t.Errorf("Robot.ConsecutiveFaultDays() = %d, want %d", *got, tt.wantDays)
			}
		})
	}
}

func TestRobot_LeaderConsecutiveFaultDays(t *testing.T) {
	faulted := LeaderStatusFaulted
	ready := LeaderStatusReady

	tests := []struct {
		name      string
		status    *LeaderStatus
		startedAt *time.Time
		wantNil   bool
		wantDays  int
	}{
		{
			name:      "returns nil when leader status is nil",
			status:    nil,
			startedAt: nil,
			wantNil:   true,
		},
		{
			name:      "returns nil when leader status is not faulted",
			status:    &ready,
			startedAt: ptrTime(time.Now().Add(-49 * time.Hour)),
			wantNil:   true,
		},
		{
			name:      "returns nil when leader fault started at is nil",
			status:    &faulted,
			startedAt: nil,
			wantNil:   true,
		},
		{
			name:      "returns floored elapsed days when leader is faulted",
			status:    &faulted,
			startedAt: ptrTime(time.Now().Add(-49 * time.Hour)),
			wantNil:   false,
			wantDays:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newValidRobot()
			r.LeaderStatus = tt.status
			r.LeaderFaultStartedAt = tt.startedAt

			got := r.LeaderConsecutiveFaultDays()

			if tt.wantNil {
				if got != nil {
					t.Errorf("Robot.LeaderConsecutiveFaultDays() = %v, want nil", *got)
				}
				return
			}

			if got == nil {
				t.Fatalf("Robot.LeaderConsecutiveFaultDays() = nil, want %d", tt.wantDays)
			}
			if *got != tt.wantDays {
				t.Errorf("Robot.LeaderConsecutiveFaultDays() = %d, want %d", *got, tt.wantDays)
			}
		})
	}
}

func ptrTime(v time.Time) *time.Time {
	return &v
}
