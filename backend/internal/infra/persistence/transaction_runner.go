package persistence

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/repository"

	"github.com/uptrace/bun"
)

type TransactionRunner struct {
	db *bun.DB
}

func NewTransactionRunner(db *bun.DB) *TransactionRunner {
	return &TransactionRunner{db: db}
}

func (r *TransactionRunner) RunInTx(ctx context.Context, fn func(ctx context.Context, conn repository.Conn) error) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return fn(ctx, tx)
	})
}
