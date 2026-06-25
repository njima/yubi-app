package usecase

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/ccontext"
	"github.com/airoa-org/yubi-app/backend/internal/event"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type SubTaskActionInput struct {
	EpisodeID string
	SubTaskID string
}

type EpisodeSubTaskUsecase interface {
	Complete(ctx context.Context, input SubTaskActionInput) error
	Skip(ctx context.Context, input SubTaskActionInput) error
}

type episodeSubTask struct {
	episodeRepo        repository.Episode
	episodeSubTaskRepo repository.EpisodeSubTask
	data               repository.DataAccess
	bus                *event.Bus
	robotBus           *event.Bus
	listBus            *event.Bus
}

func NewEpisodeSubTask(
	episodeRepo repository.Episode,
	episodeSubTaskRepo repository.EpisodeSubTask,
	data repository.DataAccess,
	bus *event.Bus,
	robotBus *event.Bus,
	listBus *event.Bus,
) EpisodeSubTaskUsecase {
	return &episodeSubTask{
		episodeRepo:        episodeRepo,
		episodeSubTaskRepo: episodeSubTaskRepo,
		data:               data,
		bus:                bus,
		robotBus:           robotBus,
		listBus:            listBus,
	}
}

func (e *episodeSubTask) Complete(ctx context.Context, input SubTaskActionInput) error {
	robotID, err := ccontext.RobotID(ctx)
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

		subtask, err := e.episodeSubTaskRepo.GetByID(ctx, conn, input.SubTaskID)
		if err != nil {
			return err
		}

		if subtask.EpisodeID != input.EpisodeID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "subtask does not belong to the episode"))
		}

		if err := subtask.Complete(); err != nil {
			return err
		}

		if err := e.episodeSubTaskRepo.Update(ctx, conn, subtask); err != nil {
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

func (e *episodeSubTask) Skip(ctx context.Context, input SubTaskActionInput) error {
	robotID, err := ccontext.RobotID(ctx)
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

		subtask, err := e.episodeSubTaskRepo.GetByID(ctx, conn, input.SubTaskID)
		if err != nil {
			return err
		}

		if subtask.EpisodeID != input.EpisodeID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "subtask does not belong to the episode"))
		}

		if err := subtask.Skip(); err != nil {
			return err
		}

		if err := e.episodeSubTaskRepo.Update(ctx, conn, subtask); err != nil {
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
