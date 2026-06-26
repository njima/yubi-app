package app

import (
	"github.com/airoa-org/yubi-app/backend/internal/config"
	"github.com/airoa-org/yubi-app/backend/internal/event"
	"github.com/airoa-org/yubi-app/backend/internal/infra/cache"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type application struct {
	conf   *config.Config
	logger zerolog.Logger

	db          *bun.DB
	redisClient *cache.Client

	datadogStarted bool
	sentryStarted  bool

	usecases

	robotStatusBus  *event.Bus
	episodeBus      *event.Bus
	robotEpisodeBus *event.Bus
	episodeListBus  *event.Bus
}
