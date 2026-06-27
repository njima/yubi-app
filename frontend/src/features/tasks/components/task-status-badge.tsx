"use client";

import { useTaskStatusLabel } from "@/shared/hooks/use-status-labels";
import { type TaskStatusValue } from "@/shared/lib/status-constants";
import { TASK_STATUS_DISPLAY } from "@/shared/lib/status-display";

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
