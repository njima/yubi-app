"use client";

import { type SortingState } from "@tanstack/react-table";
import {
  useQueryState,
  parseAsString,
  parseAsInteger,
  parseAsStringEnum,
} from "nuqs";
import { Suspense, useMemo } from "react";
import { useTranslation } from "react-i18next";

import { useFormatRelativeTime } from "@/shared/hooks/use-date-formatters";
import { useUserRoleLabel } from "@/shared/hooks/use-status-labels";
import { DEFAULT_PAGE_SIZE } from "@/shared/lib/pagination";
import { SearchableSelect } from "@/shared/ui/searchable-select";

import { useLocationSearchOptions } from "@/features/locations";
import { useSiteSearchOptions } from "@/features/sites";
import { usePermission } from "@/features/users";
import { useMeQuery, useUsersQuery } from "@/features/users";
import { CreateUserDialog } from "@/features/users/components/create-user-dialog";
import { ImportUsersDialog } from "@/features/users/components/import-users-dialog";
import { getUserColumns } from "@/features/users/components/user-columns";
import { UserDataTable } from "@/features/users/components/user-data-table";

const validSortBy = [
  "name",
  "email",
  "role",
  "location",
  "created_at",
] as const;
const validSortOrder = ["asc", "desc"] as const;

function UsersContent() {
  const { t } = useTranslation();
  const [page, setPage] = useQueryState("page", parseAsInteger.withDefault(1));
  const [limit, setLimit] = useQueryState(
    "limit",
    parseAsInteger.withDefault(DEFAULT_PAGE_SIZE)
  );
  const [siteId, setSiteId] = useQueryState("site_id", parseAsString);
  const [locationId, setLocationId] = useQueryState(
    "location_id",
    parseAsString
  );
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

  const { data, isLoading, error } = useUsersQuery({
    page,
    limit,
    site_id: siteId || undefined,
    location_id: locationId || undefined,
    sort_by: sortBy ?? undefined,
    sort_order: sortOrder ?? undefined,
  });
  const users = data?.users ?? [];
  const pagination = data?.pagination;
  const totalPages = pagination
    ? Math.ceil(pagination.count / pagination.limit)
    : 1;
  const { data: me } = useMeQuery();
  const {
    options: locationSearchOptions,
    isLoading: locationSearchLoading,
    onSearch: onLocationSearch,
    selectedLabel: locationSelectedLabel,
    onValueChange: onLocationSelectChange,
  } = useLocationSearchOptions({ site_id: siteId || undefined });
  const canUpdateRole = usePermission("user:update_role");
  const canCreateUser = usePermission("user:create");
  const formatRelativeTime = useFormatRelativeTime();
  const getRoleLabel = useUserRoleLabel();

  const columns = useMemo(
    () =>
      getUserColumns({
        canUpdateRole,
        currentUser: me,
        t,
        formatRelativeTime,
        getRoleLabel,
      }),
    [canUpdateRole, me, t, formatRelativeTime, getRoleLabel]
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
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
            {t("usersPage.title")}
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">
            {t("usersPage.subtitle")}
          </p>
        </div>
        {canCreateUser && (
          <div className="flex gap-2">
            <ImportUsersDialog />
            <CreateUserDialog />
          </div>
        )}
      </div>

      {/* Filters */}
      <div className="flex flex-wrap items-center gap-4">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
            {t("usersPage.site")}:
          </span>
          <SearchableSelect
            value={siteId || ""}
            onValueChange={(value) => {
              setSiteId(value === "" ? null : value);
              onSiteSelectChange(value);
              setLocationId(null);
              onLocationSelectChange("");
              setPage(1);
            }}
            options={[{ value: "", label: "All Sites" }, ...siteSearchOptions]}
            onSearch={onSiteSearch}
            isLoading={siteSearchLoading}
            selectedLabel={siteId ? siteSelectedLabel : undefined}
            placeholder="All Sites"
            disabled={isLoading}
            className="min-w-[200px]"
          />
        </div>

        <div className="flex items-center gap-2">
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
            {t("usersPage.location")}:
          </span>
          <SearchableSelect
            value={locationId || ""}
            onValueChange={(value) => {
              setLocationId(value === "" ? null : value);
              onLocationSelectChange(value);
              setPage(1);
            }}
            options={[
              { value: "", label: t("usersPage.allLocations") },
              ...locationSearchOptions,
            ]}
            onSearch={onLocationSearch}
            isLoading={locationSearchLoading}
            selectedLabel={locationId ? locationSelectedLabel : undefined}
            placeholder={t("usersPage.allLocations")}
            disabled={isLoading}
            className="min-w-50"
          />
        </div>
      </div>

      {/* Users Table */}
      {error ? (
        <div className="p-8 text-center text-red-600 dark:text-red-400">
          {t("usersPage.errorLoadingUsers", { message: error.message })}
        </div>
      ) : (
        <UserDataTable
          columns={columns}
          data={users}
          sorting={sorting}
          onSortingChange={handleSortingChange}
          isLoading={isLoading}
          totalCount={pagination?.count ?? 0}
          page={page}
          totalPages={totalPages}
          onPageChange={setPage}
          limit={limit}
          onLimitChange={(v) => {
            setLimit(v);
            setPage(1);
          }}
        />
      )}
    </div>
  );
}

export default function UsersPage() {
  const { t } = useTranslation();

  return (
    <Suspense
      fallback={
        <div className="p-8 text-center text-gray-600 dark:text-gray-400">
          {t("common.loading")}
        </div>
      }
    >
      <UsersContent />
    </Suspense>
  );
}
