"use client";

import {
  TASK_DIFFICULTY,
  type TaskDifficultyValue,
} from "@/lib/status/constants";
import { getTaskDifficultyLabel } from "@/lib/status/utils";
import { cn } from "@/lib/utils";

const difficultyDotColors: Record<TaskDifficultyValue, string> = {
  [TASK_DIFFICULTY.S]: "bg-red-500",
  [TASK_DIFFICULTY.A]: "bg-amber-500",
  [TASK_DIFFICULTY.B]: "bg-blue-500",
  [TASK_DIFFICULTY.C]: "bg-gray-400",
};

interface TaskDifficultyBadgeProps {
  difficulty: TaskDifficultyValue;
}

export function TaskDifficultyBadge({ difficulty }: TaskDifficultyBadgeProps) {
  return (
    <span className="inline-flex items-center gap-1.5 text-xs font-medium text-gray-700 dark:text-gray-300">
      <span
        className={cn(
          "h-2 w-2 shrink-0 rounded-full",
          difficultyDotColors[difficulty]
        )}
      />
      {getTaskDifficultyLabel(difficulty)}
    </span>
  );
}
