package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/authz"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/event"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/interfaces/http/controller"
	"github.com/airoa-org/yubi-app/backend/internal/requestctx"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

const (
	episodeHeartbeatInterval = 30 * time.Second
)

type SSEHandler struct {
	robotDeviceUsecase     usecase.RobotDeviceUsecase
	episodeUsecase         usecase.EpisodeUsecase
	taskUsecase            usecase.TaskUsecase
	taskVersionUsecase     usecase.TaskVersionUsecase
	bus                    *event.Bus
	listBus                *event.Bus
	robotStatusBroadcaster *event.Broadcaster[[]byte]
	teleopBroadcaster      *teleopBroadcaster
	logger                 zerolog.Logger
}

func NewSSEHandler(
	ctx context.Context,
	logger zerolog.Logger,
	robotDeviceUsecase usecase.RobotDeviceUsecase,
	episodeUsecase usecase.EpisodeUsecase,
	taskUsecase usecase.TaskUsecase,
	taskVersionUsecase usecase.TaskVersionUsecase,
	bus *event.Bus,
	robotBus *event.Bus,
	listBus *event.Bus,
	robotStatusBus *event.Bus,
) *SSEHandler {
	h := &SSEHandler{
		robotDeviceUsecase: robotDeviceUsecase,
		episodeUsecase:     episodeUsecase,
		taskUsecase:        taskUsecase,
		taskVersionUsecase: taskVersionUsecase,
		bus:                bus,
		listBus:            listBus,
		logger:             logger,
	}
	h.robotStatusBroadcaster = event.NewBroadcaster(ctx, robotStatusBus, h.fetchRobotStatusFrame)
	h.teleopBroadcaster = newTeleopBroadcaster(
		ctx,
		logger,
		robotStatusBus,
		robotBus,
		robotDeviceUsecase,
		episodeUsecase,
		taskUsecase,
		taskVersionUsecase,
	)
	return h
}

func (h *SSEHandler) StreamRobotStatus(c *gin.Context) {
	role, err := requestctx.UserRole(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusForbidden, openapi.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "user role not found in context",
		})
		return
	}
	if !authz.HasPermission(role, "robot:status_stream") {
		c.JSON(http.StatusForbidden, openapi.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "insufficient permissions",
		})
		return
	}

	robotID := c.Param("robotId")
	if robotID == "" {
		c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "robotId is required",
		})
		return
	}

	exists, err := h.robotDeviceUsecase.RobotExists(c.Request.Context(), robotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to check robot",
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "robot not found",
		})
		return
	}

	ch := h.robotStatusBroadcaster.Subscribe(robotID)
	defer h.robotStatusBroadcaster.Unsubscribe(robotID, ch)

	initial := func() error {
		frame, err := h.fetchRobotStatusFrame(c.Request.Context(), robotID)
		if err != nil {
			return writeSSEError(c, "failed to get robot status", err)
		}
		if _, err := c.Writer.Write(frame); err != nil {
			return err
		}
		c.Writer.Flush()
		return nil
	}
	writeFrame := func(frame []byte) error {
		if _, err := c.Writer.Write(frame); err != nil {
			return err
		}
		c.Writer.Flush()
		return nil
	}
	streamSSE(c, ch, initial, writeFrame)
}

