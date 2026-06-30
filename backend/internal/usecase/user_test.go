package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
)

var testTime = time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)

type stubUserRepo struct {
	existing     model.User
	googleSubErr error
	created      model.User
	updateCalled bool
	updateArg    model.User
}

func (s *stubUserRepo) Create(ctx context.Context, conn repository.Conn, user model.User) (model.User, error) {
	s.created = user
	return user, nil
}

func (s *stubUserRepo) Update(ctx context.Context, conn repository.Conn, user model.User) (model.User, error) {
	s.updateCalled = true
	s.updateArg = user
	return user, nil
}

func (s *stubUserRepo) GetByNaturalID(ctx context.Context, conn repository.Conn, idNatural string) (model.User, error) {
	return s.existing, nil
}

func (s *stubUserRepo) GetByGoogleSub(ctx context.Context, conn repository.Conn, googleSub string) (model.User, error) {
	if s.googleSubErr != nil {
		return model.User{}, s.googleSubErr
	}
	return s.existing, nil
}

func (s *stubUserRepo) ExistsByEmail(ctx context.Context, conn repository.Conn, email string) (bool, error) {
	return false, nil
}

func (s *stubUserRepo) ExistsByEmails(ctx context.Context, conn repository.Conn, emails []string) (map[string]bool, error) {
	return nil, nil
}

func (s *stubUserRepo) List(ctx context.Context, conn repository.Conn, filter repository.UserListFilter, limit, offset int) (model.Users, int, error) {
	return nil, 0, nil
}

func (s *stubUserRepo) Delete(ctx context.Context, conn repository.Conn, idNatural string) error {
	return nil
}

func newUserUsecaseWithStub(repo *stubUserRepo) *user {
	return &user{userRepo: repo, data: repository.NewDataAccess(nil, stubTxRunner{})}
}

type stubTxRunner struct{}

func (stubTxRunner) RunInTx(ctx context.Context, fn func(context.Context, repository.Conn) error) error {
	return fn(ctx, nil)
}

type googleProvisionOrganizationRepo struct {
	created model.Organization
	byID    map[string]model.Organization
	byName  map[string]model.Organization
}

func (s *googleProvisionOrganizationRepo) Create(ctx context.Context, conn repository.Conn, org model.Organization) (model.Organization, error) {
	if s.byName == nil {
		s.byName = map[string]model.Organization{}
	}
	if _, ok := s.byName[org.Name]; ok {
		return model.Organization{}, apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "organization name already exists"))
	}
	s.created = org
	if s.byID == nil {
		s.byID = map[string]model.Organization{}
	}
	s.byID[org.IDNatural] = org
	s.byName[org.Name] = org
	return org, nil
}

func (s *googleProvisionOrganizationRepo) GetByNaturalID(ctx context.Context, conn repository.Conn, idNatural string) (model.Organization, error) {
	if org, ok := s.byID[idNatural]; ok {
		return org, nil
	}
	return model.Organization{}, apperror.NewError(apperror.NewMessage(apperror.CodeOrganizationNotFound, "organization not found"))
}

func (s *googleProvisionOrganizationRepo) List(ctx context.Context, conn repository.Conn, limit, offset int) (model.Organizations, int, error) {
	return nil, 0, nil
}

func (s *googleProvisionOrganizationRepo) Update(ctx context.Context, conn repository.Conn, org model.Organization) (model.Organization, error) {
	return org, nil
}

func (s *googleProvisionOrganizationRepo) Delete(ctx context.Context, conn repository.Conn, idNatural string) error {
	return nil
}

type stubOrganizationMembershipRepo struct {
	created     model.OrganizationMembership
	memberships []model.OrganizationMembership
	updatedRole *model.UserRole
}

func (s *stubOrganizationMembershipRepo) Create(ctx context.Context, conn repository.Conn, membership model.OrganizationMembership) (model.OrganizationMembership, error) {
	s.created = membership
	s.memberships = append(s.memberships, membership)
	return membership, nil
}

