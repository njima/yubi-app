"use client";

import { CheckCircle, Circle } from "lucide-react";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";
import { SUBTASK_COLLECTION_STATUS } from "@/lib/status/constants";
import { cn } from "@/lib/utils";

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

interface SubtaskTimelineProps {
  subtasks: EpisodeSubTask[];
  parameterValues?: Record<string, string> | null;
  isLoading?: boolean;
}

function getRowClassName(status: number): string {
  switch (status) {
    case SUBTASK_COLLECTION_STATUS.COMPLETED:
      return "text-gray-400 dark:text-gray-500";
    case SUBTASK_COLLECTION_STATUS.IN_PROGRESS:
      return "bg-blue-50 dark:bg-blue-950/50 text-blue-600 dark:text-blue-400 font-medium";
    case SUBTASK_COLLECTION_STATUS.SKIPPED:
    case SUBTASK_COLLECTION_STATUS.CANCELLED:
      return "text-gray-400 opacity-60";
    default:
      return "";
  }
}

function StatusDisplay({
  status,
  t,
}: {
  status: number;
  t: (key: string) => string;
}) {
  switch (status) {
    case SUBTASK_COLLECTION_STATUS.COMPLETED:
      return (
        <span className="flex items-center gap-1">
          <CheckCircle className="h-3.5 w-3.5" />
          {t("subtaskTimeline.complete")}
        </span>
      );
    case SUBTASK_COLLECTION_STATUS.IN_PROGRESS:
      return (
        <span className="flex items-center gap-1">
          <Circle className="h-3.5 w-3.5 fill-current" />
          {t("status.inProgress")}
        </span>
      );
    case SUBTASK_COLLECTION_STATUS.READY:
      return (
        <span className="flex items-center gap-1">
          <Circle className="h-3.5 w-3.5" />
          {t("status.ready")}
        </span>
      );
    case SUBTASK_COLLECTION_STATUS.SKIPPED:
      return <span>{t("status.skipped")}</span>;
    case SUBTASK_COLLECTION_STATUS.CANCELLED:
      return <span>{t("status.cancelled")}</span>;
    default:
      return <span>{t("status.unknown")}</span>;
  }
}

export function SubtaskTimeline({
  subtasks,
  parameterValues,
  isLoading = false,
}: SubtaskTimelineProps) {
  const { t } = useTranslation();
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

  if (subtasks.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-8 text-center">
        <p className="text-sm text-gray-500 dark:text-gray-400">
          {t("subtaskTimeline.empty")}
        </p>
      </div>
    );
  }

  const sorted = [...subtasks].sort((a, b) => a.order_index - b.order_index);

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-10">#</TableHead>
          <TableHead>{t("subtaskTimeline.name")}</TableHead>
          <TableHead className="w-28">{t("subtaskTimeline.status")}</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {sorted.map((subtask) => (
          <TableRow
            key={subtask.subtask_id}
            className={cn(getRowClassName(subtask.status))}
          >
            <TableCell>{subtask.order_index + 1}</TableCell>
            <TableCell>
              <ParameterizedName
                name={subtask.name}
                parameterValues={parameterValues}
              />
            </TableCell>
            <TableCell>
              <StatusDisplay status={subtask.status} t={t} />
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
