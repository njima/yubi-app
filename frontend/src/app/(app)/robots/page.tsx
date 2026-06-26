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

import {
  useFormatAbsoluteTime,
  useFormatDistanceTime,
} from "@/shared/hooks/use-date-formatters";
import { DEFAULT_PAGE_SIZE } from "@/shared/lib/pagination";
import { ROBOT_STATUS, USER_ROLE } from "@/shared/lib/status-constants";
import { PaginationFooter } from "@/shared/ui/pagination-footer";
import { SearchableSelect } from "@/shared/ui/searchable-select";

import {
  useLocationsQuery,
  useLocationSearchOptions,
} from "@/features/locations";
import {
  useRobotsQuery,
  useRobotTypesQuery,
  CreateRobotDialog,
  FleetSummaryGrid,
  useRobotScope,
  isRobotInScope,
  useRobotsStatusStream,
  RobotsStatusProvider,
} from "@/features/robots";
import { getRobotColumns } from "@/features/robots/components/robot-columns";
import { RobotDataTable } from "@/features/robots/components/robot-data-table";
import { useSiteSearchOptions } from "@/features/sites";
import { usePermission } from "@/features/users";
import { useMeQuery } from "@/features/users";

const validSortBy = [
  "name",
  "location_id",
  "robot_type",
  "status",
  "leader_status",
  "last_heartbeat_at",
  "active_episode_id",
  "active_user_id",
] as const;
const validSortOrder = ["asc", "desc"] as const;

