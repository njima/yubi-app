package entity

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Organization struct {
	bun.BaseModel `bun:"table:organization,alias:o"`
	Timestamp

	// columns
	ID          int64   `bun:"id,pk,autoincrement"`                        // Auto-increment ID
	IDNatural   string  `bun:"id_natural,unique,type:varchar(36),notnull"` // UUID
	Name        string  `bun:"name,unique,type:varchar(255),notnull"`      // Organization name
	Kind        string  `bun:"kind,type:varchar(20),notnull,default:'team'"`
	Description *string `bun:"description,type:text"` // Description
}

var _ bun.BeforeAppendModelHook = (*Organization)(nil)

func (o *Organization) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if o.CreatedAt.IsZero() {
			o.CreatedAt = now
		}

		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = now
		}

	case *bun.UpdateQuery:
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = now
		}
	}

	return nil
}
