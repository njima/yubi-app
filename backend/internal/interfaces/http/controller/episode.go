package controller

import (
	"bytes"
	"context"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/ccontext"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/pagination"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

// parseDateRangeHalfOpen converts API date query parameters into half-open
// [from, to) time.Time bounds in UTC. The upper bound is the day after the
// requested to-date, so SQL "WHERE started_at < to" still includes the
// entire to-date that the user typed in. Returns a bad-request error if
// only one of from/to is specified — the date range is atomic and must
// either be fully specified or fully omitted.
func parseDateRangeHalfOpen(from, to *openapi_types.Date) (*time.Time, *time.Time, error) {
	if (from == nil) != (to == nil) {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "started_at_from and started_at_to must both be specified or both omitted"))
	}
	var fromT, toT *time.Time
	if from != nil {
		t := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
		fromT = &t
	}
	if to != nil {
		t := time.Date(to.Year(), to.Month(), to.Day()+1, 0, 0, 0, 0, time.UTC)
		toT = &t
	}
	return fromT, toT, nil
}

func groupExecutionsBySubTaskID(execs model.EpisodeSubTaskExecutions) map[string]model.EpisodeSubTaskExecutions {
	result := make(map[string]model.EpisodeSubTaskExecutions)
	for _, exec := range execs {
		result[exec.EpisodeSubTaskID] = append(result[exec.EpisodeSubTaskID], exec)
	}
	return result
}

// fetchTaskAndVersionMaps batch-fetches tasks and task versions referenced by
// the given episodes so callers can resolve display names without per-episode
// round trips.
func (c *controller) fetchTaskAndVersionMaps(ctx context.Context, eps model.Episodes) (map[string]*model.Task, map[string]*model.TaskVersion) {
	taskIDSet := make(map[string]struct{}, len(eps))
	tvIDSet := make(map[string]struct{}, len(eps))
	for _, e := range eps {
		if e.TaskID != "" {
			taskIDSet[e.TaskID] = struct{}{}
		}
		if e.TaskVersionID != "" {
			tvIDSet[e.TaskVersionID] = struct{}{}
		}
	}
	taskIDs := make([]string, 0, len(taskIDSet))
	for id := range taskIDSet {
		taskIDs = append(taskIDs, id)
	}
	tvIDs := make([]string, 0, len(tvIDSet))
	for id := range tvIDSet {
		tvIDs = append(tvIDs, id)
	}

	taskMap := make(map[string]*model.Task, len(taskIDs))
	if len(taskIDs) > 0 {
		tasks, err := c.taskUsecase.ListByIDs(ctx, taskIDs)
		if err != nil {
			c.logger.Error().Err(err).Msg("failed to fetch tasks for episode display name resolution")
		} else {
			for _, t := range tasks {
				taskMap[t.IDNatural] = t
			}
		}
	}

	tvMap := make(map[string]*model.TaskVersion, len(tvIDs))
	if len(tvIDs) > 0 {
		versions, err := c.taskVersionUsecase.ListByIDs(ctx, tvIDs)
		if err != nil {
			c.logger.Error().Err(err).Msg("failed to fetch task versions for episode display name resolution")
		} else {
			for _, v := range versions {
				tvMap[v.IDNatural] = v
			}
		}
	}
	return taskMap, tvMap
}

// applyTaskVersionDisplayName decorates the response with the server-resolved
// task version display name when both task and version are present in the
// supplied maps.
func applyTaskVersionDisplayName(resp *openapi.Episode, ep *model.Episode, taskMap map[string]*model.Task, tvMap map[string]*model.TaskVersion) {
	tk, ok := taskMap[ep.TaskID]
	if !ok {
		return
	}
	tv, ok := tvMap[ep.TaskVersionID]
	if !ok {
		return
	}
	resolved := tv.DisplayLabel(tk.Name)
	resp.TaskVersionDisplayName = &resolved
}

