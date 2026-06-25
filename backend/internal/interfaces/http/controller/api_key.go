package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/pagination"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func apiKeyToResponse(k model.APIKey) openapi.ApiKeyResponse {
	var updatedAt = k.CreatedAt
	if k.UpdatedAt != nil {
		updatedAt = *k.UpdatedAt
	}
	return openapi.ApiKeyResponse{
		Id:             k.IDNatural,
		Name:           k.Name,
		UserId:         k.UserID,
		UserName:       k.UserName,
		OrganizationId: k.OrganizationID,
		RobotId:        k.RobotID,
		RobotName:      k.RobotName,
		KeyHint:        k.KeyHint,
		ExpiresAt:      k.ExpiresAt,
		LastUsedAt:     k.LastUsedAt,
		RevokedAt:      k.RevokedAt,
		CreatedAt:      k.CreatedAt,
		UpdatedAt:      updatedAt,
	}
}

func (c *controller) ListApiKeys(ctx context.Context, request openapi.ListApiKeysRequestObject) (openapi.ListApiKeysResponseObject, error) {
	params := request.Params
	pg := pagination.Parse(params.Page, params.Limit)

	filter := usecase.APIKeyListFilter{
		RobotID: params.RobotId,
		UserID:  params.UserId,
	}
	if params.IncludeRevoked != nil {
		filter.IncludeRevoked = *params.IncludeRevoked
	}

	keys, total, err := c.apiKeyUsecase.List(ctx, filter, pg.Page, pg.Limit)
	if err != nil {
		return nil, err
	}

	respKeys := make([]openapi.ApiKeyResponse, 0, len(keys))
	for _, k := range keys {
		respKeys = append(respKeys, apiKeyToResponse(*k))
	}

	return openapi.ListApiKeys200JSONResponse{
		ApiKeys: respKeys,
		Pagination: openapi.Pagination{
			Count: total,
			Limit: pg.Limit,
			Page:  pg.Page,
		},
	}, nil
}

func (c *controller) CreateApiKey(ctx context.Context, request openapi.CreateApiKeyRequestObject) (openapi.CreateApiKeyResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	out, err := c.apiKeyUsecase.Create(ctx, usecase.APIKeyCreateInput{
		Name:      request.Body.Name,
		RobotID:   request.Body.RobotId,
		ExpiresAt: request.Body.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	resp := apiKeyToResponse(out.APIKey)
	updatedAt := resp.UpdatedAt
	return openapi.CreateApiKey201JSONResponse{
		Id:             resp.Id,
		Name:           resp.Name,
		UserId:         resp.UserId,
		UserName:       resp.UserName,
		OrganizationId: resp.OrganizationId,
		RobotId:        resp.RobotId,
		RobotName:      resp.RobotName,
		KeyHint:        resp.KeyHint,
		ExpiresAt:      resp.ExpiresAt,
		LastUsedAt:     resp.LastUsedAt,
		RevokedAt:      resp.RevokedAt,
		CreatedAt:      resp.CreatedAt,
		UpdatedAt:      updatedAt,
		Key:            out.RawKey,
	}, nil
}

func (c *controller) GetApiKey(ctx context.Context, request openapi.GetApiKeyRequestObject) (openapi.GetApiKeyResponseObject, error) {
	k, err := c.apiKeyUsecase.Get(ctx, request.ApiKeyId)
	if err != nil {
		return nil, err
	}
	return openapi.GetApiKey200JSONResponse(apiKeyToResponse(k)), nil
}

func (c *controller) UpdateApiKey(ctx context.Context, request openapi.UpdateApiKeyRequestObject) (openapi.UpdateApiKeyResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	in := usecase.APIKeyUpdateInput{
		IDNatural: request.ApiKeyId,
		Name:      request.Body.Name,
		ExpiresAt: request.Body.ExpiresAt,
	}
	if request.Body.ClearExpiresAt != nil && *request.Body.ClearExpiresAt {
		in.ClearExpiry = true
	}

	updated, err := c.apiKeyUsecase.Update(ctx, in)
	if err != nil {
		return nil, err
	}
	return openapi.UpdateApiKey200JSONResponse(apiKeyToResponse(updated)), nil
}

func (c *controller) RevokeApiKey(ctx context.Context, request openapi.RevokeApiKeyRequestObject) (openapi.RevokeApiKeyResponseObject, error) {
	if err := c.apiKeyUsecase.Revoke(ctx, request.ApiKeyId); err != nil {
		return nil, err
	}
	return openapi.RevokeApiKey204Response{}, nil
}
