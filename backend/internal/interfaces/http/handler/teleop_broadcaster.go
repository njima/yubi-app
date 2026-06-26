package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sync"

	"github.com/rs/zerolog"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/interfaces/http/controller"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/eventbus"
)

// teleopBroadcaster multiplexes robotStatusBus and robotEpisodeBus into a
// single per-robot goroutine that emits SSE *named* events to every
// connection watching that robot. Event types:
//
//	event: robot_status — live robot status (battery/connection/gates)
//	event: episode      — current episode + subtasks (full snapshot)
//	event: task         — task metadata; only on connect and on task change
//
// The per-robot goroutine caches the last emitted task_id so `task` frames
// are only sent when the current episode's task actually changes.
type teleopBroadcaster struct {
	parent             context.Context
	logger             zerolog.Logger
	robotStatusBus     *eventbus.Bus
	robotEpisodeBus    *eventbus.Bus
	robotDeviceUsecase usecase.RobotDeviceUsecase
	episodeUsecase     usecase.EpisodeUsecase
	taskUsecase        usecase.TaskUsecase
	taskVersionUsecase usecase.TaskVersionUsecase

	mu    sync.Mutex
	rooms map[string]*teleopRoom
}

type teleopRoom struct {
	cancel context.CancelFunc
	subs   map[chan []byte]struct{}

	// Latest fan-out of each event type, cached so a connection that
	// joins an already-running room gets an immediate snapshot without
	// waiting for the next bus tick.
	lastStatus  []byte
	lastEpisode []byte
	lastTask    []byte
}

const teleopSubBuffer = 16

func newTeleopBroadcaster(
	ctx context.Context,
	logger zerolog.Logger,
	robotStatusBus, robotEpisodeBus *eventbus.Bus,
	robotDeviceUsecase usecase.RobotDeviceUsecase,
	episodeUsecase usecase.EpisodeUsecase,
	taskUsecase usecase.TaskUsecase,
	taskVersionUsecase usecase.TaskVersionUsecase,
) *teleopBroadcaster {
	return &teleopBroadcaster{
		parent:             ctx,
		logger:             logger,
		robotStatusBus:     robotStatusBus,
		robotEpisodeBus:    robotEpisodeBus,
		robotDeviceUsecase: robotDeviceUsecase,
		episodeUsecase:     episodeUsecase,
		taskUsecase:        taskUsecase,
		taskVersionUsecase: taskVersionUsecase,
		rooms:              make(map[string]*teleopRoom),
	}
}

// Subscribe registers a new connection for a robot and returns a channel
// on which SSE frames (already-formatted `event:`/`data:` bytes) will be
// delivered. The first subscriber for a given robot starts a goroutine;
// the last Unsubscribe cancels it.
//
// If the room is already running, the new subscriber receives whichever
// of (status, episode, task) frames have been emitted so far, so it
// sees the current state without waiting for the next bus tick.
func (tb *teleopBroadcaster) Subscribe(robotID string) chan []byte {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	room, ok := tb.rooms[robotID]
	if !ok {
		ctx, cancel := context.WithCancel(tb.parent)
		room = &teleopRoom{
			cancel: cancel,
			subs:   make(map[chan []byte]struct{}),
		}
		tb.rooms[robotID] = room
		go tb.loop(ctx, robotID)
	}
	ch := make(chan []byte, teleopSubBuffer)
	room.subs[ch] = struct{}{}

	// Replay cached snapshot for late joiners. Non-blocking: the buffer
	// has space (we just created the channel).
	if room.lastStatus != nil {
		ch <- room.lastStatus
	}
	if room.lastEpisode != nil {
		ch <- room.lastEpisode
	}
	if room.lastTask != nil {
		ch <- room.lastTask
	}
	return ch
}

func (tb *teleopBroadcaster) Unsubscribe(robotID string, ch chan []byte) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	room, ok := tb.rooms[robotID]
	if !ok {
		return
	}
	delete(room.subs, ch)
	close(ch)
	if len(room.subs) == 0 {
		room.cancel()
		delete(tb.rooms, robotID)
	}
}

// teleopEventKind identifies which cache slot (and event type) a frame
// belongs to. The string value is for symmetry with the SSE event names
// emitted on the wire; only the enum is used internally.
type teleopEventKind int

