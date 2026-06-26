package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
	"github.com/rs/zerolog"
)

// lastUsedFlushInterval is the minimum gap between two last_used_at writes for the same key.
// All requests landing inside this window for the same key result in at most one UPDATE.
const lastUsedFlushInterval = 60 * time.Second

// APIKeyAuthResult is what the auth middleware needs to populate request context after a successful X-API-Key lookup.
type APIKeyAuthResult struct {
	APIKeyID       int64
	APIKeyIDNat    string
	UserID         string
	UserRole       model.UserRole
	RobotID        string
	OrganizationID string
}

type APIKeyCreateInput struct {
	Name      string
	RobotID   string
	ExpiresAt *time.Time
}

type APIKeyUpdateInput struct {
	IDNatural   string
	Name        *string
	ExpiresAt   *time.Time // pointer to *time.Time pattern would be ambiguous; SetExpiresAt iff non-nil pointer
	ClearExpiry bool       // explicit "set to NULL" flag
}

type APIKeyCreateResult struct {
	APIKey model.APIKey
	RawKey string // returned exactly once on creation
}

type APIKeyUsecase interface {
	// Authenticate looks up an API key by raw key (not hash). Returns auth result on success.
	Authenticate(ctx context.Context, rawKey string) (APIKeyAuthResult, error)
	// MarkUsed enqueues an async last_used_at update for the given key id.
	MarkUsed(apiKeyID int64)
	// List returns API keys in the caller's organization. Admin-only at the controller layer.
	List(ctx context.Context, filter APIKeyListFilter, page, limit int) (model.APIKeys, int, error)
	Create(ctx context.Context, input APIKeyCreateInput) (APIKeyCreateResult, error)
	Get(ctx context.Context, idNatural string) (model.APIKey, error)
	Update(ctx context.Context, input APIKeyUpdateInput) (model.APIKey, error)
	Revoke(ctx context.Context, idNatural string) error
}

type apiKey struct {
	repo      repository.APIKey
	userRepo  repository.User
	robotRepo repository.Robot
	data      repository.DataAccess
	logger    zerolog.Logger

	mu            sync.Mutex
	lastUsedFlush map[int64]time.Time // id -> last successful flush time; intentionally unbounded, cleared on process restart
	flushInterval time.Duration
	now           func() time.Time
}

func NewAPIKey(
	repo repository.APIKey,
	userRepo repository.User,
	robotRepo repository.Robot,
	data repository.DataAccess,
	logger zerolog.Logger,
) *apiKey {
	return &apiKey{
		repo:          repo,
		userRepo:      userRepo,
		robotRepo:     robotRepo,
		data:          data,
		logger:        logger,
		lastUsedFlush: make(map[int64]time.Time),
		flushInterval: lastUsedFlushInterval,
		now:           func() time.Time { return time.Now().UTC() },
	}
}

func (a *apiKey) Authenticate(ctx context.Context, rawKey string) (APIKeyAuthResult, error) {
	if rawKey == "" {
		return APIKeyAuthResult{}, apperror.NewError(apperror.NewMessage(apperror.CodeUnauthorized, "api key is empty"))
	}

	hash := model.HashAPIKey(rawKey)
	k, err := a.repo.FindActiveByHash(ctx, a.data.Conn(), hash, a.now())
	if err != nil {
		// Log the failure with the non-sensitive hint so an operator can
		// correlate 401s with the key actually in flight, without leaking the
		// raw key value into logs.
		hint := model.APIKeyHint(rawKey)
		reason := "lookup_failed"
		if apperror.SameKind(err, apperror.KindNotFound) {
			reason = "not_found_or_inactive"
		}
		a.logger.Warn().
			Err(err).
			Str("event", "api_key_auth_fail").
			Str("reason", reason).
			Str("key_hint", hint).
			Msg("api key authentication failed")
		return APIKeyAuthResult{}, err
	}

	// robot_id is required for the current "1 key = 1 robot" policy.
	if k.RobotID == nil || *k.RobotID == "" {
		a.logger.Warn().
			Str("event", "api_key_auth_fail").
			Str("reason", "no_robot_binding").
			Str("key_hint", k.KeyHint).
			Int64("api_key_id", k.ID).
			Msg("api key authentication failed")
		return APIKeyAuthResult{}, apperror.NewError(apperror.NewMessage(apperror.CodeUnauthorized, "api key has no robot binding"))
	}

	// Defensive: the owner user must exist. The FK is ON DELETE CASCADE so a
	// missing user should never resolve to a live api_key row in practice, but
	// guard anyway so a buggy join (LEFT JOIN with NULL user) never silently
	// elevates an authenticated request to UserRole=0 (Admin).
	if k.UserID == "" || k.UserName == "" {
		a.logger.Warn().
			Str("event", "api_key_auth_fail").
			Str("reason", "user_not_loaded").
			Str("key_hint", k.KeyHint).
			Int64("api_key_id", k.ID).
			Msg("api key authentication failed")
		return APIKeyAuthResult{}, apperror.NewError(apperror.NewMessage(apperror.CodeUnauthorized, "api key owner missing"))
	}

	return APIKeyAuthResult{
		APIKeyID:       k.ID,
		APIKeyIDNat:    k.IDNatural,
		UserID:         k.UserID,
		UserRole:       k.UserRole,
		RobotID:        *k.RobotID,
		OrganizationID: k.OrganizationID,
	}, nil
}

