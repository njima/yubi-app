package usecase

import (
	"context"
	"math/rand"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/eventbus"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
	"github.com/rs/zerolog"
)

type StartEpisodeInput struct {
	EpisodeID    string
	OccurredAt   time.Time
	ActiveUserID *string // Live operator from Redis heartbeat; takes priority over RecordedByID and API key user.
}

type FinishEpisodeInput struct {
	EpisodeID  string
	OccurredAt time.Time
}

type CancelEpisodeInput struct {
	EpisodeID string
}

type EpisodeUsecase interface {
	Create(ctx context.Context, input EpisodeCreateInput) (model.Episode, error)
	BulkCreate(ctx context.Context, input EpisodeCreateInput, count int) (model.Episodes, error)
	GetByID(ctx context.Context, id string) (model.Episode, error)
	GetCurrentRobotEpisode(ctx context.Context, robotID string) (*model.Episode, error)
	GetSubTasksByEpisodeID(ctx context.Context, episodeID string, taskVersionID string) (model.SubTasks, model.EpisodeSubTasks, model.EpisodeSubTaskExecutions, error)
	List(ctx context.Context, filter EpisodeListFilter, page, limit int) (model.Episodes, int, error)
	Update(ctx context.Context, input EpisodeUpdateInput) (model.Episode, error)
	Delete(ctx context.Context, id string) error
	Start(ctx context.Context, input StartEpisodeInput) error
	Finish(ctx context.Context, input FinishEpisodeInput) error
	Cancel(ctx context.Context, input CancelEpisodeInput) error
	RepeatLast(ctx context.Context) (model.Episode, error)
	GetRecordings(ctx context.Context, episodeID string) (map[string]string, error)
	GetStats(ctx context.Context, episodeID string) (model.EpisodeRecordingStats, error)
}

type EpisodeCreateInput struct {
	OrganizationID  string
	LocationID      string
	TaskID          string
	TaskVersionID   *string
	RobotID         string
	UserID          string
	RecordedByID    *string
	ParameterValues map[string]string
}

type EpisodeUpdateInput struct {
	ID           string
	StartedAt    *time.Time
	FinishedAt   *time.Time
	Status       *model.EpisodeStatus
	ErrorDetails *string
	RecordedByID *string
}

type episode struct {
	repo            repository.Episode
	gradeRepo       repository.EpisodeGrade
	logger          zerolog.Logger
	tvRepo          repository.TaskVersion
	sr              repository.SubTask
	estr            repository.EpisodeSubTask
	execr           repository.EpisodeSubTaskExecution
	rr              repository.Robot
	robotStatusRepo repository.RobotStatusRepository
	recRepo         repository.EpisodeRecording
	taskRepo        repository.Task
	locRepo         repository.Location
	siteRepo        repository.Site
	data            repository.DataAccess

	// bus and robotBus are notified after database transactions commit so
	// SSE subscribers can refetch updated state. Because the notification
	// happens outside the transaction, there is a small window where a
	// subscriber's read may return stale data (PostgreSQL MVCC). This is
	// acceptable for the current single-writer setup; consider a post-commit
	// hook if read replicas are introduced.
	bus      *eventbus.Bus
	robotBus *eventbus.Bus
	listBus  *eventbus.Bus
}

type EpisodeDependencies struct {
	Repository               repository.Episode
	GradeRepository          repository.EpisodeGrade
	Logger                   zerolog.Logger
	TaskVersionRepository    repository.TaskVersion
	SubTaskRepository        repository.SubTask
	EpisodeSubTaskRepository repository.EpisodeSubTask
	ExecutionRepository      repository.EpisodeSubTaskExecution
	RobotRepository          repository.Robot
	RobotStatusRepository    repository.RobotStatusRepository
	RecordingRepository      repository.EpisodeRecording
	TaskRepository           repository.Task
	LocationRepository       repository.Location
	SiteRepository           repository.Site
	DataAccess               repository.DataAccess
	EventBus                 *eventbus.Bus
	RobotEventBus            *eventbus.Bus
	ListEventBus             *eventbus.Bus
}

func NewEpisode(deps EpisodeDependencies) *episode {
	return &episode{
		repo:            deps.Repository,
		gradeRepo:       deps.GradeRepository,
		logger:          deps.Logger,
		tvRepo:          deps.TaskVersionRepository,
		sr:              deps.SubTaskRepository,
		estr:            deps.EpisodeSubTaskRepository,
		execr:           deps.ExecutionRepository,
		rr:              deps.RobotRepository,
		robotStatusRepo: deps.RobotStatusRepository,
		recRepo:         deps.RecordingRepository,
		taskRepo:        deps.TaskRepository,
		locRepo:         deps.LocationRepository,
		siteRepo:        deps.SiteRepository,
		data:            deps.DataAccess,
		bus:             deps.EventBus,
		robotBus:        deps.RobotEventBus,
		listBus:         deps.ListEventBus,
	}
}

