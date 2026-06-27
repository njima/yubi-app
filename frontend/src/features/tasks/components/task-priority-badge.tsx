"use client";

import { useTaskPriorityLabel } from "@/lib/hooks/use-status-labels";
import { TASK_PRIORITY, type TaskPriorityValue } from "@/lib/status/constants";
import { cn } from "@/lib/utils";

const priorityDotColors: Record<TaskPriorityValue, string> = {
  [TASK_PRIORITY.LOW]: "bg-gray-400",
  [TASK_PRIORITY.NORMAL]: "bg-blue-500",
  [TASK_PRIORITY.HIGH]: "bg-amber-500",
  [TASK_PRIORITY.URGENT]: "bg-red-500",
};

interface TaskPriorityBadgeProps {
  priority: TaskPriorityValue;
}

export function TaskPriorityBadge({ priority }: TaskPriorityBadgeProps) {
  const getPriorityLabel = useTaskPriorityLabel();

  return (
    <span className="inline-flex items-center gap-1.5 text-xs font-medium text-gray-700 dark:text-gray-300">
      <span
        className={cn(
          "h-2 w-2 shrink-0 rounded-full",
          priorityDotColors[priority]
        )}
      />
      {getPriorityLabel(priority)}
    </span>
  );
}
