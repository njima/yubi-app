package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

type GoogleAuthSessionRequest struct {
	GoogleSub string `json:"google_sub"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type GoogleAuthSessionResponse struct {
	UserID                 string         `json:"user_id"`
	Email                  string         `json:"email"`
	DisplayName            string         `json:"display_name"`
	AvatarURL              *string        `json:"avatar_url"`
	ActiveOrganizationID   string         `json:"active_organization_id"`
	ActiveOrganizationName string         `json:"active_organization_name"`
	ActiveRole             model.UserRole `json:"active_role"`
}

func (c *controller) CreateGoogleAuthSession(ctx context.Context, req GoogleAuthSessionRequest) (GoogleAuthSessionResponse, error) {
	session, err := c.userUsecase.FindOrProvisionGoogleUser(ctx, usecase.GoogleUserInput{
		GoogleSub: req.GoogleSub,
		Email:     req.Email,
		Name:      req.Name,
		AvatarURL: req.AvatarURL,
	})
	if err != nil {
		return GoogleAuthSessionResponse{}, err
	}

	return GoogleAuthSessionResponse{
		UserID:                 session.User.IDNatural,
		Email:                  session.User.Email,
		DisplayName:            session.User.Name,
		AvatarURL:              session.User.AvatarURL,
		ActiveOrganizationID:   session.ActiveOrganization.IDNatural,
		ActiveOrganizationName: session.ActiveOrganization.Name,
		ActiveRole:             session.ActiveMembership.Role,
	}, nil
}
