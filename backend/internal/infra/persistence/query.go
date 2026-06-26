package persistence

import (
	"strings"

	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

// escapeILIKE escapes ILIKE wildcards (% and _) so user-supplied search terms
// don't act as pattern matchers.
func escapeILIKE(s string) string {
	return strings.NewReplacer("%", "\\%", "_", "\\_").Replace(s)
}

func applyLocationListFilters(sel *bun.SelectQuery, filter repository.LocationListFilter) *bun.SelectQuery {
	if filter.SiteID != nil && *filter.SiteID != "" {
		sel = sel.Where("l.site_id = ?", *filter.SiteID)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		sel = sel.Where("l.name ILIKE ?", "%"+escaped+"%")
	}
	return sel
}

func applySiteListFilters(sel *bun.SelectQuery, filter repository.SiteListFilter) *bun.SelectQuery {
	if filter.OrganizationID != nil {
		sel = sel.Where("organization_id = ?", *filter.OrganizationID)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		sel = sel.Where("name ILIKE ?", "%"+escaped+"%")
	}
	return sel
}

// namedSQLArgs lets bun's NewRaw resolve ?name placeholders against a flat
// map. bun only consults a NamedArgAppender when args has length 1, so the
// whole map must be passed as the sole argument:
//
//	conn.NewRaw("SELECT ... WHERE id = ?id", namedSQLArgs{"id": x})
//
// Prefer over positional ? when a query has many parameters or reuses the
// same value multiple times.
type namedSQLArgs map[string]any

var _ schema.NamedArgAppender = namedSQLArgs(nil)

func (a namedSQLArgs) AppendNamedArg(gen schema.QueryGen, b []byte, name string) ([]byte, bool) {
	v, ok := a[name]
	if !ok {
		return b, false
	}
	return gen.Append(b, v), true
}
