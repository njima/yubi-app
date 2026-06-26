package controller

import (
	"bytes"
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

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
		taskList = append(taskList, taskResponse(*t))
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
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	input := usecase.TaskCreateInput{
		OrganizationID: body.OrganizationId,
		Name:           body.Name,
		Description:    body.Description,
		ManualURL:      body.ManualUrl,
		Priority:       model.TaskPriority(body.Priority),
		Difficulty:     model.TaskDifficulty(body.Difficulty),
		Status:         model.TaskStatus(body.Status),
		Deadline:       body.Deadline,
		RobotType:      body.RobotType,
	}
	if body.TagIds != nil {
		input.TagIDs = *body.TagIds
	}

	tk, err := c.taskUsecase.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return openapi.CreateTask201JSONResponse(taskResponse(tk)), nil
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

	return openapi.GetTaskById200JSONResponse(taskResponse(tk)), nil
}

func (c *controller) UpdateTaskById(ctx context.Context, request openapi.UpdateTaskByIdRequestObject) (openapi.UpdateTaskByIdResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	input := usecase.TaskUpdateInput{
		ID: request.TaskId,
	}
	if body.Name != nil {
		input.Name = body.Name
	}
	if body.Description != nil {
		input.Description = body.Description
	}
	if body.ManualUrl != nil {
		input.ManualURL = body.ManualUrl
	}
	if body.Priority != nil {
		input.Priority = taskPriorityPtr(*body.Priority)
	}
	if body.Difficulty != nil {
		input.Difficulty = taskDifficultyPtr(*body.Difficulty)
	}
	if body.Status != nil {
		input.Status = taskStatusPtr(*body.Status)
	}
	if body.Deadline != nil {
		input.Deadline = body.Deadline
	}
	if body.RobotType != nil {
		input.RobotType = body.RobotType
	}
	if body.TagIds != nil {
		input.TagIDs = body.TagIds
	}

	tk, err := c.taskUsecase.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateTaskById200JSONResponse(taskResponse(tk)), nil
}

func (c *controller) CreateTaskVersion(ctx context.Context, request openapi.CreateTaskVersionRequestObject) (openapi.CreateTaskVersionResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	input := usecase.TaskVersionCreateInput{
		TaskID:                          request.TaskId,
		Version:                         body.Version,
		DisplayName:                     body.DisplayName,
		BaseTaskVersionID:               body.BaseTaskVersionId,
		TargetDurationSeconds:           body.TargetDurationSeconds,
		TargetEpisodeCount:              body.TargetEpisodeCount,
		TargetDurationPerEpisodeSeconds: body.TargetDurationPerEpisodeSeconds,
	}

	tv, err := c.taskVersionUsecase.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return openapi.CreateTaskVersion201JSONResponse(taskVersionResponse(tv)), nil
}

func (c *controller) ListTaskVersions(ctx context.Context, request openapi.ListTaskVersionsRequestObject) (openapi.ListTaskVersionsResponseObject, error) {
	versions, err := c.taskVersionUsecase.ListByTaskID(ctx, request.TaskId)
	if err != nil {
		return nil, err
	}

	resp := make([]openapi.TaskVersion, 0, len(versions))
	for _, v := range versions {
		resp = append(resp, taskVersionResponse(*v))
	}
	return openapi.ListTaskVersions200JSONResponse(resp), nil
}

func (c *controller) UpdateTaskVersion(ctx context.Context, request openapi.UpdateTaskVersionRequestObject) (openapi.UpdateTaskVersionResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	tv, err := c.taskVersionUsecase.Update(ctx, request.TaskId, usecase.TaskVersionUpdateInput{
		ID:                              request.VersionId,
		DisplayName:                     body.DisplayName,
		TargetDurationSeconds:           body.TargetDurationSeconds,
		TargetEpisodeCount:              body.TargetEpisodeCount,
		TargetDurationPerEpisodeSeconds: body.TargetDurationPerEpisodeSeconds,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateTaskVersion200JSONResponse(taskVersionResponse(tv)), nil
}

func (c *controller) ApproveTaskVersion(ctx context.Context, request openapi.ApproveTaskVersionRequestObject) (openapi.ApproveTaskVersionResponseObject, error) {
	tv, err := c.taskVersionUsecase.Approve(ctx, request.TaskId, request.VersionId)
	if err != nil {
		return nil, err
	}

	return openapi.ApproveTaskVersion200JSONResponse(taskVersionResponse(tv)), nil
}

func (c *controller) UpdateTaskVersionParameters(ctx context.Context, request openapi.UpdateTaskVersionParametersRequestObject) (openapi.UpdateTaskVersionParametersResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	tv, err := c.taskVersionUsecase.UpdateParameters(ctx, usecase.TaskVersionUpdateParametersInput{
		TaskID:     request.TaskId,
		VersionID:  request.VersionId,
		Parameters: openAPIToModelParameters(body.Parameters),
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateTaskVersionParameters200JSONResponse(taskVersionResponse(tv)), nil
}

func openAPIToModelParameters(params []openapi.TaskVersionParameter) []model.TaskVersionParameter {
	result := make([]model.TaskVersionParameter, len(params))
	for i, p := range params {
		result[i] = model.TaskVersionParameter{Key: p.Key, Values: p.Values}
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
