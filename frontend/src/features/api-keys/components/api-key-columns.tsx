"use client";

import { type ColumnDef } from "@tanstack/react-table";
import { Link as LinkIcon } from "lucide-react";
import Link from "next/link";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { Badge } from "@/components/ui/badge";

import { RevokeApiKeyDialog } from "./revoke-api-key-dialog";

type ApiKey = z.infer<typeof schemas.ApiKeyResponse>;

interface ApiKeyColumnsOptions {
  canRevoke: boolean;
  t: (key: string, options?: Record<string, unknown>) => string;
  formatDateTime: (s: string | null | undefined) => string;
}

export function getApiKeyColumns({
  canRevoke,
  t,
  formatDateTime,
}: ApiKeyColumnsOptions): ColumnDef<ApiKey>[] {
  return [
    {
      accessorKey: "name",
      header: t("apiKeyColumns.name"),
      enableSorting: false,
      cell: ({ row }) => (
        <span className="font-medium">{row.original.name}</span>
      ),
    },
    {
      accessorKey: "key_hint",
      header: t("apiKeyColumns.keyHint"),
      enableSorting: false,
      cell: ({ row }) => (
        <span className="rounded bg-gray-100 dark:bg-gray-800 px-2 py-0.5 font-mono text-xs">
          {row.original.key_hint}
        </span>
      ),
    },
    {
      accessorKey: "robot_name",
      header: t("apiKeyColumns.robot"),
      enableSorting: false,
      cell: ({ row }) => {
        const { robot_id, robot_name } = row.original;
        if (!robot_id) {
          return <span className="text-gray-400">-</span>;
        }
        return (
          <Link
            href={`/robots/${robot_id}`}
            className="inline-flex items-center gap-1 text-blue-600 dark:text-blue-400 hover:underline"
          >
            <LinkIcon className="h-3 w-3" />
            {robot_name ?? robot_id}
          </Link>
        );
      },
    },
    {
      accessorKey: "user_name",
      header: t("apiKeyColumns.owner"),
      enableSorting: false,
      cell: ({ row }) => <span>{row.original.user_name}</span>,
    },
    {
      accessorKey: "created_at",
      header: t("apiKeyColumns.createdAt"),
      enableSorting: false,
      cell: ({ row }) => (
        <span className="text-sm text-gray-600 dark:text-gray-400">
          {formatDateTime(row.original.created_at)}
        </span>
      ),
    },
    {
      accessorKey: "last_used_at",
      header: t("apiKeyColumns.lastUsedAt"),
      enableSorting: false,
      cell: ({ row }) => (
        <span className="text-sm text-gray-600 dark:text-gray-400">
          {row.original.last_used_at
            ? formatDateTime(row.original.last_used_at)
            : t("apiKeyColumns.never")}
        </span>
      ),
    },
    {
      id: "status",
      header: t("apiKeyColumns.status"),
      enableSorting: false,
      cell: ({ row }) => {
        const revoked = row.original.revoked_at != null;
        return revoked ? (
          <Badge variant="destructive">
            {t("apiKeyColumns.statusRevoked")}
          </Badge>
        ) : (
          <Badge variant="default" className="bg-green-600 hover:bg-green-700">
            {t("apiKeyColumns.statusActive")}
          </Badge>
        );
      },
    },
    {
      id: "actions",
      header: () => (
        <div className="text-right">{t("apiKeyColumns.actions")}</div>
      ),
      enableSorting: false,
      cell: ({ row }) => {
        const revoked = row.original.revoked_at != null;
        if (revoked || !canRevoke) return null;
        return (
          <div className="flex justify-end">
            <RevokeApiKeyDialog
              apiKeyId={row.original.id}
              name={row.original.name}
              robotName={row.original.robot_name}
            />
          </div>
        );
      },
    },
  ];
}
