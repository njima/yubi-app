package main

import "github.com/airoa-org/yubi-app/backend/internal/usecase/eventbus"

type eventBuses struct {
	RobotStatus  *eventbus.Bus
	Episode      *eventbus.Bus
	RobotEpisode *eventbus.Bus
	EpisodeList  *eventbus.Bus
}

func newEventBuses() eventBuses {
	return eventBuses{
		RobotStatus:  eventbus.NewBus(),
		Episode:      eventbus.NewBus(),
		RobotEpisode: eventbus.NewBus(),
		EpisodeList:  eventbus.NewBus(),
	}
}