func (e *episode) Create(ctx context.Context, input EpisodeCreateInput) (model.Episode, error) {
	episodes, err := e.BulkCreate(ctx, input, 1)
	if err != nil {
		return model.Episode{}, err
	}
	if len(episodes) == 0 {
		return model.Episode{}, apperror.NewError(
			apperror.NewMessage(apperror.CodeInternal, "no episode created"),
		)
	}
	return *episodes[0], nil
}

func (e *episode) BulkCreate(ctx context.Context, input EpisodeCreateInput, count int) (model.Episodes, error) {
	if count < 1 {
		return nil, apperror.NewError(
			apperror.NewMessage(apperror.CodeBadRequest, "count must be greater than 0"),
		)
	}

	var (
		tv  model.TaskVersion
		err error
	)
	if input.TaskVersionID != nil {
		tv, err = e.tvRepo.GetByID(ctx, e.data.Conn(), *input.TaskVersionID)
		if err != nil {
			return nil, err
		}
		if tv.TaskID != input.TaskID {
			return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "task_version_id does not belong to task_id"))
		}
		if !tv.IsApproved() {
			return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "task_version is not approved"))
		}
	} else {
		tv, err = e.tvRepo.GetLatestApprovedByTaskID(ctx, e.data.Conn(), input.TaskID)
		if err != nil {
			return nil, err
		}
	}

	// Validate that location_id matches the robot's location
	robot, err := e.rr.GetByID(ctx, e.data.Conn(), input.RobotID)
	if err != nil {
		return nil, err
	}
	if robot.LocationID != input.LocationID {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "location_id does not match robot's location"))
	}

	subtasks, err := e.sr.GetByTaskVersionID(ctx, e.data.Conn(), tv.IDNatural)
	if err != nil {
		return nil, err
	}

	episodes := make(model.Episodes, 0, count)
	err = e.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		allEpisodeSubTasks := make([]model.EpisodeSubTask, 0, len(subtasks)*count)

		for i := 0; i < count; i++ {
			paramValues, err := resolveParameterValues(tv.Parameters, input.ParameterValues)
			if err != nil {
				return err
			}

			inep, err := model.InitEpisode(input.OrganizationID, input.TaskID, input.LocationID, input.RobotID, input.UserID, input.RecordedByID)
			if err != nil {
				return err
			}

			if err := inep.SetTaskVersionID(tv.IDNatural); err != nil {
				return err
			}
			inep.ParameterValues = paramValues

			ep, err := e.repo.Create(ctx, conn, inep)
			if err != nil {
				return err
			}

			epCopy := ep
			episodes = append(episodes, &epCopy)

			for _, st := range subtasks {
				est, err := model.InitEpisodeSubTask(
					input.OrganizationID,
					ep.IDNatural,
					st.IDNatural,
				)
				if err != nil {
					return err
				}
				allEpisodeSubTasks = append(allEpisodeSubTasks, est)
			}
		}

		if err := e.estr.BulkCreate(ctx, conn, allEpisodeSubTasks); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	e.robotBus.Notify(input.RobotID)
	e.listBus.Notify("list")
	return episodes, nil
}

func (e *episode) GetByID(ctx context.Context, id string) (model.Episode, error) {
	ep, err := e.repo.GetByID(ctx, e.data.Conn(), id)
	if err != nil {
		return model.Episode{}, err
	}
	e.mergeGradeAggregates(ctx, model.Episodes{&ep})
	return ep, nil
}

func (e *episode) GetCurrentRobotEpisode(ctx context.Context, robotID string) (*model.Episode, error) {
	ep, err := e.repo.GetCurrentRobotEpisode(ctx, e.data.Conn(), robotID)
	if err != nil {
		return nil, err
	}
	if ep == nil {
		return nil, nil
	}
	e.mergeGradeAggregates(ctx, model.Episodes{ep})
	return ep, nil
}

