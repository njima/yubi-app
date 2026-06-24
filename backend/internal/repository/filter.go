package repository

import "github.com/airoa-org/yubi-app/backend/internal/domain/model"

type EpisodeStatus = model.EpisodeStatus

const (
	EpisodeStatusReady     = model.EpisodeStatusReady
	EpisodeStatusRecording = model.EpisodeStatusRecording
	EpisodeStatusCancel    = model.EpisodeStatusCancel
	EpisodeStatusCompleted = model.EpisodeStatusCompleted
)

type RobotFilterStatus int

const (
	RobotFilterStatusOnline      RobotFilterStatus = 0
	RobotFilterStatusBusy        RobotFilterStatus = 1
	RobotFilterStatusOffline     RobotFilterStatus = 2
	RobotFilterStatusFaulted     RobotFilterStatus = 3
	RobotFilterStatusMaintenance RobotFilterStatus = 4
	RobotFilterStatusReady       RobotFilterStatus = 5
)

type TaskStatus = model.TaskStatus

const (
	TaskStatusPlanning  = model.TaskStatusPlanning
	TaskStatusDoing     = model.TaskStatusDoing
	TaskStatusCompleted = model.TaskStatusCompleted
	TaskStatusCanceled  = model.TaskStatusCanceled
)

type TaskPriority = model.TaskPriority

const (
	TaskPriorityLow    = model.TaskPriorityLow
	TaskPriorityNormal = model.TaskPriorityNormal
	TaskPriorityHigh   = model.TaskPriorityHigh
	TaskPriorityUrgent = model.TaskPriorityUrgent
)

type TaskDifficulty = model.TaskDifficulty

const (
	TaskDifficultyS = model.TaskDifficultyS
	TaskDifficultyA = model.TaskDifficultyA
	TaskDifficultyB = model.TaskDifficultyB
	TaskDifficultyC = model.TaskDifficultyC
)
