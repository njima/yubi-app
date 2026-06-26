"use client";

import { registerLayoutComponent } from "@/shared/lib/layout-registry";
import { isCameraItem } from "@/shared/lib/layout-types";

import { TeachMeBizCard } from "@/features/tasks";

import {
  CameraRenderer,
  GateStatusCard,
  SubtaskDetailList,
  SubTaskProgressCard,
  TaskInformationCard,
} from "./teleop-layout-components";
import { SubtaskTimeline } from "../components/subtask-timeline";
import { TeleoperationStatusCard } from "../components/teleoperation-status-card";
import { TfVisualizer } from "../components/tf-visualizer";
import { ThreeDModelCard } from "../components/three-d-model-card";

let registered = false;

/**
 * Register teleoperation components with the layout registry.
 * Safe to call multiple times — only registers once.
 */
export function registerTeleopComponents() {
  if (registered) return;
  registered = true;

  registerLayoutComponent("camera", (ctx, item) => {
    if (!isCameraItem(item)) return null;
    return <CameraRenderer item={item} context={ctx} />;
  });

  registerLayoutComponent("task-information", (ctx) => (
    <TaskInformationCard ctx={ctx} />
  ));

  registerLayoutComponent("gate-information", (ctx) => (
    <GateStatusCard ctx={ctx} />
  ));

  registerLayoutComponent("status-card", (ctx) => (
    <TeleoperationStatusCard
      robot={
        ctx.robot
          ? {
              ...ctx.robot,
              organization_id: undefined,
              location_id: undefined,
              battery_level: undefined,
            }
          : undefined
      }
      realtimeStatus={ctx.realtimeStatus}
      episodeStatus={ctx.episode?.status}
    />
  ));

  registerLayoutComponent("3d-model", () => <ThreeDModelCard />);

  registerLayoutComponent("subtask-progress", (ctx) => (
    <SubTaskProgressCard ctx={ctx} />
  ));

  registerLayoutComponent("subtask-timeline", (ctx) => (
    <div className="max-h-[400px] overflow-y-auto">
      <SubtaskTimeline
        subtasks={ctx.episode?.subtasks ?? []}
        parameterValues={ctx.episode?.parameter_values}
        isLoading={ctx.isLoadingEpisode}
      />
    </div>
  ));

  registerLayoutComponent("teach-me-card", (ctx) => (
    <TeachMeBizCard manualUrl={ctx.taskManualUrl} />
  ));

  registerLayoutComponent("subtask-detail-list", (ctx) => (
    <SubtaskDetailList ctx={ctx} />
  ));

  registerLayoutComponent("tf-visualizer", (ctx) => <TfVisualizer ctx={ctx} />);
}
