"use client";

import {
  AlertTriangle,
  CheckCircle,
  Circle,
  Clock,
  ExternalLink,
} from "lucide-react";
import Link from "next/link";
import { useTranslation } from "react-i18next";

import type { schemas } from "@/lib/api/generated/api";

import { SUBTASK_COLLECTION_STATUS } from "@/shared/lib/status-constants";
import { cn } from "@/shared/lib/utils";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import type {
  LayoutContext,
  LayoutCamera,
} from "@/features/robots/lib/teleop-layout-registry";
import type { CameraLayoutItem } from "@/features/robots/lib/teleop-layout-types";
import { ParameterizedName } from "@/features/tasks";

import { CameraView } from "../components/camera-view";
import {
  GateStatusBadge,
  GateGroupCell,
} from "../components/gate-status-badge";
import { SubtaskTimeline } from "../components/subtask-timeline";

import type { LucideIcon } from "lucide-react";
import type { z } from "zod";

type EpisodeSubTask = z.infer<typeof schemas.EpisodeSubTask>;

function resolveCamera(
  ref: string,
  cameras: LayoutCamera[]
): LayoutCamera | undefined {
  const lower = (s: string | undefined) => s?.toLowerCase() ?? "";
  if (ref === "*main*") {
    return cameras.find(
      (c) =>
        lower(c.name).includes("head") || lower(c.namespace).includes("head")
    );
  }
  if (ref === "*left*") {
    return cameras.find(
      (c) =>
        lower(c.name).includes("left") || lower(c.namespace).includes("left")
    );
  }
  if (ref === "*right*") {
    return cameras.find(
      (c) =>
        lower(c.name).includes("right") || lower(c.namespace).includes("right")
    );
  }
  return cameras.find((c) => c.namespace === ref);
}

export function CameraRenderer({
  item,
  context,
}: {
  item: CameraLayoutItem;
  context: LayoutContext;
}) {
  const { cameras, host, port } = context;
  const camera = cameras?.length ? resolveCamera(item.ref, cameras) : undefined;

  return (
    <CameraView
      camera={camera}
      host={host}
      port={port}
      robotName={context.robot?.name}
      placeholderLabel={item.ref}
      showOverlays={item.overlay}
      episodeStatus={context.episode?.status}
      errorDetails={context.episode?.error_details}
      currentSubtask={context.currentSubtask}
      nextSubtask={context.nextSubtask}
      parameterValues={context.episode?.parameter_values}
      gateLevel={context.gateLevel}
      streamConfig={context.streamConfig}
    />
  );
}

// --- Gate Information Card ---

