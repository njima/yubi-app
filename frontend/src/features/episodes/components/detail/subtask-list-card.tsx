"use client";

import { ChevronRight } from "lucide-react";
import { Fragment, useState } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { useSubtaskCollectionStatusConfig } from "@/shared/hooks/use-subtask-collection-status-config";
import { formatDateTime, formatDuration } from "@/shared/lib/date-utils";
import { SUBTASK_COLLECTION_STATUS } from "@/shared/lib/status-constants";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

import { ParameterizedName } from "@/features/tasks";

type EpisodeSubTask = z.infer<typeof schemas.EpisodeSubTask>;

interface SubtaskListCardProps {
  subtasks: EpisodeSubTask[];
  parameterValues?: Record<string, string> | null;
}

export function SubtaskListCard({
  subtasks,
  parameterValues,
}: SubtaskListCardProps) {
  const { t } = useTranslation();
  const statusConfig = useSubtaskCollectionStatusConfig();
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());

  if (subtasks.length === 0) {
    return null;
  }

  const toggleRow = (subtaskId: string) => {
    setExpandedRows((prev) => {
      const next = new Set(prev);
      if (next.has(subtaskId)) {
        next.delete(subtaskId);
      } else {
        next.add(subtaskId);
      }
      return next;
    });
  };

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">
          {t("episodeSubtasks.subtasks")}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-10"></TableHead>
              <TableHead className="w-12">#</TableHead>
              <TableHead>{t("episodeSubtasks.name")}</TableHead>
              <TableHead className="w-32">
                {t("episodeSubtasks.status")}
              </TableHead>
              <TableHead className="w-20">
                {t("episodeSubtasks.runs")}
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {subtasks.map((subtask, index) => {
              const statusKey = subtask.status as keyof ReturnType<
                typeof useSubtaskCollectionStatusConfig
              >;
              const config =
                statusConfig[statusKey] ??
                statusConfig[SUBTASK_COLLECTION_STATUS.READY];
              const IconComponent = config.icon;
              const executions = subtask.executions ?? [];
              const isExpanded = expandedRows.has(subtask.id);
              const hasExecutions = executions.length > 0;

              return (
                <Fragment key={subtask.id}>
                  <TableRow
                    className={
                      hasExecutions ? "cursor-pointer hover:bg-muted/50" : ""
                    }
                    onClick={() => hasExecutions && toggleRow(subtask.id)}
                  >
                    <TableCell className="px-2">
                      {hasExecutions && (
                        <ChevronRight
                          className={`h-4 w-4 text-muted-foreground transition-transform duration-200 ${
                            isExpanded ? "rotate-90" : ""
                          }`}
                        />
                      )}
                    </TableCell>
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
                          className={`h-3 w-3 ${
                            statusKey === SUBTASK_COLLECTION_STATUS.IN_PROGRESS
                              ? "animate-spin"
                              : ""
                          }`}
                        />
                        {config.label}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {executions.length > 0 ? (
                        <Badge variant="outline">{executions.length}</Badge>
                      ) : (
                        <span className="text-muted-foreground">-</span>
                      )}
                    </TableCell>
                  </TableRow>
                  {isExpanded && (
                    <TableRow key={`${subtask.id}-executions`}>
                      <TableCell colSpan={5} className="bg-muted/30 py-2 px-4">
                        <div className="space-y-1 text-sm">
                          {executions.map((execution, execIndex) => (
                            <div
                              key={execution.id}
                              className="flex items-center gap-2 text-muted-foreground"
                            >
                              <span className="font-medium text-foreground">
                                #{execIndex + 1}
                              </span>
                              <span>
                                {formatDateTime(execution.started_at)}
                              </span>
                              <span>→</span>
                              <span>
                                {formatDateTime(execution.finished_at)}
                              </span>
                              <span className="text-xs">
                                (
                                {formatDuration(
                                  execution.started_at,
                                  execution.finished_at
                                )}
                                )
                              </span>
                            </div>
                          ))}
                        </div>
                      </TableCell>
                    </TableRow>
                  )}
                </Fragment>
              );
            })}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
