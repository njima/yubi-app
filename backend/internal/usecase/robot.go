package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/pagination"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/uptrace/bun"
)

type RobotUsecase interface {
	Create(ctx context.Context, input RobotCreateInput) (model.Robot, error)
	GetByID(ctx context.Context, id string) (model.Robot, error)
	List(ctx context.Context, filter repository.RobotListFilter, page, limit int) (model.Robots, int, error)
	ListTypes(ctx context.Context, filter repository.RobotTypeFilter) ([]string, error)
	Update(ctx context.Context, input RobotUpdateInput) (model.Robot, error)
	Delete(ctx context.Context, id string) error
}

type RobotCreateInput struct {
	OrganizationID string
	LocationID     string
	Name           string
	RobotType      *string
	LeaderStatus   *openapi.LeaderStatus
	RobotConfig    *json.RawMessage
}

type RobotUpdateInput struct {
	ID              string
	OrganizationID  string
	LocationID      string
	Name            *string
	RobotType       *string
	Status          *openapi.RobotStatus
	LeaderStatus    *openapi.LeaderStatus
	HasLeaderStatus bool // true when leader_status is explicitly provided (including null)
	LastHeartbeatAt *time.Time
	OfflineReason   *string
	RobotConfig     *json.RawMessage
}

type robot struct {
	repo            repository.Robot
	robotStatusRepo repository.RobotStatusRepository
	uptimeDeltaRepo repository.RobotUptimeDeltaRepository
	db              *bun.DB
}

func NewRobot(repo repository.Robot, robotStatusRepo repository.RobotStatusRepository, uptimeDeltaRepo repository.RobotUptimeDeltaRepository, db *bun.DB) *robot {
	return &robot{repo: repo, robotStatusRepo: robotStatusRepo, uptimeDeltaRepo: uptimeDeltaRepo, db: db}
}

func (r *robot) Create(ctx context.Context, input RobotCreateInput) (model.Robot, error) {
	ro, err := model.InitRobot(input.OrganizationID, input.LocationID, input.Name, input.RobotType, input.RobotConfig)
	if err != nil {
		return model.Robot{}, err
	}

	if input.LeaderStatus != nil {
		ro.SetLeaderStatus(input.LeaderStatus)
	}

	cro, err := r.repo.Create(ctx, r.db, ro)
	if err != nil {
		return model.Robot{}, err
	}

	return cro, nil
}

func (r *robot) GetByID(ctx context.Context, id string) (model.Robot, error) {
	rob, err := r.repo.GetByID(ctx, r.db, id)
	if err != nil {
		return model.Robot{}, err
	}
	status, err := r.robotStatusRepo.GetByRobotID(ctx, rob.IDNatural)
	if err != nil {
		return model.Robot{}, err
	}
	rob.ResolvedStatus(status != nil)
	return rob, nil
}

