package model

import (
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/google/uuid"
)

func InitID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeInternal, "failed to generate UUID: %v", err))
	}

	return id.String(), nil
}
