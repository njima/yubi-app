package entity

import (
	"context"

	"github.com/uptrace/bun"
)

// OrgIDFromContext is set at application startup to avoid circular imports
// between the entity and requestctx packages.
// Returns the organization ID and true if present, or empty string and false if not.
var OrgIDFromContext func(ctx context.Context) (string, bool)

// OrgScoped is an embeddable struct that automatically applies
// WHERE <alias>.organization_id = ? to SELECT, UPDATE, and DELETE queries via bun hooks.
//
// Usage: embed this struct in any entity that has an organization_id column.
//
//	type Robot struct {
//	    bun.BaseModel `bun:"table:robot,alias:r"`
//	    OrgScoped
//	    ...
//	}
//
// No-op when org is absent from context (background jobs, auth middleware internal calls).
//
// IMPORTANT: hooks only fire when the entity is used as the primary model via Model().
// Queries built with TableExpr() bypass hooks entirely and MUST add the org filter manually:
//
//	if orgID, err := requestctx.OrganizationID(ctx); err == nil {
//	    q = q.Where("<alias>.organization_id = ?", orgID)
//	}
type OrgScoped struct{}

var _ bun.BeforeSelectHook = OrgScoped{}
var _ bun.BeforeUpdateHook = OrgScoped{}
var _ bun.BeforeDeleteHook = OrgScoped{}

func (OrgScoped) BeforeSelect(ctx context.Context, q *bun.SelectQuery) error {
	if orgID, ok := orgIDFromCtx(ctx); ok {
		if tm, ok := q.GetModel().(bun.TableModel); ok {
			alias := tm.Table().Alias
			if alias != "" {
				q.Where(alias+".organization_id = ?", orgID)
				return nil
			}
		}
		q.Where("organization_id = ?", orgID)
	}
	return nil
}

func (OrgScoped) BeforeUpdate(ctx context.Context, q *bun.UpdateQuery) error {
	if orgID, ok := orgIDFromCtx(ctx); ok {
		// UPDATE operates on a single table — no JOIN ambiguity, alias not needed.
		q.Where("organization_id = ?", orgID)
	}
	return nil
}

func (OrgScoped) BeforeDelete(ctx context.Context, q *bun.DeleteQuery) error {
	if orgID, ok := orgIDFromCtx(ctx); ok {
		// DELETE operates on a single table — no JOIN ambiguity, alias not needed.
		q.Where("organization_id = ?", orgID)
	}
	return nil
}

func orgIDFromCtx(ctx context.Context) (string, bool) {
	if OrgIDFromContext == nil {
		return "", false
	}
	return OrgIDFromContext(ctx)
}
