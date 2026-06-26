package persistence

import (
	"github.com/airoa-org/yubi-app/backend/internal/repository"

	"github.com/uptrace/bun"
)

type bunConnection interface {
	NewInsert() *bun.InsertQuery
	NewSelect() *bun.SelectQuery
	NewRaw(query string, args ...any) *bun.RawQuery
	NewUpdate() *bun.UpdateQuery
	NewDelete() *bun.DeleteQuery
}

func bunConn(conn repository.Conn) bunConnection {
	c, ok := conn.(bunConnection)
	if !ok {
		panic("persistence requires a bun-compatible database connection")
	}
	return c
}
