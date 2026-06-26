package controller

import (
	"testing"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
)

func TestLocationResponse(t *testing.T) {
	loc := model.Location{
		IDNatural: "loc-1",
		SiteID:    "site-1",
		SiteName:  "Main Site",
		Name:      "Dock",
	}

	got := locationResponse(loc)

	if got.Id != loc.IDNatural || got.SiteId != loc.SiteID || got.SiteName != loc.SiteName || got.Name != loc.Name {
		t.Fatalf("locationResponse() = %+v, want values from %+v", got, loc)
	}
}

func TestSiteResponse(t *testing.T) {
	createdAt := time.Date(2026, 6, 26, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	site := model.Site{
		IDNatural:      "site-1",
		OrganizationID: "org-1",
		Name:           "Main Site",
		CreatedAt:      createdAt,
		UpdatedAt:      &updatedAt,
	}

	got := siteResponse(site)

	if got.Id != site.IDNatural || got.OrganizationId != site.OrganizationID || got.Name != site.Name {
		t.Fatalf("siteResponse() = %+v, want values from %+v", got, site)
	}
	if got.CreatedAt == nil || !got.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, createdAt)
	}
	if got.UpdatedAt == nil || !got.UpdatedAt.Equal(updatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, updatedAt)
	}
}

func TestOrganizationResponse(t *testing.T) {
	createdAt := time.Date(2026, 6, 26, 10, 0, 0, 0, time.UTC)
	description := "robotics"
	org := model.Organization{
		IDNatural:   "org-1",
		Name:        "Airoa",
		Description: &description,
		CreatedAt:   createdAt,
	}

	got := organizationResponse(org)

	if got.OrganizationId != org.IDNatural || got.DisplayName != org.Name {
		t.Fatalf("organizationResponse() = %+v, want values from %+v", got, org)
	}
	if got.Description == nil || *got.Description != description {
		t.Errorf("Description = %v, want %q", got.Description, description)
	}
	if got.CreatedAt == nil || !got.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, createdAt)
	}
}

