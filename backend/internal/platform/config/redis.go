package config

import (
	"context"

	"github.com/kelseyhightower/envconfig"
)

type Redis struct {
	Host     string `envconfig:"REDIS_HOST" default:"localhost"`
	Port     string `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}

func NewRedis(ctx context.Context) (*Redis, error) {
	r := Redis{}
	if err := envconfig.Process("", &r); err != nil {
		return nil, newError(ErrorKindProcess, err)
	}

	return &r, nil
}