func (r *robot) List(ctx context.Context, filter repository.RobotListFilter, page, limit int) (model.Robots, int, error) {
	if limit <= 0 {
		limit = pagination.DefaultLimit
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	dbFilter := filter
	if filter.Status != nil &&
		(*filter.Status == repository.RobotFilterStatusOnline || *filter.Status == repository.RobotFilterStatusOffline) {
		dbFilter.Status = nil

		onlineIDs, err := r.robotStatusRepo.GetAllOnlineRobotIDs(ctx)
		if err != nil {
			return nil, 0, err
		}

		if *filter.Status == repository.RobotFilterStatusOnline {
			if len(onlineIDs) == 0 {
				return model.Robots{}, 0, nil
			}
			dbFilter.OnlineRobotIDs = &onlineIDs
		} else {
			dbFilter.OnlineRobotIDs = &onlineIDs
			dbFilter.ExcludeOnline = true
		}
	}

	robs, total, err := r.repo.List(ctx, r.db, dbFilter, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make(model.Robots, 0, len(robs))
	for _, rob := range robs {
		status, err := r.robotStatusRepo.GetByRobotID(ctx, rob.IDNatural)
		if err != nil {
			return nil, 0, err
		}
		rob.ResolvedStatus(status != nil)
		result = append(result, rob)
	}
	return result, total, nil
}

func (r *robot) Update(ctx context.Context, input RobotUpdateInput) (model.Robot, error) {
	rob, err := r.repo.GetByID(ctx, r.db, input.ID)
	if err != nil {
		return model.Robot{}, err
	}

	rob, err = r.update(ctx, rob, input)
	if err != nil {
		return model.Robot{}, err
	}

	urob, err := r.repo.Update(ctx, r.db, rob)
	if err != nil {
		return model.Robot{}, err
	}

	status, err := r.robotStatusRepo.GetByRobotID(ctx, urob.IDNatural)
	if err != nil {
		return model.Robot{}, err
	}
	urob.ResolvedStatus(status != nil)

	return urob, nil
}

func (r *robot) update(ctx context.Context, robot model.Robot, input RobotUpdateInput) (model.Robot, error) {
	if input.Name != nil {
		if err := robot.SetName(*input.Name); err != nil {
			return model.Robot{}, err
		}
	}
	if input.RobotType != nil {
		if err := robot.SetRobotType(*input.RobotType); err != nil {
			return model.Robot{}, err
		}
	}
	if input.Status != nil {
		wasFaulted := robot.Status == openapi.RobotStatusFaulted

		// Only Ready(5), Faulted(3), Maintenance(4) can be manually set
		switch *input.Status {
		case openapi.RobotStatusReady, openapi.RobotStatusFaulted, openapi.RobotStatusMaintenance:
			// allowed
		default:
			return model.Robot{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "status must be Ready, Faulted, or Maintenance"))
		}
		// Validate status transition based on current DB status
		switch robot.Status {
		case openapi.RobotStatusBusy:
			return model.Robot{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "cannot manually change status while robot is Busy"))
		case openapi.RobotStatusFaulted:
			if *input.Status != openapi.RobotStatusReady && *input.Status != openapi.RobotStatusMaintenance {
				return model.Robot{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "Faulted robot can only be set to Ready or Maintenance"))
			}
		case openapi.RobotStatusMaintenance:
			if *input.Status != openapi.RobotStatusReady && *input.Status != openapi.RobotStatusFaulted {
				return model.Robot{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "Maintenance robot can only be set to Ready or Faulted"))
			}
		}
		if err := robot.SetStatus(*input.Status); err != nil {
			return model.Robot{}, err
		}
		if *input.Status == openapi.RobotStatusFaulted {
			now := time.Now()
			robot.SetFaultStartedAt(&now)
		}
		if wasFaulted && *input.Status != openapi.RobotStatusFaulted {
			robot.SetFaultStartedAt(nil)
		}
	}
	if input.LastHeartbeatAt != nil {
		if err := robot.SetLastHeartbeatAt(*input.LastHeartbeatAt); err != nil {
			return model.Robot{}, err
		}
	}
	if input.OfflineReason != nil {
		if err := robot.SetOfflineReason(*input.OfflineReason); err != nil {
			return model.Robot{}, err
		}
	}
	if input.RobotConfig != nil {
		if err := robot.SetRobotConfig(*input.RobotConfig); err != nil {
			return model.Robot{}, err
		}
	}
	if input.HasLeaderStatus {
		wasLeaderFaulted := robot.LeaderStatus != nil && *robot.LeaderStatus == openapi.LeaderFaulted
		robot.SetLeaderStatus(input.LeaderStatus)
		if !wasLeaderFaulted && input.LeaderStatus != nil && *input.LeaderStatus == openapi.LeaderFaulted {
			now := time.Now()
			robot.SetLeaderFaultStartedAt(&now)
		} else if wasLeaderFaulted && (input.LeaderStatus == nil || *input.LeaderStatus != openapi.LeaderFaulted) {
			robot.SetLeaderFaultStartedAt(nil)
		}
	}

	return robot, nil
}

func (r *robot) ListTypes(ctx context.Context, filter repository.RobotTypeFilter) ([]string, error) {
	if filter.Status != nil &&
		(*filter.Status == repository.RobotFilterStatusOnline || *filter.Status == repository.RobotFilterStatusOffline) {
		requestedStatus := *filter.Status
		filter.Status = nil

		onlineIDs, err := r.robotStatusRepo.GetAllOnlineRobotIDs(ctx)
		if err != nil {
			return nil, err
		}

		if requestedStatus == repository.RobotFilterStatusOnline {
			if len(onlineIDs) == 0 {
				return []string{}, nil
			}
			filter.OnlineRobotIDs = &onlineIDs
		} else {
			filter.OnlineRobotIDs = &onlineIDs
			filter.ExcludeOnline = true
		}
	}
	return r.repo.ListTypes(ctx, r.db, filter)
}

func (r *robot) Delete(ctx context.Context, id string) error {
	rob, err := r.repo.GetByID(ctx, r.db, id)
	if err != nil {
		return err
	}
	if rob.Status == openapi.RobotStatusBusy {
		return apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "cannot delete robot while it is Busy"))
	}
	if err := r.repo.Delete(ctx, r.db, id); err != nil {
		return err
	}
	// Clean up Redis keys for the deleted robot. Best-effort: a failure here leaves
	// orphaned keys in Redis but does not affect correctness — the flush loop only
	// processes robots that exist in the DB.
	_ = r.uptimeDeltaRepo.Delete(ctx, id)
	return nil
}