func (h *SSEHandler) StreamRobotStatusByIds(c *gin.Context) {
	ctx := c.Request.Context()

	role, err := requestctx.UserRole(ctx)
	if err != nil {
		c.JSON(http.StatusForbidden, openapi.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "user role not found in context",
		})
		return
	}

	if !authz.HasPermission(role, "robot:status_stream") {
		c.JSON(http.StatusForbidden, openapi.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "insufficient permissions",
		})
		return
	}

	robotIDsParam := c.Query("robotIds")
	if robotIDsParam == "" {
		c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "robotIds is required",
		})
		return
	}

	robotIDs := strings.Split(robotIDsParam, ",")

	seen := make(map[string]struct{})
	uniqueRobotIDs := make([]string, 0, len(robotIDs))

	for _, id := range robotIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}

		if _, ok := seen[id]; ok {
			continue
		}

		seen[id] = struct{}{}
		uniqueRobotIDs = append(uniqueRobotIDs, id)
	}

	const maxRobotIDs = 40
	if len(uniqueRobotIDs) > maxRobotIDs {
		c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("too many robotIds: max %d", maxRobotIDs),
		})
		return
	}

	validIDs := make([]string, 0, len(uniqueRobotIDs))

	for _, id := range uniqueRobotIDs {
		exists, err := h.robotDeviceUsecase.RobotExists(ctx, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "failed to check robot existence",
			})
			return
		}
		if exists {
			validIDs = append(validIDs, id)
		} else {
			h.logger.Warn().
				Str("robot_id", id).
				Msg("requested robot id does not exist, skipping")
		}
	}

	if len(validIDs) == 0 {
		c.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "robot not found",
		})
		return
	}

	var wg sync.WaitGroup
	var once sync.Once

	// merged channel
	merged := make(chan []byte, 10)

	// ctx is stopped, then merged close
	go func() {
		<-ctx.Done()

		wg.Wait()

		once.Do(func() {
			close(merged)
		})
	}()

	for _, id := range validIDs {
		id := id

		wg.Add(1)

		ch := h.robotStatusBroadcaster.Subscribe(id)

		go func(id string, ch chan []byte) {
			defer wg.Done()
			defer h.robotStatusBroadcaster.Unsubscribe(id, ch)

			defer func() {
				if r := recover(); r != nil {
					h.logger.Error().
						Interface("panic", r).
						Str("robot_id", id).
						Bytes("stack", debug.Stack()).
						Msg("panic recovered in robot status stream")
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return

				case msg, ok := <-ch:
					if !ok {
						return
					}

					select {
					case merged <- msg:
					case <-ctx.Done():
						return
					}
				}
			}
		}(id, ch)
	}

	initial := func() error {
		for _, id := range validIDs {
			frame, err := h.fetchRobotStatusFrame(ctx, id)
			if err != nil {
				return writeSSEError(c, "failed to get robot status", err)
			}
			if _, err := c.Writer.Write(frame); err != nil {
				return err
			}
		}
		c.Writer.Flush()
		return nil
	}

	writeFrame := func(frame []byte) error {
		if _, err := c.Writer.Write(frame); err != nil {
			return err
		}
		c.Writer.Flush()
		return nil
	}

	streamSSE(c, merged, initial, writeFrame)
}

func (h *SSEHandler) fetchRobotStatusFrame(ctx context.Context, robotID string) ([]byte, error) {
	status, err := h.robotDeviceUsecase.GetRobotStatus(ctx, robotID)
	if err != nil {
		return nil, err
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
			RobotId:   robotID,
			RobotType: status.RobotType,
			Status:    detail,
		}
	} else {
		response = openapi.RobotStatusStreamResponse{
			RobotId:   robotID,
			RobotType: "",
			Status: openapi.RobotStatusStreamDetail{
				BatteryPct:    0,
				ConnectionPct: 0,
				UptimeSec:     0,
			},
		}
	}

	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return fmt.Appendf(nil, "data: %s\n\n", data), nil
}

// streamSSE sets up SSE headers and runs the common event loop: send an
// initial snapshot, then select on ctx.Done / ch / heartbeat.
func streamSSE[T any](c *gin.Context, ch <-chan T, initial func() error, onEvent func(T) error) {
	rc := http.NewResponseController(c.Writer)
	_ = rc.SetWriteDeadline(time.Time{})

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.Flush()

	heartbeat := time.NewTicker(episodeHeartbeatInterval)
	defer heartbeat.Stop()

	if err := initial(); err != nil {
		return
	}

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case v := <-ch:
			if err := onEvent(v); err != nil {
				return
			}
		case <-heartbeat.C:
			if _, err := fmt.Fprintf(c.Writer, ": heartbeat\n\n"); err != nil {
				return
			}
			c.Writer.Flush()
		}
	}
}

