package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

var testTime = time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)

type stubUserRepo struct {
	existing     model.User
	updateCalled bool
	updateArg    model.User
}

func (s *stubUserRepo) Create(ctx context.Context, conn repository.DBConn, user model.User) (model.User, error) {
	return user, nil
}

func (s *stubUserRepo) Update(ctx context.Context, conn repository.DBConn, user model.User) (model.User, error) {
	s.updateCalled = true
	s.updateArg = user
	return user, nil
}

func (s *stubUserRepo) UpdateRole(ctx context.Context, conn repository.DBConn, idNatural string, role model.UserRole) (model.User, error) {
	return model.User{}, nil
}

func (s *stubUserRepo) GetByNaturalID(ctx context.Context, conn repository.DBConn, idNatural string) (model.User, error) {
	return s.existing, nil
}

func (s *stubUserRepo) ExistsByEmail(ctx context.Context, conn repository.DBConn, email string) (bool, error) {
	return false, nil
}

func (s *stubUserRepo) ExistsByEmails(ctx context.Context, conn repository.DBConn, emails []string) (map[string]bool, error) {
	return nil, nil
}

func (s *stubUserRepo) List(ctx context.Context, conn repository.DBConn, filter repository.UserListFilter, limit, offset int) (model.Users, int, error) {
	return nil, 0, nil
}

func (s *stubUserRepo) Delete(ctx context.Context, conn repository.DBConn, idNatural string) error {
	return nil
}

func newUserUsecaseWithStub(repo *stubUserRepo) *user {
	return &user{userRepo: repo, data: repository.NewDataAccess(nil, nil)}
}

func TestUserUsecase_Update_UsesPointerFieldsForPartialUpdate(t *testing.T) {
	repo := &stubUserRepo{
		existing: model.NewUser(
			1,
			"user-1",
			"org-1",
			"Alice",
			"alice@example.com",
			model.UserRoleOperator,
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
			"org-1",
			"Alice",
			"alice@example.com",
			model.UserRoleOperator,
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
		"org-1",
		"Alice",
		"alice@example.com",
		model.UserRoleOperator,
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