function RobotsContent() {
  const { t } = useTranslation();
  const formatDistanceTime = useFormatDistanceTime();
  const formatAbsoluteTime = useFormatAbsoluteTime();
  const [siteId, setSiteId] = useQueryState("site_id", parseAsString);
  const [locationId, setLocationId] = useQueryState(
    "location_id",
    parseAsString
  );
  const [status, setStatus] = useQueryState("status", parseAsInteger);
  const [robotType, setRobotType] = useQueryState("robot_type", parseAsString);
  const [sortBy, setSortBy] = useQueryState(
    "sort_by",
    parseAsStringEnum([...validSortBy])
  );
  const [sortOrder, setSortOrder] = useQueryState(
    "sort_order",
    parseAsStringEnum([...validSortOrder])
  );

  const [page, setPage] = useQueryState("page", parseAsInteger.withDefault(1));
  const [limit, setLimit] = useQueryState(
    "limit",
    parseAsInteger.withDefault(DEFAULT_PAGE_SIZE)
  );

  const parsedStatus =
    status === ROBOT_STATUS.ONLINE ||
    status === ROBOT_STATUS.BUSY ||
    status === ROBOT_STATUS.OFFLINE ||
    status === ROBOT_STATUS.FAULTED ||
    status === ROBOT_STATUS.MAINTENANCE
      ? status
      : undefined;

  const {
    data: robotsData,
    isLoading,
    error,
  } = useRobotsQuery({
    site_id: siteId || undefined,
    location_id: locationId || undefined,
    status: parsedStatus,
    robot_type: robotType || undefined,
    page,
    limit,
    sort_by: sortBy ?? undefined,
    sort_order: sortOrder ?? undefined,
  });
  const robots = useMemo(() => robotsData?.robots ?? [], [robotsData?.robots]);

  // Scope + user
  const { scopeIds } = useRobotScope();
  const { data: me } = useMeQuery();

  // Filter robots for OPERATOR role
  const filteredRobots = useMemo(
    () =>
      me?.role === USER_ROLE.OPERATOR
        ? robots.filter((robot) => isRobotInScope(robot.id, scopeIds))
        : robots,
    [robots, scopeIds, me?.role]
  );

  const filteredRobotIds = useMemo(
    () => filteredRobots.map((r) => r.id),
    [filteredRobots]
  );

  // Fetch realtime status for filtered robots
  // This hook returns a map of { robot_id: realtimeStatus } and connection state.
  const { data: realtimeStatusForFilteredRobots, isConnected } =
    useRobotsStatusStream(filteredRobotIds);

  // Fetch distinct robot types for filter dropdown (filtered by current site/location/status)
  const { data: robotTypes } = useRobotTypesQuery({
    site_id: siteId || undefined,
    location_id: locationId || undefined,
    status: parsedStatus,
  });

  // Async search for Site filter dropdown
  const {
    options: siteSearchOptions,
    isLoading: siteSearchLoading,
    onSearch: onSiteSearch,
    selectedLabel: siteSelectedLabel,
    onValueChange: onSiteSelectChange,
  } = useSiteSearchOptions();

  // Async search for Location filter dropdown
  const {
    options: locationSearchOptions,
    isLoading: locationSearchLoading,
    onSearch: onLocationSearch,
    selectedLabel: locationSelectedLabel,
    onValueChange: onLocationSelectChange,
  } = useLocationSearchOptions({ site_id: siteId || undefined });

  // TODO: Technical debt - Server should JOIN and include location names in response
  // Currently making multiple API calls on client side as a temporary workaround
  const { data: locationsData } = useLocationsQuery({ limit: 1000 });
  const locations = locationsData?.locations;

  const canCreate = usePermission("robot:create");
  const canEdit = usePermission("robot:update");
  const canDelete = usePermission("robot:delete");

  const pagination = robotsData?.pagination;
  const totalPages = pagination
    ? Math.ceil(pagination.count / pagination.limit)
    : 1;

  const locationNameById = useMemo(
    () => new Map((locations ?? []).map((l) => [l.id, l.name])),
    [locations]
  );

  const columns = useMemo(
    () =>
      getRobotColumns({
        canEdit,
        canDelete,
        locationNameById,
        meRole: me?.role,
        isInScope: (robotId: string) =>
          me?.role === USER_ROLE.ADMIN || isRobotInScope(robotId, scopeIds),
        t,
        formatDistanceTime,
        formatAbsoluteTime,
      }),
    [
      canEdit,
      canDelete,
      locationNameById,
      me?.role,
      scopeIds,
      t,
      formatDistanceTime,
      formatAbsoluteTime,
    ]
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
      return;
    }

    const first = newSorting[0];
    if (!first) return;
    setSortBy(first.id as (typeof validSortBy)[number]);
    setSortOrder(first.desc ? "desc" : "asc");
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
          {t("robotsPage.title")}
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-2">
          {t("robotsPage.subtitle")}
        </p>
      </div>

      <FleetSummaryGrid />

      {/* Robot List section */}
      <div className="border-t border-gray-200 dark:border-gray-700 pt-6">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
          {t("robotsPage.robotList")}
        </h2>

        {/* Actions */}
        <div className="flex items-center justify-between gap-4">
          <div className="flex flex-wrap items-center gap-4">
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
                {t("robotsPage.site")}:
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
                options={[
                  { value: "", label: t("robotsPage.allSites") },
                  ...siteSearchOptions,
                ]}
                onSearch={onSiteSearch}
                isLoading={siteSearchLoading}
                selectedLabel={siteId ? siteSelectedLabel : undefined}
                placeholder={t("robotsPage.allSites")}
                disabled={isLoading}
                className="min-w-[200px]"
              />
            </div>

            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
                {t("robotsPage.location")}:
              </span>
              <SearchableSelect
                value={locationId || ""}
                onValueChange={(value) => {
                  setLocationId(value === "" ? null : value);
                  onLocationSelectChange(value);
                  setPage(1);
                }}
                options={[
                  { value: "", label: t("robotsPage.allLocations") },
                  ...locationSearchOptions,
                ]}
                onSearch={onLocationSearch}
                isLoading={locationSearchLoading}
                selectedLabel={locationId ? locationSelectedLabel : undefined}
                placeholder={t("robotsPage.allLocations")}
                disabled={isLoading}
                className="min-w-50"
              />
            </div>

            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
                {t("robotsPage.followerStatus")}:
              </span>
              <SearchableSelect
                value={status !== null ? String(status) : ""}
                onValueChange={(value) => {
                  setStatus(value === "" ? null : parseInt(value, 10));
                  setPage(1);
                }}
                options={[
                  { value: "", label: t("robotsPage.allStatuses") },
                  {
                    value: String(ROBOT_STATUS.ONLINE),
                    label: t("status.online"),
                  },
                  { value: String(ROBOT_STATUS.BUSY), label: t("status.busy") },
                  {
                    value: String(ROBOT_STATUS.OFFLINE),
                    label: t("status.offline"),
                  },
                  {
                    value: String(ROBOT_STATUS.FAULTED),
                    label: t("status.faulted"),
                  },
                  {
                    value: String(ROBOT_STATUS.MAINTENANCE),
                    label: t("status.maintenance"),
                  },
                ]}
                placeholder={t("robotsPage.allStatuses")}
                disabled={isLoading}
                className="min-w-40"
              />
            </div>

            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
                {t("robotsPage.robotType")}:
              </span>
              <SearchableSelect
                value={robotType || ""}
                onValueChange={(value) => {
                  setRobotType(value === "" ? null : value);
                  setPage(1);
                }}
                options={[
                  { value: "", label: t("robotsPage.allRobotTypes") },
                  ...(robotTypes ?? []).map((m) => ({ value: m, label: m })),
                ]}
                placeholder={t("robotsPage.allRobotTypes")}
                disabled={isLoading}
                className="min-w-40"
              />
            </div>
          </div>

          {canCreate && <CreateRobotDialog />}
        </div>

        {/* Robots Table */}
        <div className="mt-4 rounded-lg border bg-white dark:bg-gray-800 dark:border-gray-700">
          {error ? (
            <div className="p-8 text-center text-red-600 dark:text-red-400">
              {t("robotsPage.errorLoadingRobots", { message: error.message })}
            </div>
          ) : (
            <>
              <RobotsStatusProvider
                statusMap={realtimeStatusForFilteredRobots}
                isConnected={isConnected}
              >
                <RobotDataTable
                  columns={columns}
                  data={filteredRobots}
                  sorting={sorting}
                  onSortingChange={handleSortingChange}
                  isLoading={isLoading}
                />
              </RobotsStatusProvider>

              <PaginationFooter
                page={page}
                totalPages={totalPages}
                totalCount={pagination?.count ?? 0}
                onPageChange={setPage}
                itemLabel={t("topNav.robots").toLowerCase()}
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
    </div>
  );
}

export default function RobotsPage() {
  const { t } = useTranslation();

  return (
    <Suspense
      fallback={
        <div className="p-8 text-center text-gray-600 dark:text-gray-400">
          {t("common.loading")}
        </div>
      }
    >
      <RobotsContent />
    </Suspense>
  );
}
