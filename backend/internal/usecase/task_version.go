package usecase

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type TaskVersionCreateInput struct {
	TaskID                          string
	Version                         string
	DisplayName                     *string
	BaseTaskVersionID               string
	TargetDurationSeconds           *int
	TargetEpisodeCount              *int
	TargetDurationPerEpisodeSeconds *int
}

type TaskVersionUpdateInput struct {
	ID                              string
	DisplayName                     *string
	TargetDurationSeconds           *int
	TargetEpisodeCount              *int
	TargetDurationPerEpisodeSeconds *int
}

type TaskVersionUpdateParametersInput struct {
	TaskID     string
	VersionID  string
	Parameters []model.TaskVersionParameter
}

type TaskVersionUsecase interface {
	ListByTaskID(ctx context.Context, taskID string) (model.TaskVersions, error)
	ListByIDs(ctx context.Context, ids []string) (model.TaskVersions, error)
	GetByID(ctx context.Context, id string) (model.TaskVersion, error)
	Create(ctx context.Context, input TaskVersionCreateInput) (model.TaskVersion, error)
	Update(ctx context.Context, taskID string, input TaskVersionUpdateInput) (model.TaskVersion, error)
	Approve(ctx context.Context, taskID, versionID string) (model.TaskVersion, error)
	UpdateParameters(ctx context.Context, input TaskVersionUpdateParametersInput) (model.TaskVersion, error)
}

type taskVersion struct {
	repo        repository.TaskVersion
	taskRepo    repository.Task
	subtaskRepo repository.SubTask
	episodeRepo repository.Episode
	db          *bun.DB
}

func NewTaskVersion(repo repository.TaskVersion, taskRepo repository.Task, subtaskRepo repository.SubTask, episodeRepo repository.Episode, db *bun.DB) TaskVersionUsecase {
	return &taskVersion{repo: repo, taskRepo: taskRepo, subtaskRepo: subtaskRepo, episodeRepo: episodeRepo, db: db}
}

func (u *taskVersion) GetByID(ctx context.Context, id string) (model.TaskVersion, error) {
	return u.repo.GetByID(ctx, u.db, id)
}

func (u *taskVersion) ListByIDs(ctx context.Context, ids []string) (model.TaskVersions, error) {
	return u.repo.ListByIDs(ctx, u.db, ids)
}

func (u *taskVersion) ListByTaskID(ctx context.Context, taskID string) (model.TaskVersions, error) {
	exists, err := u.taskRepo.Exists(ctx, u.db, taskID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeTaskNotFound, "task not found: id=%s", taskID))
	}

	versions, err := u.repo.ListByTaskID(ctx, u.db, taskID)
	if err != nil {
		return nil, err
	}

	// Mark the latest approved version as current (versions are ordered by created_at DESC)
	for _, v := range versions {
		if v.IsApproved() {
			v.IsCurrent = true
			break
		}
	}

	return versions, nil
}

func (u *taskVersion) Create(ctx context.Context, input TaskVersionCreateInput) (model.TaskVersion, error) {
	exists, err := u.taskRepo.Exists(ctx, u.db, input.TaskID)
	if err != nil {
		return model.TaskVersion{}, err
	}
	if !exists {
		return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskNotFound, "task not found: id=%s", input.TaskID))
	}

	var result model.TaskVersion
	err = u.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Get all versions for this task (includes OrganizationID)
		versions, err := u.repo.ListByTaskID(ctx, tx, input.TaskID)
		if err != nil {
			return err
		}
		if len(versions) == 0 {
			return apperror.NewError(apperror.NewMessage(apperror.CodeTaskVersionNotFound, "no versions found for task: id=%s", input.TaskID))
		}

		orgID := versions[0].OrganizationID

		// Validate base version belongs to this task and copy its parameters
		var baseVersion *model.TaskVersion
		for _, v := range versions {
			if v.IDNatural == input.BaseTaskVersionID {
				baseVersion = v
				break
			}
		}
		if baseVersion == nil {
			return apperror.NewError(apperror.NewMessage(apperror.CodeTaskVersionNotFound, "base task version not found: id=%s", input.BaseTaskVersionID))
		}

		// Check version string uniqueness
		for _, v := range versions {
			if v.Version == input.Version {
				return apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "version already exists: %s", input.Version))
			}
		}

		// Create new version
		tv, err := model.InitTaskVersion(orgID, input.TaskID, input.Version, input.DisplayName, input.TargetDurationSeconds, input.TargetEpisodeCount, input.TargetDurationPerEpisodeSeconds, baseVersion.Parameters)
		if err != nil {
			return err
		}
		result, err = u.repo.Create(ctx, tx, tv)
		if err != nil {
			return err
		}

		// Copy subtasks from base version
		baseSubtasks, err := u.subtaskRepo.GetByTaskVersionID(ctx, tx, input.BaseTaskVersionID)
		if err != nil {
			return err
		}
		for _, st := range baseSubtasks {
			newSt, err := model.InitSubTask(orgID, result.IDNatural, st.Name, st.OrderIndex, st.Description, st.TargetDurationSeconds)
			if err != nil {
				return err
			}
			if _, err := u.subtaskRepo.Create(ctx, tx, newSt); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return model.TaskVersion{}, err
	}

	result.IsCurrent = false
	return result, nil
}

