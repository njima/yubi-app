package usecase

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
)

type RobotOperatorUsecase interface {
	Set(ctx context.Context, robotID string, operator model.RobotOperator) (*model.RobotOperator, error)
	Get(ctx context.Context, robotID string) (*model.RobotOperator, error)
	Clear(ctx context.Context, robotID string, userID string) error
}

type robotOperator struct {
	repo repository.RobotOperatorRepository
}

func NewRobotOperator(repo repository.RobotOperatorRepository) RobotOperatorUsecase {
	return &robotOperator{repo: repo}
}

// Set registers or refreshes the active operator. Uses atomic SET NX for
// initial lock acquisition to prevent race conditions. Returns the existing
// operator if locked by someone else (caller should return 409).
func (u *robotOperator) Set(ctx context.Context, robotID string, operator model.RobotOperator) (*model.RobotOperator, error) {
	// Try atomic acquire (SET NX EX)
	acquired, err := u.repo.SaveNX(ctx, robotID, operator)
	if err != nil {
		return nil, err
	}
	if acquired {
		return nil, nil
	}

	// Key exists — check if same user (heartbeat refresh)
	existing, err := u.repo.GetByRobotID(ctx, robotID)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.UserID != operator.UserID {
		return existing, apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "robot is locked by another operator"))
	}

	// Same user — refresh TTL
	if err := u.repo.Save(ctx, robotID, operator); err != nil {
		return nil, err
	}
	return nil, nil
}

// Get returns the current operator or nil if none.
func (u *robotOperator) Get(ctx context.Context, robotID string) (*model.RobotOperator, error) {
	return u.repo.GetByRobotID(ctx, robotID)
}

// Clear releases the operator lock. Only the active operator can clear.
func (u *robotOperator) Clear(ctx context.Context, robotID string, userID string) error {
	existing, err := u.repo.GetByRobotID(ctx, robotID)
	if err != nil {
		return err
	}
	if existing != nil && existing.UserID != userID {
		return apperror.NewError(apperror.NewMessage(apperror.CodeForbidden, "only the active operator can release the lock"))
	}
	return u.repo.Delete(ctx, robotID)
}
