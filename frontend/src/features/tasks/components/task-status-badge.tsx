"use client";

import { useTaskStatusLabel } from "@/lib/hooks/use-status-labels";
import { type TaskStatusValue } from "@/lib/status/constants";
import { TASK_STATUS_DISPLAY } from "@/lib/status/display";

import { Badge } from "@/components/ui/badge";

interface TaskStatusBadgeProps {
  status: TaskStatusValue;
}

export function TaskStatusBadge({ status }: TaskStatusBadgeProps) {
  const getStatusLabel = useTaskStatusLabel();

  return (
    <Badge variant="outline" className={TASK_STATUS_DISPLAY[status].className}>
      {getStatusLabel(status)}
    </Badge>
  );
}
