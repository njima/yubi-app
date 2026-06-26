package config

import (
	"context"

	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Driver   string `envconfig:"DB_DRIVER" required:"true"`
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     string `envconfig:"DB_PORT" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
	SSLMode  string `envconfig:"DB_SSL_MODE" default:"disable"`
}

func NewDatabase(ctx context.Context) (*Database, error) {
	db := Database{}
	if err := envconfig.Process("", &db); err != nil {
		return nil, newError(ErrorKindProcess, err)
	}

	return &db, nil
}
