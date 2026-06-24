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

type TaskStatus int

const (
	TaskStatusPlanning  TaskStatus = 0
	TaskStatusDoing     TaskStatus = 1
	TaskStatusCompleted TaskStatus = 2
	TaskStatusCanceled  TaskStatus = 3
)

type TaskPriority int

const (
	TaskPriorityLow    TaskPriority = 0
	TaskPriorityNormal TaskPriority = 1
	TaskPriorityHigh   TaskPriority = 2
	TaskPriorityUrgent TaskPriority = 3
)

type TaskDifficulty int

const (
	TaskDifficultyS TaskDifficulty = 0
	TaskDifficultyA TaskDifficulty = 1
	TaskDifficultyB TaskDifficulty = 2
	TaskDifficultyC TaskDifficulty = 3
)