func episodeToResponse(ep model.Episode) openapi.Episode {
	resp := openapi.Episode{
		Id:            ep.IDNatural,
		LocationId:    ep.LocationID,
		UserId:        ep.UserID,
		RobotId:       ep.RobotID,
		Status:        openAPIEpisodeStatus(ep.Status),
		TaskId:        ep.TaskID,
		TaskVersionId: ep.TaskVersionID,
		StartedAt:     ep.StartedAt,
		EndedAt:       ep.FinishedAt,
		ErrorDetails:  ep.ErrorDetails,
		CreatedAt:     ep.CreatedAt,
		RecordedBy:    ep.RecordedByID,
		AverageGrade:  ep.AverageGrade,
		GradeCount:    &ep.GradeCount,
	}
	if len(ep.ParameterValues) > 0 {
		resp.ParameterValues = &ep.ParameterValues
	}
	return resp
}

// BuildEpisodeSubTasks assembles the API response for episode subtasks by
// joining subtask master definitions with their episode-specific records and
// executions. Exported for reuse by the SSE handler.
func BuildEpisodeSubTasks(subtaskMasters model.SubTasks, records model.EpisodeSubTasks, executions model.EpisodeSubTaskExecutions, parameterValues map[string]string) []openapi.EpisodeSubTask {
	recordMap := make(map[string]*model.EpisodeSubTask)
	for _, r := range records {
		recordMap[r.SubTaskID] = r
	}

	executionsMap := groupExecutionsBySubTaskID(executions)

	subtasks := make([]openapi.EpisodeSubTask, 0, len(subtaskMasters))
	for _, st := range subtaskMasters {
		apiSubtask := openapi.EpisodeSubTask{
			SubtaskId:  st.IDNatural,
			Name:       model.InterpolateSubTaskName(st.Name, parameterValues),
			OrderIndex: st.OrderIndex,
		}
		if record, ok := recordMap[st.IDNatural]; ok {
			apiSubtask.Id = record.IDNatural
			apiSubtask.Status = record.CollectionStatus

			if execs, ok := executionsMap[record.IDNatural]; ok {
				apiExecs := make([]openapi.EpisodeSubTaskExecution, 0, len(execs))
				for _, exec := range execs {
					apiExecs = append(apiExecs, openapi.EpisodeSubTaskExecution{
						Id:         exec.IDNatural,
						Status:     exec.ExecutionStatus,
						StartedAt:  exec.StartedAt,
						FinishedAt: exec.FinishedAt,
					})
				}
				apiSubtask.Executions = &apiExecs
			}
		}
		subtasks = append(subtasks, apiSubtask)
	}

	return subtasks
}

func (c *controller) ExportEpisodes(ctx context.Context, request openapi.ExportEpisodesRequestObject) (openapi.ExportEpisodesResponseObject, error) {
	fromT, toT, err := parseDateRangeHalfOpen(request.Params.StartedAtFrom, request.Params.StartedAtTo)
	if err != nil {
		return nil, err
	}
	listFilter := repository.EpisodeListFilter{
		TaskID:        request.Params.TaskId,
		TaskVersionID: request.Params.TaskVersionId,
		RobotID:       request.Params.RobotId,
		UserID:        request.Params.UserId,
		StartedAtFrom: fromT,
		StartedAtTo:   toT,
	}
	if request.Params.Status != nil {
		statuses, err := episodeStatuses(*request.Params.Status)
		if err != nil {
			return nil, err
		}
		listFilter.Statuses = statuses
	}

	filter := repository.EpisodeExportFilter{EpisodeListFilter: listFilter}

	csvBytes, err := c.episodeExportUsecase.Export(ctx, filter)
	if err != nil {
		return nil, err
	}

	return openapi.ExportEpisodes200TextcsvResponse{
		Body: bytes.NewReader(csvBytes),
		Headers: openapi.ExportEpisodes200ResponseHeaders{
			ContentDisposition: `attachment; filename="episodes_export.csv"`,
		},
	}, nil
}

