package config

import (
	"context"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Environment string `envconfig:"ENVIRONMENT" default:"local"`
	AppName     string `envconfig:"APP_NAME" default:"yubi-app-backend"`
	AppPort     string `envconfig:"APP_PORT" default:"8000"`
	Database    Database
	Auth        Auth
	Redis       Redis
	Sentry      Sentry
	Datadog     Datadog
	S3          S3
}

func NewConfig(ctx context.Context) (*Config, error) {
	conf := Config{}
	if err := envconfig.Process("", &conf); err != nil {
		return nil, newError(ErrorKindProcess, err)
	}

	return &conf, nil
}

func (c *Config) IsLocal() bool {
	return strings.ToLower(c.Environment) == "local"
}

// var envKeywords = []string{
// 	"local",
// 	"development",
// 	"production",
// }