func (s *stubOrganizationMembershipRepo) GetByUserAndOrganization(ctx context.Context, conn repository.Conn, userID, organizationID string) (model.OrganizationMembership, error) {
	for _, membership := range s.memberships {
		if membership.UserID == userID && membership.OrganizationID == organizationID {
			return membership, nil
		}
	}
	return model.OrganizationMembership{}, apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "organization membership not found"))
}

func (s *stubOrganizationMembershipRepo) ListByUser(ctx context.Context, conn repository.Conn, userID string) ([]model.OrganizationMembership, error) {
	var memberships []model.OrganizationMembership
	for _, membership := range s.memberships {
		if membership.UserID == userID {
			memberships = append(memberships, membership)
		}
	}
	return memberships, nil
}

func (s *stubOrganizationMembershipRepo) CountByUser(ctx context.Context, conn repository.Conn, userID string) (int, error) {
	memberships, err := s.ListByUser(ctx, conn, userID)
	if err != nil {
		return 0, err
	}
	return len(memberships), nil
}

func (s *stubOrganizationMembershipRepo) UpdateRole(ctx context.Context, conn repository.Conn, userID, organizationID string, role model.UserRole) (model.OrganizationMembership, error) {
	for i, membership := range s.memberships {
		if membership.UserID == userID && membership.OrganizationID == organizationID {
			s.memberships[i].Role = role
			s.updatedRole = &s.memberships[i].Role
			return s.memberships[i], nil
		}
	}
	return model.OrganizationMembership{}, apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "organization membership not found"))
}

type stubUserLocationRepo struct {
	organizationID string
}

func (s *stubUserLocationRepo) SetUserLocations(ctx context.Context, conn repository.Conn, userID string, organizationID string, locationIDs []string) error {
	s.organizationID = organizationID
	return nil
}

type stubUserSiteRepo struct {
	organizationID string
}

func (s *stubUserSiteRepo) SetUserSites(ctx context.Context, conn repository.Conn, userID string, organizationID string, siteIDs []string) error {
	s.organizationID = organizationID
	return nil
}

func TestUserUsecase_FindOrProvisionGoogleUser_CreatesIdentityPersonalOrganizationAndAdminMembership(t *testing.T) {
	userRepo := &stubUserRepo{
		googleSubErr: apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "user not found")),
	}
	orgRepo := &googleProvisionOrganizationRepo{}
	membershipRepo := &stubOrganizationMembershipRepo{}
	uc := &user{
		userRepo:       userRepo,
		orgRepo:        orgRepo,
		membershipRepo: membershipRepo,
		data:           repository.NewDataAccess(nil, stubTxRunner{}),
	}

	input := GoogleUserInput{
		GoogleSub: "google-oauth2|123",
		Email:     "ada@example.com",
		Name:      "Ada Lovelace",
		AvatarURL: "https://example.com/a.png",
	}
	got, err := uc.FindOrProvisionGoogleUser(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.User.Email != "ada@example.com" {
		t.Errorf("got.User.Email = %q, want %q", got.User.Email, "ada@example.com")
	}
	if got.User.GoogleSub != "google-oauth2|123" {
		t.Errorf("got.User.GoogleSub = %q, want %q", got.User.GoogleSub, "google-oauth2|123")
	}
	if got.User.Name != "Ada Lovelace" {
		t.Errorf("got.User.Name = %q, want %q", got.User.Name, "Ada Lovelace")
	}
	if got.User.AvatarURL == nil || *got.User.AvatarURL != "https://example.com/a.png" {
		t.Errorf("got.User.AvatarURL = %v, want %q", got.User.AvatarURL, "https://example.com/a.png")
	}
	if got.ActiveOrganization.Kind != model.OrganizationKindPersonal {
		t.Errorf("got.ActiveOrganization.Kind = %q, want %q", got.ActiveOrganization.Kind, model.OrganizationKindPersonal)
	}
	if got.ActiveOrganization.Name == "Ada Lovelace's Workspace" {
		t.Errorf("got.ActiveOrganization.Name = %q, want unique suffix", got.ActiveOrganization.Name)
	}
	if got.ActiveOrganization.Name == "" {
		t.Errorf("got.ActiveOrganization.Name is empty")
	}
	if got.ActiveMembership.UserID != got.User.IDNatural {
		t.Errorf("got.ActiveMembership.UserID = %q, want %q", got.ActiveMembership.UserID, got.User.IDNatural)
	}
	if got.ActiveMembership.OrganizationID != got.ActiveOrganization.IDNatural {
		t.Errorf("got.ActiveMembership.OrganizationID = %q, want %q", got.ActiveMembership.OrganizationID, got.ActiveOrganization.IDNatural)
	}
	if got.ActiveMembership.Role != model.UserRoleAdmin {
		t.Errorf("got.ActiveMembership.Role = %v, want %v", got.ActiveMembership.Role, model.UserRoleAdmin)
	}
}

