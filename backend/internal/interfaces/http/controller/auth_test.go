package controller

import (
	"context"
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func TestCreateAuthGoogleSessionReturnsProvisionedSession(t *testing.T) {
	session := usecase.AuthenticatedUserSession{
		User: model.User{
			IDNatural: "user-1",
			GoogleSub: "google-sub-1",
			Name:      "Ada Lovelace",
			Email:     "ada@example.com",
		},
		ActiveOrganization: model.Organization{
			IDNatural: "org-1",
			Name:      "Ada's Workspace",
			Kind:      model.OrganizationKindPersonal,
		},
		ActiveMembership: model.OrganizationMembership{
			UserID:         "user-1",
			OrganizationID: "org-1",
			Role:           model.UserRoleAdmin,
		},
	}
	userUC := &stubAuthGoogleUserUsecase{session: session}
	c := NewController(Dependencies{UserUsecase: userUC})

	got, err := c.CreateGoogleAuthSession(context.Background(), GoogleAuthSessionRequest{
		GoogleSub: "google-sub-1",
		Email:     "ada@example.com",
		Name:      "Ada Lovelace",
		AvatarURL: "https://example.com/avatar.png",
	})
	if err != nil {
		t.Fatalf("CreateGoogleAuthSession() error = %v", err)
	}
	if got.UserID != "user-1" || got.ActiveOrganizationID != "org-1" || got.ActiveRole != model.UserRoleAdmin {
		t.Fatalf("response = %+v, want provisioned user session", got)
	}
	if userUC.input.GoogleSub != "google-sub-1" || userUC.input.Email != "ada@example.com" {
		t.Fatalf("usecase input = %+v, want google profile", userUC.input)
	}
}

type stubAuthGoogleUserUsecase struct {
	input   usecase.GoogleUserInput
	session usecase.AuthenticatedUserSession
}

func (s *stubAuthGoogleUserUsecase) FindOrProvisionGoogleUser(ctx context.Context, input usecase.GoogleUserInput) (usecase.AuthenticatedUserSession, error) {
	s.input = input
	return s.session, nil
}

func (s *stubAuthGoogleUserUsecase) Create(ctx context.Context, input usecase.CreateInput) (model.User, error) {
	panic("not implemented")
}

func (s *stubAuthGoogleUserUsecase) Update(ctx context.Context, input usecase.UserUpdateInput) (model.User, error) {
	panic("not implemented")
}

func (s *stubAuthGoogleUserUsecase) UpdateRole(ctx context.Context, input usecase.UserRoleUpdateInput) (model.User, error) {
	panic("not implemented")
}

func (s *stubAuthGoogleUserUsecase) SetLocations(ctx context.Context, userID string, locationIDs []string) (model.User, error) {
	panic("not implemented")
}

func (s *stubAuthGoogleUserUsecase) SetSites(ctx context.Context, userID string, siteIDs []string) (model.User, error) {
	panic("not implemented")
}

func (s *stubAuthGoogleUserUsecase) GetByNaturalID(ctx context.Context, idNatural string) (model.User, error) {
	panic("not implemented")
}

func (s *stubAuthGoogleUserUsecase) List(ctx context.Context, filter usecase.UserListFilter, page, limit int) (model.Users, int, error) {
	panic("not implemented")
}

func (s *stubAuthGoogleUserUsecase) Delete(ctx context.Context, idNatural string) error {
	panic("not implemented")
}

func (s *stubAuthGoogleUserUsecase) GetAuthenticatedSession(ctx context.Context, userID string, organizationID *string) (usecase.AuthenticatedUserSession, error) {
	panic("not implemented")
}

func (s *stubAuthGoogleUserUsecase) ResolveActiveMembership(ctx context.Context, userID string, organizationID *string) (model.OrganizationMembership, error) {
	panic("not implemented")
}
