"use client";

import { type ColumnDef } from "@tanstack/react-table";
import { Pencil, Trash2 } from "lucide-react";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { truncateUuid } from "@/shared/lib/format";

import { Button } from "@/components/ui/button";

import { DataTableColumnHeader } from "./data-table-column-header";
import { DeleteLocationDialog } from "./delete-location-dialog";
import { EditLocationDialog } from "./edit-location-dialog";

type Location = z.infer<typeof schemas.Location>;

interface LocationColumnsOptions {
  canUpdate: boolean;
  canDelete: boolean;
  t: (key: string) => string;
}

export function getLocationColumns({
  canUpdate,
  canDelete,
  t,
}: LocationColumnsOptions): ColumnDef<Location>[] {
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
        <DataTableColumnHeader
          column={column}
          title={t("locationColumns.name")}
        />
      ),
      cell: ({ row }) => (
        <span className="font-medium">{row.original.name}</span>
      ),
    },
    {
      accessorKey: "site_name",
      header: t("locationColumns.site"),
      enableSorting: false,
      cell: ({ row }) => <span>{row.original.site_name}</span>,
    },
    ...(canUpdate || canDelete
      ? [
          {
            id: "actions",
            header: () => (
              <div className="text-right">{t("locationColumns.actions")}</div>
            ),
            enableSorting: false,
            cell: ({ row }: { row: { original: Location } }) => {
              const location = row.original;
              return (
                <div className="flex justify-end gap-1">
                  {canUpdate && (
                    <EditLocationDialog location={location}>
                      <Button size="sm" variant="ghost">
                        <Pencil className="h-4 w-4" />
                      </Button>
                    </EditLocationDialog>
                  )}
                  {canDelete && (
                    <DeleteLocationDialog
                      locationId={location.id}
                      name={location.name}
                    >
                      <Button size="sm" variant="ghost">
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </DeleteLocationDialog>
                  )}
                </div>
              );
            },
          } as ColumnDef<Location>,
        ]
      : []),
  ];
}
