package main

import "github.com/airoa-org/yubi-app/backend/internal/event"

type eventBuses struct {
	RobotStatus  *event.Bus
	Episode      *event.Bus
	RobotEpisode *event.Bus
	EpisodeList  *event.Bus
}

func newEventBuses() eventBuses {
	return eventBuses{
		RobotStatus:  event.NewBus(),
		Episode:      event.NewBus(),
		RobotEpisode: event.NewBus(),
		EpisodeList:  event.NewBus(),
	}
}
