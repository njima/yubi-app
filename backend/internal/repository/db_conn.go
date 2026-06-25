package repository

import (
	"context"

	"github.com/uptrace/bun"
)

type DBConn interface {
	NewInsert() *bun.InsertQuery
	NewSelect() *bun.SelectQuery
	NewRaw(query string, args ...any) *bun.RawQuery
	NewUpdate() *bun.UpdateQuery
	NewDelete() *bun.DeleteQuery
}

type TxRunner interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context, conn DBConn) error) error
}