func (e *episode) GetSubTasksByEpisodeID(ctx context.Context, episodeID string, taskVersionID string) (model.SubTasks, model.EpisodeSubTasks, model.EpisodeSubTaskExecutions, error) {
	subtasks, err := e.sr.GetByTaskVersionID(ctx, e.data.Conn(), taskVersionID)
	if err != nil {
		return nil, nil, nil, err
	}

	records, err := e.estr.GetByEpisodeID(ctx, e.data.Conn(), episodeID)
	if err != nil {
		return nil, nil, nil, err
	}

	if len(records) == 0 {
		return subtasks, records, nil, nil
	}

	subTaskIDs := make([]string, 0, len(records))
	for _, r := range records {
		subTaskIDs = append(subTaskIDs, r.IDNatural)
	}

	executions, err := e.execr.GetByEpisodeSubTaskIDs(ctx, e.data.Conn(), subTaskIDs)
	if err != nil {
		return nil, nil, nil, err
	}

	return subtasks, records, executions, nil
}

func (e *episode) List(ctx context.Context, filter EpisodeListFilter, page, limit int) (model.Episodes, int, error) {
	if filter.TaskID != nil && filter.TaskVersionID != nil {
		tv, err := e.tvRepo.GetByID(ctx, e.data.Conn(), *filter.TaskVersionID)
		if err != nil {
			return nil, 0, err
		}
		if tv.TaskID != *filter.TaskID {
			return nil, 0, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "task_version_id does not belong to task_id: task_version_id=%s task_id=%s", *filter.TaskVersionID, *filter.TaskID))
		}
	}

	pg := pagination.Normalize(page, limit)
	episodes, total, err := e.repo.List(ctx, e.data.Conn(), filter.repositoryFilter(), pg.Limit, pg.Offset)
	if err != nil {
		return nil, 0, err
	}

	e.mergeGradeAggregates(ctx, episodes)

	return episodes, total, nil
}

// mergeGradeAggregates swallows repo errors so the core episode read path
// stays up when episode_grade is unavailable (migration skew, replica lag).
func (e *episode) mergeGradeAggregates(ctx context.Context, episodes model.Episodes) {
	if len(episodes) == 0 {
		return
	}

	ids := make([]string, 0, len(episodes))
	for _, ep := range episodes {
		ids = append(ids, ep.IDNatural)
	}

	aggMap, err := e.gradeRepo.GetAverageMap(ctx, e.data.Conn(), ids)
	if err != nil {
		e.logger.Warn().Err(err).Int("episode_count", len(episodes)).Msg("episode: grade aggregate fetch failed, returning episodes without grades")
		return
	}

	for _, ep := range episodes {
		if agg, ok := aggMap[ep.IDNatural]; ok {
			avg := agg.Average
			ep.AverageGrade = &avg
			ep.GradeCount = agg.Count
		}
	}
}

func (e *episode) Update(ctx context.Context, input EpisodeUpdateInput) (model.Episode, error) {
	ep, err := e.repo.GetByID(ctx, e.data.Conn(), input.ID)
	if err != nil {
		return model.Episode{}, err
	}

	ep, err = e.update(ctx, ep, input)
	if err != nil {
		return model.Episode{}, err
	}

	result, err := e.repo.Update(ctx, e.data.Conn(), ep)
	if err != nil {
		return model.Episode{}, err
	}

	e.bus.Notify(input.ID)
	e.robotBus.Notify(ep.RobotID)
	e.listBus.Notify("list")
	return result, nil
}

func (e *episode) update(ctx context.Context, ep model.Episode, input EpisodeUpdateInput) (model.Episode, error) {
	if input.StartedAt != nil {
		if err := ep.SetStartedAt(*input.StartedAt); err != nil {
			return model.Episode{}, err
		}
	}
	if input.FinishedAt != nil {
		if err := ep.SetFinishedAt(*input.FinishedAt); err != nil {
			return model.Episode{}, err
		}
	}
	if input.Status != nil {
		if err := ep.SetStatus(*input.Status); err != nil {
			return model.Episode{}, err
		}
	}
	if input.ErrorDetails != nil {
		if err := ep.SetErrorDetails(*input.ErrorDetails); err != nil {
			return model.Episode{}, err
		}
	}
	if input.RecordedByID != nil {
		if err := ep.SetRecordedByID(*input.RecordedByID); err != nil {
			return model.Episode{}, err
		}
	}

	return ep, nil
}

func (e *episode) Delete(ctx context.Context, id string) error {
	return e.repo.Delete(ctx, e.data.Conn(), id)
}

