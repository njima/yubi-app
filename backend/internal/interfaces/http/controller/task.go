package controller

import (
	"bytes"
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

func toOpenAPITag(t model.TaskTag) openapi.TaskTag {
	return openapi.TaskTag{
		Id:               t.ID,
		Name:             t.Name,
		CategoryTypeId:   t.CategoryTypeID,
		CategoryTypeName: t.CategoryTypeName,
	}
}

func toOpenAPITags(tags model.TaskTags) *[]openapi.TaskTag {
	if len(tags) == 0 {
		return nil
	}
	result := make([]openapi.TaskTag, 0, len(tags))
	for _, t := range tags {
		result = append(result, toOpenAPITag(*t))
	}
	return &result
}

func toOpenAPITask(t model.Task) openapi.Task {
	var priority openapi.TaskPriority
	if t.Priority != nil {
		priority = openapi.TaskPriority(*t.Priority)
	}
	var difficulty openapi.TaskDifficulty
	if t.Difficulty != nil {
		difficulty = openapi.TaskDifficulty(*t.Difficulty)
	}
	var status openapi.TaskStatus
	if t.Status != nil {
		status = openapi.TaskStatus(*t.Status)
	}
	task := openapi.Task{
		Id:                    t.IDNatural,
		Name:                  t.Name,
		Description:           t.Description,
		ManualUrl:             t.ManualURL,
		Priority:              priority,
		Difficulty:            difficulty,
		Status:                status,
		Deadline:              t.Deadline,
		RobotType:             t.RobotType,
		TargetDurationSeconds: t.TargetDurationSeconds,
		TargetEpisodeCount:    t.TargetEpisodeCount,
		ActualEpisodeCount:    t.ActualEpisodeCount,
		Tags:                  toOpenAPITags(t.Tags),
	}
	if t.Version != "" {
		task.Version = &t.Version
		tv := model.TaskVersion{Version: t.Version, DisplayName: t.VersionDisplayName}
		resolved := tv.DisplayLabel(t.Name)
		task.VersionDisplayName = &resolved
	}
	return task
}

func (c *controller) ListTasks(ctx context.Context, request openapi.ListTasksRequestObject) (openapi.ListTasksResponseObject, error) {
	pg := pagination.Parse(request.Params.Page, request.Params.Limit)

	filter := usecase.TaskListFilter{
		HasApprovedVersion: request.Params.HasApprovedVersion,
		SortBy:             taskSortBy(request.Params.SortBy),
		SortOrder:          sortOrder(request.Params.SortOrder),
		RobotType:          request.Params.RobotType,
		Search:             request.Params.Search,
	}
	if request.Params.Status != nil {
		filter.Statuses = taskStatuses(*request.Params.Status)
	}
	if request.Params.Priority != nil {
		filter.Priorities = taskPriorities(*request.Params.Priority)
	}
	if request.Params.Difficulty != nil {
		filter.Difficulties = taskDifficulties(*request.Params.Difficulty)
	}

	tasks, total, err := c.taskUsecase.List(ctx, filter, pg.Page, pg.Limit)
	if err != nil {
		return nil, err
	}

	taskList := make([]openapi.Task, 0, len(tasks))
	for _, t := range tasks {
		taskList = append(taskList, toOpenAPITask(*t))
	}

	return openapi.ListTasks200JSONResponse{
		Tasks: taskList,
		Pagination: openapi.Pagination{
			Count: total,
			Page:  pg.Page,
			Limit: pg.Limit,
		},
	}, nil
}

func (c *controller) CreateTask(ctx context.Context, request openapi.CreateTaskRequestObject) (openapi.CreateTaskResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	input := usecase.TaskCreateInput{
		OrganizationID: request.Body.OrganizationId,
		Name:           request.Body.Name,
		Description:    request.Body.Description,
		ManualURL:      request.Body.ManualUrl,
		Priority:       model.TaskPriority(request.Body.Priority),
		Difficulty:     model.TaskDifficulty(request.Body.Difficulty),
		Status:         model.TaskStatus(request.Body.Status),
		Deadline:       request.Body.Deadline,
		RobotType:      request.Body.RobotType,
	}
	if request.Body.TagIds != nil {
		input.TagIDs = *request.Body.TagIds
	}

	tk, err := c.taskUsecase.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return openapi.CreateTask201JSONResponse(toOpenAPITask(tk)), nil
}

func (c *controller) DeleteTaskById(ctx context.Context, request openapi.DeleteTaskByIdRequestObject) (openapi.DeleteTaskByIdResponseObject, error) {
	if err := c.taskUsecase.Delete(ctx, request.TaskId); err != nil {
		return nil, err
	}
	return openapi.DeleteTaskById204Response{}, nil
}

func (c *controller) GetTaskById(ctx context.Context, request openapi.GetTaskByIdRequestObject) (openapi.GetTaskByIdResponseObject, error) {
	tk, err := c.taskUsecase.GetByID(ctx, request.TaskId)
	if err != nil {
		return nil, err
	}

	return openapi.GetTaskById200JSONResponse(toOpenAPITask(tk)), nil
}

func (c *controller) UpdateTaskById(ctx context.Context, request openapi.UpdateTaskByIdRequestObject) (openapi.UpdateTaskByIdResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	input := usecase.TaskUpdateInput{
		ID: request.TaskId,
	}
	if request.Body.Name != nil {
		input.Name = request.Body.Name
	}
	if request.Body.Description != nil {
		input.Description = request.Body.Description
	}
	if request.Body.ManualUrl != nil {
		input.ManualURL = request.Body.ManualUrl
	}
	if request.Body.Priority != nil {
		input.Priority = taskPriorityPtr(*request.Body.Priority)
	}
	if request.Body.Difficulty != nil {
		input.Difficulty = taskDifficultyPtr(*request.Body.Difficulty)
	}
	if request.Body.Status != nil {
		input.Status = taskStatusPtr(*request.Body.Status)
	}
	if request.Body.Deadline != nil {
		input.Deadline = request.Body.Deadline
	}
	if request.Body.RobotType != nil {
		input.RobotType = request.Body.RobotType
	}
	if request.Body.TagIds != nil {
		input.TagIDs = request.Body.TagIds
	}

	tk, err := c.taskUsecase.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateTaskById200JSONResponse(toOpenAPITask(tk)), nil
}

func (c *controller) CreateTaskVersion(ctx context.Context, request openapi.CreateTaskVersionRequestObject) (openapi.CreateTaskVersionResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	input := usecase.TaskVersionCreateInput{
		TaskID:                          request.TaskId,
		Version:                         request.Body.Version,
		DisplayName:                     request.Body.DisplayName,
		BaseTaskVersionID:               request.Body.BaseTaskVersionId,
		TargetDurationSeconds:           request.Body.TargetDurationSeconds,
		TargetEpisodeCount:              request.Body.TargetEpisodeCount,
		TargetDurationPerEpisodeSeconds: request.Body.TargetDurationPerEpisodeSeconds,
	}

	tv, err := c.taskVersionUsecase.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return openapi.CreateTaskVersion201JSONResponse(taskVersionToResponse(tv)), nil
}

func (c *controller) ListTaskVersions(ctx context.Context, request openapi.ListTaskVersionsRequestObject) (openapi.ListTaskVersionsResponseObject, error) {
	versions, err := c.taskVersionUsecase.ListByTaskID(ctx, request.TaskId)
	if err != nil {
		return nil, err
	}

	resp := make([]openapi.TaskVersion, 0, len(versions))
	for _, v := range versions {
		resp = append(resp, taskVersionToResponse(*v))
	}
	return openapi.ListTaskVersions200JSONResponse(resp), nil
}

func (c *controller) UpdateTaskVersion(ctx context.Context, request openapi.UpdateTaskVersionRequestObject) (openapi.UpdateTaskVersionResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	tv, err := c.taskVersionUsecase.Update(ctx, request.TaskId, usecase.TaskVersionUpdateInput{
		ID:                              request.VersionId,
		DisplayName:                     request.Body.DisplayName,
		TargetDurationSeconds:           request.Body.TargetDurationSeconds,
		TargetEpisodeCount:              request.Body.TargetEpisodeCount,
		TargetDurationPerEpisodeSeconds: request.Body.TargetDurationPerEpisodeSeconds,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateTaskVersion200JSONResponse(taskVersionToResponse(tv)), nil
}

func (c *controller) ApproveTaskVersion(ctx context.Context, request openapi.ApproveTaskVersionRequestObject) (openapi.ApproveTaskVersionResponseObject, error) {
	tv, err := c.taskVersionUsecase.Approve(ctx, request.TaskId, request.VersionId)
	if err != nil {
		return nil, err
	}

	return openapi.ApproveTaskVersion200JSONResponse(taskVersionToResponse(tv)), nil
}

func (c *controller) UpdateTaskVersionParameters(ctx context.Context, request openapi.UpdateTaskVersionParametersRequestObject) (openapi.UpdateTaskVersionParametersResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	tv, err := c.taskVersionUsecase.UpdateParameters(ctx, usecase.TaskVersionUpdateParametersInput{
		TaskID:     request.TaskId,
		VersionID:  request.VersionId,
		Parameters: openAPIToModelParameters(request.Body.Parameters),
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateTaskVersionParameters200JSONResponse(taskVersionToResponse(tv)), nil
}

func taskVersionToResponse(tv model.TaskVersion) openapi.TaskVersion {
	resp := openapi.TaskVersion{
		Id:                              tv.IDNatural,
		TaskId:                          tv.TaskID,
		Version:                         tv.Version,
		DisplayName:                     tv.DisplayName,
		IsCurrent:                       tv.IsCurrent,
		ApprovalStatus:                  openAPIApprovalStatus(tv.ApprovalStatus),
		CreatedAt:                       tv.CreatedAt,
		TargetDurationSeconds:           tv.TargetDurationSeconds,
		TargetEpisodeCount:              tv.TargetEpisodeCount,
		TargetDurationPerEpisodeSeconds: tv.TargetDurationPerEpisodeSeconds,
	}
	if tv.ActualDurationSeconds != nil {
		v := int(*tv.ActualDurationSeconds)
		resp.ActualDurationSeconds = &v
	}
	if tv.ActualEpisodeCount != nil {
		resp.ActualEpisodeCount = tv.ActualEpisodeCount
	}
	if len(tv.Parameters) > 0 {
		params := modelToOpenAPIParameters(tv.Parameters)
		resp.Parameters = &params
	}
	return resp
}

func openAPIToModelParameters(params []openapi.TaskVersionParameter) []model.TaskVersionParameter {
	result := make([]model.TaskVersionParameter, len(params))
	for i, p := range params {
		result[i] = model.TaskVersionParameter{Key: p.Key, Values: p.Values}
	}
	return result
}

func modelToOpenAPIParameters(params []model.TaskVersionParameter) []openapi.TaskVersionParameter {
	result := make([]openapi.TaskVersionParameter, len(params))
	for i, p := range params {
		result[i] = openapi.TaskVersionParameter{Key: p.Key, Values: p.Values}
	}
	return result
}

func buildTaskSummaryFilter(robotTypes *[]string, categoryTypeID *string, tagID *[]string) usecase.TaskSummaryFilter {
	filter := usecase.TaskSummaryFilter{}
	if robotTypes != nil {
		filter.RobotTypes = *robotTypes
	}
	if categoryTypeID != nil {
		filter.CategoryTypeID = categoryTypeID
	}
	if tagID != nil {
		filter.TagIDs = *tagID
	}
	return filter
}

func (c *controller) GetTaskSummary(ctx context.Context, request openapi.GetTaskSummaryRequestObject) (openapi.GetTaskSummaryResponseObject, error) {
	filter := buildTaskSummaryFilter(request.Params.RobotType, request.Params.CategoryTypeId, request.Params.TagId)
	if request.Params.From != nil {
		t := request.Params.From.Time
		filter.DeadlineFrom = &t
	}
	if request.Params.To != nil {
		// to is inclusive, add 1 day for exclusive end
		t := request.Params.To.Time.AddDate(0, 0, 1)
		filter.DeadlineTo = &t
	}

	summary, err := c.taskUsecase.GetSummary(ctx, filter)
	if err != nil {
		return nil, err
	}

	return openapi.GetTaskSummary200JSONResponse{
		TotalTasks:            summary.TotalTasks,
		TargetDurationSeconds: summary.TargetDurationSeconds,
		TargetEpisodeCount:    summary.TargetEpisodeCount,
	}, nil
}

func (c *controller) GetTaskCompletionTrend(ctx context.Context, request openapi.GetTaskCompletionTrendRequestObject) (openapi.GetTaskCompletionTrendResponseObject, error) {
	filter := buildTaskSummaryFilter(request.Params.RobotType, request.Params.CategoryTypeId, request.Params.TagId)

	// Date range defaults: 2 months ago → 2 months ahead
	now := time.Now().UTC()
	from := now.AddDate(0, -2, 0)
	to := now.AddDate(0, 2, 0)
	if request.Params.From != nil {
		from = request.Params.From.Time
	}
	if request.Params.To != nil {
		to = request.Params.To.Time
	}

	// Only filter by DeadlineTo for SQL efficiency.
	// DeadlineFrom is not set because overdue tasks (deadline < from) need to be fetched too.
	toExclusive := to.AddDate(0, 0, 1)
	filter.DeadlineTo = &toExclusive

	interval := "2week"
	if request.Params.Interval != nil {
		interval = string(*request.Params.Interval)
	}

	trend, err := c.taskUsecase.GetCompletionTrend(ctx, filter, string(request.Params.GroupBy), from, to, interval)
	if err != nil {
		return nil, err
	}

	periods := make([]openapi.TaskTrendPeriod, 0, len(trend.Periods))
	for _, p := range trend.Periods {
		periods = append(periods, toOpenAPITrendPeriod(p))
	}

	return openapi.GetTaskCompletionTrend200JSONResponse(openapi.TaskCompletionTrend{
		Periods: periods,
	}), nil
}

func toOpenAPITrendPeriod(p model.TrendPeriod) openapi.TaskTrendPeriod {
	groups := make([]openapi.TaskTrendGroup, 0, len(p.Groups))
	for _, g := range p.Groups {
		groups = append(groups, openapi.TaskTrendGroup{
			Label:          g.Label,
			TargetTasks:    g.TargetTasks,
			ActualTasks:    g.ActualTasks,
			TargetDuration: g.TargetDuration,
			ActualDuration: g.ActualDuration,
			TargetEpisodes: g.TargetEpisodes,
			ActualEpisodes: g.ActualEpisodes,
		})
	}

	result := openapi.TaskTrendPeriod{Groups: groups}
	if !p.Start.IsZero() {
		result.Start = &p.Start
	}
	if !p.End.IsZero() {
		result.End = &p.End
	}
	return result
}

func (c *controller) GetTaskAvailableTags(ctx context.Context, request openapi.GetTaskAvailableTagsRequestObject) (openapi.GetTaskAvailableTagsResponseObject, error) {
	var robotTypes []string
	if request.Params.RobotType != nil {
		robotTypes = *request.Params.RobotType
	}

	tags, err := c.taskTagUsecase.GetAvailableTags(ctx, robotTypes, request.Params.CategoryTypeId)
	if err != nil {
		return nil, err
	}

	resp := make([]openapi.TaskTag, 0, len(tags))
	for _, t := range tags {
		resp = append(resp, openapi.TaskTag{
			Id:               t.ID,
			Name:             t.Name,
			CategoryTypeId:   t.CategoryTypeID,
			CategoryTypeName: t.CategoryTypeName,
		})
	}

	return openapi.GetTaskAvailableTags200JSONResponse(resp), nil
}

func (c *controller) ExportTasks(ctx context.Context, request openapi.ExportTasksRequestObject) (openapi.ExportTasksResponseObject, error) {
	filter := usecase.TaskListFilter{}
	if request.Params.Status != nil {
		filter.Statuses = taskStatuses(*request.Params.Status)
	}
	if request.Params.Priority != nil {
		filter.Priorities = taskPriorities(*request.Params.Priority)
	}
	if request.Params.Difficulty != nil {
		filter.Difficulties = taskDifficulties(*request.Params.Difficulty)
	}
	filter.RobotType = request.Params.RobotType

	csvBytes, err := c.taskExportUsecase.Export(ctx, filter)
	if err != nil {
		return nil, err
	}

	return openapi.ExportTasks200TextcsvResponse{
		Body: bytes.NewReader(csvBytes),
		Headers: openapi.ExportTasks200ResponseHeaders{
			ContentDisposition: `attachment; filename="tasks_export.csv"`,
		},
	}, nil
}
