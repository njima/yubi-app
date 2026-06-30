package controller

import (
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func locationResponse(loc model.Location) openapi.Location {
	return openapi.Location{
		Id:       loc.IDNatural,
		Name:     loc.Name,
		SiteId:   loc.SiteID,
		SiteName: loc.SiteName,
	}
}

func apiKeyResponse(key model.APIKey) openapi.ApiKeyResponse {
	updatedAt := key.CreatedAt
	if key.UpdatedAt != nil {
		updatedAt = *key.UpdatedAt
	}
	return openapi.ApiKeyResponse{
		Id:             key.IDNatural,
		Name:           key.Name,
		UserId:         key.UserID,
		UserName:       key.UserName,
		OrganizationId: key.OrganizationID,
		RobotId:        key.RobotID,
		RobotName:      key.RobotName,
		KeyHint:        key.KeyHint,
		ExpiresAt:      key.ExpiresAt,
		LastUsedAt:     key.LastUsedAt,
		RevokedAt:      key.RevokedAt,
		CreatedAt:      key.CreatedAt,
		UpdatedAt:      updatedAt,
	}
}

func episodeGradeResponse(grade model.EpisodeGrade, userName string) openapi.EpisodeGrade {
	return openapi.EpisodeGrade{
		EpisodeId: grade.EpisodeID,
		UserId:    grade.UserID,
		UserName:  userName,
		Grade:     grade.Grade,
		Comment:   grade.Comment,
		GradedAt:  grade.GradedAt,
		CreatedAt: grade.CreatedAt,
		UpdatedAt: grade.UpdatedAt,
	}
}

func episodeResponse(ep model.Episode) openapi.Episode {
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

func siteResponse(site model.Site) openapi.Site {
	return openapi.Site{
		Id:             site.IDNatural,
		Name:           site.Name,
		OrganizationId: site.OrganizationID,
		CreatedAt:      &site.CreatedAt,
		UpdatedAt:      site.UpdatedAt,
	}
}

func subTaskResponse(subtask model.SubTask) openapi.SubTask {
	return openapi.SubTask{
		Id:                    subtask.IDNatural,
		Name:                  subtask.Name,
		Description:           subtask.Description,
		TargetDurationSeconds: subtask.TargetDurationSeconds,
	}
}

func subTaskResponses(subtasks model.SubTasks) []openapi.SubTask {
	result := make([]openapi.SubTask, 0, len(subtasks))
	for _, subtask := range subtasks {
		result = append(result, subTaskResponse(*subtask))
	}
	return result
}

func taskTagResponse(tag model.TaskTag) openapi.TaskTag {
	return openapi.TaskTag{
		Id:               tag.ID,
		Name:             tag.Name,
		CategoryTypeId:   tag.CategoryTypeID,
		CategoryTypeName: tag.CategoryTypeName,
	}
}

func taskTagResponses(tags model.TaskTags) *[]openapi.TaskTag {
	if len(tags) == 0 {
		return nil
	}
	result := make([]openapi.TaskTag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, taskTagResponse(*tag))
	}
	return &result
}

func taskResponse(task model.Task) openapi.Task {
	var priority openapi.TaskPriority
	if task.Priority != nil {
		priority = openapi.TaskPriority(*task.Priority)
	}
	var difficulty openapi.TaskDifficulty
	if task.Difficulty != nil {
		difficulty = openapi.TaskDifficulty(*task.Difficulty)
	}
	var status openapi.TaskStatus
	if task.Status != nil {
		status = openapi.TaskStatus(*task.Status)
	}
	resp := openapi.Task{
		Id:                    task.IDNatural,
		Name:                  task.Name,
		Description:           task.Description,
		ManualUrl:             task.ManualURL,
		Priority:              priority,
		Difficulty:            difficulty,
		Status:                status,
		Deadline:              task.Deadline,
		RobotType:             task.RobotType,
		TargetDurationSeconds: task.TargetDurationSeconds,
		TargetEpisodeCount:    task.TargetEpisodeCount,
		ActualEpisodeCount:    task.ActualEpisodeCount,
		Tags:                  taskTagResponses(task.Tags),
	}
	if task.Version != "" {
		resp.Version = &task.Version
		tv := model.TaskVersion{Version: task.Version, DisplayName: task.VersionDisplayName}
		resolved := tv.DisplayLabel(task.Name)
		resp.VersionDisplayName = &resolved
	}
	return resp
}

func taskVersionResponse(taskVersion model.TaskVersion) openapi.TaskVersion {
	resp := openapi.TaskVersion{
		Id:                              taskVersion.IDNatural,
		TaskId:                          taskVersion.TaskID,
		Version:                         taskVersion.Version,
		DisplayName:                     taskVersion.DisplayName,
		IsCurrent:                       taskVersion.IsCurrent,
		ApprovalStatus:                  openAPIApprovalStatus(taskVersion.ApprovalStatus),
		CreatedAt:                       taskVersion.CreatedAt,
		TargetDurationSeconds:           taskVersion.TargetDurationSeconds,
		TargetEpisodeCount:              taskVersion.TargetEpisodeCount,
		TargetDurationPerEpisodeSeconds: taskVersion.TargetDurationPerEpisodeSeconds,
	}
	if taskVersion.ActualDurationSeconds != nil {
		v := int(*taskVersion.ActualDurationSeconds)
		resp.ActualDurationSeconds = &v
	}
	if taskVersion.ActualEpisodeCount != nil {
		resp.ActualEpisodeCount = taskVersion.ActualEpisodeCount
	}
	if len(taskVersion.Parameters) > 0 {
		params := taskVersionParameterResponses(taskVersion.Parameters)
		resp.Parameters = &params
	}
	return resp
}

func taskVersionParameterResponses(params []model.TaskVersionParameter) []openapi.TaskVersionParameter {
	result := make([]openapi.TaskVersionParameter, len(params))
	for i, param := range params {
		result[i] = openapi.TaskVersionParameter{Key: param.Key, Values: param.Values}
	}
	return result
}

func organizationResponse(org model.Organization) openapi.OrganizationResponse {
	return openapi.OrganizationResponse{
		OrganizationId: org.IDNatural,
		DisplayName:    org.Name,
		Description:    org.Description,
		CreatedAt:      &org.CreatedAt,
		UpdatedAt:      org.UpdatedAt,
	}
}

func robotResponse(robot model.Robot) openapi.Robot {
	status, leaderStatus := openAPIRobotStatus(robot.Status), openAPILeaderStatus(robot.LeaderStatus)
	return openapi.Robot{
		Id:                         robot.IDNatural,
		OrganizationId:             &robot.OrganizationID,
		OrganizationName:           &robot.OrganizationName,
		SiteId:                     &robot.SiteID,
		SiteName:                   &robot.SiteName,
		LocationId:                 &robot.LocationID,
		LocationName:               &robot.LocationName,
		Name:                       robot.Name,
		RobotType:                  robot.RobotType,
		Status:                     &status,
		LeaderStatus:               leaderStatus,
		ConsecutiveFaultDays:       robot.ConsecutiveFaultDays(),
		LeaderConsecutiveFaultDays: robot.LeaderConsecutiveFaultDays(),
		LeaderFaultStartedAt:       robot.LeaderFaultStartedAt,
		LastHeartbeatAt:            robot.LastHeartbeatAt,
		OfflineReason:              robot.OfflineReason,
		RobotConfig:                mapPtrFromRawMessagePtr(robot.RobotConfig),
		ActiveEpisodeId:            robot.ActiveEpisodeID,
		ActiveUserId:               robot.ActiveUserID,
	}
}

func robotResponses(robots model.Robots) []openapi.Robot {
	result := make([]openapi.Robot, 0, len(robots))
	for _, robot := range robots {
		result = append(result, robotResponse(*robot))
	}
	return result
}

func robotOperatorResponse(operator model.RobotOperator) openapi.RobotOperator {
	return openapi.RobotOperator{
		UserId:           operator.UserID,
		DisplayName:      operator.DisplayName,
		OrganizationName: operator.OrganizationName,
	}
}

func userResponse(user model.User) openapi.UserResponse {
	return openapi.UserResponse{
		UserId:      user.IDNatural,
		Email:       user.Email,
		DisplayName: user.Name,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Locations:   locationSummaries(user.Locations),
		Sites:       siteSummaries(user.Sites),
	}
}

func userResponseWithWorkspace(user model.User, org model.Organization, membership model.OrganizationMembership) openapi.UserResponse {
	resp := userResponse(user)
	resp.Role = openAPIUserRolePtr(membership.Role)
	resp.OrganizationId = org.IDNatural
	resp.OrganizationName = org.Name
	return resp
}

func locationSummaries(locs []model.LocationSummary) []openapi.LocationSummary {
	result := make([]openapi.LocationSummary, 0, len(locs))
	for _, loc := range locs {
		result = append(result, openapi.LocationSummary{
			LocationId: loc.LocationID,
			Name:       loc.Name,
		})
	}
	return result
}

func siteSummaries(sites []model.SiteSummary) []openapi.SiteSummary {
	result := make([]openapi.SiteSummary, 0, len(sites))
	for _, site := range sites {
		result = append(result, openapi.SiteSummary{
			SiteId: site.SiteID,
			Name:   site.Name,
		})
	}
	return result
}

func meResponse(session usecase.AuthenticatedUserSession) openapi.MeResponse {
	return openapi.MeResponse{
		UserId:                 session.User.IDNatural,
		Email:                  session.User.Email,
		DisplayName:            session.User.Name,
		AvatarUrl:              session.User.AvatarURL,
		ActiveOrganizationId:   session.ActiveOrganization.IDNatural,
		ActiveOrganizationName: session.ActiveOrganization.Name,
		ActiveRole:             openAPIUserRole(session.ActiveMembership.Role),
	}
}
