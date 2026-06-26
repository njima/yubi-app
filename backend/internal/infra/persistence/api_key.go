package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type apiKey struct{}

func NewAPIKey() *apiKey { return &apiKey{} }

func apiKeyEntityToModel(e entity.APIKey) model.APIKey {
	var userRole model.UserRole
	var userName string
	if e.User != nil {
		userRole = model.UserRole(e.User.Role)
		userName = e.User.Name
	}
	var robotName *string
	if e.Robot != nil {
		n := e.Robot.Name
		robotName = &n
	}
	updatedAt := e.UpdatedAt
	return model.APIKey{
		ID:             e.ID,
		IDNatural:      e.IDNatural,
		OrganizationID: e.OrganizationID,
		UserID:         e.UserID,
		UserName:       userName,
		UserRole:       userRole,
		RobotID:        e.RobotID,
		RobotName:      robotName,
		Name:           e.Name,
		KeyHash:        e.KeyHash,
		KeyHint:        e.KeyHint,
		ExpiresAt:      e.ExpiresAt,
		LastUsedAt:     e.LastUsedAt,
		RevokedAt:      e.RevokedAt,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      &updatedAt,
	}
}

func (a *apiKey) FindActiveByHash(ctx context.Context, conn repository.DBConn, keyHashHex string, now time.Time) (model.APIKey, error) {
	// Cross-organization lookup: auth middleware runs before any organization is known,
	// so OrgIDFromContext returns false and the OrgScoped BeforeSelect hook is a no-op.
	// We use Model() to keep relations available.
	var e entity.APIKey
	err := conn.NewSelect().
		Model(&e).
		Relation("User").
		Where("ak.key_hash = ?", keyHashHex).
		Where("ak.revoked_at IS NULL").
		Where("(ak.expires_at IS NULL OR ak.expires_at > ?)", now).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.APIKey{}, apperror.NewError(apperror.NewMessage(apperror.CodeAPIKeyNotFound, "api key not found or inactive"))
		}
		return model.APIKey{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to lookup api key: %v", err))
	}
	if e.User == nil {
		return model.APIKey{}, apperror.NewError(
			apperror.NewMessage(apperror.CodeDatabaseError, "api key owner not found: user relation not loaded"),
		)
	}
	return apiKeyEntityToModel(e), nil
}

func (a *apiKey) TouchLastUsedAt(ctx context.Context, conn repository.DBConn, id int64, at time.Time) error {
	// Bypass OrgScoped by using TableExpr so this can run without an organization context
	// (e.g. from a background flusher started after the request returned).
	_, err := conn.NewUpdate().
		TableExpr(`"api_key"`).
		Set("last_used_at = ?", at).
		Set("updated_at = ?", at).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to touch last_used_at: %v", err))
	}
	return nil
}

func (a *apiKey) Create(ctx context.Context, conn repository.DBConn, k model.APIKey) (model.APIKey, error) {
	db := entity.APIKey{
		IDNatural:      k.IDNatural,
		OrganizationID: k.OrganizationID,
		UserID:         k.UserID,
		RobotID:        k.RobotID,
		Name:           k.Name,
		KeyHash:        k.KeyHash,
		KeyHint:        k.KeyHint,
		ExpiresAt:      k.ExpiresAt,
	}

	var inserted entity.APIKey
	if err := conn.NewInsert().
		Model(&db).
		Returning("*").
		Scan(ctx, &inserted); err != nil {
		return model.APIKey{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create api key: %v", err))
	}

	out := apiKeyEntityToModel(inserted)
	out.UserRole = k.UserRole
	return out, nil
}

func (a *apiKey) GetByNaturalID(ctx context.Context, conn repository.DBConn, idNatural string) (model.APIKey, error) {
	var e entity.APIKey
	err := conn.NewSelect().
		Model(&e).
		Relation("User").
		Relation("Robot").
		Where("ak.id_natural = ?", idNatural).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.APIKey{}, apperror.NewError(apperror.NewMessage(apperror.CodeAPIKeyNotFound, "api key not found: id_natural=%s", idNatural))
		}
		return model.APIKey{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get api key: %v", err))
	}
	return apiKeyEntityToModel(e), nil
}

func (a *apiKey) List(ctx context.Context, conn repository.DBConn, filter repository.APIKeyListFilter, limit, offset int) (model.APIKeys, int, error) {
	var rows []entity.APIKey
	q := conn.NewSelect().
		Model(&rows).
		Relation("User").
		Relation("Robot")

	if filter.RobotID != nil {
		q = q.Where("ak.robot_id = ?", *filter.RobotID)
	}
	if filter.UserID != nil {
		q = q.Where("ak.user_id = ?", *filter.UserID)
	}
	if !filter.IncludeRevoked {
		q = q.Where("ak.revoked_at IS NULL")
	}

	q = q.Order("ak.created_at DESC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	total, err := q.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list api keys: %v", err))
	}

	out := make(model.APIKeys, 0, len(rows))
	for i := range rows {
		m := apiKeyEntityToModel(rows[i])
		out = append(out, &m)
	}
	return out, total, nil
}

func (a *apiKey) Update(ctx context.Context, conn repository.DBConn, k model.APIKey) (model.APIKey, error) {
	upd := conn.NewUpdate().Model((*entity.APIKey)(nil))

	// name is always set (validation enforced upstream); allow expires_at to be nilable
	upd = upd.Set("name = ?", k.Name)
	upd = upd.Set("expires_at = ?", k.ExpiresAt)
	upd = upd.Set("updated_at = ?", time.Now().UTC())

	res, err := upd.Where("id_natural = ?", k.IDNatural).Exec(ctx)
	if err != nil {
		return model.APIKey{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update api key: %v", err))
	}
	n, err := res.RowsAffected()
	if err != nil {
		return model.APIKey{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update api key: %v", err))
	}
	if n == 0 {
		return model.APIKey{}, apperror.NewError(apperror.NewMessage(apperror.CodeAPIKeyNotFound, "api key not found: id_natural=%s", k.IDNatural))
	}

	// Re-fetch with relations so user_name/robot_name are populated in the response.
	return a.GetByNaturalID(ctx, conn, k.IDNatural)
}

func (a *apiKey) Revoke(ctx context.Context, conn repository.DBConn, idNatural string, at time.Time) (model.APIKey, error) {
	var updated entity.APIKey
	if err := conn.NewUpdate().
		Model((*entity.APIKey)(nil)).
		Set("revoked_at = COALESCE(revoked_at, ?)", at).
		Set("updated_at = ?", at).
		Where("id_natural = ?", idNatural).
		Returning("*").
		Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.APIKey{}, apperror.NewError(apperror.NewMessage(apperror.CodeAPIKeyNotFound, "api key not found: id_natural=%s", idNatural))
		}
		return model.APIKey{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to revoke api key: %v", err))
	}
	return apiKeyEntityToModel(updated), nil
}
