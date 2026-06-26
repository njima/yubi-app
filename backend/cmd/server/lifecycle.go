package main

import (
	"time"

	"github.com/getsentry/sentry-go"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func (a *application) Close() {
	if a.redisClient != nil {
		if err := a.redisClient.Close(); err != nil {
			a.logger.Error().Err(err).Msg("failed to close redis")
		}
	}
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.logger.Error().Err(err).Msg("failed to close db")
		}
	}
	if a.sentryStarted {
		sentry.Flush(2 * time.Second)
	}
	if a.datadogStarted {
		ddtracer.Stop()
	}
}