func TestUserResponse(t *testing.T) {
	createdAt := time.Date(2026, 6, 26, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	user := model.User{
		IDNatural:        "user-1",
		OrganizationID:   "org-1",
		OrganizationName: "Airoa",
		Name:             "Operator",
		Email:            "operator@example.com",
		Role:             model.UserRoleOperator,
		CreatedAt:        createdAt,
		UpdatedAt:        &updatedAt,
		Locations: []model.LocationSummary{
			{LocationID: "loc-1", Name: "Dock"},
		},
		Sites: []model.SiteSummary{
			{SiteID: "site-1", Name: "Tokyo"},
		},
	}

	got := userResponse(user)

	if got.UserId != user.IDNatural || got.Email != user.Email || got.DisplayName != user.Name {
		t.Fatalf("userResponse() = %+v, want values from %+v", got, user)
	}
	if got.OrganizationId != user.OrganizationID || got.OrganizationName != user.OrganizationName {
		t.Errorf("organization fields = (%q, %q), want (%q, %q)", got.OrganizationId, got.OrganizationName, user.OrganizationID, user.OrganizationName)
	}
	if got.Role == nil || *got.Role != openapi.Operator {
		t.Errorf("Role = %v, want operator", got.Role)
	}
	if !got.CreatedAt.Equal(createdAt) || got.UpdatedAt == nil || !got.UpdatedAt.Equal(updatedAt) {
		t.Errorf("timestamps = (%v, %v), want (%v, %v)", got.CreatedAt, got.UpdatedAt, createdAt, updatedAt)
	}
	if len(got.Locations) != 1 || got.Locations[0].LocationId != "loc-1" {
		t.Errorf("Locations = %+v, want loc-1", got.Locations)
	}
	if len(got.Sites) != 1 || got.Sites[0].SiteId != "site-1" {
		t.Errorf("Sites = %+v, want site-1", got.Sites)
	}
}

func TestRobotResponse(t *testing.T) {
	robotType := "arm"
	activeEpisodeID := "episode-1"
	activeUserID := "user-1"
	robot := model.Robot{
		IDNatural:        "robot-1",
		OrganizationID:   "org-1",
		OrganizationName: "Airoa",
		SiteID:           "site-1",
		SiteName:         "Tokyo",
		LocationID:       "loc-1",
		LocationName:     "Dock",
		Name:             "Yubi",
		RobotType:        &robotType,
		Status:           model.RobotStatusReady,
		ActiveEpisodeID:  &activeEpisodeID,
		ActiveUserID:     &activeUserID,
	}

	got := robotResponse(robot)

	if got.Id != robot.IDNatural || got.Name != robot.Name {
		t.Fatalf("robotResponse() = %+v, want values from %+v", got, robot)
	}
	if got.OrganizationId == nil || *got.OrganizationId != robot.OrganizationID {
		t.Errorf("OrganizationId = %v, want %q", got.OrganizationId, robot.OrganizationID)
	}
	if got.SiteId == nil || *got.SiteId != robot.SiteID {
		t.Errorf("SiteId = %v, want %q", got.SiteId, robot.SiteID)
	}
	if got.LocationId == nil || *got.LocationId != robot.LocationID {
		t.Errorf("LocationId = %v, want %q", got.LocationId, robot.LocationID)
	}
	if got.Status == nil || *got.Status != openapi.RobotStatusReady {
		t.Errorf("Status = %v, want ready", got.Status)
	}
	if got.ActiveEpisodeId == nil || *got.ActiveEpisodeId != activeEpisodeID {
		t.Errorf("ActiveEpisodeId = %v, want %q", got.ActiveEpisodeId, activeEpisodeID)
	}
	if got.ActiveUserId == nil || *got.ActiveUserId != activeUserID {
		t.Errorf("ActiveUserId = %v, want %q", got.ActiveUserId, activeUserID)
	}
}

func TestSubTaskResponse(t *testing.T) {
	description := "Pick an item"
	targetDurationSeconds := 120
	subtask := model.SubTask{
		IDNatural:             "subtask-1",
		Name:                  "Pick",
		Description:           &description,
		TargetDurationSeconds: &targetDurationSeconds,
	}

	got := subTaskResponse(subtask)

	if got.Id != subtask.IDNatural || got.Name != subtask.Name {
		t.Fatalf("subTaskResponse() = %+v, want values from %+v", got, subtask)
	}
	if got.Description == nil || *got.Description != description {
		t.Errorf("Description = %v, want %q", got.Description, description)
	}
	if got.TargetDurationSeconds == nil || *got.TargetDurationSeconds != targetDurationSeconds {
		t.Errorf("TargetDurationSeconds = %v, want %d", got.TargetDurationSeconds, targetDurationSeconds)
	}
}

func TestAPIKeyResponse(t *testing.T) {
	createdAt := time.Date(2026, 6, 26, 10, 0, 0, 0, time.UTC)
	expiresAt := createdAt.Add(24 * time.Hour)
	robotID := "robot-1"
	robotName := "Yubi"
	key := model.APIKey{
		IDNatural:      "key-1",
		OrganizationID: "org-1",
		UserID:         "user-1",
		UserName:       "Operator",
		RobotID:        &robotID,
		RobotName:      &robotName,
		Name:           "Robot key",
		KeyHint:        "abcd...wxyz",
		ExpiresAt:      &expiresAt,
		CreatedAt:      createdAt,
	}

	got := apiKeyResponse(key)

	if got.Id != key.IDNatural || got.Name != key.Name || got.KeyHint != key.KeyHint {
		t.Fatalf("apiKeyResponse() = %+v, want values from %+v", got, key)
	}
	if got.UserId != key.UserID || got.UserName != key.UserName || got.OrganizationId != key.OrganizationID {
		t.Errorf("owner fields = (%q, %q, %q), want (%q, %q, %q)", got.UserId, got.UserName, got.OrganizationId, key.UserID, key.UserName, key.OrganizationID)
	}
	if got.RobotId == nil || *got.RobotId != robotID || got.RobotName == nil || *got.RobotName != robotName {
		t.Errorf("robot fields = (%v, %v), want (%q, %q)", got.RobotId, got.RobotName, robotID, robotName)
	}
	if !got.UpdatedAt.Equal(createdAt) {
		t.Errorf("UpdatedAt = %v, want CreatedAt fallback %v", got.UpdatedAt, createdAt)
	}
}

func TestEpisodeGradeResponse(t *testing.T) {
	gradedAt := time.Date(2026, 6, 26, 10, 0, 0, 0, time.UTC)
	createdAt := gradedAt.Add(-time.Hour)
	updatedAt := gradedAt.Add(time.Hour)
	comment := "good"
	grade := model.EpisodeGrade{
		EpisodeID: "episode-1",
		UserID:    "user-1",
		Grade:     0.8,
		Comment:   &comment,
		GradedAt:  gradedAt,
		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
	}

	got := episodeGradeResponse(grade, "Operator")

	if got.EpisodeId != grade.EpisodeID || got.UserId != grade.UserID || got.UserName != "Operator" {
		t.Fatalf("episodeGradeResponse() = %+v, want values from %+v", got, grade)
	}
	if got.Comment == nil || *got.Comment != comment {
		t.Errorf("Comment = %v, want %q", got.Comment, comment)
	}
	if got.UpdatedAt == nil || !got.UpdatedAt.Equal(updatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, updatedAt)
	}
}