func (e *episode) Start(ctx context.Context, input StartEpisodeInput) error {
	robotID, err := requestctx.RobotID(ctx)
	if err != nil {
		return err
	}
	userID, err := requestctx.UserID(ctx)
	if err != nil {
		return err
	}

	err = e.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		episode, err := e.repo.GetByID(ctx, conn, input.EpisodeID)
		if err != nil {
			return err
		}

		if episode.RobotID != robotID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeForbidden, "robot is not authorized to operate this episode"))
		}

		if err := episode.Start(input.OccurredAt); err != nil {
			return err
		}

		robot, err := e.rr.GetByID(ctx, conn, robotID)
		if err != nil {
			return err
		}

		robotStatus, err := e.robotStatusRepo.GetByRobotID(ctx, robot.IDNatural)
		if err != nil {
			return err
		}
		robot.ResolvedStatus(robotStatus != nil)

		// Priority: live teleop operator > episode.RecordedByID > API key user
		activeUserID := userID
		if input.ActiveUserID != nil {
			activeUserID = *input.ActiveUserID
			// Stamp the episode with the actual operator so the episode
			// record is self-contained for audit/analytics.
			if episode.RecordedByID == nil {
				if err := episode.SetRecordedByID(activeUserID); err != nil {
					return err
				}
			}
		} else if episode.RecordedByID != nil {
			activeUserID = *episode.RecordedByID
		}
		if err := robot.StartTeleoperation(input.EpisodeID, activeUserID); err != nil {
			return err
		}

		if _, err := e.repo.Update(ctx, conn, episode); err != nil {
			return err
		}

		if _, err := e.rr.Update(ctx, conn, robot); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	e.bus.Notify(input.EpisodeID)
	e.robotBus.Notify(robotID)
	e.listBus.Notify("list")
	return nil
}

func (e *episode) Finish(ctx context.Context, input FinishEpisodeInput) error {
	robotID, err := requestctx.RobotID(ctx)
	if err != nil {
		return err
	}

	err = e.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		episode, err := e.repo.GetByID(ctx, conn, input.EpisodeID)
		if err != nil {
			return err
		}

		if episode.RobotID != robotID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeForbidden, "robot is not authorized to operate this episode"))
		}

		if err := episode.Finish(input.OccurredAt); err != nil {
			return err
		}

		robot, err := e.rr.GetByID(ctx, conn, robotID)
		if err != nil {
			return err
		}

		if err := robot.EndTeleoperation(); err != nil {
			return err
		}

		if _, err := e.repo.Update(ctx, conn, episode); err != nil {
			return err
		}

		if _, err := e.rr.Update(ctx, conn, robot); err != nil {
			return err
		}

		// Auto-update task status based on collection progress
		actual, err := e.repo.SumDurationByTaskID(ctx, conn, episode.TaskID)
		if err != nil {
			return err
		}
		target, err := e.tvRepo.SumTargetByTaskID(ctx, conn, episode.TaskID)
		if err != nil {
			return err
		}
		tk, err := e.taskRepo.GetByID(ctx, conn, episode.TaskID)
		if err != nil {
			return err
		}
		if tk.Status != nil && *tk.Status == model.TaskStatusCanceled {
			return nil
		}
		newStatus := model.DetermineTaskStatus(actual, target)
		if tk.Status != nil && *tk.Status == newStatus {
			return nil
		}
		updateTask := model.Task{IDNatural: episode.TaskID}
		if err := updateTask.SetStatus(&newStatus); err != nil {
			return err
		}
		if _, err := e.taskRepo.Update(ctx, conn, updateTask); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	e.bus.Notify(input.EpisodeID)
	e.robotBus.Notify(robotID)
	e.listBus.Notify("list")
	return nil
}

func (e *episode) RepeatLast(ctx context.Context) (model.Episode, error) {
	robotID, err := requestctx.RobotID(ctx)
	if err != nil {
		return model.Episode{}, err
	}

	robot, err := e.rr.GetByID(ctx, e.data.Conn(), robotID)
	if err != nil {
		return model.Episode{}, err
	}
	if robot.ActiveEpisodeID != nil {
		return model.Episode{}, apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "robot already has an active episode: %s", *robot.ActiveEpisodeID))
	}

	// Workaround: repo.List does not populate OrganizationID on the returned
	// model (field is omitted in the manual mapping), so we read it from the
	// auth context instead of copying from the last episode.
	organizationID, err := requestctx.OrganizationID(ctx)
	if err != nil {
		return model.Episode{}, err
	}

	statuses := []repository.EpisodeStatus{
		repository.EpisodeStatusCompleted,
		repository.EpisodeStatusCancel,
	}
	eps, _, err := e.repo.List(ctx, e.data.Conn(), repository.EpisodeListFilter{
		RobotID:  &robotID,
		Statuses: statuses,
	}, 1, 0)
	if err != nil {
		return model.Episode{}, err
	}
	if len(eps) == 0 {
		return model.Episode{}, apperror.NewError(
			apperror.NewMessage(apperror.CodeEpisodeNotFound, "no completed or cancelled episodes found for robot"))
	}

	last := eps[0]
	return e.Create(ctx, EpisodeCreateInput{
		OrganizationID:  organizationID,
		LocationID:      last.LocationID,
		TaskID:          last.TaskID,
		RobotID:         robotID,
		UserID:          last.UserID,
		ParameterValues: last.ParameterValues,
	})
}

