package repository

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

// APIKeyListFilter is the set of optional filters supported by APIKey.List.
type APIKeyListFilter struct {
	RobotID        *string
	UserID         *string
	IncludeRevoked bool
}

// APIKey is the persistence-layer port for managing API keys.
//
// All methods operate within the organization scope of the context (enforced by
// the OrgScoped bun hook on the entity), except for FindActiveByHash which is
// intentionally cross-organization — the request is authenticated before an
// organization is known, so the lookup must scan all rows.
type APIKey interface {
	// FindActiveByHash looks up an API key by its SHA-256 hex hash, ignoring the
	// OrgScoped filter. Returns CodeAPIKeyNotFound when no row matches, the row is
	// revoked, or the row has expired.
	FindActiveByHash(ctx context.Context, conn Conn, keyHashHex string, now time.Time) (model.APIKey, error)

	// TouchLastUsedAt updates last_used_at for the given API key id (internal int64 id, not natural id).
	// Used by the async debounced flusher; does not enforce OrgScoped.
	TouchLastUsedAt(ctx context.Context, conn Conn, id int64, at time.Time) error

	Create(ctx context.Context, conn Conn, k model.APIKey) (model.APIKey, error)
	GetByNaturalID(ctx context.Context, conn Conn, idNatural string) (model.APIKey, error)
	List(ctx context.Context, conn Conn, filter APIKeyListFilter, limit, offset int) (model.APIKeys, int, error)
	Update(ctx context.Context, conn Conn, k model.APIKey) (model.APIKey, error)
	Revoke(ctx context.Context, conn Conn, idNatural string, at time.Time) (model.APIKey, error)
}