// StreamEpisodeListUpdates sends a lightweight SSE ping whenever any episode
// is created or mutated, allowing the frontend to invalidate its list cache.
func (h *SSEHandler) StreamEpisodeListUpdates(c *gin.Context) {
	role, err := requestctx.UserRole(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusForbidden, openapi.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "user role not found in context",
		})
		return
	}
	if !authz.HasPermission(role, "episode:list_stream") {
		c.JSON(http.StatusForbidden, openapi.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "insufficient permissions",
		})
		return
	}

	ch := h.listBus.Subscribe("list")
	defer h.listBus.Unsubscribe("list", ch)

	ping := func() error {
		if _, err := fmt.Fprint(c.Writer, "data: {}\n\n"); err != nil {
			return err
		}
		c.Writer.Flush()
		return nil
	}
	streamSSE(c, ch, ping, func(struct{}) error { return ping() })
}

// StreamEpisodeUpdates streams episode state via SSE, pushing updates
// whenever the event bus signals a change.
func (h *SSEHandler) StreamEpisodeUpdates(c *gin.Context) {
	role, err := requestctx.UserRole(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusForbidden, openapi.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "user role not found in context",
		})
		return
	}
	if !authz.HasPermission(role, "episode:stream") {
		c.JSON(http.StatusForbidden, openapi.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "insufficient permissions",
		})
		return
	}

	episodeID := c.Param("episodeId")
	if episodeID == "" {
		c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "episodeId is required",
		})
		return
	}

	initialEp, err := h.episodeUsecase.GetByID(c.Request.Context(), episodeID)
	if err != nil {
		if apperror.SameKind(err, apperror.KindNotFound) {
			c.JSON(http.StatusNotFound, openapi.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "episode not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "failed to get episode",
			})
		}
		return
	}

	// Note: task name/description are cached for the connection's lifetime.
	// Renames mid-connection won't show until the client reconnects.
	var task *model.Task
	if tk, err := h.taskUsecase.GetByID(c.Request.Context(), initialEp.TaskID); err != nil {
		h.logger.Error().Err(err).Str("task_id", initialEp.TaskID).Msg("failed to fetch task for episode stream")
	} else {
		task = &tk
	}

	ch := h.bus.Subscribe(episodeID)
	defer h.bus.Unsubscribe(episodeID, ch)

	ctx := c.Request.Context()
	send := func() error { return h.sendEpisodeSSEEvent(ctx, c, episodeID, task) }
	streamSSE(c, ch, send, func(struct{}) error { return send() })
}

func (h *SSEHandler) sendEpisodeSSEEvent(ctx context.Context, c *gin.Context, episodeID string, task *model.Task) error {
	ep, err := h.episodeUsecase.GetByID(ctx, episodeID)
	if err != nil {
		return writeSSEError(c, "failed to get episode", err)
	}

	return h.writeEpisodeSSE(ctx, c, &ep, task)
}

// resolveTaskVersionDisplayName returns the resolved display name for the
// given episode's task version using the same fallback rule as the REST
// controller: explicit display_name, otherwise "{task name} {version}".
func (h *SSEHandler) resolveTaskVersionDisplayName(ctx context.Context, taskID, taskVersionID string) (string, bool) {
	if taskID == "" || taskVersionID == "" {
		return "", false
	}
	tk, err := h.taskUsecase.GetByID(ctx, taskID)
	if err != nil {
		h.logger.Error().Err(err).Str("task_id", taskID).Msg("failed to fetch task for SSE display name")
		return "", false
	}
	tv, err := h.taskVersionUsecase.GetByID(ctx, taskVersionID)
	if err != nil {
		h.logger.Error().Err(err).Str("task_version_id", taskVersionID).Msg("failed to fetch task version for SSE display name")
		return "", false
	}
	return tv.DisplayLabel(tk.Name), true
}

