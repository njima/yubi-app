"use client";

import { useTranslation } from "react-i18next";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { useSubtaskCollectionStatusConfig } from "@/shared/hooks/use-subtask-collection-status-config";
import { formatDuration, formatTime } from "@/shared/lib/date-utils";
import { SUBTASK_COLLECTION_STATUS } from "@/shared/lib/status-constants";

import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

import { ParameterizedName } from "@/features/tasks";

type EpisodeSubTask = z.infer<typeof schemas.EpisodeSubTask>;

interface SubtaskControlTableProps {
  subtasks: EpisodeSubTask[];
  parameterValues?: Record<string, string> | null;
  isLoading?: boolean;
}

export function SubtaskControlTable({
  subtasks,
  parameterValues,
  isLoading = false,
}: SubtaskControlTableProps) {
  const { t } = useTranslation();
  const statusConfig = useSubtaskCollectionStatusConfig();
  // Loading state: show skeleton
  if (isLoading) {
    return (
      <div className="space-y-2">
        {[1, 2, 3].map((i) => (
          <div
            key={i}
            className="h-10 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"
          />
        ))}
      </div>
    );
  }

  // Empty state: show message
  if (subtasks.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-8 text-center">
        <p className="text-sm text-gray-500 dark:text-gray-400">
          {t("subtaskControl.empty")}
        </p>
      </div>
    );
  }

  // Normal state: show table
  return (
    <div className="max-h-64 overflow-y-auto">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-12">#</TableHead>
            <TableHead>{t("subtaskControl.name")}</TableHead>
            <TableHead className="w-32">{t("subtaskControl.status")}</TableHead>
            <TableHead className="w-20">{t("subtaskControl.runs")}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {subtasks.map((subtask, index) => {
            const statusKey = subtask.status ?? SUBTASK_COLLECTION_STATUS.READY;
            const config = statusConfig[statusKey];
            const IconComponent = config.icon;
            const executions = subtask.executions ?? [];

            return (
              <TableRow key={subtask.subtask_id}>
                <TableCell className="font-medium">
                  {(subtask.order_index ?? index) + 1}
                </TableCell>
                <TableCell>
                  <ParameterizedName
                    name={subtask.name}
                    parameterValues={parameterValues}
                  />
                </TableCell>
                <TableCell>
                  <Badge variant={config.variant} className="gap-1">
                    <IconComponent
                      className={`h-3 w-3 ${statusKey === SUBTASK_COLLECTION_STATUS.IN_PROGRESS ? "animate-spin" : ""}`}
                    />
                    {config.label}
                  </Badge>
                </TableCell>
                <TableCell>
                  {executions.length > 0 ? (
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <Badge variant="outline" className="cursor-default">
                            {executions.length}
                          </Badge>
                        </TooltipTrigger>
                        <TooltipContent side="top" className="max-w-xs">
                          <div className="space-y-1 text-xs">
                            {executions.map((execution, execIndex) => (
                              <div key={execution.id}>
                                <span className="font-medium">
                                  #{execIndex + 1}
                                </span>{" "}
                                {formatTime(execution.started_at)}
                                {" — "}
                                {formatDuration(
                                  execution.started_at,
                                  execution.finished_at
                                )}
                              </div>
                            ))}
                          </div>
                        </TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                  ) : (
                    <span className="text-muted-foreground">-</span>
                  )}
                </TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
}