func (u *taskVersion) Update(ctx context.Context, taskID string, input TaskVersionUpdateInput) (model.TaskVersion, error) {
	tv, err := u.repo.GetByID(ctx, u.db, input.ID)
	if err != nil {
		return model.TaskVersion{}, err
	}
	if tv.TaskID != taskID {
		return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskVersionNotFound, "task version not found: id=%s for task=%s", input.ID, taskID))
	}
	if !tv.IsDraft() {
		return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "cannot edit approved task version: id=%s", input.ID))
	}

	if input.TargetDurationSeconds != nil {
		tv.TargetDurationSeconds = input.TargetDurationSeconds
	}
	if input.TargetEpisodeCount != nil {
		tv.TargetEpisodeCount = input.TargetEpisodeCount
	}
	if input.TargetDurationPerEpisodeSeconds != nil {
		tv.TargetDurationPerEpisodeSeconds = input.TargetDurationPerEpisodeSeconds
	}
	if input.DisplayName != nil {
		tv.DisplayName = input.DisplayName
	}

	return u.repo.Update(ctx, u.db, tv)
}

func (u *taskVersion) Approve(ctx context.Context, taskID, versionID string) (model.TaskVersion, error) {
	// Validate the version exists and belongs to this task
	tv, err := u.repo.GetByID(ctx, u.db, versionID)
	if err != nil {
		return model.TaskVersion{}, err
	}
	if tv.TaskID != taskID {
		return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskVersionNotFound, "task version not found: id=%s for task=%s", versionID, taskID))
	}
	if tv.IsApproved() {
		return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "task version already approved: id=%s", versionID))
	}

	var result model.TaskVersion
	err = u.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var txErr error
		result, txErr = u.repo.Approve(ctx, tx, versionID)
		if txErr != nil {
			return txErr
		}

		// Auto-update task status based on collection progress
		actual, err := u.episodeRepo.SumDurationByTaskID(ctx, tx, taskID)
		if err != nil {
			return err
		}
		target, err := u.repo.SumTargetByTaskID(ctx, tx, taskID)
		if err != nil {
			return err
		}
		tk, err := u.taskRepo.GetByID(ctx, tx, taskID)
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
		updateTask := model.Task{IDNatural: taskID}
		if err := updateTask.SetStatus(&newStatus); err != nil {
			return err
		}
		if _, err := u.taskRepo.Update(ctx, tx, updateTask); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return model.TaskVersion{}, err
	}

	result.IsCurrent = true
	return result, nil
}

func (u *taskVersion) UpdateParameters(ctx context.Context, input TaskVersionUpdateParametersInput) (model.TaskVersion, error) {
	tv, err := u.repo.GetByID(ctx, u.db, input.VersionID)
	if err != nil {
		return model.TaskVersion{}, err
	}
	if tv.TaskID != input.TaskID {
		return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskVersionNotFound, "task version not found: id=%s for task=%s", input.VersionID, input.TaskID))
	}
	if tv.IsApproved() {
		return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "cannot update parameters on approved version: id=%s", input.VersionID))
	}

	if err := model.ValidateParameterDefinitions(input.Parameters); err != nil {
		return model.TaskVersion{}, err
	}

	result, err := u.repo.UpdateParameters(ctx, u.db, input.VersionID, input.Parameters)
	if err != nil {
		return model.TaskVersion{}, err
	}

	return result, nil
}