export function GateStatusCard({ ctx }: { ctx: LayoutContext }) {
  const { t } = useTranslation();
  const gateConditions = ctx.realtimeStatus?.gate_conditions;

  return (
    <Card>
      <CardHeader className="px-4 py-3">
        <CardTitle className="text-xl font-medium text-gray-600 dark:text-gray-300">
          {t("teleop.gate")}
        </CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-3 pt-0">
        {gateConditions && (
          <div className="flex items-center gap-2">
            <GateStatusBadge
              gateConditions={gateConditions}
              className="text-xl"
            />
            {Object.entries(gateConditions.groups).map(([name, group]) => (
              <div key={name} className="flex-1">
                <GateGroupCell name={name} group={group} className="text-xl" />
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

// -- Task Progressing Card ---

export function SubTaskProgressCard({ ctx }: { ctx: LayoutContext }) {
  const { t } = useTranslation();
  const subtasks = [...(ctx.episode?.subtasks ?? [])].sort(
    (a, b) => a.order_index - b.order_index
  );

  if (subtasks.length === 0) {
    return null;
  }

  const completedCount = subtasks.filter(
    (s) => s.status === SUBTASK_COLLECTION_STATUS.COMPLETED
  ).length;

  const currentTask = subtasks.find(
    (s) => s.status === SUBTASK_COLLECTION_STATUS.IN_PROGRESS
  )?.name;

  return (
    <Card>
      <CardHeader className="px-4 py-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-xl font-medium text-gray-600 dark:text-gray-300">
            {t("teleop.subtaskProgress")}
          </CardTitle>
          <div className="flex items-center gap-2">
            <span className="text-xl font-medium text-gray-700 dark:text-gray-200">
              {t("teleop.episodeId")}:{" "}
            </span>
            <span className="text-xl font-medium text-gray-700 dark:text-gray-200">
              {ctx.activeEpisodeId ?? "-"}
            </span>
          </div>
        </div>
      </CardHeader>
      <CardContent className="flex flex-col gap-2 px-4 pb-3 pt-0">
        {/* subtasks dot */}
        <div className="flex items-center gap-4">
          <div className="flex gap-1">
            {subtasks.map((s) => {
              const completed =
                s.status === SUBTASK_COLLECTION_STATUS.COMPLETED;
              const inProgress =
                s.status === SUBTASK_COLLECTION_STATUS.IN_PROGRESS;
              const skipped = s.status === SUBTASK_COLLECTION_STATUS.SKIPPED;
              const cancelled =
                s.status === SUBTASK_COLLECTION_STATUS.CANCELLED;

              return (
                <span
                  key={s.id}
                  title={s.name}
                  className={cn(
                    "h-4 w-4 rounded-full transition-colors",
                    completed && "bg-green-500",
                    inProgress && "bg-blue-500 ring-4 ring-blue-300",
                    skipped && "bg-gray-300 dark:bg-gray-600",
                    cancelled && "bg-red-500",
                    !completed &&
                      !inProgress &&
                      !skipped &&
                      !cancelled &&
                      "bg-gray-200 dark:bg-gray-700"
                  )}
                />
              );
            })}
          </div>
          <span className="text-xl font-medium text-gray-700 dark:text-gray-200">
            {t("teleop.subtaskCompleted", {
              completed: completedCount,
              total: subtasks.length,
            })}
          </span>
        </div>

        {/* Current Task */}
        <div className="flex items-center gap-2">
          <span className="text-xl font-medium text-gray-700 dark:text-gray-200">
            {t("teleop.currentTask")}:{" "}
          </span>
          <span className="text-xl font-medium text-gray-700 dark:text-gray-200">
            {currentTask ?? "-"}
          </span>
        </div>
      </CardContent>
    </Card>
  );
}

// --- Task Information Card ---

export function TaskInformationCard({ ctx }: { ctx: LayoutContext }) {
  const { t } = useTranslation();
  const robotId = ctx.robot?.id;

  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="text-base font-medium">
            {t("teleopTaskInfo.title")}
          </CardTitle>
          {robotId && (
            <Link href={`/robots/${robotId}/teleoperation/subtasks`}>
              <ExternalLink className="h-4 w-4 text-gray-400 hover:text-gray-600" />
            </Link>
          )}
        </div>
        {(ctx.taskName || ctx.activeEpisodeId) && (
          <div className="text-sm text-gray-500 dark:text-gray-400 space-y-0.5 mt-1">
            {ctx.taskName && (
              <p>
                {ctx.taskName}
                {ctx.taskVersion && (
                  <span className="ml-1 text-gray-400">
                    ({ctx.taskVersion})
                  </span>
                )}
              </p>
            )}
            {ctx.activeEpisodeId && (
              <p className="text-xs font-mono text-gray-400 dark:text-gray-500">
                {t("teleopTaskInfo.episode")}:{" "}
                {ctx.activeEpisodeId.substring(0, 8)}
              </p>
            )}
          </div>
        )}
      </CardHeader>
      <CardContent>
        <div className="max-h-[400px] overflow-y-auto">
          <SubtaskTimeline
            subtasks={ctx.episode?.subtasks ?? []}
            parameterValues={ctx.episode?.parameter_values}
            isLoading={ctx.isLoadingEpisode}
          />
        </div>
      </CardContent>
    </Card>
  );
}

// --- Subtask Detail List (from subtasks page) ---

export function SubtaskDetailList({ ctx }: { ctx: LayoutContext }) {
  const { t } = useTranslation();
  const subtasks = [...(ctx.episode?.subtasks ?? [])].sort(
    (a, b) => a.order_index - b.order_index
  );

  if (ctx.isLoadingEpisode) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="space-y-2">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="h-16 bg-gray-200 dark:bg-gray-700 rounded-lg animate-pulse"
              />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!ctx.episode) {
    return (
      <Card>
        <CardContent className="py-16 text-center">
          <p className="text-sm text-gray-400">
            {t("teleopTaskInfo.noActiveEpisode")}
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardContent className="pt-6">
        {ctx.taskDescription && (
          <p className="text-sm text-gray-700 dark:text-gray-300 mb-4">
            {ctx.taskDescription}
          </p>
        )}

        <h3 className="text-sm font-medium text-gray-500 mb-3">
          {t("teleopTaskInfo.subtasks")} ({subtasks.length})
        </h3>

        {subtasks.length === 0 ? (
          <p className="text-sm text-gray-400 text-center py-8">
            {t("teleopTaskInfo.noSubtasksDefined")}
          </p>
        ) : (
          <div className="space-y-2">
            {subtasks.map((subtask) => (
              <SubtaskItem
                key={subtask.id}
                subtask={subtask}
                parameterValues={ctx.episode?.parameter_values}
              />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

// --- SubtaskItem + ExecutionStatusBadge ---

function useExecutionStatusConfig() {
  const { t } = useTranslation();

  return {
    0: {
      label: t("status.ready"),
      className: "bg-gray-100 text-gray-700",
      icon: Circle,
    },
    1: {
      label: t("status.started"),
      className: "bg-blue-100 text-blue-700",
      icon: Clock,
    },
    2: {
      label: t("status.cancel"),
      className: "bg-yellow-100 text-yellow-700",
      icon: AlertTriangle,
    },
    3: {
      label: t("status.finished"),
      className: "bg-green-100 text-green-700",
      icon: CheckCircle,
    },
  } as Record<number, { label: string; className: string; icon: LucideIcon }>;
}

function ExecutionStatusBadge({ status }: { status: number }) {
  const { t } = useTranslation();
  const statusConfig = useExecutionStatusConfig();
  const config = statusConfig[status] ?? {
    label: t("status.ready"),
    className: "bg-gray-100 text-gray-700",
    icon: Circle,
  };
  const Icon = config.icon;
  return (
    <Badge variant="outline" className={cn(config.className, "gap-1")}>
      <Icon className="h-3 w-3" /> {config.label}
    </Badge>
  );
}

function SubtaskItem({
  subtask,
  parameterValues,
}: {
  subtask: EpisodeSubTask;
  parameterValues?: Record<string, string> | null;
}) {
  const isInProgress = subtask.status === SUBTASK_COLLECTION_STATUS.IN_PROGRESS;
  const isSkippedOrCancelled =
    subtask.status === SUBTASK_COLLECTION_STATUS.SKIPPED ||
    subtask.status === SUBTASK_COLLECTION_STATUS.CANCELLED;
  const executions = subtask.executions ?? [];

  return (
    <div
      className={cn(
        "border rounded-lg p-4",
        isInProgress &&
          "bg-blue-50 dark:bg-blue-950/50 border-blue-200 dark:border-blue-800",
        isSkippedOrCancelled && "opacity-60"
      )}
    >
      <p className="font-medium">
        {subtask.order_index + 1}.{" "}
        <ParameterizedName
          name={subtask.name}
          parameterValues={parameterValues}
        />
      </p>

      {executions.length > 0 && (
        <div className="mt-2 space-y-1">
          {executions.map((exec) => (
            <div
              key={exec.id}
              className="flex items-center justify-between gap-2"
            >
              <span className="text-xs text-gray-500 font-mono truncate">
                {exec.id}
              </span>
              <ExecutionStatusBadge status={exec.status} />
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
