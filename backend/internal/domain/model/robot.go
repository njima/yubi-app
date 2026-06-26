package model

import (
	"encoding/json"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RobotStatus int

const (
	RobotStatusOnline      RobotStatus = 0
	RobotStatusBusy        RobotStatus = 1
	RobotStatusOffline     RobotStatus = 2
	RobotStatusFaulted     RobotStatus = 3
	RobotStatusMaintenance RobotStatus = 4
	RobotStatusReady       RobotStatus = 5
)

type LeaderStatus int

const (
	LeaderStatusReady       LeaderStatus = 0
	LeaderStatusFaulted     LeaderStatus = 1
	LeaderStatusMaintenance LeaderStatus = 2
)

type Robot struct {
	ID                   int64
	IDNatural            string
	OrganizationID       string
	OrganizationName     string
	SiteID               string
	SiteName             string
	LocationID           string
	LocationName         string
	Name                 string
	RobotType            *string
	Status               RobotStatus
	LeaderStatus         *LeaderStatus
	LeaderFaultStartedAt *time.Time
	FaultStartedAt       *time.Time
	LastHeartbeatAt      *time.Time
	OfflineReason        *string
	RobotConfig          *json.RawMessage
	ActiveEpisodeID      *string
	ActiveUserID         *string
	CreatedAt            time.Time
	UpdatedAt            *time.Time
}

type Robots []*Robot

func InitRobot(
	organizationID,
	locationID,
	name string,
	robotType *string,
	robotConfig *json.RawMessage,
) (Robot, error) {
	ID, err := InitID()
	if err != nil {
		return Robot{}, err
	}

	if robotType == nil {
		emptyStr := ""
		robotType = &emptyStr
	}

	if robotConfig == nil {
		emptyJSON := json.RawMessage(`{}`)
		robotConfig = &emptyJSON
	}

	rob := Robot{
		IDNatural:      ID,
		OrganizationID: organizationID,
		LocationID:     locationID,
		Name:           name,
		RobotType:      robotType,
		Status:         RobotStatusReady,
		RobotConfig:    robotConfig,
		CreatedAt:      time.Now(),
	}

	if err := rob.validate(); err != nil {
		return Robot{}, err
	}

	return rob, nil
}

func NewRobot(idNatural, organizationID, locationID, name string) *Robot {
	return &Robot{
		IDNatural:      idNatural,
		OrganizationID: organizationID,
		LocationID:     locationID,
		Name:           name,
		CreatedAt:      time.Now(),
	}
}

func (r Robot) validate() error {
	if err := (validation.Errors{
		"id_natural":      validation.Validate(r.IDNatural, validation.Required.Error("id_natural is required")),
		"organization_id": validation.Validate(r.OrganizationID, validation.Required.Error("organization_id is required")),
		"location_id":     validation.Validate(r.LocationID, validation.Required.Error("location_id is required")),
		"name": validation.Validate(
			r.Name,
			validation.Required.Error("name is required"),
			validation.RuneLength(1, 100).Error("name must be between 1 and 100 characters"),
		),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "robot validation failed: %v", err))
	}

	return nil
}

func (r *Robot) SetName(name string) error {
	r.Name = name
	return r.validate()
}

func (r *Robot) SetRobotType(robotType string) error {
	r.RobotType = &robotType
	return r.validate()
}

func (r *Robot) SetStatus(status RobotStatus) error {
	r.Status = status
	return r.validate()
}

func (r *Robot) SetLastHeartbeatAt(lastHeartbeatAt time.Time) error {
	r.LastHeartbeatAt = &lastHeartbeatAt
	return r.validate()
}

func (r *Robot) SetOfflineReason(offlineReason string) error {
	r.OfflineReason = &offlineReason
	return r.validate()
}

func (r *Robot) SetLeaderStatus(leaderStatus *LeaderStatus) {
	r.LeaderStatus = leaderStatus
}

func (r *Robot) SetLeaderFaultStartedAt(leaderFaultStartedAt *time.Time) {
	r.LeaderFaultStartedAt = leaderFaultStartedAt
}

func (r *Robot) SetFaultStartedAt(faultStartedAt *time.Time) {
	r.FaultStartedAt = faultStartedAt
}

func (r *Robot) SetRobotConfig(robotConfig json.RawMessage) error {
	r.RobotConfig = &robotConfig
	return r.validate()
}

// CanStartTeleoperation checks if the robot can start teleoperation.
// The caller must resolve the robot's status via ResolvedStatus before calling this.
func (r *Robot) CanStartTeleoperation() error {
	if r.Status != RobotStatusOnline {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "robot must be online to start teleoperation, current status: %d", r.Status),
		)
	}
	return nil
}

// StartTeleoperation transitions the robot to Busy state and sets active episode/user
func (r *Robot) StartTeleoperation(episodeID, userID string) error {
	if err := r.CanStartTeleoperation(); err != nil {
		return err
	}
	r.Status = RobotStatusBusy
	r.ActiveEpisodeID = &episodeID
	r.ActiveUserID = &userID
	return nil
}

// CanEndTeleoperation checks if the robot can end teleoperation
func (r *Robot) CanEndTeleoperation() error {
	if r.Status != RobotStatusBusy {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "robot status must be Busy to end teleoperation, current: %d", r.Status),
		)
	}
	return nil
}

// EndTeleoperation transitions the robot back to Ready state and clears active episode/user
func (r *Robot) EndTeleoperation() error {
	if err := r.CanEndTeleoperation(); err != nil {
		return err
	}
	r.Status = RobotStatusReady
	r.ActiveEpisodeID = nil
	r.ActiveUserID = nil
	return nil
}

// ResolvedStatus updates r.Status by combining DB status with Redis heartbeat state.
// Ready or Online + heartbeat alive → Online
// Ready or Online + heartbeat absent → Offline
// Busy / Faulted / Maintenance → unchanged
func (r *Robot) ResolvedStatus(heartbeatAlive bool) {
	if r.Status == RobotStatusReady || r.Status == RobotStatusOnline {
		if heartbeatAlive {
			r.Status = RobotStatusOnline
		} else {
			r.Status = RobotStatusOffline
		}
	}
}

func (r Robot) ConsecutiveFaultDays() *int {
	if r.Status != RobotStatusFaulted || r.FaultStartedAt == nil {
		return nil
	}
	days := int(time.Since(*r.FaultStartedAt).Hours() / 24)
	return &days
}

func (r Robot) LeaderConsecutiveFaultDays() *int {
	if r.LeaderStatus == nil || *r.LeaderStatus != LeaderStatusFaulted || r.LeaderFaultStartedAt == nil {
		return nil
	}
	days := int(time.Since(*r.LeaderFaultStartedAt).Hours() / 24)
	return &days
}

// Robot operator identity (stored in Redis with TTL).

type RobotOperator struct {
	UserID           string `json:"user_id"`
	DisplayName      string `json:"display_name"`
	OrganizationName string `json:"organization_name"`
}

// Gate condition types for recording gates.

type GateConditionStatus struct {
	GateLevel int                        `json:"gate_level"`
	Groups    map[string]GateGroupStatus `json:"groups"`
}

type GateGroupStatus struct {
	Level      int             `json:"level"`
	Settled    bool            `json:"settled"`
	Conditions []GateCondition `json:"conditions"`
}

type GateCondition struct {
	Name       string `json:"name"`
	Passed     bool   `json:"passed"`
	Reason     string `json:"reason"`
	Escalation int    `json:"escalation"`
}
