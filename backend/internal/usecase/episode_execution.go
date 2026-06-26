package usecase

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/requestctx"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/eventbus"
)

type CreateExecutionInput struct {
	EpisodeID string
	SubTaskID string
}

type ExecutionActionInput struct {
	EpisodeID   string
	SubTaskID   string
	ExecutionID string
	OccurredAt  time.Time
}

type CancelExecutionInput struct {
	EpisodeID   string
	SubTaskID   string
	ExecutionID string
}

type EpisodeExecutionUsecase interface {
	Create(ctx context.Context, input CreateExecutionInput) (string, error)
	Start(ctx context.Context, input ExecutionActionInput) error
	Finish(ctx context.Context, input ExecutionActionInput) error
	Cancel(ctx context.Context, input CancelExecutionInput) error
}

type episodeExecution struct {
	episodeRepo        repository.Episode
	episodeSubTaskRepo repository.EpisodeSubTask
	executionRepo      repository.EpisodeSubTaskExecution
	data               repository.DataAccess
	bus                *eventbus.Bus
	robotBus           *eventbus.Bus
	listBus            *eventbus.Bus
}

func NewEpisodeExecution(
	episodeRepo repository.Episode,
	episodeSubTaskRepo repository.EpisodeSubTask,
	executionRepo repository.EpisodeSubTaskExecution,
	data repository.DataAccess,
	bus *eventbus.Bus,
	robotBus *eventbus.Bus,
	listBus *eventbus.Bus,
) EpisodeExecutionUsecase {
	return &episodeExecution{
		episodeRepo:        episodeRepo,
		episodeSubTaskRepo: episodeSubTaskRepo,
		executionRepo:      executionRepo,
		data:               data,
		bus:                bus,
		robotBus:           robotBus,
		listBus:            listBus,
	}
}

func (e *episodeExecution) Create(ctx context.Context, input CreateExecutionInput) (string, error) {
	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil {
		return "", err
	}

	robotID, err := requestctx.RobotID(ctx)
	if err != nil {
		return "", err
	}

	var executionID string
	err = e.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		episode, err := e.episodeRepo.GetByID(ctx, conn, input.EpisodeID)
		if err != nil {
			return err
		}

		if episode.RobotID != robotID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeForbidden, "robot is not authorized to operate this episode"))
		}

		subtask, err := e.episodeSubTaskRepo.GetByID(ctx, conn, input.SubTaskID)
		if err != nil {
			return err
		}

		if subtask.EpisodeID != input.EpisodeID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "subtask does not belong to the episode"))
		}

		if err := subtask.StartProgress(); err != nil {
			return err
		}
		if err := e.episodeSubTaskRepo.Update(ctx, conn, subtask); err != nil {
			return err
		}

		execution, err := model.InitEpisodeSubTaskExecution(orgID, input.SubTaskID)
		if err != nil {
			return err
		}

		created, err := e.executionRepo.Create(ctx, conn, execution)
		if err != nil {
			return err
		}

		executionID = created.IDNatural
		return nil
	})

	if err != nil {
		return "", err
	}

	e.bus.Notify(input.EpisodeID)
	e.robotBus.Notify(robotID)
	e.listBus.Notify("list")
	return executionID, nil
}

func (e *episodeExecution) Start(ctx context.Context, input ExecutionActionInput) error {
	robotID, err := requestctx.RobotID(ctx)
	if err != nil {
		return err
	}

	err = e.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		episode, err := e.episodeRepo.GetByID(ctx, conn, input.EpisodeID)
		if err != nil {
			return err
		}

		if episode.RobotID != robotID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeForbidden, "robot is not authorized to operate this episode"))
		}

		execution, err := e.executionRepo.GetByID(ctx, conn, input.ExecutionID)
		if err != nil {
			return err
		}

		subtask, err := e.episodeSubTaskRepo.GetByID(ctx, conn, input.SubTaskID)
		if err != nil {
			return err
		}
		if subtask.EpisodeID != input.EpisodeID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "subtask does not belong to the episode"))
		}

		if execution.EpisodeSubTaskID != input.SubTaskID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "execution does not belong to the subtask"))
		}

		count, err := e.executionRepo.CountStartedBySubTaskID(ctx, conn, input.SubTaskID)
		if err != nil {
			return err
		}
		if count > 0 {
			return apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "another execution is already running for this subtask"))
		}

		if err := execution.Start(input.OccurredAt); err != nil {
			return err
		}

		if err := e.executionRepo.Update(ctx, conn, execution); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	e.bus.Notify(input.EpisodeID)
	e.robotBus.Notify(robotID)
	e.listBus.Notify("list")
	return nil
}

func (e *episodeExecution) Finish(ctx context.Context, input ExecutionActionInput) error {
	robotID, err := requestctx.RobotID(ctx)
	if err != nil {
		return err
	}

	err = e.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		episode, err := e.episodeRepo.GetByID(ctx, conn, input.EpisodeID)
		if err != nil {
			return err
		}

		if episode.RobotID != robotID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeForbidden, "robot is not authorized to operate this episode"))
		}

		execution, err := e.executionRepo.GetByID(ctx, conn, input.ExecutionID)
		if err != nil {
			return err
		}

		subtask, err := e.episodeSubTaskRepo.GetByID(ctx, conn, input.SubTaskID)
		if err != nil {
			return err
		}
		if subtask.EpisodeID != input.EpisodeID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "subtask does not belong to the episode"))
		}

		if execution.EpisodeSubTaskID != input.SubTaskID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "execution does not belong to the subtask"))
		}

		if err := execution.Finish(input.OccurredAt); err != nil {
			return err
		}

		if err := e.executionRepo.Update(ctx, conn, execution); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	e.bus.Notify(input.EpisodeID)
	e.robotBus.Notify(robotID)
	e.listBus.Notify("list")
	return nil
}

func (e *episodeExecution) Cancel(ctx context.Context, input CancelExecutionInput) error {
	robotID, err := requestctx.RobotID(ctx)
	if err != nil {
		return err
	}

	err = e.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		episode, err := e.episodeRepo.GetByID(ctx, conn, input.EpisodeID)
		if err != nil {
			return err
		}

		if episode.RobotID != robotID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeForbidden, "robot is not authorized to operate this episode"))
		}

		execution, err := e.executionRepo.GetByID(ctx, conn, input.ExecutionID)
		if err != nil {
			return err
		}

		subtask, err := e.episodeSubTaskRepo.GetByID(ctx, conn, input.SubTaskID)
		if err != nil {
			return err
		}
		if subtask.EpisodeID != input.EpisodeID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "subtask does not belong to the episode"))
		}

		if execution.EpisodeSubTaskID != input.SubTaskID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "execution does not belong to the subtask"))
		}

		if err := execution.Cancel(); err != nil {
			return err
		}

		if err := e.executionRepo.Update(ctx, conn, execution); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	e.bus.Notify(input.EpisodeID)
	e.robotBus.Notify(robotID)
	e.listBus.Notify("list")
	return nil
}
