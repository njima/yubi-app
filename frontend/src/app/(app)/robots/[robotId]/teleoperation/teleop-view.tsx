"use client";

import {
  ArrowLeft,
  ChevronRight,
  Maximize,
  Minimize,
  Play,
  Signal,
  SignalZero,
  Square,
} from "lucide-react";
import Link from "next/link";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";

import type { LayoutContext } from "@/shared/lib/layout-registry";
import { DEFAULT_TELEOP_VIEWS } from "@/shared/lib/layout-types";
import {
  SUBTASK_COLLECTION_STATUS,
  USER_ROLE,
} from "@/shared/lib/status-constants";
import { cn } from "@/shared/lib/utils";

import { LayoutRenderer } from "@/components/layout/layout-renderer";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

import {
  useRobotQuery,
  useRobotScope,
  isRobotInScope,
} from "@/features/robots";
import { useOperatorHeartbeat } from "@/features/robots/hooks/use-operator-heartbeat";
import { useTeleopStream } from "@/features/robots/hooks/use-teleop-stream";
import {
  extractCameras,
  extractHostPort,
  extractRobotLayout,
  extractStreamConfig,
} from "@/features/robots/lib/camera-utils";
import { registerTeleopComponents } from "@/features/robots/lib/register-teleop-components";
import { useMeQuery } from "@/features/users";

registerTeleopComponents();

interface TeleopViewProps {
  robotId: string;
  viewName: string;
}

