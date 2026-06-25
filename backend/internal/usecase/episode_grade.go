package usecase

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/pagination"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type EpisodeGradeUsecase interface {
	Upsert(ctx context.Context, input EpisodeGradeUpsertInput) (model.EpisodeGrade, error)
	GetMyGrade(ctx context.Context, episodeID, userID string) (*model.EpisodeGrade, error)
	List(ctx context.Context, episodeID string, page, limit int) ([]repository.EpisodeGradeListItem, int, error)
}

type EpisodeGradeUpsertInput struct {
	EpisodeID      string
	UserID         string
	OrganizationID string
	Grade          float64
	Comment        *string
}

type episodeGrade struct {
	repo repository.EpisodeGrade
	db   repository.DBConn
}

func NewEpisodeGrade(repo repository.EpisodeGrade, db repository.DBConn) *episodeGrade {
	return &episodeGrade{repo: repo, db: db}
}

func (u *episodeGrade) Upsert(ctx context.Context, input EpisodeGradeUpsertInput) (model.EpisodeGrade, error) {
	grade, err := model.InitEpisodeGrade(input.OrganizationID, input.EpisodeID, input.UserID, input.Grade, input.Comment)
	if err != nil {
		return model.EpisodeGrade{}, err
	}

	return u.repo.Upsert(ctx, u.db, grade)
}

func (u *episodeGrade) GetMyGrade(ctx context.Context, episodeID, userID string) (*model.EpisodeGrade, error) {
	return u.repo.GetMyGrade(ctx, u.db, episodeID, userID)
}

func (u *episodeGrade) List(ctx context.Context, episodeID string, page, limit int) ([]repository.EpisodeGradeListItem, int, error) {
	if limit <= 0 {
		limit = pagination.DefaultLimit
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return u.repo.ListByEpisodeID(ctx, u.db, episodeID, limit, offset)
}
