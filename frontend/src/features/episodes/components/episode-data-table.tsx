"use client";

import {
  type ColumnDef,
  type OnChangeFn,
  type SortingState,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { useTranslation } from "react-i18next";

import { PaginationFooter } from "@/components/ui/pagination-footer";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

interface EpisodeDataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  sorting: SortingState;
  onSortingChange: OnChangeFn<SortingState>;
  isLoading?: boolean;
  totalCount?: number;
  page?: number;
  totalPages?: number;
  onPageChange?: (page: number) => void;
  limit?: number;
  onLimitChange?: (limit: number) => void;
}

export function EpisodeDataTable<TData, TValue>({
  columns,
  data,
  sorting,
  onSortingChange,
  isLoading = false,
  totalCount,
  page,
  totalPages,
  onPageChange,
  limit,
  onLimitChange,
}: EpisodeDataTableProps<TData, TValue>) {
  const { t } = useTranslation();

  // eslint-disable-next-line react-hooks/incompatible-library -- TanStack Table is not yet compatible with React Compiler
  const table = useReactTable({
    data,
    columns,
    state: { sorting },
    onSortingChange,
    manualSorting: true,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="rounded-lg border bg-white dark:bg-gray-800 dark:border-gray-700">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {isLoading ? (
            Array.from({ length: 5 }).map((_, i) => (
              <TableRow key={`skeleton-${i}`}>
                {columns.map((_, j) => (
                  <TableCell key={`skeleton-${i}-${j}`}>
                    <div className="h-4 w-24 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : table.getRowModel().rows.length > 0 ? (
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell
                colSpan={columns.length}
                className="h-24 text-center text-gray-600 dark:text-gray-400"
              >
                {t("episodesPage.noEpisodes")}
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>

      {totalCount != null &&
        page != null &&
        totalPages != null &&
        onPageChange && (
          <PaginationFooter
            page={page}
            totalPages={totalPages}
            totalCount={totalCount}
            onPageChange={onPageChange}
            itemLabel={t("topNav.episodes").toLowerCase()}
            limit={limit}
            onLimitChange={onLimitChange}
          />
        )}
    </div>
  );
}
