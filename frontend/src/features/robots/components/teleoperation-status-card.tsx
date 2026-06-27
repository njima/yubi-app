"use client";

import { Activity, Battery, Clock, Shield, Wifi } from "lucide-react";
import { useTranslation } from "react-i18next";

import { formatUptime } from "@/lib/date-utils";
import { EPISODE_COLLECTION_STATUS } from "@/lib/status/constants";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { GateGroupGrid, GateStatusBadge } from "./gate-status-badge";
import { BATTERY_THRESHOLDS } from "../schemas/robot";

import type { Robot, RobotStatusStreamDetail } from "../schemas/robot";

interface TeleoperationStatusCardProps {
  robot?: Robot;
  realtimeStatus?: RobotStatusStreamDetail | null;
  episodeStatus?: number;
}

function StatusRow({
  icon: Icon,
  label,
  children,
}: {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
        <Icon className="h-4 w-4" />
        <span>{label}</span>
      </div>
      <div className="text-sm font-medium">{children}</div>
    </div>
  );
}

function EpisodeStateBadge({ status }: { status?: number }) {
  const { t } = useTranslation();

  switch (status) {
    case EPISODE_COLLECTION_STATUS.READY:
      return (
        <span className="inline-block px-2 py-0.5 rounded text-xs font-medium text-yellow-600 bg-yellow-100">
          {t("status.ready")}
        </span>
      );
    case EPISODE_COLLECTION_STATUS.RECORDING:
      return (
        <span className="inline-block px-2 py-0.5 rounded text-xs font-medium text-green-600 bg-green-100 animate-pulse">
          {t("status.recording")}
        </span>
      );
    case EPISODE_COLLECTION_STATUS.CANCEL:
      return (
        <span className="inline-block px-2 py-0.5 rounded text-xs font-medium text-red-600 bg-red-100">
          {t("status.cancelled")}
        </span>
      );
    case EPISODE_COLLECTION_STATUS.COMPLETED:
      return (
        <span className="inline-block px-2 py-0.5 rounded text-xs font-medium text-gray-600 bg-gray-100">
          {t("status.completed")}
        </span>
      );
    default:
      return <span className="text-sm text-gray-400">-</span>;
  }
}

export function TeleoperationStatusCard({
  robot,
  realtimeStatus,
  episodeStatus,
}: TeleoperationStatusCardProps) {
  const { t } = useTranslation();

  const batteryLevel = realtimeStatus?.battery_pct ?? robot?.battery_level;
  const hasBatteryData = batteryLevel !== undefined && batteryLevel !== null;
  const batteryColor = !hasBatteryData
    ? "text-gray-500 dark:text-gray-400"
    : batteryLevel < BATTERY_THRESHOLDS.CRITICAL
      ? "text-red-600 dark:text-red-400"
      : batteryLevel < BATTERY_THRESHOLDS.LOW
        ? "text-yellow-600 dark:text-yellow-400"
        : "text-green-600 dark:text-green-400";

  const connectionPct = realtimeStatus?.connection_pct;
  const uptimeSec = realtimeStatus?.uptime_sec;

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-base font-medium">
          {t("teleopStatus.title")}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        {/* Battery */}
        <StatusRow icon={Battery} label={t("teleopStatus.battery")}>
          <span className={batteryColor}>
            {hasBatteryData ? `${batteryLevel}%` : "-"}
          </span>
        </StatusRow>

        {/* Connection */}
        <StatusRow icon={Wifi} label={t("teleopStatus.connection")}>
          {connectionPct !== undefined && connectionPct !== null
            ? `${connectionPct}%`
            : "-"}
        </StatusRow>

        {/* Uptime */}
        <StatusRow icon={Clock} label={t("teleopStatus.uptime")}>
          {uptimeSec !== undefined && uptimeSec !== null
            ? formatUptime(uptimeSec)
            : "-"}
        </StatusRow>

        {/* Gate Status */}
        <StatusRow icon={Shield} label={t("teleopStatus.gate")}>
          <GateStatusBadge gateConditions={realtimeStatus?.gate_conditions} />
        </StatusRow>

        {realtimeStatus?.gate_conditions && (
          <GateGroupGrid gateConditions={realtimeStatus.gate_conditions} />
        )}

        {/* Episode State */}
        <StatusRow icon={Activity} label={t("teleopStatus.episode")}>
          <EpisodeStateBadge status={episodeStatus} />
        </StatusRow>
      </CardContent>
    </Card>
  );
}
