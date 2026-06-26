package entity

import (
	"time"

	"github.com/uptrace/bun"
)

// Timestamp is an embeddable struct that provides the created_at and updated_at columns
// shared by all entities. Timestamps are set automatically via bun's BeforeAppendModel hook
// on INSERT and UPDATE.
type Timestamp struct {
	CreatedAt time.Time `bun:"created_at,type:timestamptz,notnull"`
	UpdatedAt time.Time `bun:"updated_at,type:timestamptz,notnull"`
}

// TableQueryCreator is a function type that accepts a bun.DB and returns a CREATE TABLE query.
// Instances are registered in the TableCreators slice and executed during migration.
type TableQueryCreator func(db *bun.DB) *bun.CreateTableQuery

// IndexQueryCreator is a function type that accepts a bun.DB and returns a CREATE INDEX query.
// Instances are registered in the IdxCreators slice and executed during migration.
type IndexQueryCreator func(db *bun.DB) *bun.CreateIndexQuery
