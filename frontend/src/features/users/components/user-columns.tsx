"use client";

import { type ColumnDef } from "@tanstack/react-table";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";
import { type UserRoleValue } from "@/lib/status/constants";

import { DataTableColumnHeader } from "@/features/tasks/components/data-table-column-header";

import { EditUserDialog } from "./edit-user-dialog";

type User = z.infer<typeof schemas.UserResponse>;
type MeUser = User | undefined;

interface UserColumnsOptions {
  canUpdateRole: boolean;
  currentUser: MeUser;
  t: (key: string, options?: Record<string, unknown>) => string;
  formatRelativeTime: (dateString: string | null | undefined) => string;
  getRoleLabel: (role: UserRoleValue) => string;
}

export function getUserColumns({
  canUpdateRole,
  currentUser,
  t,
  formatRelativeTime,
  getRoleLabel,
}: UserColumnsOptions): ColumnDef<User>[] {
  return [
    {
      accessorKey: "user_id",
      header: "ID",
      enableSorting: false,
      cell: ({ row }) => (
        <span className="text-gray-500 dark:text-gray-400 text-sm font-mono">
          {row.original.user_id}
        </span>
      ),
    },
    {
      accessorKey: "name",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("userColumns.displayName")}
        />
      ),
      cell: ({ row }) => (
        <span className="font-medium">{row.original.display_name}</span>
      ),
    },
    {
      accessorKey: "email",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t("userColumns.email")} />
      ),
      cell: ({ row }) => (
        <span className="text-gray-600 dark:text-gray-400">
          {row.original.email}
        </span>
      ),
    },
    {
      id: "organization",
      header: t("userColumns.organization"),
      enableSorting: false,
      cell: ({ row }) => (
        <span className="text-gray-600 dark:text-gray-400">
          {row.original.organization_name}
        </span>
      ),
    },
    {
      accessorKey: "location",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("userColumns.locations")}
        />
      ),
      cell: ({ row }) => {
        const locations = row.original.locations;
        if (!locations || locations.length === 0) {
          return <span className="text-gray-400 dark:text-gray-500">-</span>;
        }
        return (
          <div className="flex flex-wrap gap-1">
            {locations.slice(0, 2).map((l) => (
              <span
                key={l.location_id}
                className="inline-flex items-center rounded-full bg-gray-100 dark:bg-gray-700 px-2 py-0.5 text-xs text-gray-700 dark:text-gray-300"
              >
                {l.name}
              </span>
            ))}
            {locations.length > 2 && (
              <span className="inline-flex items-center rounded-full bg-gray-200 dark:bg-gray-600 px-2 py-0.5 text-xs text-gray-500 dark:text-gray-400">
                {t("userColumns.more", { count: locations.length - 2 })}
              </span>
            )}
          </div>
        );
      },
    },
    {
      accessorKey: "role",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t("userColumns.role")} />
      ),
      cell: ({ row }) => (
        <span className="text-gray-600 dark:text-gray-400">
          {row.original.role !== undefined
            ? getRoleLabel(row.original.role)
            : "-"}
        </span>
      ),
    },
    {
      accessorKey: "created_at",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t("userColumns.created")}
        />
      ),
      cell: ({ row }) => (
        <span className="text-gray-600 dark:text-gray-400 text-sm">
          {formatRelativeTime(row.original.created_at)}
        </span>
      ),
    },
    {
      id: "updated_at",
      header: t("userColumns.updated"),
      enableSorting: false,
      cell: ({ row }) => (
        <span className="text-gray-600 dark:text-gray-400 text-sm">
          {formatRelativeTime(row.original.updated_at)}
        </span>
      ),
    },
    {
      id: "total_working_hours",
      header: t("userColumns.totalWorkingHours"),
      enableSorting: false,
      cell: () => (
        <span className="text-gray-600 dark:text-gray-400 text-sm">-</span>
      ),
    },
    ...(canUpdateRole
      ? [
          {
            id: "actions",
            header: () => (
              <div className="text-right">{t("userColumns.actions")}</div>
            ),
            enableSorting: false,
            cell: ({ row }: { row: { original: User } }) => (
              <div className="flex justify-end gap-2">
                <EditUserDialog user={row.original} currentUser={currentUser} />
              </div>
            ),
          } as ColumnDef<User>,
        ]
      : []),
  ];
}