// MarkUsed updates last_used_at asynchronously with per-key debouncing.
// Safe to call from request hot path. Errors are logged and swallowed.
func (a *apiKey) MarkUsed(apiKeyID int64) {
	now := a.now()
	a.mu.Lock()
	last, ok := a.lastUsedFlush[apiKeyID]
	if ok && now.Sub(last) < a.flushInterval {
		a.mu.Unlock()
		return
	}
	a.lastUsedFlush[apiKeyID] = now
	a.mu.Unlock()

	go func() {
		// Detach from request context: this runs after the response is sent.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := a.repo.TouchLastUsedAt(ctx, a.data.Conn(), apiKeyID, now); err != nil {
			a.logger.Warn().Err(err).Int64("api_key_id", apiKeyID).Msg("failed to flush last_used_at")
		}
	}()
}

func (a *apiKey) List(ctx context.Context, filter APIKeyListFilter, page, limit int) (model.APIKeys, int, error) {
	if limit <= 0 {
		limit = 50
	}
	offset := 0
	if page > 1 {
		offset = (page - 1) * limit
	}
	return a.repo.List(ctx, a.data.Conn(), filter.repositoryFilter(), limit, offset)
}

func (a *apiKey) Create(ctx context.Context, input APIKeyCreateInput) (APIKeyCreateResult, error) {
	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil {
		return APIKeyCreateResult{}, err
	}
	userID, err := requestctx.UserID(ctx)
	if err != nil {
		return APIKeyCreateResult{}, err
	}
	userRole, err := requestctx.UserRole(ctx)
	if err != nil {
		return APIKeyCreateResult{}, err
	}

	if input.RobotID == "" {
		return APIKeyCreateResult{}, apperror.NewError(apperror.NewMessage(apperror.CodeValidationError, "robot_id is required"))
	}
	// Verify the robot exists in the caller's organization. The OrgScoped hook
	// filters by org automatically, so a not-found error here means cross-org access too.
	if _, err := a.robotRepo.GetByID(ctx, a.data.Conn(), input.RobotID); err != nil {
		return APIKeyCreateResult{}, err
	}

	if input.ExpiresAt != nil && !input.ExpiresAt.After(a.now()) {
		return APIKeyCreateResult{}, apperror.NewError(
			apperror.NewMessage(apperror.CodeValidationError, "expires_at must be in the future"),
		)
	}

	robotID := input.RobotID
	domainKey, raw, err := model.InitAPIKey(orgID, userID, userRole, &robotID, input.Name, input.ExpiresAt)
	if err != nil {
		return APIKeyCreateResult{}, err
	}

	if _, err := a.repo.Create(ctx, a.data.Conn(), domainKey); err != nil {
		return APIKeyCreateResult{}, err
	}

	// Reload with relations so user_name / robot_name are populated for the
	// response. Create itself does not load relations.
	created, err := a.repo.GetByNaturalID(ctx, a.data.Conn(), domainKey.IDNatural)
	if err != nil {
		return APIKeyCreateResult{}, err
	}
	created.UserRole = userRole

	return APIKeyCreateResult{APIKey: created, RawKey: raw}, nil
}

func (a *apiKey) Get(ctx context.Context, idNatural string) (model.APIKey, error) {
	return a.repo.GetByNaturalID(ctx, a.data.Conn(), idNatural)
}

func (a *apiKey) Update(ctx context.Context, input APIKeyUpdateInput) (model.APIKey, error) {
	existing, err := a.repo.GetByNaturalID(ctx, a.data.Conn(), input.IDNatural)
	if err != nil {
		return model.APIKey{}, err
	}

	if input.Name != nil {
		if err := existing.SetName(*input.Name); err != nil {
			return model.APIKey{}, err
		}
	}
	if input.ClearExpiry {
		existing.SetExpiresAt(nil)
	} else if input.ExpiresAt != nil {
		existing.SetExpiresAt(input.ExpiresAt)
	}

	return a.repo.Update(ctx, a.data.Conn(), existing)
}

func (a *apiKey) Revoke(ctx context.Context, idNatural string) error {
	_, err := a.repo.Revoke(ctx, a.data.Conn(), idNatural, a.now())
	return err
}