func TestUserUsecase_FindOrProvisionGoogleUser_AllowsSameDisplayNamePersonalWorkspaces(t *testing.T) {
	orgRepo := &googleProvisionOrganizationRepo{}

	firstUserRepo := &stubUserRepo{
		googleSubErr: apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "user not found")),
	}
	firstMembershipRepo := &stubOrganizationMembershipRepo{}
	firstUC := &user{
		userRepo:       firstUserRepo,
		orgRepo:        orgRepo,
		membershipRepo: firstMembershipRepo,
		data:           repository.NewDataAccess(nil, stubTxRunner{}),
	}
	first, err := firstUC.FindOrProvisionGoogleUser(context.Background(), GoogleUserInput{
		GoogleSub: "google-oauth2|123",
		Email:     "ada@example.com",
		Name:      "Ada Lovelace",
	})
	if err != nil {
		t.Fatalf("first provision error: %v", err)
	}

	secondUserRepo := &stubUserRepo{
		googleSubErr: apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "user not found")),
	}
	secondMembershipRepo := &stubOrganizationMembershipRepo{}
	secondUC := &user{
		userRepo:       secondUserRepo,
		orgRepo:        orgRepo,
		membershipRepo: secondMembershipRepo,
		data:           repository.NewDataAccess(nil, stubTxRunner{}),
	}
	second, err := secondUC.FindOrProvisionGoogleUser(context.Background(), GoogleUserInput{
		GoogleSub: "google-oauth2|456",
		Email:     "ada.two@example.com",
		Name:      "Ada Lovelace",
	})
	if err != nil {
		t.Fatalf("second provision error: %v", err)
	}

	if first.ActiveOrganization.Name == second.ActiveOrganization.Name {
		t.Fatalf("personal organization names should differ, both got %q", first.ActiveOrganization.Name)
	}
}

