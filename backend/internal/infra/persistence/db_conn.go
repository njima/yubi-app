package persistence

import (
	"github.com/airoa-org/yubi-app/backend/internal/repository"

	"github.com/uptrace/bun"
)

type bunDBConn interface {
	NewInsert() *bun.InsertQuery
	NewSelect() *bun.SelectQuery
	NewRaw(query string, args ...any) *bun.RawQuery
	NewUpdate() *bun.UpdateQuery
	NewDelete() *bun.DeleteQuery
}

func bunConn(conn repository.DBConn) bunDBConn {
	c, ok := conn.(bunDBConn)
	if !ok {
		panic("persistence requires a bun-compatible database connection")
	}
	return c
}
