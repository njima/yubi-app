package repository

import "context"

type UserLocation interface {
	SetUserLocations(ctx context.Context, conn Conn, userID string, organizationID string, locationIDs []string) error
}
