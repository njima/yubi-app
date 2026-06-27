"use client";

import { type ColumnDef } from "@tanstack/react-table";
import { AlertCircle, ExternalLink, Pencil, Trash2 } from "lucide-react";
import Link from "next/link";

import { truncateUuid } from "@/shared/lib/format";
import { ROBOT_STATUS, USER_ROLE } from "@/shared/lib/status-constants";

import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

import {
  ConsecutiveFaultDaysBadge,
  DeleteRobotDialog,
  EditRobotDialog,
  LeaderStatusBadge,
  RobotStatusBadge,
  StartTeleoperationButton,
  QueueTaskButton,
} from "@/features/robots";

import { DataTableColumnHeader } from "./data-table-column-header";
import { GateStatusBadgeInRobotList } from "./gate-status-badge";
import { useRobotsStatus } from "../contexts/robots-status-context";

import type { Robot } from "../schemas/robot";

interface RobotColumnsOptions {
  canEdit: boolean;
  canDelete: boolean;
  locationNameById: Map<string, string>;
  meRole?: number;
  isInScope?: (robotId: string) => boolean;
  t: (key: string) => string;
  formatDistanceTime: (dateString: string | null | undefined) => string;
  formatAbsoluteTime: (dateString: string | null | undefined) => string;
}

// Reads SSE-driven status from RobotsStatusContext at render time so the
// column array can stay referentially stable. Column-closure access to
// frequently-changing state would invalidate the columns useMemo on every
// SSE event, force flexRender to see new cell function identities, and
// remount every cell — wiping local useState in nested components like
// EditRobotDialog (the dialog snaps shut while the user is editing).
function GateCell({ robotId }: { robotId: string }) {
  const { statusMap, isConnected } = useRobotsStatus();
  if (!isConnected) {
    return <span className="text-sm text-gray-400">Connecting...</span>;
  }
  const realtimeStatus = statusMap[robotId];
  if (!realtimeStatus) {
    return <span className="text-sm text-gray-500">-</span>;
  }
  return (
    <GateStatusBadgeInRobotList
      gateConditions={realtimeStatus.gate_conditions}
    />
  );
}

