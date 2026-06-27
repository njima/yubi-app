"use client";

import { type ColumnDef } from "@tanstack/react-table";
import { AlertCircle, ExternalLink, Pencil } from "lucide-react";
import Link from "next/link";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";
import { formatDateTime } from "@/lib/date-utils";

import { Button } from "@/components/ui/button";

import { DataTableColumnHeader } from "@/features/tasks/components/data-table-column-header";

import { EditEpisodeDialog } from "./edit-episode-dialog";
import { EpisodeStatusBadge } from "./episode-status-badge";
import { GradeBar } from "./grade-bar";

type Episode = z.infer<typeof schemas.Episode>;

interface EpisodeColumnsOptions {
  canEdit: boolean;
  taskNameById: Map<string, string>;
  robotNameById: Map<string, string>;
  userNameById: Map<string, string>;
  locationNameById: Map<string, string>;
  t: (key: string) => string;
}

export function getEpisodeColumns({
  canEdit,
  taskNameById,
  robotNameById,
  userNameById,
  locationNameById,
  t,
}: EpisodeColumnsOptions): ColumnDef<Episode>[] {
  return [
    {
      accessorKey: "id",
      header: "ID",
      enableSorting: false,
      cell: ({ row }) => (
        <span className="font-mono text-sm">
          {row.original.id.substring(0, 8)}
        </span>
      ),
    },
    {
      accessorKey: "task",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t("episodesPage.task")} />
      ),
      cell: ({ row }) => (
        <span className="text-sm">
          {row.original.task_id
            ? (taskNameById.get(row.original.task_id) ?? row.original.task_id)
            : "-"}
        </span>
      ),
    },
    {
      accessorKey: "task_version",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("episodesPage.taskVersion")}
        />
      ),
      enableSorting: false,
      cell: ({ row }) => (
        <span
          className="text-sm truncate max-w-[200px] inline-block align-bottom"
          title={
            row.original.task_version_display_name ??
            row.original.task_version_id
          }
        >
          {row.original.task_version_display_name ??
            row.original.task_version_id ??
            "-"}
        </span>
      ),
    },
    {
      id: "status",
      header: () => t("episodesPage.status"),
      enableSorting: false,
      cell: ({ row }) => <EpisodeStatusBadge status={row.original.status} />,
    },
    {
      accessorKey: "robot",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("episodesPage.robot")}
        />
      ),
      cell: ({ row }) => (
        <span className="text-sm">
          {row.original.robot_id
            ? (robotNameById.get(row.original.robot_id) ??
              row.original.robot_id)
            : "-"}
        </span>
      ),
    },
    {
      accessorKey: "recorded_by",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("episodesPage.recordedBy")}
        />
      ),
      cell: ({ row }) => (
        <span className="text-sm">
          {row.original.recorded_by
            ? (userNameById.get(row.original.recorded_by) ??
              row.original.recorded_by)
            : "-"}
        </span>
      ),
    },
    {
      id: "location",
      header: () => t("topNav.locations"),
      enableSorting: false,
      cell: ({ row }) => (
        <span className="text-sm">
          {row.original.location_id
            ? (locationNameById.get(row.original.location_id) ??
              row.original.location_id)
            : "-"}
        </span>
      ),
    },
    {
      accessorKey: "started_at",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("episodesPage.startedAt")}
        />
      ),
      cell: ({ row }) => (
        <span className="text-sm text-gray-600 dark:text-gray-400">
          {row.original.started_at
            ? formatDateTime(row.original.started_at)
            : "-"}
        </span>
      ),
    },
    {
      accessorKey: "ended_at",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("episodesPage.endedAt")}
        />
      ),
      cell: ({ row }) => (
        <span className="text-sm text-gray-600 dark:text-gray-400">
          {row.original.ended_at ? formatDateTime(row.original.ended_at) : "-"}
        </span>
      ),
    },
    {
      id: "grade",
      header: () => t("episodesPage.grade"),
      enableSorting: false,
      cell: ({ row }) => (
        <GradeBar
          value={row.original.average_grade}
          count={row.original.grade_count}
          size="sm"
        />
      ),
    },
    {
      accessorKey: "error",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t("common.error")} />
      ),
      cell: ({ row }) => {
        const hasError =
          row.original.error_details && row.original.error_details.length > 0;
        return hasError ? (
          <AlertCircle className="h-4 w-4 text-red-500" />
        ) : (
          <span className="text-gray-400">-</span>
        );
      },
    },
    {
      id: "actions",
      header: () => (
        <div className="text-right">{t("episodesPage.actions")}</div>
      ),
      enableSorting: false,
      cell: ({ row }) => (
        <div className="flex justify-end gap-2">
          <Link href={`/episodes/${row.original.id}`}>
            <Button variant="ghost" size="sm">
              <ExternalLink className="h-4 w-4" />
            </Button>
          </Link>
          {canEdit && (
            <EditEpisodeDialog episodeId={row.original.id}>
              <Button variant="ghost" size="sm">
                <Pencil className="h-4 w-4" />
              </Button>
            </EditEpisodeDialog>
          )}
        </div>
      ),
    },
  ];
}