func TestTaskResponse(t *testing.T) {
	description := "Move boxes"
	priority := model.TaskPriorityHigh
	difficulty := model.TaskDifficultyA
	status := model.TaskStatusDoing
	robotType := "arm"
	displayName := "Pilot"
	task := model.Task{
		IDNatural:          "task-1",
		Name:               "Move",
		Description:        &description,
		ManualURL:          "https://example.com/manual",
		Priority:           &priority,
		Difficulty:         &difficulty,
		Status:             &status,
		RobotType:          &robotType,
		Version:            "v2",
		VersionDisplayName: &displayName,
		Tags: model.TaskTags{
			{ID: "tag-1", Name: "Picking", CategoryTypeID: "cat-1", CategoryTypeName: "Work"},
		},
	}

	got := taskResponse(task)

	if got.Id != task.IDNatural || got.Name != task.Name || got.ManualUrl != task.ManualURL {
		t.Fatalf("taskResponse() = %+v, want values from %+v", got, task)
	}
	if got.Version == nil || *got.Version != task.Version {
		t.Errorf("Version = %v, want %q", got.Version, task.Version)
	}
	if got.VersionDisplayName == nil || *got.VersionDisplayName != displayName {
		t.Errorf("VersionDisplayName = %v, want %q", got.VersionDisplayName, displayName)
	}
	if got.Tags == nil || len(*got.Tags) != 1 || (*got.Tags)[0].Id != "tag-1" {
		t.Errorf("Tags = %+v, want tag-1", got.Tags)
	}
}

func TestTaskVersionResponse(t *testing.T) {
	createdAt := time.Date(2026, 6, 26, 10, 0, 0, 0, time.UTC)
	displayName := "Pilot"
	actualDurationSeconds := int64(90)
	actualEpisodeCount := 3
	taskVersion := model.TaskVersion{
		IDNatural:             "version-1",
		TaskID:                "task-1",
		Version:               "v2",
		DisplayName:           &displayName,
		ApprovalStatus:        model.ApprovalStatusApproved,
		IsCurrent:             true,
		CreatedAt:             createdAt,
		ActualDurationSeconds: &actualDurationSeconds,
		ActualEpisodeCount:    &actualEpisodeCount,
		Parameters: []model.TaskVersionParameter{
			{Key: "color", Values: []string{"red"}},
		},
	}

	got := taskVersionResponse(taskVersion)

	if got.Id != taskVersion.IDNatural || got.TaskId != taskVersion.TaskID || got.Version != taskVersion.Version {
		t.Fatalf("taskVersionResponse() = %+v, want values from %+v", got, taskVersion)
	}
	if got.ActualDurationSeconds == nil || *got.ActualDurationSeconds != int(actualDurationSeconds) {
		t.Errorf("ActualDurationSeconds = %v, want %d", got.ActualDurationSeconds, actualDurationSeconds)
	}
	if got.Parameters == nil || len(*got.Parameters) != 1 || (*got.Parameters)[0].Key != "color" {
		t.Errorf("Parameters = %+v, want color", got.Parameters)
	}
}