const (
	teleopEventStatus teleopEventKind = iota
	teleopEventEpisode
	teleopEventTask
)

// fanOut updates the room's cache for the given event kind and delivers
// the frame to every current subscriber. Slow subscribers drop the
// frame rather than block the per-robot loop.
func (tb *teleopBroadcaster) fanOut(robotID string, kind teleopEventKind, frame []byte) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	room, ok := tb.rooms[robotID]
	if !ok {
		return
	}
	switch kind {
	case teleopEventStatus:
		room.lastStatus = frame
	case teleopEventEpisode:
		room.lastEpisode = frame
	case teleopEventTask:
		room.lastTask = frame
	}
	for ch := range room.subs {
		select {
		case ch <- frame:
		default:
			// Slow subscriber — drop this frame rather than block the
			// per-robot loop. Mirrors the behaviour of the existing
			// eventbus.Broadcaster fan-out.
		}
	}
}

// loop is the per-robot goroutine. It subscribes to both input buses,
// pushes an initial snapshot (status + episode + task) and then reacts
// to wakes on either bus, emitting only the event type that matches
// the wake source.
func (tb *teleopBroadcaster) loop(ctx context.Context, robotID string) {
	statusCh := tb.robotStatusBus.Subscribe(robotID)
	defer tb.robotStatusBus.Unsubscribe(robotID, statusCh)
	episodeCh := tb.robotEpisodeBus.Subscribe(robotID)
	defer tb.robotEpisodeBus.Unsubscribe(robotID, episodeCh)

	var lastTaskID string

	tb.pushStatus(ctx, robotID)
	lastTaskID = tb.pushEpisodeAndMaybeTask(ctx, robotID, lastTaskID)

	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-statusCh:
			if !ok {
				return
			}
			tb.pushStatus(ctx, robotID)
		case _, ok := <-episodeCh:
			if !ok {
				return
			}
			lastTaskID = tb.pushEpisodeAndMaybeTask(ctx, robotID, lastTaskID)
		}
	}
}

// pushStatus fetches the robot's live status and fans out an
// `event: robot_status` frame. Errors are logged at WARN (so persistent
// backend failures are observable) but don't tear down the loop — the
// next bus tick will retry.
func (tb *teleopBroadcaster) pushStatus(ctx context.Context, robotID string) {
	status, err := tb.robotDeviceUsecase.GetRobotStatus(ctx, robotID)
	if err != nil {
		tb.logger.Warn().Err(err).Str("robot_id", robotID).Msg("teleop: pushStatus GetRobotStatus failed")
		return
	}

	var response openapi.RobotStatusStreamResponse
	if status != nil {
		detail := openapi.RobotStatusStreamDetail{
			BatteryPct:    status.Status.Battery.Pct,
			ConnectionPct: status.Status.Connection.QualityPct,
			UptimeSec:     int(math.Round(status.Status.UptimeSec)),
		}
		if g := status.Status.GateConditions; g != nil {
			oGate := convertGateToOpenAPI(g)
			detail.GateConditions = &oGate
		}
		response = openapi.RobotStatusStreamResponse{
			RobotType: status.RobotType,
			Status:    detail,
		}
	}

	data, err := json.Marshal(response)
	if err != nil {
		tb.logger.Warn().Err(err).Str("robot_id", robotID).Msg("teleop: pushStatus json.Marshal failed")
		return
	}
	tb.fanOut(robotID, teleopEventStatus, fmt.Appendf(nil, "event: robot_status\ndata: %s\n\n", data))
}

