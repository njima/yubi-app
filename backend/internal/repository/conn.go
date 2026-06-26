package repository

import (
	"context"
)

type Conn interface{}

type TransactionRunner interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context, conn Conn) error) error
}

type DataAccess struct {
	conn Conn
	tx   TransactionRunner
}

func NewDataAccess(conn Conn, tx TransactionRunner) DataAccess {
	return DataAccess{conn: conn, tx: tx}
}

func (d DataAccess) Conn() Conn {
	return d.conn
}

func (d DataAccess) RunInTx(ctx context.Context, fn func(ctx context.Context, data DataAccess) error) error {
	return d.tx.RunInTx(ctx, func(ctx context.Context, conn Conn) error {
		return fn(ctx, NewDataAccess(conn, d.tx))
	})
}
