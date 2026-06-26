"use client";

import { useRobotStatusLabel } from "@/shared/hooks/use-status-labels";
import {
  ROBOT_STATUS,
  type RobotStatusValue,
} from "@/shared/lib/status-constants";
import { ROBOT_STATUS_DISPLAY } from "@/shared/lib/status-display";
import { Badge } from "@/shared/ui/badge";

interface RobotStatusBadgeProps {
  statusCode: number;
}

export function RobotStatusBadge({ statusCode }: RobotStatusBadgeProps) {
  const getStatusLabel = useRobotStatusLabel();
  const status = statusCode as RobotStatusValue;
  const display =
    ROBOT_STATUS_DISPLAY[status] ?? ROBOT_STATUS_DISPLAY[ROBOT_STATUS.OFFLINE];

  return (
    <Badge variant="outline" className={display.className}>
      {getStatusLabel(status)}
    </Badge>
  );
}