// pushEpisodeAndMaybeTask fetches the current episode for the robot and
// fans out an `event: episode` frame. If the episode's task_id differs
// from the previously-emitted one, it also fetches + emits a fresh
// `event: task` frame. Returns the task_id that should be remembered for
// the next call (unchanged if no task event was emitted).
func (tb *teleopBroadcaster) pushEpisodeAndMaybeTask(ctx context.Context, robotID, prevTaskID string) string {
	ep, err := tb.episodeUsecase.GetCurrentRobotEpisode(ctx, robotID)
	if err != nil {
		tb.logger.Warn().Err(err).Str("robot_id", robotID).Msg("teleop: GetCurrentRobotEpisode failed")
		return prevTaskID
	}
	if ep == nil {
		tb.fanOut(robotID, teleopEventEpisode, []byte("event: episode\ndata: null\n\n"))
		return ""
	}

	subtaskMasters, records, executions, err := tb.episodeUsecase.GetSubTasksByEpisodeID(ctx, ep.IDNatural, ep.TaskVersionID)
	if err != nil {
		tb.logger.Warn().Err(err).
			Str("robot_id", robotID).
			Str("episode_id", ep.IDNatural).
			Msg("teleop: GetSubTasksByEpisodeID failed")
		return prevTaskID
	}
	subtasks := controller.BuildEpisodeSubTasks(subtaskMasters, records, executions, ep.ParameterValues)

	episodeResp := openapi.Episode{
		Id:            ep.IDNatural,
		LocationId:    ep.LocationID,
		UserId:        ep.UserID,
		RobotId:       ep.RobotID,
		Status:        openapi.EpisodeCollectionStatus(ep.Status),
		TaskId:        ep.TaskID,
		TaskVersionId: ep.TaskVersionID,
		StartedAt:     ep.StartedAt,
		EndedAt:       ep.FinishedAt,
		ErrorDetails:  ep.ErrorDetails,
		Subtasks:      &subtasks,
		CreatedAt:     ep.CreatedAt,
		RecordedBy:    ep.RecordedByID,
		AverageGrade:  ep.AverageGrade,
		GradeCount:    &ep.GradeCount,
	}
	if len(ep.ParameterValues) > 0 {
		episodeResp.ParameterValues = &ep.ParameterValues
	}

	var tk *model.Task
	var tv *model.TaskVersion
	if ep.TaskID != "" && ep.TaskVersionID != "" {
		if t, err := tb.taskUsecase.GetByID(ctx, ep.TaskID); err == nil {
			tk = &t
		} else {
			tb.logger.Warn().Err(err).
				Str("robot_id", robotID).
				Str("task_id", ep.TaskID).
				Msg("teleop: pushTask GetByID(task) failed")
		}
		if v, err := tb.taskVersionUsecase.GetByID(ctx, ep.TaskVersionID); err == nil {
			tv = &v
		} else {
			tb.logger.Warn().Err(err).
				Str("robot_id", robotID).
				Str("task_version_id", ep.TaskVersionID).
				Msg("teleop: pushTask GetByID(task_version) failed")
		}
	}
	if tk != nil && tv != nil {
		resolved := tv.DisplayLabel(tk.Name)
		episodeResp.TaskVersionDisplayName = &resolved
	}

	data, err := json.Marshal(episodeResp)
	if err != nil {
		tb.logger.Warn().Err(err).
			Str("robot_id", robotID).
			Str("episode_id", ep.IDNatural).
			Msg("teleop: episode json.Marshal failed")
		return prevTaskID
	}
	tb.fanOut(robotID, teleopEventEpisode, fmt.Appendf(nil, "event: episode\ndata: %s\n\n", data))

	if ep.TaskID == prevTaskID || tk == nil || tv == nil {
		return prevTaskID
	}
	if tb.pushTaskFrame(robotID, tk, tv) {
		return ep.TaskID
	}
	return prevTaskID
}

// TeleopTaskMeta is the payload shape for `event: task`. Kept as a
// private type here because it is not exposed through openapi — only the
// teleop stream emits it and only useTeleopStream consumes it.
type teleopTaskMeta struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	ManualURL   string  `json:"manual_url"`
	Version     string  `json:"version"`
}

// pushTaskFrame fans out an `event: task` frame from already-fetched task
// and task_version objects. Returns true if the frame was successfully
// emitted so the caller can update its cached task_id.
func (tb *teleopBroadcaster) pushTaskFrame(robotID string, task *model.Task, version *model.TaskVersion) bool {
	meta := teleopTaskMeta{
		Id:          task.IDNatural,
		Name:        task.Name,
		Description: task.Description,
		ManualURL:   task.ManualURL,
		Version:     version.Version,
	}
	data, err := json.Marshal(meta)
	if err != nil {
		tb.logger.Warn().Err(err).
			Str("robot_id", robotID).
			Str("task_id", task.IDNatural).
			Msg("teleop: pushTaskFrame json.Marshal failed")
		return false
	}
	tb.fanOut(robotID, teleopEventTask, fmt.Appendf(nil, "event: task\ndata: %s\n\n", data))
	return true
}