export function TeleopView({ robotId, viewName }: TeleopViewProps) {
  const { t } = useTranslation();
  const contentRef = useRef<HTMLDivElement>(null);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const isDefault = viewName === "default";

  // Operator lock — explicit start/stop
  const {
    isActive: isOperating,
    isLocked,
    lockedBy,
    start: startOperating,
    stop: stopOperating,
  } = useOperatorHeartbeat(robotId);

  const { data: me } = useMeQuery();
  const { scopeIds } = useRobotScope();

  const toggleFullscreen = useCallback(() => {
    if (!contentRef.current) return;
    if (document.fullscreenElement) {
      document.exitFullscreen();
    } else {
      contentRef.current.requestFullscreen();
    }
  }, []);

  useEffect(() => {
    const handler = () => setIsFullscreen(!!document.fullscreenElement);
    document.addEventListener("fullscreenchange", handler);
    return () => document.removeEventListener("fullscreenchange", handler);
  }, []);

  // Robot config (REST one-shot — camera layout is static per robot).
  const { data: robot, isLoading: isLoadingRobot } = useRobotQuery(robotId);

  // Auto-start: if the robot has no active operator and the current user
  // is an in-scope Operator, start once. Resume after page reload is
  // handled by useOperatorHeartbeat's checkStatus (detects own lock in
  // Redis and resumes heartbeat automatically).
  const autoStarted = useRef(false);
  useEffect(() => {
    if (autoStarted.current || !robot || !me) return;
    if (
      !robot.active_operator &&
      me.role === USER_ROLE.OPERATOR &&
      isRobotInScope(robotId, scopeIds) &&
      !isOperating &&
      !isLocked
    ) {
      autoStarted.current = true;
      startOperating();
    }
  }, [robot, me, robotId, scopeIds, isOperating, isLocked, startOperating]);

  // Single combined SSE stream for live robot status, current episode,
  // and task metadata. Replaces three SSE hooks + two REST hooks plus the
  // old `streamEpisode ?? robotEpisode` merge race.
  const {
    status: realtimeStatus,
    episode,
    task,
    isConnected: isRobotConnected,
  } = useTeleopStream(robotId, !!robot);
  const activeEpisodeId = episode?.id;
  // "Loading" whenever we don't yet have an episode but the stream is
  // connected (or still in its initial handshake). Going from connected
  // to disconnected should show a disconnected indicator, not a loader.
  const isLoadingEpisode = !episode && isRobotConnected !== false;

  const sortedSubtasks = [...(episode?.subtasks ?? [])].sort(
    (a, b) => a.order_index - b.order_index
  );
  const currentSubtask = sortedSubtasks.find(
    (s) => s.status === SUBTASK_COLLECTION_STATUS.IN_PROGRESS
  );
  const nextSubtask = sortedSubtasks.find(
    (s) =>
      s.status === SUBTASK_COLLECTION_STATUS.READY &&
      s.order_index > (currentSubtask?.order_index ?? -1)
  );

  // Layout config
  const robotLayout = extractRobotLayout(robot?.robot_config);
  const views = robotLayout?.teleoperation ?? DEFAULT_TELEOP_VIEWS;
  const layout = views[viewName];
  const defaultLayout = views.default ?? DEFAULT_TELEOP_VIEWS.default!;
  const { host, port, rosbridgePort } = extractHostPort(robot?.robot_config);
  const cameras = extractCameras(robot?.robot_config);
  const streamConfig = extractStreamConfig(robot?.robot_config);

  // Nav: other views (exclude current)
  const otherViews = Object.entries(views).filter(
    ([key]) => key !== viewName && key !== "default"
  );

  // View title
  const viewTitle =
    layout?.title ??
    (isDefault
      ? t("teleop.teleoperation")
      : viewName.charAt(0).toUpperCase() + viewName.slice(1));

  // Build context
  const layoutContext: LayoutContext = useMemo(
    () => ({
      robot: robot
        ? {
            id: robot.id,
            name: robot.name,
            robot_type: robot.robot_type ?? undefined,
          }
        : undefined,
      realtimeStatus,
      isRobotConnected,
      isLoadingRobot,
      episode,
      isLoadingEpisode,
      activeEpisodeId,
      currentSubtask: currentSubtask ?? null,
      nextSubtask: nextSubtask ?? null,
      taskName: task?.name,
      taskVersion: task?.version,
      taskManualUrl: task?.manual_url,
      taskDescription: task?.description,
      cameras,
      host,
      port,
      rosbridgePort,
      gateLevel: realtimeStatus?.gate_conditions?.gate_level,
      streamConfig,
    }),
    [
      robot,
      realtimeStatus,
      isRobotConnected,
      isLoadingRobot,
      episode,
      isLoadingEpisode,
      activeEpisodeId,
      currentSubtask,
      nextSubtask,
      task,
      cameras,
      host,
      port,
      rosbridgePort,
      streamConfig,
    ]
  );

  return (
    <div className="space-y-4">
      {/* Breadcrumb */}
      <nav className="flex items-center text-sm text-gray-500">
        <Link
          href="/robots"
          className="hover:text-gray-700 dark:hover:text-gray-300"
        >
          {t("teleop.robots")}
        </Link>
        <ChevronRight className="h-4 w-4 mx-2" />
        {robot ? (
          <Link
            href={`/robots/${robot.id}`}
            className="hover:text-gray-700 dark:hover:text-gray-300"
          >
            {robot.name}
          </Link>
        ) : (
          <div className="h-4 w-24 bg-gray-200 dark:bg-gray-700 rounded animate-pulse" />
        )}
        <ChevronRight className="h-4 w-4 mx-2" />
        {isDefault ? (
          <span className="text-gray-900 dark:text-gray-100 font-medium">
            {t("teleop.teleoperation")}
          </span>
        ) : (
          <>
            <Link
              href={`/robots/${robotId}/teleoperation`}
              className="hover:text-gray-700 dark:hover:text-gray-300"
            >
              {t("teleop.teleoperation")}
            </Link>
            <ChevronRight className="h-4 w-4 mx-2" />
            <span className="text-gray-900 dark:text-gray-100 font-medium">
              {viewTitle}
            </span>
          </>
        )}
      </nav>

      {/* Back */}
      <Link
        href={
          isDefault ? `/robots/${robotId}` : `/robots/${robotId}/teleoperation`
        }
      >
        <Button variant="ghost" size="sm">
          <ArrowLeft className="h-4 w-4 mr-1" />
          {t("teleop.back")}
        </Button>
      </Link>

      {/* Main Content */}
      <div
        ref={contentRef}
        className={cn(
          "space-y-4",
          isFullscreen && "bg-background p-4 h-screen overflow-auto text-lg"
        )}
      >
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              {viewTitle}
            </h1>
            <p className="text-gray-600 dark:text-gray-400">
              {robot?.robot_type && <span>{robot.robot_type} · </span>}
              {robot?.name ?? t("teleop.loadingRobotName")}
            </p>
          </div>

          <div className="flex items-center gap-3">
            {/* Connection Indicator */}
            <div
              className={cn(
                "flex items-center gap-1 text-sm",
                isRobotConnected ? "text-green-600" : "text-red-500"
              )}
            >
              {isRobotConnected ? (
                <Signal className="h-4 w-4" />
              ) : (
                <SignalZero className="h-4 w-4" />
              )}
              <span>
                {isRobotConnected
                  ? t("teleop.connected")
                  : t("teleop.disconnected")}
              </span>
            </div>

            {/* Nav to default (if not on default) */}
            {!isDefault && (
              <Link href={`/robots/${robotId}/teleoperation`}>
                <Button variant="outline" size="sm">
                  {defaultLayout.title ?? t("teleop.teleoperation")}
                </Button>
              </Link>
            )}

            {/* Nav to other views */}
            {otherViews.map(([vName, vConfig]) => (
              <Link
                key={vName}
                href={`/robots/${robotId}/teleoperation/${vName}`}
              >
                <Button variant="outline" size="sm">
                  {vConfig.title ?? vName}
                </Button>
              </Link>
            ))}

            {/* Start/Stop Teleoperation */}
            {isOperating ? (
              <Button
                variant="outline"
                size="sm"
                onClick={stopOperating}
                className="text-red-600 border-red-200 hover:bg-red-50 dark:text-red-400 dark:border-red-800 dark:hover:bg-red-950"
              >
                <Square className="h-4 w-4 mr-1" />
                {t("teleop.stopTeleoperation")}
              </Button>
            ) : (
              <Button
                variant="outline"
                size="sm"
                onClick={startOperating}
                disabled={isLocked}
                className={
                  isLocked
                    ? ""
                    : "text-green-600 border-green-200 hover:bg-green-50 dark:text-green-400 dark:border-green-800 dark:hover:bg-green-950"
                }
                title={
                  isLocked && lockedBy
                    ? t("teleop.operatedByOrg", {
                        name: lockedBy.display_name,
                        org: lockedBy.organization_name,
                      })
                    : undefined
                }
              >
                <Play className="h-4 w-4 mr-1" />
                {isLocked && lockedBy
                  ? t("teleop.operatedBy", { name: lockedBy.display_name })
                  : t("teleop.startTeleoperation")}
              </Button>
            )}

            {/* Fullscreen Toggle */}
            <Button variant="outline" size="sm" onClick={toggleFullscreen}>
              {isFullscreen ? (
                <>
                  <Minimize className="h-4 w-4 mr-1" />
                  {t("teleop.exitFullscreen")}
                </>
              ) : (
                <>
                  <Maximize className="h-4 w-4 mr-1" />
                  {t("teleop.fullscreen")}
                </>
              )}
            </Button>
          </div>
        </div>

        {/* Layout or not-found */}
        {layout ? (
          <div className={undefined}>
            <LayoutRenderer
              layout={layout}
              context={layoutContext}
              fillHeight={isFullscreen}
            />
          </div>
        ) : (
          <Card>
            <CardContent className="py-16 text-center">
              <p className="text-gray-500 dark:text-gray-400">
                {t("teleop.viewNotDefined", { viewName })}
              </p>
              <Link href={`/robots/${robotId}/teleoperation`}>
                <Button variant="outline" size="sm" className="mt-4">
                  {t("teleop.backToTeleoperation")}
                </Button>
              </Link>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