func (e *episode) Cancel(ctx context.Context, input CancelEpisodeInput) error {
	robotID, err := requestctx.RobotID(ctx)
	if err != nil {
		return err
	}

	err = e.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		episode, err := e.repo.GetByID(ctx, conn, input.EpisodeID)
		if err != nil {
			return err
		}

		if episode.RobotID != robotID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeForbidden, "robot is not authorized to operate this episode"))
		}

		if err := episode.Cancel(); err != nil {
			return err
		}

		robot, err := e.rr.GetByID(ctx, conn, robotID)
		if err != nil {
			return err
		}

		if err := robot.EndTeleoperation(); err != nil {
			return err
		}

		if err := e.estr.BulkCancelByEpisodeID(ctx, conn, input.EpisodeID); err != nil {
			return err
		}

		if err := e.execr.BulkCancelByEpisodeID(ctx, conn, input.EpisodeID); err != nil {
			return err
		}

		if _, err := e.repo.Update(ctx, conn, episode); err != nil {
			return err
		}

		if _, err := e.rr.Update(ctx, conn, robot); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	e.bus.Notify(input.EpisodeID)
	e.robotBus.Notify(robotID)
	e.listBus.Notify("list")
	return nil
}

func (e *episode) buildPreviewPath(ctx context.Context, ep model.Episode) (model.EpisodePreviewPath, error) {
	if ep.StartedAt == nil {
		return model.EpisodePreviewPath{}, apperror.NewError(apperror.NewMessage(apperror.CodeInternal, "episode has no started_at"))
	}

	robot, err := e.rr.GetByID(ctx, e.data.Conn(), ep.RobotID)
	if err != nil {
		return model.EpisodePreviewPath{}, err
	}

	loc, err := e.locRepo.GetByID(ctx, e.data.Conn(), ep.LocationID)
	if err != nil {
		return model.EpisodePreviewPath{}, err
	}

	site, err := e.siteRepo.GetByID(ctx, e.data.Conn(), loc.SiteID)
	if err != nil {
		return model.EpisodePreviewPath{}, err
	}

	robotType := ""
	if robot.RobotType != nil {
		robotType = *robot.RobotType
	}

	return model.EpisodePreviewPath{
		UUID:         ep.IDNatural,
		Organization: robot.OrganizationName,
		Site:         site.Name,
		Location:     loc.Name,
		RobotType:    robotType,
		RobotID:      robot.IDNatural,
		StartedAt:    *ep.StartedAt,
	}, nil
}

func (e *episode) GetRecordings(ctx context.Context, episodeID string) (map[string]string, error) {
	ep, err := e.repo.GetByID(ctx, e.data.Conn(), episodeID)
	if err != nil {
		return nil, err
	}
	path, err := e.buildPreviewPath(ctx, ep)
	if err != nil {
		return nil, err
	}
	return e.recRepo.GetRecordingURLs(ctx, path)
}

func (e *episode) GetStats(ctx context.Context, episodeID string) (model.EpisodeRecordingStats, error) {
	ep, err := e.repo.GetByID(ctx, e.data.Conn(), episodeID)
	if err != nil {
		return nil, err
	}
	path, err := e.buildPreviewPath(ctx, ep)
	if err != nil {
		return nil, err
	}
	return e.recRepo.GetStats(ctx, path)
}

// resolveParameterValues validates provided values and fills in missing keys randomly.
func resolveParameterValues(params []model.TaskVersionParameter, provided map[string]string) (map[string]string, error) {
	if len(params) == 0 {
		return nil, nil
	}

	if err := model.ValidateParameterValues(params, provided); err != nil {
		return nil, err
	}

	resolved := make(map[string]string, len(params))
	for k, v := range provided {
		resolved[k] = v
	}
	for _, p := range params {
		if _, ok := resolved[p.Key]; !ok {
			resolved[p.Key] = p.Values[rand.Intn(len(p.Values))]
		}
	}
	return resolved, nil
}