export function getRobotColumns({
  canEdit,
  canDelete,
  locationNameById,
  meRole,
  isInScope,
  t,
  formatDistanceTime,
  formatAbsoluteTime,
}: RobotColumnsOptions): ColumnDef<Robot>[] {
  return [
    {
      accessorKey: "id",
      header: "ID",
      enableSorting: false,
      cell: ({ row }) => (
        <span className="font-mono text-sm text-gray-500 dark:text-gray-400">
          {truncateUuid(row.original.id)}
        </span>
      ),
    },
    {
      accessorKey: "name",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t("robotColumns.name")} />
      ),
      cell: ({ row }) => (
        <span className="font-medium">{row.original.name}</span>
      ),
    },
    {
      accessorKey: "location_id",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("robotColumns.location")}
        />
      ),
      cell: ({ row }) => {
        const id = row.original.location_id;
        const name = id ? (locationNameById.get(id) ?? id) : "-";
        return <span className="text-sm">{name}</span>;
      },
    },
    {
      accessorKey: "robot_type",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("robotColumns.robotType")}
        />
      ),
      cell: ({ row }) => (
        <span className="text-sm">{row.original.robot_type || "-"}</span>
      ),
    },
    {
      accessorKey: "status",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("robotColumns.followerStatus")}
        />
      ),
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
          <RobotStatusBadge statusCode={row.original.status ?? 2} />
          <ConsecutiveFaultDaysBadge
            days={row.original.consecutive_fault_days}
          />
        </div>
      ),
    },
    {
      accessorKey: "leader_status",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("robotColumns.leaderStatus")}
        />
      ),
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
          <LeaderStatusBadge statusCode={row.original.leader_status} />
          <ConsecutiveFaultDaysBadge
            days={row.original.leader_consecutive_fault_days}
          />
        </div>
      ),
    },
    {
      accessorKey: "gate",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Gate" />
      ),
      cell: ({ row }) => <GateCell robotId={row.original.id} />,
    },
    {
      accessorKey: "last_heartbeat_at",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("robotColumns.lastHeartbeat")}
        />
      ),
      cell: ({ row }) => {
        const robot = row.original;
        if (!robot.last_heartbeat_at) {
          return <span className="text-sm text-gray-500">-</span>;
        }
        return (
          <div className="flex items-center gap-2">
            <Tooltip>
              <TooltipTrigger asChild>
                <span className="text-sm cursor-help">
                  {formatDistanceTime(robot.last_heartbeat_at)}
                </span>
              </TooltipTrigger>
              <TooltipContent>
                {formatAbsoluteTime(robot.last_heartbeat_at)}
              </TooltipContent>
            </Tooltip>
            {robot.offline_reason && (
              <Tooltip>
                <TooltipTrigger asChild>
                  <AlertCircle className="h-4 w-4 text-red-500 cursor-help" />
                </TooltipTrigger>
                <TooltipContent className="max-w-xs">
                  <p className="font-semibold">
                    {t("robotColumns.offlineReason")}:
                  </p>
                  <p>{robot.offline_reason}</p>
                </TooltipContent>
              </Tooltip>
            )}
          </div>
        );
      },
    },
    {
      accessorKey: "active_episode_id",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("robotColumns.activeEpisode")}
        />
      ),
      cell: ({ row }) => {
        const episodeId = row.original.active_episode_id;
        if (!episodeId) {
          return <span className="text-sm text-gray-500">-</span>;
        }
        return (
          <Link
            href={`/episodes/${episodeId}`}
            className="inline-flex items-center gap-1 px-2 py-1 bg-blue-50 dark:bg-blue-950 rounded text-xs font-mono text-blue-600 dark:text-blue-400 hover:bg-blue-100 dark:hover:bg-blue-900 transition-colors"
          >
            {truncateUuid(episodeId)}
            <ExternalLink className="h-3 w-3" />
          </Link>
        );
      },
    },
    {
      accessorKey: "active_user_id",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("robotColumns.operator")}
        />
      ),
      cell: ({ row }) => (
        <span className="text-sm">
          {row.original.active_operator?.display_name ?? "—"}
        </span>
      ),
    },
    {
      id: "actions",
      header: () => (
        <div className="text-right">{t("robotColumns.actions")}</div>
      ),
      enableSorting: false,
      cell: ({ row }) => {
        const robot = row.original;
        return (
          <div className="flex justify-end gap-2">
            <StartTeleoperationButton
              robot={robot}
              inScope={isInScope ? isInScope(robot.id) : true}
            />
            {meRole !== undefined && meRole <= USER_ROLE.MANAGER && (
              <QueueTaskButton robot={robot} />
            )}
            <Link href={`/robots/${robot.id}`}>
              <Button variant="ghost" size="sm">
                <ExternalLink className="h-4 w-4" />
              </Button>
            </Link>
            {canEdit && (
              <EditRobotDialog robotId={robot.id}>
                <Button variant="ghost" size="sm">
                  <Pencil className="h-4 w-4" />
                </Button>
              </EditRobotDialog>
            )}
            {canDelete &&
              (robot.status === ROBOT_STATUS.BUSY ? (
                <Tooltip>
                  <TooltipTrigger asChild>
                    <span>
                      <Button variant="ghost" size="sm" disabled>
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </span>
                  </TooltipTrigger>
                  <TooltipContent>
                    {t("robotColumns.cannotDeleteBusy")}
                  </TooltipContent>
                </Tooltip>
              ) : (
                <DeleteRobotDialog robotId={robot.id} name={robot.name}>
                  <Button variant="ghost" size="sm">
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </DeleteRobotDialog>
              ))}
          </div>
        );
      },
    },
  ];
}