func TestUserUsecase_SetLocations_UsesActiveOrganizationFromContext(t *testing.T) {
	repo := &stubUserRepo{
		existing: model.NewUser(
			1,
			"user-1",
			"google-oauth2|alice",
			"Alice",
			"alice@example.com",
			nil,
			testTime,
			nil,
		),
	}
	locationRepo := &stubUserLocationRepo{}
	membershipRepo := &stubOrganizationMembershipRepo{
		memberships: []model.OrganizationMembership{
			{UserID: "user-1", OrganizationID: "org-first", Role: model.UserRoleViewer},
			{UserID: "user-1", OrganizationID: "org-active", Role: model.UserRoleAdmin},
		},
	}
	uc := &user{
		userRepo:         repo,
		membershipRepo:   membershipRepo,
		userLocationRepo: locationRepo,
		data:             repository.NewDataAccess(nil, stubTxRunner{}),
	}
	ctx := requestctx.SetOrganizationID(context.Background(), "org-active")

	_, err := uc.SetLocations(ctx, "user-1", []string{"loc-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if locationRepo.organizationID != "org-active" {
		t.Errorf("location organizationID = %q, want %q", locationRepo.organizationID, "org-active")
	}
}

func TestUserUsecase_SetSites_UsesActiveOrganizationFromContext(t *testing.T) {
	repo := &stubUserRepo{
		existing: model.NewUser(
			1,
			"user-1",
			"google-oauth2|alice",
			"Alice",
			"alice@example.com",
			nil,
			testTime,
			nil,
		),
	}
	siteRepo := &stubUserSiteRepo{}
	membershipRepo := &stubOrganizationMembershipRepo{
		memberships: []model.OrganizationMembership{
			{UserID: "user-1", OrganizationID: "org-first", Role: model.UserRoleViewer},
			{UserID: "user-1", OrganizationID: "org-active", Role: model.UserRoleAdmin},
		},
	}
	uc := &user{
		userRepo:       repo,
		membershipRepo: membershipRepo,
		userSiteRepo:   siteRepo,
		data:           repository.NewDataAccess(nil, stubTxRunner{}),
	}
	ctx := requestctx.SetOrganizationID(context.Background(), "org-active")

	_, err := uc.SetSites(ctx, "user-1", []string{"site-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if siteRepo.organizationID != "org-active" {
		t.Errorf("site organizationID = %q, want %q", siteRepo.organizationID, "org-active")
	}
}

func TestUserUsecase_SetLocations_RejectsActiveOrganizationWithoutMembership(t *testing.T) {
	repo := &stubUserRepo{
		existing: model.NewUser(
			1,
			"user-1",
			"google-oauth2|alice",
			"Alice",
			"alice@example.com",
			nil,
			testTime,
			nil,
		),
	}
	locationRepo := &stubUserLocationRepo{}
	membershipRepo := &stubOrganizationMembershipRepo{
		memberships: []model.OrganizationMembership{
			{UserID: "user-1", OrganizationID: "org-first", Role: model.UserRoleViewer},
		},
	}
	uc := &user{
		userRepo:         repo,
		membershipRepo:   membershipRepo,
		userLocationRepo: locationRepo,
		data:             repository.NewDataAccess(nil, stubTxRunner{}),
	}
	ctx := requestctx.SetOrganizationID(context.Background(), "org-missing")

	_, err := uc.SetLocations(ctx, "user-1", []string{"loc-1"})
	if err == nil {
		t.Fatalf("expected error for active organization without membership")
	}
	if locationRepo.organizationID != "" {
		t.Errorf("location organizationID = %q, want no write", locationRepo.organizationID)
	}
}

func TestUserUsecase_UpdateRole_UpdatesMembershipForActiveOrganization(t *testing.T) {
	repo := &stubUserRepo{
		existing: model.NewUser(
			1,
			"user-1",
			"google-oauth2|alice",
			"Alice",
			"alice@example.com",
			nil,
			testTime,
			nil,
		),
	}
	membershipRepo := &stubOrganizationMembershipRepo{
		memberships: []model.OrganizationMembership{
			{UserID: "user-1", OrganizationID: "org-first", Role: model.UserRoleViewer},
			{UserID: "user-1", OrganizationID: "org-active", Role: model.UserRoleViewer},
		},
	}
	uc := &user{
		userRepo:       repo,
		membershipRepo: membershipRepo,
		data:           repository.NewDataAccess(nil, stubTxRunner{}),
	}
	ctx := requestctx.SetOrganizationID(context.Background(), "org-active")

	_, err := uc.UpdateRole(ctx, UserRoleUpdateInput{
		UserID: "user-1",
		Role:   model.UserRoleManager,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	active, err := membershipRepo.GetByUserAndOrganization(context.Background(), nil, "user-1", "org-active")
	if err != nil {
		t.Fatalf("active membership lookup error: %v", err)
	}
	if active.Role != model.UserRoleManager {
		t.Errorf("active membership role = %v, want %v", active.Role, model.UserRoleManager)
	}
	first, err := membershipRepo.GetByUserAndOrganization(context.Background(), nil, "user-1", "org-first")
	if err != nil {
		t.Fatalf("first membership lookup error: %v", err)
	}
	if first.Role != model.UserRoleViewer {
		t.Errorf("first membership role = %v, want %v", first.Role, model.UserRoleViewer)
	}
}

func TestUserUsecase_GetAuthenticatedSession_UsesActiveOrganizationMembership(t *testing.T) {
	userRepo := &stubUserRepo{
		existing: model.User{
			IDNatural: "user-1",
			Name:      "Operator",
			Email:     "operator@example.com",
		},
	}
	orgRepo := &googleProvisionOrganizationRepo{
		byID: map[string]model.Organization{
			"org-active": {
				IDNatural: "org-active",
				Name:      "Active Org",
				Kind:      model.OrganizationKindTeam,
			},
		},
	}
	membershipRepo := &stubOrganizationMembershipRepo{
		memberships: []model.OrganizationMembership{
			{UserID: "user-1", OrganizationID: "org-other", Role: model.UserRoleViewer},
			{UserID: "user-1", OrganizationID: "org-active", Role: model.UserRoleOperator},
		},
	}
	uc := &user{
		userRepo:       userRepo,
		orgRepo:        orgRepo,
		membershipRepo: membershipRepo,
		data:           repository.NewDataAccess(nil, stubTxRunner{}),
	}

	orgID := "org-active"
	got, err := uc.GetAuthenticatedSession(context.Background(), "user-1", &orgID)

	if err != nil {
		t.Fatalf("GetAuthenticatedSession() error = %v", err)
	}
	if got.User.IDNatural != "user-1" {
		t.Errorf("User.IDNatural = %q, want user-1", got.User.IDNatural)
	}
	if got.ActiveOrganization.IDNatural != "org-active" || got.ActiveOrganization.Name != "Active Org" {
		t.Errorf("ActiveOrganization = %+v, want org-active Active Org", got.ActiveOrganization)
	}
	if got.ActiveMembership.Role != model.UserRoleOperator {
		t.Errorf("ActiveMembership.Role = %v, want operator", got.ActiveMembership.Role)
	}
}

func TestUserUsecase_Update_UsesPointerFieldsForPartialUpdate(t *testing.T) {
	repo := &stubUserRepo{
		existing: model.NewUser(
			1,
			"user-1",
			"google-oauth2|alice",
			"Alice",
			"alice@example.com",
			nil,
			testTime,
			nil,
		),
	}
	uc := newUserUsecaseWithStub(repo)
	name := "Alice Smith"

	got, err := uc.Update(context.Background(), UserUpdateInput{
		UserID: "user-1",
		Name:   &name,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.updateCalled {
		t.Fatalf("expected repo.Update to be called")
	}
	if repo.updateArg.Name != name {
		t.Errorf("repo update Name = %q, want %q", repo.updateArg.Name, name)
	}
	if repo.updateArg.Email != "alice@example.com" {
		t.Errorf("repo update Email = %q, want existing email", repo.updateArg.Email)
	}
	if got.Name != name {
		t.Errorf("returned Name = %q, want %q", got.Name, name)
	}
}

func TestUserUsecase_Update_RejectsExplicitEmptyName(t *testing.T) {
	repo := &stubUserRepo{
		existing: model.NewUser(
			1,
			"user-1",
			"google-oauth2|alice",
			"Alice",
			"alice@example.com",
			nil,
			testTime,
			nil,
		),
	}
	uc := newUserUsecaseWithStub(repo)
	emptyName := ""

	_, err := uc.Update(context.Background(), UserUpdateInput{
		UserID: "user-1",
		Name:   &emptyName,
	})
	if err == nil {
		t.Fatalf("expected validation error for explicit empty name")
	}
	if repo.updateCalled {
		t.Fatalf("invalid update should not reach repo.Update")
	}
}

func TestUserUsecase_Update_NoFieldsReturnsExistingUser(t *testing.T) {
	existing := model.NewUser(
		1,
		"user-1",
		"google-oauth2|alice",
		"Alice",
		"alice@example.com",
		nil,
		testTime,
		nil,
	)
	repo := &stubUserRepo{existing: existing}
	uc := newUserUsecaseWithStub(repo)

	got, err := uc.Update(context.Background(), UserUpdateInput{
		UserID: "user-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.updateCalled {
		t.Fatalf("expected no repo.Update call when no fields are provided")
	}
	if got.Name != existing.Name || got.Email != existing.Email {
		t.Errorf("got %+v, want existing user %+v", got, existing)
	}
}