func (c *controller) ListEpisodes(ctx context.Context, request openapi.ListEpisodesRequestObject) (openapi.ListEpisodesResponseObject, error) {
	pg := pagination.Parse(request.Params.Page, request.Params.Limit)

	fromT, toT, err := parseDateRangeHalfOpen(request.Params.StartedAtFrom, request.Params.StartedAtTo)
	if err != nil {
		return nil, err
	}
	filter := repository.EpisodeListFilter{
		TaskID:        request.Params.TaskId,
		TaskVersionID: request.Params.TaskVersionId,
		RobotID:       request.Params.RobotId,
		UserID:        request.Params.UserId,
		StartedAtFrom: fromT,
		StartedAtTo:   toT,
		SortBy:        episodeSortBy(request.Params.SortBy),
		SortOrder:     sortOrder(request.Params.SortOrder),
	}
	if request.Params.Status != nil {
		statuses, err := episodeStatuses(*request.Params.Status)
		if err != nil {
			return nil, err
		}
		filter.Statuses = statuses
	}

	eps, total, err := c.episodeUsecase.List(ctx, filter, pg.Page, pg.Limit)
	if err != nil {
		return nil, err
	}

	taskMap, tvMap := c.fetchTaskAndVersionMaps(ctx, eps)
	episodes := make([]openapi.Episode, 0, len(eps))
	for _, e := range eps {
		resp := episodeToResponse(*e)
		applyTaskVersionDisplayName(&resp, e, taskMap, tvMap)
		episodes = append(episodes, resp)
	}

	return openapi.ListEpisodes200JSONResponse{
		Episodes: episodes,
		Pagination: openapi.Pagination{
			Count: total,
			Page:  pg.Page,
			Limit: pg.Limit,
		},
	}, nil
}

func (c *controller) CreateEpisode(ctx context.Context, request openapi.CreateEpisodeRequestObject) (openapi.CreateEpisodeResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	userID, err := ccontext.UserID(ctx)
	if err != nil {
		return nil, err
	}

	body := *request.Body

	input := usecase.EpisodeCreateInput{
		OrganizationID: body.OrganizationId,
		LocationID:     body.LocationId,
		TaskID:         body.TaskId,
		TaskVersionID:  body.TaskVersionId,
		RobotID:        body.RobotId,
		UserID:         userID,
		RecordedByID:   body.RecordedBy,
	}
	if body.ParameterValues != nil {
		input.ParameterValues = *body.ParameterValues
	}

	ep, err := c.episodeUsecase.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	resp := episodeToResponse(ep)
	taskMap, tvMap := c.fetchTaskAndVersionMaps(ctx, model.Episodes{&ep})
	applyTaskVersionDisplayName(&resp, &ep, taskMap, tvMap)
	return openapi.CreateEpisode201JSONResponse(resp), nil
}

func (c *controller) CreateEpisodesBulk(ctx context.Context, request openapi.CreateEpisodesBulkRequestObject) (openapi.CreateEpisodesBulkResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	body := *request.Body
	if body.Count < 1 || body.Count > 100 {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "count must be between 1 and 100"))
	}

	userID, err := ccontext.UserID(ctx)
	if err != nil {
		return nil, err
	}

	input := usecase.EpisodeCreateInput{
		OrganizationID: body.OrganizationId,
		LocationID:     body.LocationId,
		TaskID:         body.TaskId,
		TaskVersionID:  body.TaskVersionId,
		RobotID:        body.RobotId,
		UserID:         userID,
		RecordedByID:   body.RecordedBy,
	}
	if body.ParameterValues != nil {
		input.ParameterValues = *body.ParameterValues
	}

	eps, err := c.episodeUsecase.BulkCreate(ctx, input, body.Count)
	if err != nil {
		return nil, err
	}

	taskMap, tvMap := c.fetchTaskAndVersionMaps(ctx, eps)
	resp := make([]openapi.Episode, 0, len(eps))
	for _, ep := range eps {
		r := episodeToResponse(*ep)
		applyTaskVersionDisplayName(&r, ep, taskMap, tvMap)
		resp = append(resp, r)
	}

	return openapi.CreateEpisodesBulk201JSONResponse(resp), nil
}

