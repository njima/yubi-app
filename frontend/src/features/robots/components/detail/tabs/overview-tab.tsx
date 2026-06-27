"use client";

import { Battery, Clock, MapPin, Building2, Signal, User } from "lucide-react";
import { useTranslation } from "react-i18next";

import { formatUptime } from "@/shared/lib/date-utils";
import { LEADER_STATUS, ROBOT_STATUS } from "@/shared/lib/status-constants";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { useLocationQuery } from "@/features/locations";
import { useOrganizationQuery } from "@/features/organizations";

import { BATTERY_THRESHOLDS } from "../../../schemas/robot";
import { ConsecutiveFaultDaysBadge } from "../../consecutive-fault-days-badge";
import { GateStatusBadge } from "../../gate-status-badge";

import type { Robot, RobotStatusStreamDetail } from "../../../schemas/robot";

interface OverviewTabProps {
  robot: Robot;
  realtimeStatus?: RobotStatusStreamDetail | null;
  isConnected?: boolean;
}

export function OverviewTab({
  robot,
  realtimeStatus,
  isConnected,
}: OverviewTabProps) {
  const { t } = useTranslation();
  const { data: organization } = useOrganizationQuery(
    robot.organization_id ?? ""
  );
  const { data: location } = useLocationQuery(robot.location_id ?? "");
  const activeOperator = robot.active_operator;

  const batteryLevel = realtimeStatus?.battery_pct ?? robot.battery_level;

  return (
    <div className="space-y-4">
      {/* Real-time Status */}
      <Card>
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between">
            <CardTitle className="text-base font-medium">
              {t("robotOverview.realTimeStatus")}
            </CardTitle>
            <ConnectionIndicator isConnected={isConnected} t={t} />
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            <div>
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                {t("robotOverview.followerStatus")}
              </p>
              <div className="flex items-center gap-2">
                <StatusBadge status={robot.status} />
                <GateStatusBadge
                  gateConditions={realtimeStatus?.gate_conditions}
                />
              </div>
            </div>
            <div>
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                {t("robotOverview.battery")}
              </p>
              <BatteryValue batteryLevel={batteryLevel} />
            </div>
            <div>
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                {t("robotOverview.connection")}
              </p>
              <ConnectionValue connectionPct={realtimeStatus?.connection_pct} />
            </div>
            <div>
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                {t("robotOverview.uptime")}
              </p>
              <UptimeValue uptimeSec={realtimeStatus?.uptime_sec} />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Information */}
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-base font-medium">
            {t("robotOverview.information")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-5 gap-6">
            <div>
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                {t("robotOverview.leaderStatus")}
              </p>
              <LeaderStatusBadgeLocal status={robot.leader_status} />
            </div>
            <div>
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                {t("robotOverview.consecutiveFaultDays")}
              </p>
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <span className="text-sm text-gray-500 dark:text-gray-400">
                    {t("robotOverview.follower")}
                  </span>
                  <ConsecutiveFaultDaysBadge
                    days={robot.consecutive_fault_days}
                  />
                  {robot.consecutive_fault_days == null && (
                    <p className="text-lg font-semibold text-gray-400 dark:text-gray-500">
                      -
                    </p>
                  )}
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-sm text-gray-500 dark:text-gray-400">
                    {t("robotOverview.leader")}
                  </span>
                  <ConsecutiveFaultDaysBadge
                    days={robot.leader_consecutive_fault_days}
                    className="text-sm"
                  />
                  {robot.leader_consecutive_fault_days == null && (
                    <p className="text-lg font-semibold text-gray-400 dark:text-gray-500">
                      -
                    </p>
                  )}
                </div>
              </div>
            </div>
            <div>
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                {t("robotOverview.operator")}
              </p>
              {activeOperator ? (
                <div className="flex items-center gap-2">
                  <User className="h-5 w-5 text-gray-600 dark:text-gray-400" />
                  <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                    {activeOperator.display_name}
                  </p>
                  <span className="text-sm text-gray-500 dark:text-gray-400">
                    ({activeOperator.organization_name})
                  </span>
                </div>
              ) : (
                <p className="text-lg font-semibold text-gray-400 dark:text-gray-500">
                  {t("robotOverview.noActiveOperator")}
                </p>
              )}
            </div>
            <div>
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                {t("robotOverview.location")}
              </p>
              <div className="flex items-center gap-2">
                <MapPin className="h-5 w-5 text-gray-600 dark:text-gray-400" />
                <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  {location?.name || "-"}
                </p>
              </div>
            </div>
            <div>
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                {t("robotOverview.organization")}
              </p>
              <div className="flex items-center gap-2">
                <Building2 className="h-5 w-5 text-gray-600 dark:text-gray-400" />
                <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  {organization?.display_name || "-"}
                </p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function ConnectionIndicator({
  isConnected,
  t,
}: {
  isConnected?: boolean;
  t: (key: string) => string;
}) {
  return (
    <div className="flex items-center gap-2">
      <span className="relative flex h-2.5 w-2.5">
        {isConnected && (
          <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75" />
        )}
        <span
          className={`relative inline-flex rounded-full h-2.5 w-2.5 ${
            isConnected ? "bg-green-500" : "bg-gray-400 dark:bg-gray-500"
          }`}
        />
      </span>
      <span
        className={`text-xs font-medium ${
          isConnected
            ? "text-green-600 dark:text-green-400"
            : "text-gray-500 dark:text-gray-400"
        }`}
      >
        {isConnected
          ? t("robotOverview.live")
          : t("robotOverview.disconnected")}
      </span>
    </div>
  );
}

function StatusBadge({ status }: { status?: number }) {
  const { t } = useTranslation();

  const getLabel = () => {
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

  const getStyle = () => {
    switch (status) {
      case ROBOT_STATUS.ONLINE:
        return "text-green-600 bg-green-50 border-green-200 dark:text-green-400 dark:bg-green-950 dark:border-green-800";
      case ROBOT_STATUS.BUSY:
        return "text-yellow-600 bg-yellow-50 border-yellow-200 dark:text-yellow-400 dark:bg-yellow-950 dark:border-yellow-800";
      case ROBOT_STATUS.OFFLINE:
        return "text-orange-600 bg-orange-50 border-orange-200 dark:text-orange-400 dark:bg-orange-950 dark:border-orange-800";
      case ROBOT_STATUS.FAULTED:
        return "text-red-600 bg-red-50 border-red-200 dark:text-red-400 dark:bg-red-950 dark:border-red-800";
      case ROBOT_STATUS.MAINTENANCE:
        return "text-blue-600 bg-blue-50 border-blue-200 dark:text-blue-400 dark:bg-blue-950 dark:border-blue-800";
      default:
        return "text-gray-600 bg-gray-50 border-gray-200 dark:text-gray-400 dark:bg-gray-950 dark:border-gray-800";
    }
  };

  return (
    <span
      className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium border ${getStyle()}`}
    >
      {getLabel()}
    </span>
  );
}

const CONNECTION_THRESHOLDS = {
  GOOD: 70,
  MODERATE: 30,
} as const;

function BatteryValue({ batteryLevel }: { batteryLevel?: number }) {
  const hasData = batteryLevel !== undefined && batteryLevel !== null;
  const isCritical = hasData && batteryLevel < BATTERY_THRESHOLDS.CRITICAL;
  const isLow = hasData && batteryLevel < BATTERY_THRESHOLDS.LOW;

  const getColorClass = () => {
    if (!hasData) return "text-gray-500 dark:text-gray-400";
    if (isCritical) return "text-red-600 dark:text-red-400";
    if (isLow) return "text-yellow-600 dark:text-yellow-400";
    return "text-green-600 dark:text-green-400";
  };

  return (
    <div className={`flex items-center gap-2 ${getColorClass()}`}>
      <Battery className="h-5 w-5" />
      <p className="text-lg font-semibold">
        {hasData ? `${batteryLevel}%` : "-"}
      </p>
    </div>
  );
}

function ConnectionValue({ connectionPct }: { connectionPct?: number }) {
  const hasData = connectionPct !== undefined && connectionPct !== null;

  const getColorClass = () => {
    if (!hasData) return "text-gray-500 dark:text-gray-400";
    if (connectionPct >= CONNECTION_THRESHOLDS.GOOD)
      return "text-green-600 dark:text-green-400";
    if (connectionPct >= CONNECTION_THRESHOLDS.MODERATE)
      return "text-yellow-600 dark:text-yellow-400";
    return "text-red-600 dark:text-red-400";
  };

  return (
    <div className={`flex items-center gap-2 ${getColorClass()}`}>
      <Signal className="h-5 w-5" />
      <p className="text-lg font-semibold">
        {hasData ? `${connectionPct}%` : "-"}
      </p>
    </div>
  );
}

function UptimeValue({ uptimeSec }: { uptimeSec?: number }) {
  const hasData = uptimeSec !== undefined && uptimeSec !== null;

  return (
    <div className="flex items-center gap-2 text-gray-900 dark:text-gray-100">
      <Clock className="h-5 w-5 text-gray-600 dark:text-gray-400" />
      <p className="text-lg font-semibold">
        {hasData ? formatUptime(uptimeSec) : "-"}
      </p>
    </div>
  );
}

function LeaderStatusBadgeLocal({ status }: { status?: number | null }) {
  if (status == null) {
    return (
      <p className="text-lg font-semibold text-gray-400 dark:text-gray-500">
        -
      </p>
    );
  }

  const getLabel = () => {
    switch (status) {
      case LEADER_STATUS.READY:
        return "Ready";
      case LEADER_STATUS.FAULTED:
        return "Faulted";
      case LEADER_STATUS.MAINTENANCE:
        return "Maintenance";
      default:
        return "Unknown";
    }
  };

  const getStyle = () => {
    switch (status) {
      case LEADER_STATUS.READY:
        return "text-green-600 bg-green-50 border-green-200 dark:text-green-400 dark:bg-green-950 dark:border-green-800";
      case LEADER_STATUS.FAULTED:
        return "text-red-600 bg-red-50 border-red-200 dark:text-red-400 dark:bg-red-950 dark:border-red-800";
      case LEADER_STATUS.MAINTENANCE:
        return "text-blue-600 bg-blue-50 border-blue-200 dark:text-blue-400 dark:bg-blue-950 dark:border-blue-800";
      default:
        return "text-gray-600 bg-gray-50 border-gray-200 dark:text-gray-400 dark:bg-gray-950 dark:border-gray-800";
    }
  };

  return (
    <span
      className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium border ${getStyle()}`}
    >
      {getLabel()}
    </span>
  );
}
