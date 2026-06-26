package config

type Auth struct {
	DefaultOrganizationID string `envconfig:"AUTH_DEFAULT_ORGANIZATION_ID"`
	DefaultUserRole       uint   `envconfig:"AUTH_DEFAULT_USER_ROLE" default:"0"` // 0 = Admin
}
