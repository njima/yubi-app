package main

import (
	"github.com/airoa-org/yubi-app/backend/internal/config"
	"github.com/airoa-org/yubi-app/backend/internal/infra/cache"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/eventbus"

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

	robotStatusBus  *eventbus.Bus
	episodeBus      *eventbus.Bus
	robotEpisodeBus *eventbus.Bus
	episodeListBus  *eventbus.Bus
}
