"use client";

import { type ColumnDef } from "@tanstack/react-table";
import { ExternalLink, Pencil } from "lucide-react";
import Link from "next/link";

import { toDateTimeLocalValue } from "@/shared/lib/date-utils";
import { truncateUuid } from "@/shared/lib/format";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

import { DataTableColumnHeader } from "./data-table-column-header";
import { EditTaskDialog } from "./edit-task-dialog";
import { TagCategoryBadge } from "./tag-category-badge";
import { TaskDifficultyBadge } from "./task-difficulty-badge";
import { TaskPriorityBadge } from "./task-priority-badge";
import { TaskStatusBadge } from "./task-status-badge";
import { secondsToHoursMinutes } from "../lib/duration";

import type { Task } from "../schemas";

interface TaskColumnsOptions {
  canEdit: boolean;
  t: (key: string) => string;
}

export function getTaskColumns({
  canEdit,
  t,
}: TaskColumnsOptions): ColumnDef<Task>[] {
  return [
    {
      accessorKey: "id",
      header: "ID",
      enableSorting: false,
      cell: ({ row }) => (
        <Tooltip>
          <TooltipTrigger asChild>
            <span className="cursor-default text-gray-500 dark:text-gray-400 text-sm font-mono">
              {truncateUuid(row.original.id)}
            </span>
          </TooltipTrigger>
          <TooltipContent>{row.original.id}</TooltipContent>
        </Tooltip>
      ),
    },
    {
      accessorKey: "name",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t("taskColumns.name")} />
      ),
      cell: ({ row }) => (
        <span className="font-medium">{row.original.name}</span>
      ),
    },
    {
      accessorKey: "robot_type",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("taskColumns.robotType")}
        />
      ),
      cell: ({ row }) => (
        <span className="text-sm text-gray-600 dark:text-gray-400">
          {row.original.robot_type ?? "-"}
        </span>
      ),
    },
    {
      id: "tags",
      header: t("taskColumns.categoryTags"),
      enableSorting: false,
      cell: ({ row }) => {
        const tags = row.original.tags ?? [];
        return (
          <div className="flex flex-wrap gap-1">
            {tags.slice(0, 2).map((tag) => (
              <TagCategoryBadge
                key={tag.id}
                categoryTypeName={tag.category_type_name}
                name={tag.name}
              />
            ))}
            {tags.length > 2 && (
              <Badge variant="outline" className="text-xs">
                +{tags.length - 2}
              </Badge>
            )}
          </div>
        );
      },
    },
    {
      accessorKey: "priority",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("taskColumns.priority")}
        />
      ),
      cell: ({ row }) => <TaskPriorityBadge priority={row.original.priority} />,
    },
    {
      accessorKey: "difficulty",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("taskColumns.difficulty")}
        />
      ),
      cell: ({ row }) => (
        <TaskDifficultyBadge difficulty={row.original.difficulty} />
      ),
    },
    {
      accessorKey: "target_duration_seconds",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("taskColumns.collectionTimeTarget")}
        />
      ),
      cell: ({ row }) => {
        const seconds = row.original.target_duration_seconds;
        if (seconds == null)
          return (
            <span className="text-sm text-gray-400 dark:text-gray-500">—</span>
          );
        const { hours, minutes } = secondsToHoursMinutes(seconds);
        return (
          <span className="text-sm text-gray-600 dark:text-gray-400">
            {hours}h {minutes}m
          </span>
        );
      },
    },
    {
      id: "episodes",
      header: t("taskColumns.episodes"),
      enableSorting: false,
      cell: ({ row }) => {
        const target = row.original.target_episode_count;
        const actual = row.original.actual_episode_count ?? 0;

        if (target != null && target > 0) {
          const done = actual >= target;
          const percent = Math.min(Math.round((actual / target) * 100), 100);
          return (
            <div className="min-w-[100px]">
              <div className="flex items-center justify-between text-xs mb-1">
                <span className="font-medium">
                  {actual} / {target}
                </span>
                {done && (
                  <Badge
                    variant="outline"
                    className="bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300"
                  >
                    {t("taskColumns.done")}
                  </Badge>
                )}
              </div>
              <div className="h-1.5 w-full rounded-full bg-gray-200 dark:bg-gray-700">
                <div
                  className={`h-full rounded-full ${done ? "bg-green-500" : "bg-blue-500"}`}
                  style={{ width: `${percent}%` }}
                />
              </div>
            </div>
          );
        }

        if (actual > 0) {
          return (
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {actual}
            </span>
          );
        }

        return (
          <span className="text-sm text-gray-400 dark:text-gray-500">—</span>
        );
      },
    },
    {
      accessorKey: "status",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("taskColumns.status")}
        />
      ),
      cell: ({ row }) => <TaskStatusBadge status={row.original.status} />,
    },
    {
      accessorKey: "version",
      header: t("taskColumns.version"),
      enableSorting: false,
      cell: ({ row }) => {
        const versionText = row.original.version ?? "-";
        const display = row.original.version_display_name;
        const fallback = row.original.version
          ? `${row.original.name} ${row.original.version}`
          : null;
        // Hide the secondary label when the backend resolved to the trivial
        // "{name} {version}" fallback to avoid duplicate text.
        const showDisplay = display && display !== fallback;
        return (
          <div className="flex flex-col">
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {versionText}
            </span>
            {showDisplay && (
              <span
                className="text-xs text-gray-500 dark:text-gray-500 truncate max-w-[200px]"
                title={display}
              >
                {display}
              </span>
            )}
          </div>
        );
      },
    },
    {
      id: "actions",
      header: () => (
        <div className="text-right">{t("taskColumns.actions")}</div>
      ),
      enableSorting: false,
      cell: ({ row }) => {
        const task = row.original;
        return (
          <div className="flex justify-end gap-2">
            <Link href={`/tasks/${task.id}`}>
              <Button variant="ghost" size="sm">
                <ExternalLink className="h-4 w-4" />
              </Button>
            </Link>
            {canEdit && (
              <EditTaskDialog
                taskId={task.id}
                name={task.name}
                description={task.description}
                manual_url={task.manual_url}
                priority={task.priority}
                difficulty={task.difficulty}
                status={task.status}
                deadline={toDateTimeLocalValue(task.deadline)}
                robot_type={task.robot_type ?? undefined}
                tags={task.tags}
              >
                <Button variant="ghost" size="sm">
                  <Pencil className="h-4 w-4" />
                </Button>
              </EditTaskDialog>
            )}
          </div>
        );
      },
    },
  ];
}
