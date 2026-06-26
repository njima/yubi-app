package repository

import "context"

type UserSite interface {
	SetUserSites(ctx context.Context, conn Conn, userID string, organizationID string, siteIDs []string) error
}
