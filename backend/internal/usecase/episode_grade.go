package usecase

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

type EpisodeGradeUsecase interface {
	Upsert(ctx context.Context, input EpisodeGradeUpsertInput) (model.EpisodeGrade, error)
	GetMyGrade(ctx context.Context, episodeID, userID string) (*model.EpisodeGrade, error)
	List(ctx context.Context, episodeID string, page, limit int) ([]EpisodeGradeListItem, int, error)
}

type EpisodeGradeListItem struct {
	Grade    model.EpisodeGrade
	UserName string
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
	data repository.DataAccess
}

func NewEpisodeGrade(repo repository.EpisodeGrade, data repository.DataAccess) *episodeGrade {
	return &episodeGrade{repo: repo, data: data}
}

func (u *episodeGrade) Upsert(ctx context.Context, input EpisodeGradeUpsertInput) (model.EpisodeGrade, error) {
	grade, err := model.InitEpisodeGrade(input.OrganizationID, input.EpisodeID, input.UserID, input.Grade, input.Comment)
	if err != nil {
		return model.EpisodeGrade{}, err
	}

	return u.repo.Upsert(ctx, u.data.Conn(), grade)
}

func (u *episodeGrade) GetMyGrade(ctx context.Context, episodeID, userID string) (*model.EpisodeGrade, error) {
	return u.repo.GetMyGrade(ctx, u.data.Conn(), episodeID, userID)
}

func (u *episodeGrade) List(ctx context.Context, episodeID string, page, limit int) ([]EpisodeGradeListItem, int, error) {
	if limit <= 0 {
		limit = pagination.DefaultLimit
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	items, total, err := u.repo.ListByEpisodeID(ctx, u.data.Conn(), episodeID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return episodeGradeListItems(items), total, nil
}

func episodeGradeListItems(items []repository.EpisodeGradeListItem) []EpisodeGradeListItem {
	result := make([]EpisodeGradeListItem, len(items))
	for i, item := range items {
		result[i] = EpisodeGradeListItem{
			Grade:    item.Grade,
			UserName: item.UserName,
		}
	}
	return result
}
