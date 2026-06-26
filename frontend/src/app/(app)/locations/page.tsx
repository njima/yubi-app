"use client";

import { type SortingState } from "@tanstack/react-table";
import {
  useQueryState,
  parseAsInteger,
  parseAsString,
  parseAsStringEnum,
} from "nuqs";
import { Suspense, useMemo } from "react";
import { useTranslation } from "react-i18next";

import { DEFAULT_PAGE_SIZE } from "@/shared/lib/pagination";
import { PaginationFooter } from "@/shared/ui/pagination-footer";
import { SearchableSelect } from "@/shared/ui/searchable-select";

import {
  useLocationsQuery,
  CreateLocationDialog,
  LocationDataTable,
  getLocationColumns,
} from "@/features/locations";
import { useSiteSearchOptions } from "@/features/sites";
import { usePermission } from "@/features/users";

const validSortBy = ["name"] as const;
const validSortOrder = ["asc", "desc"] as const;

function LocationsContent() {
  const { t } = useTranslation();
  const [page, setPage] = useQueryState("page", parseAsInteger.withDefault(1));
  const [limit, setLimit] = useQueryState(
    "limit",
    parseAsInteger.withDefault(DEFAULT_PAGE_SIZE)
  );
  const [siteId, setSiteId] = useQueryState("site_id", parseAsString);
  const [sortBy, setSortBy] = useQueryState(
    "sort_by",
    parseAsStringEnum([...validSortBy])
  );
  const [sortOrder, setSortOrder] = useQueryState(
    "sort_order",
    parseAsStringEnum([...validSortOrder])
  );

  const {
    options: siteSearchOptions,
    isLoading: siteSearchLoading,
    onSearch: onSiteSearch,
    selectedLabel: siteSelectedLabel,
    onValueChange: onSiteSelectChange,
  } = useSiteSearchOptions();

  const {
    data: locationsData,
    isLoading,
    error,
  } = useLocationsQuery({
    page,
    limit,
    site_id: siteId ?? undefined,
    sort_by: sortBy ?? undefined,
    sort_order: sortOrder ?? undefined,
  });
  const locations = locationsData?.locations ?? [];
  const pagination = locationsData?.pagination;
  const totalPages = pagination
    ? Math.ceil(pagination.count / pagination.limit)
    : 1;

  const canCreate = usePermission("location:create");
  const canUpdate = usePermission("location:update");
  const canDelete = usePermission("location:delete");

  const columns = useMemo(
    () => getLocationColumns({ canUpdate, canDelete, t }),
    [canUpdate, canDelete, t]
  );

  // URL → TanStack SortingState
  const sorting: SortingState = sortBy
    ? [{ id: sortBy, desc: sortOrder === "desc" }]
    : [];

  // TanStack onSortingChange → URL
  const handleSortingChange = (
    updaterOrValue: SortingState | ((old: SortingState) => SortingState)
  ) => {
    const newSorting =
      typeof updaterOrValue === "function"
        ? updaterOrValue(sorting)
        : updaterOrValue;

    if (newSorting.length === 0) {
      setSortBy(null);
      setSortOrder(null);
      setPage(1);
      return;
    }

    const first = newSorting[0];
    if (!first) return;
    setSortBy(first.id as (typeof validSortBy)[number]);
    setSortOrder(first.desc ? "desc" : "asc");
    setPage(1);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
            {t("locationsPage.title")}
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">
            {t("locationsPage.subtitle")}
          </p>
        </div>
        {canCreate && <CreateLocationDialog />}
      </div>

      <div className="flex flex-wrap items-center gap-4">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
            {t("locationsPage.site")}:
          </span>
          <SearchableSelect
            value={siteId || ""}
            onValueChange={(value) => {
              setSiteId(value === "" ? null : value);
              onSiteSelectChange(value);
              setPage(1);
            }}
            options={[
              { value: "", label: t("locationsPage.allSites") },
              ...siteSearchOptions,
            ]}
            onSearch={onSiteSearch}
            isLoading={siteSearchLoading}
            selectedLabel={siteId ? siteSelectedLabel : undefined}
            placeholder={t("locationsPage.allSites")}
            disabled={isLoading}
            className="min-w-[200px]"
          />
        </div>
      </div>

      {/* Locations Table */}
      <div className="rounded-lg border bg-white dark:bg-gray-800 dark:border-gray-700">
        {error ? (
          <div className="p-8 text-center text-red-600 dark:text-red-400">
            {t("locationsPage.errorLoadingLocations", {
              message: error.message,
            })}
          </div>
        ) : (
          <>
            <LocationDataTable
              columns={columns}
              data={locations}
              sorting={sorting}
              onSortingChange={handleSortingChange}
              isLoading={isLoading}
            />

            <PaginationFooter
              page={page}
              totalPages={totalPages}
              totalCount={pagination?.count ?? 0}
              onPageChange={setPage}
              itemLabel={t("topNav.locations").toLowerCase()}
              limit={limit}
              onLimitChange={(v) => {
                setLimit(v);
                setPage(1);
              }}
            />
          </>
        )}
      </div>
    </div>
  );
}

export default function LocationsPage() {
  const { t } = useTranslation();

  return (
    <Suspense
      fallback={
        <div className="p-8 text-center text-gray-600 dark:text-gray-400">
          {t("common.loading")}
        </div>
      }
    >
      <LocationsContent />
    </Suspense>
  );
}
