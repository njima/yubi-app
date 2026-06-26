package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

type RobotUsecase interface {
	Create(ctx context.Context, input RobotCreateInput) (model.Robot, error)
	GetByID(ctx context.Context, id string) (model.Robot, error)
	List(ctx context.Context, filter RobotListFilter, page, limit int) (model.Robots, int, error)
	ListTypes(ctx context.Context, filter RobotTypeFilter) ([]string, error)
	Update(ctx context.Context, input RobotUpdateInput) (model.Robot, error)
	Delete(ctx context.Context, id string) error
}

type RobotCreateInput struct {
	OrganizationID string
	LocationID     string
	Name           string
	RobotType      *string
	LeaderStatus   *model.LeaderStatus
	RobotConfig    *json.RawMessage
}

type RobotUpdateInput struct {
	ID              string
	OrganizationID  string
	LocationID      string
	Name            *string
	RobotType       *string
	Status          *model.RobotStatus
	LeaderStatus    *model.LeaderStatus
	HasLeaderStatus bool // true when leader_status is explicitly provided (including null)
	LastHeartbeatAt *time.Time
	OfflineReason   *string
	RobotConfig     *json.RawMessage
}

type robot struct {
	repo            repository.Robot
	robotStatusRepo repository.RobotStatusRepository
	uptimeDeltaRepo repository.RobotUptimeDeltaRepository
	data            repository.DataAccess
}

func NewRobot(repo repository.Robot, robotStatusRepo repository.RobotStatusRepository, uptimeDeltaRepo repository.RobotUptimeDeltaRepository, data repository.DataAccess) *robot {
	return &robot{repo: repo, robotStatusRepo: robotStatusRepo, uptimeDeltaRepo: uptimeDeltaRepo, data: data}
}

func (r *robot) Create(ctx context.Context, input RobotCreateInput) (model.Robot, error) {
	ro, err := model.InitRobot(input.OrganizationID, input.LocationID, input.Name, input.RobotType, input.RobotConfig)
	if err != nil {
		return model.Robot{}, err
	}

	if input.LeaderStatus != nil {
		ro.SetLeaderStatus(input.LeaderStatus)
	}

	cro, err := r.repo.Create(ctx, r.data.Conn(), ro)
	if err != nil {
		return model.Robot{}, err
	}

	return cro, nil
}

func (r *robot) GetByID(ctx context.Context, id string) (model.Robot, error) {
	rob, err := r.repo.GetByID(ctx, r.data.Conn(), id)
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

func (r *robot) List(ctx context.Context, filter RobotListFilter, page, limit int) (model.Robots, int, error) {
	pg := pagination.Normalize(page, limit)

	dbFilter := filter.repositoryFilter()
	if dbFilter.Status != nil &&
		(*dbFilter.Status == repository.RobotFilterStatusOnline || *dbFilter.Status == repository.RobotFilterStatusOffline) {
		requestedStatus := *dbFilter.Status
		dbFilter.Status = nil

		onlineIDs, err := r.robotStatusRepo.GetAllOnlineRobotIDs(ctx)
		if err != nil {
			return nil, 0, err
		}

		if requestedStatus == repository.RobotFilterStatusOnline {
			if len(onlineIDs) == 0 {
				return model.Robots{}, 0, nil
			}
			dbFilter.OnlineRobotIDs = &onlineIDs
		} else {
			dbFilter.OnlineRobotIDs = &onlineIDs
			dbFilter.ExcludeOnline = true
		}
	}

	robs, total, err := r.repo.List(ctx, r.data.Conn(), dbFilter, pg.Limit, pg.Offset)
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
	rob, err := r.repo.GetByID(ctx, r.data.Conn(), input.ID)
	if err != nil {
		return model.Robot{}, err
	}

	rob, err = r.update(ctx, rob, input)
	if err != nil {
		return model.Robot{}, err
	}

	urob, err := r.repo.Update(ctx, r.data.Conn(), rob)
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
		wasFaulted := robot.Status == model.RobotStatusFaulted

		// Only Ready(5), Faulted(3), Maintenance(4) can be manually set
		switch *input.Status {
		case model.RobotStatusReady, model.RobotStatusFaulted, model.RobotStatusMaintenance:
			// allowed
		default:
			return model.Robot{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "status must be Ready, Faulted, or Maintenance"))
		}
		// Validate status transition based on current DB status
		switch robot.Status {
		case model.RobotStatusBusy:
			return model.Robot{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "cannot manually change status while robot is Busy"))
		case model.RobotStatusFaulted:
			if *input.Status != model.RobotStatusReady && *input.Status != model.RobotStatusMaintenance {
				return model.Robot{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "Faulted robot can only be set to Ready or Maintenance"))
			}
		case model.RobotStatusMaintenance:
			if *input.Status != model.RobotStatusReady && *input.Status != model.RobotStatusFaulted {
				return model.Robot{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "Maintenance robot can only be set to Ready or Faulted"))
			}
		}
		if err := robot.SetStatus(*input.Status); err != nil {
			return model.Robot{}, err
		}
		if *input.Status == model.RobotStatusFaulted {
			now := time.Now()
			robot.SetFaultStartedAt(&now)
		}
		if wasFaulted && *input.Status != model.RobotStatusFaulted {
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
		wasLeaderFaulted := robot.LeaderStatus != nil && *robot.LeaderStatus == model.LeaderStatusFaulted
		robot.SetLeaderStatus(input.LeaderStatus)
		if !wasLeaderFaulted && input.LeaderStatus != nil && *input.LeaderStatus == model.LeaderStatusFaulted {
			now := time.Now()
			robot.SetLeaderFaultStartedAt(&now)
		} else if wasLeaderFaulted && (input.LeaderStatus == nil || *input.LeaderStatus != model.LeaderStatusFaulted) {
			robot.SetLeaderFaultStartedAt(nil)
		}
	}

	return robot, nil
}

func (r *robot) ListTypes(ctx context.Context, filter RobotTypeFilter) ([]string, error) {
	dbFilter := filter.repositoryFilter()
	if dbFilter.Status != nil &&
		(*dbFilter.Status == repository.RobotFilterStatusOnline || *dbFilter.Status == repository.RobotFilterStatusOffline) {
		requestedStatus := *dbFilter.Status
		dbFilter.Status = nil

		onlineIDs, err := r.robotStatusRepo.GetAllOnlineRobotIDs(ctx)
		if err != nil {
			return nil, err
		}

		if requestedStatus == repository.RobotFilterStatusOnline {
			if len(onlineIDs) == 0 {
				return []string{}, nil
			}
			dbFilter.OnlineRobotIDs = &onlineIDs
		} else {
			dbFilter.OnlineRobotIDs = &onlineIDs
			dbFilter.ExcludeOnline = true
		}
	}
	return r.repo.ListTypes(ctx, r.data.Conn(), dbFilter)
}

func (r *robot) Delete(ctx context.Context, id string) error {
	rob, err := r.repo.GetByID(ctx, r.data.Conn(), id)
	if err != nil {
		return err
	}
	if rob.Status == model.RobotStatusBusy {
		return apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "cannot delete robot while it is Busy"))
	}
	if err := r.repo.Delete(ctx, r.data.Conn(), id); err != nil {
		return err
	}
	// Clean up Redis keys for the deleted robot. Best-effort: a failure here leaves
	// orphaned keys in Redis but does not affect correctness — the flush loop only
	// processes robots that exist in the DB.
	_ = r.uptimeDeltaRepo.Delete(ctx, id)
	return nil
}
