package model

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SubTaskCollectionStatus int

const (
	SubTaskCollectionStatusReady      SubTaskCollectionStatus = 0
	SubTaskCollectionStatusInProgress SubTaskCollectionStatus = 1
	SubTaskCollectionStatusCompleted  SubTaskCollectionStatus = 2
	SubTaskCollectionStatusSkipped    SubTaskCollectionStatus = 3
	SubTaskCollectionStatusCancelled  SubTaskCollectionStatus = 4
)

type TaskResult int

const (
	TaskResultUndetermined TaskResult = 0
	TaskResultSuccess      TaskResult = 1
	TaskResultFailed       TaskResult = 2
)

func (s SubTaskCollectionStatus) IsTerminal() bool {
	switch s {
	case SubTaskCollectionStatusCompleted, SubTaskCollectionStatusSkipped, SubTaskCollectionStatusCancelled:
		return true
	default:
		return false
	}
}

func (s SubTaskCollectionStatus) IsWorkflowResolved() bool {
	return s == SubTaskCollectionStatusCompleted || s == SubTaskCollectionStatusSkipped
}

func (s SubTaskCollectionStatus) IsSuccessfulCompletion() bool {
	return s == SubTaskCollectionStatusCompleted
}

type EpisodeSubTask struct {
	ID               int64
	IDNatural        string
	OrganizationID   string
	EpisodeID        string
	SubTaskID        string
	CollectionStatus SubTaskCollectionStatus
	CreatedAt        time.Time
	UpdatedAt        *time.Time
}

type EpisodeSubTasks []*EpisodeSubTask

func InitEpisodeSubTask(organizationID, episodeID, subtaskID string) (EpisodeSubTask, error) {
	idNatural, err := InitID()
	if err != nil {
		return EpisodeSubTask{}, err
	}

	est := EpisodeSubTask{
		IDNatural:        idNatural,
		OrganizationID:   organizationID,
		EpisodeID:        episodeID,
		SubTaskID:        subtaskID,
		CollectionStatus: SubTaskCollectionStatusReady,
		CreatedAt:        time.Now(),
	}

	if err := est.validate(); err != nil {
		return EpisodeSubTask{}, err
	}

	return est, nil
}

func NewEpisodeSubTask(
	id int64,
	idNatural,
	organizationID,
	episodeID,
	subtaskID string,
	collectionStatus SubTaskCollectionStatus,
	createdAt time.Time,
	updatedAt *time.Time,
) EpisodeSubTask {
	return EpisodeSubTask{
		ID:               id,
		IDNatural:        idNatural,
		OrganizationID:   organizationID,
		EpisodeID:        episodeID,
		SubTaskID:        subtaskID,
		CollectionStatus: collectionStatus,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

func (est EpisodeSubTask) validate() error {
	if err := (validation.Errors{
		"id_natural":      validation.Validate(est.IDNatural, validation.Required.Error("id_natural is required")),
		"organization_id": validation.Validate(est.OrganizationID, validation.Required.Error("organization_id is required")),
		"episode_id":      validation.Validate(est.EpisodeID, validation.Required.Error("episode_id is required")),
		"subtask_id":      validation.Validate(est.SubTaskID, validation.Required.Error("subtask_id is required")),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "episode_sub_task validation failed: %v", err))
	}
	return nil
}

func (est EpisodeSubTask) IsTerminal() bool {
	return est.CollectionStatus.IsTerminal()
}

func (est EpisodeSubTask) IsWorkflowResolved() bool {
	return est.CollectionStatus.IsWorkflowResolved()
}

func (est EpisodeSubTask) IsSuccessfulCompletion() bool {
	return est.CollectionStatus.IsSuccessfulCompletion()
}

func (est *EpisodeSubTask) CanStartProgress() error {
	if est.CollectionStatus != SubTaskCollectionStatusReady {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "subtask status must be Ready to start progress, current: %v", est.CollectionStatus),
		)
	}
	return nil
}

func (est *EpisodeSubTask) StartProgress() error {
	if est.CollectionStatus == SubTaskCollectionStatusInProgress {
		return nil // Already in progress
	}
	if err := est.CanStartProgress(); err != nil {
		return err
	}
	est.CollectionStatus = SubTaskCollectionStatusInProgress
	return nil
}

func (est *EpisodeSubTask) Complete() error {
	if est.CollectionStatus == SubTaskCollectionStatusCompleted {
		return nil // Already completed
	}
	if est.CollectionStatus == SubTaskCollectionStatusCancelled ||
		est.CollectionStatus == SubTaskCollectionStatusSkipped {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "cannot complete subtask with status: %v", est.CollectionStatus),
		)
	}
	est.CollectionStatus = SubTaskCollectionStatusCompleted
	return nil
}

func (est *EpisodeSubTask) Skip() error {
	if est.CollectionStatus == SubTaskCollectionStatusSkipped {
		return nil // Already skipped
	}
	if est.CollectionStatus == SubTaskCollectionStatusCancelled ||
		est.CollectionStatus == SubTaskCollectionStatusCompleted {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "cannot skip subtask with status: %v", est.CollectionStatus),
		)
	}
	est.CollectionStatus = SubTaskCollectionStatusSkipped
	return nil
}

func (est *EpisodeSubTask) Cancel() error {
	if est.CollectionStatus == SubTaskCollectionStatusCancelled {
		return nil // Already cancelled
	}
	if est.CollectionStatus == SubTaskCollectionStatusCompleted {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "cannot cancel completed subtask"),
		)
	}
	est.CollectionStatus = SubTaskCollectionStatusCancelled
	return nil
}