// writeEpisodeSSE fetches subtasks for the given episode, builds the response,
// and writes it as an SSE data frame.
func (h *SSEHandler) writeEpisodeSSE(ctx context.Context, c *gin.Context, ep *model.Episode, task *model.Task) error {
	subtaskMasters, records, executions, err := h.episodeUsecase.GetSubTasksByEpisodeID(ctx, ep.IDNatural, ep.TaskVersionID)
	if err != nil {
		return writeSSEError(c, "failed to get episode subtasks", err)
	}

	subtasks := controller.BuildEpisodeSubTasks(subtaskMasters, records, executions, ep.ParameterValues)

	var taskName, taskDescription *string
	if task != nil {
		taskName = &task.Name
		taskDescription = task.Description
	}

	response := openapi.Episode{
		Id:              ep.IDNatural,
		LocationId:      ep.LocationID,
		UserId:          ep.UserID,
		RobotId:         ep.RobotID,
		Status:          openapi.EpisodeCollectionStatus(ep.Status),
		TaskId:          ep.TaskID,
		TaskName:        taskName,
		TaskDescription: taskDescription,
		TaskVersionId:   ep.TaskVersionID,
		StartedAt:       ep.StartedAt,
		EndedAt:         ep.FinishedAt,
		ErrorDetails:    ep.ErrorDetails,
		Subtasks:        &subtasks,
		CreatedAt:       ep.CreatedAt,
		RecordedBy:      ep.RecordedByID,
		AverageGrade:    ep.AverageGrade,
		GradeCount:      &ep.GradeCount,
	}
	if len(ep.ParameterValues) > 0 {
		response.ParameterValues = &ep.ParameterValues
	}
	if resolved, ok := h.resolveTaskVersionDisplayName(ctx, ep.TaskID, ep.TaskVersionID); ok {
		response.TaskVersionDisplayName = &resolved
	}

	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", data); err != nil {
		return err
	}
	c.Writer.Flush()
	return nil
}

func convertGateToOpenAPI(g *model.GateConditionStatus) openapi.GateConditionStatus {
	groups := make(map[string]openapi.GateGroupStatus, len(g.Groups))
	for name, grp := range g.Groups {
		conditions := make([]openapi.GateCondition, len(grp.Conditions))
		for i, c := range grp.Conditions {
			conditions[i] = openapi.GateCondition{
				Name:       c.Name,
				Passed:     c.Passed,
				Reason:     c.Reason,
				Escalation: c.Escalation,
			}
		}
		groups[name] = openapi.GateGroupStatus{
			Level:      grp.Level,
			Settled:    grp.Settled,
			Conditions: conditions,
		}
	}
	return openapi.GateConditionStatus{
		GateLevel: g.GateLevel,
		Groups:    groups,
	}
}

// StreamRobotTeleop streams a consolidated teleop view for a robot over a
// single SSE connection, emitting three named event types:
//
//	event: robot_status — battery/connection/gates (per heartbeat)
//	event: episode      — current episode + subtasks
//	event: task         — task name/description/manual_url/version,
//	                      only on connect and when the task_id changes
//
// Per-robot fan-out: a single goroutine drives every subscriber watching
// the same robot, matching the pattern of StreamRobotStatus.
func (h *SSEHandler) StreamRobotTeleop(c *gin.Context) {
	role, err := requestctx.UserRole(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusForbidden, openapi.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "user role not found in context",
		})
		return
	}
	if !authz.HasPermission(role, "robot:teleop_stream") {
		c.JSON(http.StatusForbidden, openapi.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "insufficient permissions",
		})
		return
	}

	robotID := c.Param("robotId")
	if robotID == "" {
		c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "robotId is required",
		})
		return
	}

	exists, err := h.robotDeviceUsecase.RobotExists(c.Request.Context(), robotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to check robot",
		})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "robot not found",
		})
		return
	}

	ch := h.teleopBroadcaster.Subscribe(robotID)
	defer h.teleopBroadcaster.Unsubscribe(robotID, ch)

	// The broadcaster drives the initial snapshot on its own goroutine
	// (first subscriber); late joiners are served from the room's
	// per-event cache inside teleopBroadcaster.Subscribe, so streamSSE
	// doesn't need an initial() callback here.
	noopInitial := func() error { return nil }
	writeFrame := func(frame []byte) error {
		if _, err := c.Writer.Write(frame); err != nil {
			return err
		}
		c.Writer.Flush()
		return nil
	}
	streamSSE(c, ch, noopInitial, writeFrame)
}

// writeSSEError writes an error event to the SSE stream.
func writeSSEError(c *gin.Context, message string, cause error) error {
	data, _ := json.Marshal(openapi.ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
	})
	if _, writeErr := fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", data); writeErr != nil {
		return writeErr
	}
	c.Writer.Flush()
	return cause
}