func (c *controller) DeleteEpisodeById(ctx context.Context, request openapi.DeleteEpisodeByIdRequestObject) (openapi.DeleteEpisodeByIdResponseObject, error) {
	if err := c.episodeUsecase.Delete(ctx, request.EpisodeId); err != nil {
		return nil, err
	}
	return openapi.DeleteEpisodeById204Response{}, nil
}

func (c *controller) GetEpisodeById(ctx context.Context, request openapi.GetEpisodeByIdRequestObject) (openapi.GetEpisodeByIdResponseObject, error) {
	ep, err := c.episodeUsecase.GetByID(ctx, request.EpisodeId)
	if err != nil {
		return nil, err
	}

	subtaskMasters, records, executions, err := c.episodeUsecase.GetSubTasksByEpisodeID(ctx, ep.IDNatural, ep.TaskVersionID)
	if err != nil {
		return nil, err
	}

	subtasks := BuildEpisodeSubTasks(subtaskMasters, records, executions, ep.ParameterValues)

	var taskName *string
	var taskDescription *string
	if tk, err := c.taskUsecase.GetByID(ctx, ep.TaskID); err != nil {
		c.logger.Error().Err(err).Str("task_id", ep.TaskID).Msg("failed to fetch task for episode")
	} else {
		taskName = &tk.Name
		taskDescription = tk.Description
	}

	resp := episodeToResponse(ep)
	resp.Subtasks = &subtasks
	resp.TaskName = taskName
	resp.TaskDescription = taskDescription
	taskMap, tvMap := c.fetchTaskAndVersionMaps(ctx, model.Episodes{&ep})
	applyTaskVersionDisplayName(&resp, &ep, taskMap, tvMap)
	return openapi.GetEpisodeById200JSONResponse(resp), nil
}

func (c *controller) UpdateEpisodeById(ctx context.Context, request openapi.UpdateEpisodeByIdRequestObject) (openapi.UpdateEpisodeByIdResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	body := *request.Body
	status, err := episodeStatusModel(body.Status)
	if err != nil {
		return nil, err
	}

	input := usecase.EpisodeUpdateInput{
		ID:           request.EpisodeId,
		StartedAt:    body.StartTime,
		FinishedAt:   body.EndTime,
		Status:       status,
		ErrorDetails: body.ErrorDetails,
		RecordedByID: body.RecordedBy,
	}

	ep, err := c.episodeUsecase.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	resp := episodeToResponse(ep)
	taskMap, tvMap := c.fetchTaskAndVersionMaps(ctx, model.Episodes{&ep})
	applyTaskVersionDisplayName(&resp, &ep, taskMap, tvMap)
	return openapi.UpdateEpisodeById200JSONResponse(resp), nil
}

