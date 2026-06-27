"use client";

import { Battery } from "lucide-react";
import { useTranslation } from "react-i18next";

import { formatUptime } from "@/lib/date-utils";
import { ROBOT_STATUS } from "@/lib/status/constants";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { BATTERY_THRESHOLDS } from "../schemas/robot";

import type { Robot, RobotStatusStreamDetail } from "../schemas/robot";

interface RobotStatusCardProps {
  robot: Robot;
  realtimeStatus?: RobotStatusStreamDetail | null;
}

export function RobotStatusCard({
  robot,
  realtimeStatus,
}: RobotStatusCardProps) {
  const { t } = useTranslation();
  const batteryLevel = realtimeStatus?.battery_pct ?? robot.battery_level;
  const hasBatteryData = batteryLevel !== undefined && batteryLevel !== null;
  const isCritical =
    hasBatteryData && batteryLevel < BATTERY_THRESHOLDS.CRITICAL;
  const isLow = hasBatteryData && batteryLevel < BATTERY_THRESHOLDS.LOW;

  const getBatteryColorClass = () => {
    if (!hasBatteryData) return "text-gray-500 dark:text-gray-400";
    if (isCritical) return "text-red-600 dark:text-red-400";
    if (isLow) return "text-yellow-600 dark:text-yellow-400";
    return "text-green-600 dark:text-green-400";
  };

  const getStatusLabel = (status?: number) => {
    switch (status) {
      case ROBOT_STATUS.ONLINE:
        return t("status.online");
      case ROBOT_STATUS.BUSY:
        return t("status.busy");
      case ROBOT_STATUS.OFFLINE:
        return t("status.offline");
      case ROBOT_STATUS.FAULTED:
        return t("status.faulted");
      case ROBOT_STATUS.MAINTENANCE:
        return t("status.maintenance");
      default:
        return t("status.unknown");
    }
  };

  const getStatusStyle = (status?: number) => {
    switch (status) {
      case ROBOT_STATUS.ONLINE:
        return "text-green-600 bg-green-50 border-green-200 dark:text-green-400 dark:bg-green-950 dark:border-green-800";
      case ROBOT_STATUS.BUSY:
        return "text-yellow-600 bg-yellow-50 border-yellow-200 dark:text-yellow-400 dark:bg-yellow-950 dark:border-yellow-800";
      case ROBOT_STATUS.OFFLINE:
        return "text-red-600 bg-red-50 border-red-200 dark:text-red-400 dark:bg-red-950 dark:border-red-800";
      case ROBOT_STATUS.FAULTED:
        return "text-orange-600 bg-orange-50 border-orange-200 dark:text-orange-400 dark:bg-orange-950 dark:border-orange-800";
      case ROBOT_STATUS.MAINTENANCE:
        return "text-blue-600 bg-blue-50 border-blue-200 dark:text-blue-400 dark:bg-blue-950 dark:border-blue-800";
      default:
        return "text-gray-600 bg-gray-50 border-gray-200 dark:text-gray-400 dark:bg-gray-950 dark:border-gray-800";
    }
  };

  const uptimeSec = realtimeStatus?.uptime_sec;
  const hasUptimeData = uptimeSec !== undefined && uptimeSec !== null;

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-base font-medium">
          {t("robotStatusCard.title")}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Battery */}
        <div className="flex items-center justify-between">
          <span className="text-sm text-gray-600 dark:text-gray-400">
            {t("robotStatusCard.battery")}
          </span>
          <span className="text-sm font-medium">
            {hasBatteryData ? `${batteryLevel}%` : "-"}
          </span>
        </div>
        <div className={`flex items-center gap-2 ${getBatteryColorClass()}`}>
          <Battery className="h-4 w-4" />
          <span className="text-sm font-medium">
            {hasBatteryData ? `${batteryLevel}%` : "-"}
          </span>
        </div>

        <div className="border-t border-gray-200 dark:border-gray-700" />

        {/* Status */}
        <div className="space-y-2">
          <span className="text-sm text-gray-600 dark:text-gray-400">
            {t("robotStatusCard.status")}
          </span>
          <div>
            <span
              className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${getStatusStyle(robot.status)}`}
            >
              {getStatusLabel(robot.status)}
            </span>
          </div>
        </div>

        <div className="border-t border-gray-200 dark:border-gray-700" />

        {/* Uptime */}
        <div className="space-y-1">
          <span className="text-sm text-gray-600 dark:text-gray-400">
            {t("robotStatusCard.uptime")}
          </span>
          <p className="text-sm font-medium">
            {hasUptimeData ? formatUptime(uptimeSec) : "-"}
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
