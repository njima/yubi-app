// Components
export { RobotListPage } from "./components/robot-list-page";
export { CreateRobotDialog } from "./components/create-robot-dialog";
export { CreateRobotForm } from "./components/create-robot-form";
export { DeleteRobotDialog } from "./components/delete-robot-dialog";
export { EditRobotDialog } from "./components/edit-robot-dialog";
export { EditRobotForm } from "./components/edit-robot-form";
export { RobotStatusBadge } from "./components/robot-status-badge";
export { LeaderStatusBadge } from "./components/leader-status-badge";
export { GateStatusBadgeInRobotList } from "./components/gate-status-badge";
export { ConsecutiveFaultDaysBadge } from "./components/consecutive-fault-days-badge";
export { FleetSummaryGrid } from "./components/fleet-summary-grid";
export { FleetStatsPanel } from "./components/fleet-stats-table";
export { BatteryLevelIndicator } from "./components/battery-level-indicator";
export { RobotDetailPage } from "./components/detail";

// Teleoperation Components
export { StartTeleoperationButton } from "./components/start-teleoperation-button";
export { QueueTaskButton } from "./components/queue-task-button";
export { StartTeleoperationDialog } from "./components/start-teleoperation-dialog";
export { StartTeleoperationForm } from "./components/start-teleoperation-form";
export { EStopButton } from "./components/e-stop-button";
export { SubtaskControlTable } from "./components/subtask-control-table";
export { SubtaskTimeline } from "./components/subtask-timeline";
export { TeleoperationStatusCard } from "./components/teleoperation-status-card";
export { ThreeDModelCard } from "./components/three-d-model-card";

// Status Cards (shared with Teleoperation)
export { RobotStatusCard } from "./components/robot-status-card";
export { ForceTorqueCard } from "./components/force-torque-card";
export { JointTemperaturesCard } from "./components/joint-temperatures-card";

// Camera
export { RobotCameraViewer } from "./components/robot-camera-viewer";
export { MjpegViewer } from "./components/mjpeg-viewer";

// Hooks
export { useRobotsQuery, useRobotQuery } from "./hooks/use-robots-query";
export { useRobotTypesQuery } from "./hooks/use-robot-types-query";
export { useRobotSearchOptions } from "./hooks/use-robot-search-options";
export { useCreateRobotMutation } from "./hooks/use-create-robot-mutation";
export { useUpdateRobotMutation } from "./hooks/use-update-robot-mutation";
export { useDeleteRobotMutation } from "./hooks/use-delete-robot-mutation";
export { useRobotsStatusStream } from "./hooks/use-robots-status-stream";
export { useRobotStatusStream } from "./hooks/use-robot-status-stream";
export {
  useFleetSummaryQuery,
  fleetSummaryQueryKeys,
} from "./hooks/use-fleet-summary-query";
export {
  useFleetStatsQuery,
  useCollectionTrendQuery,
  fleetStatsQueryKeys,
} from "./hooks/use-fleet-stats-query";

// Scope
export { useRobotScope, isRobotInScope } from "./hooks/use-robot-scope";

// Contexts
export {
  RobotsStatusProvider,
  useRobotsStatus,
} from "./contexts/robots-status-context";

// Lib
export { computeFleetTotals, computeStatsTotals } from "./lib/fleet-utils";

// Fleet Schemas
export type {
  FleetSiteSummary,
  FleetRobotTypeSummary,
  FleetStatusCount,
  FleetRobotTypeStats,
  FleetSiteStats,
  CollectionTrend,
  TrendSeries,
  ColoredTrendSeries,
  ColoredCollectionTrend,
  Granularity,
  FleetTotals,
  StatsTotals,
} from "./schemas/fleet";
