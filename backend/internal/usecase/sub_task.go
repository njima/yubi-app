package usecase

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

type SubTaskUsecase interface {
	Create(ctx context.Context, input SubTaskCreateInput) (model.SubTask, error)
	GetByID(ctx context.Context, id string) (model.SubTask, error)
	List(ctx context.Context, taskID, taskVersionID *string, page, limit int) (model.SubTasks, int, error)
	Update(ctx context.Context, input SubTaskUpdateInput) (model.SubTask, error)
	Reorder(ctx context.Context, input SubTaskReorderInput) (model.SubTasks, error)
	Delete(ctx context.Context, id string) error
}

type SubTaskCreateInput struct {
	OrganizationID        string
	TaskID                string
	TaskVersionID         string
	Name                  string
	Description           *string
	TargetDurationSeconds *int
}

type SubTaskUpdateInput struct {
	ID                    string
	Name                  *string
	OrderIndex            *int
	Description           *string
	TargetDurationSeconds *int
}

type SubTaskReorderInput struct {
	TaskVersionID string
	SubTaskIDs    []string
}

type subtask struct {
	repo   repository.SubTask
	rt     repository.Task
	tvRepo repository.TaskVersion
	data   repository.DataAccess
}

func NewSubTask(repo repository.SubTask, rt repository.Task, tvRepo repository.TaskVersion, data repository.DataAccess) *subtask {
	return &subtask{repo: repo, rt: rt, tvRepo: tvRepo, data: data}
}

func (s *subtask) Create(ctx context.Context, input SubTaskCreateInput) (model.SubTask, error) {
	var cst model.SubTask
	err := s.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		tv, err := s.tvRepo.GetByIDForUpdate(ctx, conn, input.TaskVersionID)
		if err != nil {
			return err
		}
		if !tv.IsDraft() {
			return apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "cannot edit approved task version: id=%s", input.TaskVersionID))
		}

		maxIndex, err := s.repo.GetMaxOrderIndex(ctx, conn, input.TaskVersionID)
		if err != nil {
			return err
		}

		st, err := model.InitSubTask(input.OrganizationID, input.TaskVersionID, input.Name, model.NextOrderIndex(maxIndex), input.Description, input.TargetDurationSeconds)
		if err != nil {
			return err
		}

		cst, err = s.repo.Create(ctx, conn, st)
		return err
	})
	if err != nil {
		return model.SubTask{}, err
	}

	return cst, nil
}

func (s *subtask) GetByID(ctx context.Context, id string) (model.SubTask, error) {
	return s.repo.GetByID(ctx, s.data.Conn(), id)
}

func (s *subtask) List(ctx context.Context, taskID, taskVersionID *string, page, limit int) (model.SubTasks, int, error) {
	if limit <= 0 {
		limit = pagination.DefaultLimit
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	filter := repository.SubTaskListFilter{
		TaskID:        taskID,
		TaskVersionID: taskVersionID,
	}
	return s.repo.List(ctx, s.data.Conn(), filter, limit, offset)
}

func (s *subtask) Update(ctx context.Context, input SubTaskUpdateInput) (model.SubTask, error) {
	existing, err := s.repo.GetByID(ctx, s.data.Conn(), input.ID)
	if err != nil {
		return model.SubTask{}, err
	}

	tv, err := s.tvRepo.GetByID(ctx, s.data.Conn(), existing.TaskVersionID)
	if err != nil {
		return model.SubTask{}, err
	}
	if !tv.IsDraft() {
		return model.SubTask{}, apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "cannot edit approved task version: id=%s", existing.TaskVersionID))
	}

	st := model.SubTask{IDNatural: input.ID}
	if input.Name != nil {
		st.Name = *input.Name
	}
	if input.OrderIndex != nil {
		st.OrderIndex = *input.OrderIndex
	}
	if input.Description != nil {
		st.Description = input.Description
	}
	if input.TargetDurationSeconds != nil {
		st.TargetDurationSeconds = input.TargetDurationSeconds
	}

	return s.repo.Update(ctx, s.data.Conn(), st)
}

func (s *subtask) Reorder(ctx context.Context, input SubTaskReorderInput) (model.SubTasks, error) {
	var result model.SubTasks
	err := s.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		tv, err := s.tvRepo.GetByIDForUpdate(ctx, conn, input.TaskVersionID)
		if err != nil {
			return err
		}
		if !tv.IsDraft() {
			return apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "cannot edit approved task version: id=%s", input.TaskVersionID))
		}

		existing, err := s.repo.GetByTaskVersionID(ctx, conn, input.TaskVersionID)
		if err != nil {
			return err
		}
		existingIDs := make(map[string]struct{}, len(existing))
		for _, st := range existing {
			existingIDs[st.IDNatural] = struct{}{}
		}
		if len(input.SubTaskIDs) != len(existingIDs) {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "subtask IDs count mismatch: expected %d, got %d", len(existingIDs), len(input.SubTaskIDs)))
		}
		for _, id := range input.SubTaskIDs {
			if _, ok := existingIDs[id]; !ok {
				return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "subtask ID does not belong to task version: id=%s", id))
			}
		}

		if err := s.repo.UpdateOrderIndices(ctx, conn, input.SubTaskIDs); err != nil {
			return err
		}

		result, err = s.repo.GetByTaskVersionID(ctx, conn, input.TaskVersionID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *subtask) Delete(ctx context.Context, id string) error {
	existing, err := s.repo.GetByID(ctx, s.data.Conn(), id)
	if err != nil {
		return err
	}

	tv, err := s.tvRepo.GetByID(ctx, s.data.Conn(), existing.TaskVersionID)
	if err != nil {
		return err
	}
	if !tv.IsDraft() {
		return apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "cannot edit approved task version: id=%s", existing.TaskVersionID))
	}

	return s.repo.Delete(ctx, s.data.Conn(), id)
}
