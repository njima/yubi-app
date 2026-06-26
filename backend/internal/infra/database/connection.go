package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// NewDatabase connects to PostgreSQL and returns a bun.DB.
// Only the "postgres" driver is supported.
//
// If sslMode is empty, it defaults to "disable".
// The connection pool is initialized with the following settings:
//   - MaxIdleConns: 10
//   - MaxOpenConns: 50
//   - ConnMaxLifetime: 5 minutes
func NewDatabase(
	driver,
	user,
	pass,
	host,
	port,
	dbName,
	sslMode string,
) (*bun.DB, error) {
	if driver != "postgres" {
		return nil, newError(ErrorKindConnect, fmt.Errorf("unsupported driver: %s, only postgres is supported", driver))
	}

	if sslMode == "" {
		sslMode = "disable"
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, dbName, sslMode)

	baseDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, newError(ErrorKindConnect, err)
	}

	if err := baseDB.Ping(); err != nil {
		return nil, newError(ErrorKindPing, err)
	}

	baseDB.SetMaxIdleConns(10)
	baseDB.SetMaxOpenConns(50)
	baseDB.SetConnMaxLifetime(5 * time.Minute)

	db := bun.NewDB(baseDB, pgdialect.New())

	return db, nil
}
