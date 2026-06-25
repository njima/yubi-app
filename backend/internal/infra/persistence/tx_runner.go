package persistence

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/repository"

	"github.com/uptrace/bun"
)

type TxRunner struct {
	db *bun.DB
}

func NewTxRunner(db *bun.DB) *TxRunner {
	return &TxRunner{db: db}
}

func (r *TxRunner) RunInTx(ctx context.Context, fn func(ctx context.Context, conn repository.DBConn) error) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return fn(ctx, tx)
	})
}
