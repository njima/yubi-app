package controller

import "github.com/airoa-org/yubi-app/backend/internal/shared/apperror"

func requiredBody[T any](body *T) (*T, error) {
	if body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}
	return body, nil
}
