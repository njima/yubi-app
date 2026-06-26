package repository

import (
	"context"
)

type DBConn interface{}

type TxRunner interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context, conn DBConn) error) error
}

type DataAccess struct {
	conn DBConn
	tx   TxRunner
}

func NewDataAccess(conn DBConn, tx TxRunner) DataAccess {
	return DataAccess{conn: conn, tx: tx}
}

func (d DataAccess) Conn() DBConn {
	return d.conn
}

func (d DataAccess) RunInTx(ctx context.Context, fn func(ctx context.Context, data DataAccess) error) error {
	return d.tx.RunInTx(ctx, func(ctx context.Context, conn DBConn) error {
		return fn(ctx, NewDataAccess(conn, d.tx))
	})
}
