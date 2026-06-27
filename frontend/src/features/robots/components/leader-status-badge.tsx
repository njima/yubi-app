"use client";

import {
  LEADER_STATUS,
  type LeaderStatusValue,
} from "@/shared/lib/status-constants";

import { Badge } from "@/components/ui/badge";

const statusStyles: Record<LeaderStatusValue, string> = {
  [LEADER_STATUS.READY]:
    "bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300",
  [LEADER_STATUS.FAULTED]:
    "bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300",
  [LEADER_STATUS.MAINTENANCE]:
    "bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300",
};

const statusLabels: Record<LeaderStatusValue, string> = {
  [LEADER_STATUS.READY]: "Ready",
  [LEADER_STATUS.FAULTED]: "Faulted",
  [LEADER_STATUS.MAINTENANCE]: "Maintenance",
};

interface LeaderStatusBadgeProps {
  statusCode?: number | null;
}

export function LeaderStatusBadge({ statusCode }: LeaderStatusBadgeProps) {
  if (statusCode == null) {
    return <span className="text-sm text-gray-500">-</span>;
  }

  const status = statusCode as LeaderStatusValue;
  return (
    <Badge
      variant="outline"
      className={statusStyles[status] ?? statusStyles[LEADER_STATUS.READY]}
    >
      {statusLabels[status] ?? "Unknown"}
    </Badge>
  );
}
