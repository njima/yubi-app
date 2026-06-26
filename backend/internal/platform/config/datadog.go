package config

type Datadog struct {
	Enabled bool `envconfig:"DD_ENABLED" default:"false"`
}
