package model

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	apiKeyRandomBytes = 32 // 256 bit of entropy
	apiKeyHintHead    = 4
	apiKeyHintTail    = 4
)

// APIKey is the domain representation of a stored API key.
// The raw key string is never persisted; only its SHA-256 hash and a short hint are kept.
//
// UserName / RobotName are denormalized fields populated by repository JOINs
// for UI display. They are not stored on the api_key row itself.
type APIKey struct {
	ID             int64
	IDNatural      string
	OrganizationID string
	UserID         string
	UserName       string
	UserRole       UserRole
	RobotID        *string
	RobotName      *string
	Name           string
	KeyHash        string
	KeyHint        string
	ExpiresAt      *time.Time
	LastUsedAt     *time.Time
	RevokedAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type APIKeys []*APIKey

// IsActive reports whether the key can be used for authentication right now.
func (a APIKey) IsActive(now time.Time) bool {
	if a.RevokedAt != nil {
		return false
	}
	if a.ExpiresAt != nil && !a.ExpiresAt.After(now) {
		return false
	}
	return true
}

// GenerateAPIKey returns a freshly generated raw key and its SHA-256 hex hash.
// The raw key is 32 bytes of crypto/rand entropy encoded as base64url (no padding, 43 chars).
func GenerateAPIKey() (raw string, hashHex string, hint string, err error) {
	buf := make([]byte, apiKeyRandomBytes)
	if _, err = rand.Read(buf); err != nil {
		return "", "", "", apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeInternal, "failed to generate random bytes: %v", err))
	}
	raw = base64.RawURLEncoding.EncodeToString(buf)
	hashHex = HashAPIKey(raw)
	hint = APIKeyHint(raw)
	return raw, hashHex, hint, nil
}

// HashAPIKey returns the SHA-256 hex digest of the raw key.
func HashAPIKey(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// CompareHashes returns true if two hex SHA-256 hashes are equal in constant time.
func CompareHashes(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// APIKeyHint returns a short display hint for the raw key.
// Format: "<first 4>…<last 4>". Falls back to the whole string if too short.
func APIKeyHint(raw string) string {
	if len(raw) <= apiKeyHintHead+apiKeyHintTail {
		return raw
	}
	return raw[:apiKeyHintHead] + "…" + raw[len(raw)-apiKeyHintTail:]
}

// InitAPIKey creates a new API key domain object with a freshly generated raw key.
// Returns the model plus the raw key (which must be returned to the caller exactly once
// and never persisted).
func InitAPIKey(
	organizationID, userID string,
	userRole UserRole,
	robotID *string,
	name string,
	expiresAt *time.Time,
) (key APIKey, raw string, err error) {
	id, err := InitID()
	if err != nil {
		return APIKey{}, "", err
	}

	raw, hashHex, hint, err := GenerateAPIKey()
	if err != nil {
		return APIKey{}, "", err
	}

	key = APIKey{
		IDNatural:      id,
		OrganizationID: organizationID,
		UserID:         userID,
		UserRole:       userRole,
		RobotID:        robotID,
		Name:           name,
		KeyHash:        hashHex,
		KeyHint:        hint,
		ExpiresAt:      expiresAt,
		CreatedAt:      time.Now().UTC(),
	}

	if err := key.validate(); err != nil {
		return APIKey{}, "", err
	}

	return key, raw, nil
}

func (a APIKey) validate() error {
	if err := (validation.Errors{
		"id_natural":      validation.Validate(a.IDNatural, validation.Required.Error("id_natural is required")),
		"organization_id": validation.Validate(a.OrganizationID, validation.Required.Error("organization_id is required")),
		"user_id":         validation.Validate(a.UserID, validation.Required.Error("user_id is required")),
		"name": validation.Validate(
			a.Name,
			validation.Required.Error("name is required"),
			validation.RuneLength(1, 255).Error("name must be between 1 and 255 characters"),
		),
		"key_hash": validation.Validate(
			a.KeyHash,
			validation.Required.Error("key_hash is required"),
			validation.RuneLength(64, 64).Error("key_hash must be 64 hex characters"),
		),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "api_key validation failed: %v", err))
	}
	return nil
}

// SetName updates the human-readable label, validating length.
func (a *APIKey) SetName(name string) error {
	a.Name = name
	return a.validate()
}

// SetExpiresAt updates the expiration timestamp. Pass nil for "no expiry".
func (a *APIKey) SetExpiresAt(t *time.Time) {
	a.ExpiresAt = t
}

// Revoke marks the key as revoked at the given time. Idempotent.
func (a *APIKey) Revoke(at time.Time) {
	if a.RevokedAt != nil {
		return
	}
	a.RevokedAt = &at
}
