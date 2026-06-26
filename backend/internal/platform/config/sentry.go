package config

import (
	"context"

	"github.com/kelseyhightower/envconfig"
)

type Sentry struct {
	DSN              string  `envconfig:"SENTRY_DSN" default:""`
	Environment      string  `envconfig:"SENTRY_ENVIRONMENT" default:"local"`
	TracesSampleRate float64 `envconfig:"SENTRY_TRACES_SAMPLE_RATE" default:"0.0"`
}

func NewSentry(ctx context.Context) (*Sentry, error) {
	s := Sentry{}
	if err := envconfig.Process("", &s); err != nil {
		return nil, newError(ErrorKindProcess, err)
	}

	return &s, nil
}