func (c *controller) GetRobotEpisodeById(ctx context.Context, request openapi.GetRobotEpisodeByIdRequestObject) (openapi.GetRobotEpisodeByIdResponseObject, error) {
	robotID, err := ccontext.RobotID(ctx)
	if err != nil {
		return nil, err
	}

	ep, err := c.episodeUsecase.GetByID(ctx, request.EpisodeId)
	if err != nil {
		return nil, err
	}

	if ep.RobotID != robotID {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeForbidden, "robot is not authorized to access this episode"))
	}

	subtaskMasters, records, executions, err := c.episodeUsecase.GetSubTasksByEpisodeID(ctx, ep.IDNatural, ep.TaskVersionID)
	if err != nil {
		return nil, err
	}

	subtasks := BuildEpisodeSubTasks(subtaskMasters, records, executions, ep.ParameterValues)

	var taskName *string
	var taskDescription *string
	if tk, err := c.taskUsecase.GetByID(ctx, ep.TaskID); err != nil {
		c.logger.Error().Err(err).Str("task_id", ep.TaskID).Msg("failed to fetch task for robot episode")
	} else {
		taskName = &tk.Name
		taskDescription = tk.Description
	}

	resp := episodeToResponse(ep)
	resp.Subtasks = &subtasks
	resp.TaskName = taskName
	resp.TaskDescription = taskDescription
	taskMap, tvMap := c.fetchTaskAndVersionMaps(ctx, model.Episodes{&ep})
	applyTaskVersionDisplayName(&resp, &ep, taskMap, tvMap)
	return openapi.GetRobotEpisodeById200JSONResponse(resp), nil
}

func (c *controller) RepeatLastRobotEpisode(ctx context.Context, request openapi.RepeatLastRobotEpisodeRequestObject) (openapi.RepeatLastRobotEpisodeResponseObject, error) {
	ep, err := c.episodeUsecase.RepeatLast(ctx)
	if err != nil {
		return nil, err
	}

	resp := episodeToResponse(ep)
	taskMap, tvMap := c.fetchTaskAndVersionMaps(ctx, model.Episodes{&ep})
	applyTaskVersionDisplayName(&resp, &ep, taskMap, tvMap)
	return openapi.RepeatLastRobotEpisode201JSONResponse(resp), nil
}

func (c *controller) ListRobotEpisodes(ctx context.Context, request openapi.ListRobotEpisodesRequestObject) (openapi.ListRobotEpisodesResponseObject, error) {
	robotID, err := ccontext.RobotID(ctx)
	if err != nil {
		return nil, err
	}

	filter := repository.EpisodeListFilter{
		RobotID:  &robotID,
		Statuses: []repository.EpisodeStatus{repository.EpisodeStatusReady},
	}

	eps, _, err := c.episodeUsecase.List(ctx, filter, 1, 100)
	if err != nil {
		return nil, err
	}

	taskMap, tvMap := c.fetchTaskAndVersionMaps(ctx, eps)
	resp := make([]openapi.Episode, 0, len(eps))
	for _, e := range eps {
		ep := episodeToResponse(*e)

		if tk, ok := taskMap[e.TaskID]; ok {
			ep.TaskName = &tk.Name
			ep.TaskDescription = tk.Description
		}
		applyTaskVersionDisplayName(&ep, e, taskMap, tvMap)

		resp = append(resp, ep)
	}

	return openapi.ListRobotEpisodes200JSONResponse(resp), nil
}

func (c *controller) GetEpisodeRecordings(ctx context.Context, request openapi.GetEpisodeRecordingsRequestObject) (openapi.GetEpisodeRecordingsResponseObject, error) {
	recordings, err := c.episodeUsecase.GetRecordings(ctx, request.EpisodeId)
	if err != nil {
		return nil, err
	}
	return openapi.GetEpisodeRecordings200JSONResponse{
		Recordings: recordings,
	}, nil
}

func (c *controller) GetEpisodeStats(ctx context.Context, request openapi.GetEpisodeStatsRequestObject) (openapi.GetEpisodeStatsResponseObject, error) {
	stats, err := c.episodeUsecase.GetStats(ctx, request.EpisodeId)
	if err != nil {
		return nil, err
	}

	apiStats := make(map[string]openapi.EpisodeFeatureStats, len(stats))
	for feature, s := range stats {
		apiStats[feature] = openapi.EpisodeFeatureStats{
			Min:   s.Min,
			Max:   s.Max,
			Mean:  s.Mean,
			Std:   s.Std,
			Count: s.Count,
		}
	}
	return openapi.GetEpisodeStats200JSONResponse{
		Stats: apiStats,
	}, nil
}
