"use client";

import { type SortingState } from "@tanstack/react-table";
import { ChevronDown } from "lucide-react";
import {
  useQueryState,
  parseAsString,
  parseAsArrayOf,
  parseAsInteger,
  parseAsStringEnum,
} from "nuqs";
import { Suspense, useEffect, useMemo } from "react";
import { useTranslation } from "react-i18next";

import { useEpisodeCollectionStatusLabel } from "@/lib/hooks/use-status-labels";
import { DEFAULT_PAGE_SIZE } from "@/lib/pagination";
import { parseEpisodeCollectionStatus } from "@/lib/status/utils";

import {
  DateRangePicker,
  type DateRange,
} from "@/components/ui/date-range-picker";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { SearchableSelect } from "@/components/ui/searchable-select";

import { useLocationsQuery } from "@/features/locations";
import { useRobotsQuery, useRobotSearchOptions } from "@/features/robots";
import {
  useTasksQuery,
  useTaskVersionsQuery,
  useTaskSearchOptions,
} from "@/features/tasks";
import { usePermission } from "@/features/users";
import { useUsersQuery, useUserSearchOptions } from "@/features/users";

import { CreateEpisodeDialog } from "./create-episode-dialog";
import { getEpisodeColumns } from "./episode-columns";
import { EpisodeDataTable } from "./episode-data-table";
import { ExportMenu } from "./export-menu";
import { episodeStatusOptions } from "../constants";
import { useEpisodesListStream } from "../hooks/use-episodes-list-stream";
import { useEpisodesQuery } from "../hooks/use-episodes-query";

const validSortBy = [
  "task",
  "robot",
  "recorded_by",
  "started_at",
  "ended_at",
  "error",
] as const;
const validSortOrder = ["asc", "desc"] as const;

function EpisodesContent() {
  const { t } = useTranslation();
  const getEpisodeCollectionStatusLabel = useEpisodeCollectionStatusLabel();
  const [taskId, setTaskId] = useQueryState("task_id", parseAsString);

  const [statuses, setStatuses] = useQueryState(
    "status",
    parseAsArrayOf(parseAsInteger).withDefault([])
  );

  const [robotId, setRobotId] = useQueryState("robot_id", parseAsString);
  const [taskVersionId, setTaskVersionId] = useQueryState(
    "task_version_id",
    parseAsString
  );
  const [userId, setUserId] = useQueryState("user_id", parseAsString);
  const [startedAtFrom, setStartedAtFrom] = useQueryState(
    "started_at_from",
    parseAsString
  );
  const [startedAtTo, setStartedAtTo] = useQueryState(
    "started_at_to",
    parseAsString
  );

  const [page, setPage] = useQueryState("page", parseAsInteger.withDefault(1));
  const [limit, setLimit] = useQueryState(
    "limit",
    parseAsInteger.withDefault(DEFAULT_PAGE_SIZE)
  );

  const [sortBy, setSortBy] = useQueryState(
    "sort_by",
    parseAsStringEnum([...validSortBy])
  );
  const [sortOrder, setSortOrder] = useQueryState(
    "sort_order",
    parseAsStringEnum([...validSortOrder])
  );

  const parsedStatuses = (statuses || [])
    .map((s) => parseEpisodeCollectionStatus(String(s)))
    .filter((s): s is NonNullable<typeof s> => s !== undefined);

  // Date range is atomic: both bounds must be present to apply the filter.
  // A one-sided URL param is treated as no filter so the picker (which only
  // shows a value when both bounds exist) stays consistent with the result set.
  const hasBothDateBounds = !!(startedAtFrom && startedAtTo);
  const { data, isLoading, error } = useEpisodesQuery({
    task_id: taskId || undefined,
    task_version_id: taskVersionId || undefined,
    robot_id: robotId || undefined,
    user_id: userId || undefined,
    status: parsedStatuses.length > 0 ? parsedStatuses : undefined,
    started_at_from: hasBothDateBounds
      ? (startedAtFrom ?? undefined)
      : undefined,
    started_at_to: hasBothDateBounds ? (startedAtTo ?? undefined) : undefined,
    page,
    limit,
    sort_by: sortBy || undefined,
    sort_order: sortOrder || undefined,
  });

  useEpisodesListStream();

  // Async search for filter dropdowns
  const {
    options: taskSearchOptions,
    isLoading: taskSearchLoading,
    onSearch: onTaskSearch,
    selectedLabel: taskSelectedLabel,
    onValueChange: onTaskSelectChange,
  } = useTaskSearchOptions();
  const {
    options: userSearchOptions,
    isLoading: userSearchLoading,
    onSearch: onUserSearch,
    selectedLabel: userSelectedLabel,
    onValueChange: onUserSelectChange,
  } = useUserSearchOptions();
  const {
    options: robotSearchOptions,
    isLoading: robotSearchLoading,
    onSearch: onRobotSearch,
    selectedLabel: robotSelectedLabel,
    onValueChange: onRobotSelectChange,
  } = useRobotSearchOptions();
  // Name resolution for table display
  const { data: tasksData } = useTasksQuery({ limit: 1000 });
  const tasks = tasksData?.tasks;
  const {
    data: taskVersions,
    isLoading: isTaskVersionsLoading,
    isError: isTaskVersionsError,
  } = useTaskVersionsQuery(taskId ?? "", {
    enabled: !!taskId,
  });
  const { data: usersData } = useUsersQuery({ limit: 1000 });
  const users = usersData?.users;

  const canCreate = usePermission("episode:create");
  const canEdit = usePermission("episode:update");
  const { data: robotsData } = useRobotsQuery({ limit: 1000 });
  const robots = robotsData?.robots;
  const { data: locationsData } = useLocationsQuery({ limit: 1000 });
  const locations = locationsData?.locations;

  const taskNameById = useMemo(
    () => new Map((tasks || []).map((task) => [task.id, task.name])),
    [tasks]
  );
  const userNameById = useMemo(
    () =>
      new Map((users || []).map((user) => [user.user_id, user.display_name])),
    [users]
  );
  const robotNameById = useMemo(
    () => new Map((robots || []).map((robot) => [robot.id, robot.name])),
    [robots]
  );
  const locationNameById = useMemo(
    () =>
      new Map(
        (locations || []).map((location) => [location.id, location.name])
      ),
    [locations]
  );

  const columns = useMemo(
    () =>
      getEpisodeColumns({
        canEdit,
        taskNameById,
        robotNameById,
        userNameById,
        locationNameById,
        t,
      }),
    [canEdit, taskNameById, robotNameById, userNameById, locationNameById, t]
  );

  const dateRange: DateRange | undefined =
    startedAtFrom && startedAtTo
      ? { from: startedAtFrom, to: startedAtTo }
      : undefined;

  const exportInitialFilters = useMemo(
    () => ({
      taskId: taskId ?? undefined,
      taskLabel: taskId ? taskNameById.get(taskId) : undefined,
      taskVersionId: taskVersionId ?? undefined,
      taskVersionLabel:
        taskVersionId && taskVersions
          ? taskVersions.find((v) => v.id === taskVersionId)?.version
          : undefined,
      robotId: robotId ?? undefined,
      robotLabel: robotId ? robotNameById.get(robotId) : undefined,
      userId: userId ?? undefined,
      userLabel: userId ? userNameById.get(userId) : undefined,
      statuses: parsedStatuses.length > 0 ? parsedStatuses : undefined,
      startedAtFrom: hasBothDateBounds
        ? (startedAtFrom ?? undefined)
        : undefined,
      startedAtTo: hasBothDateBounds ? (startedAtTo ?? undefined) : undefined,
    }),
    [
      taskId,
      taskNameById,
      taskVersionId,
      taskVersions,
      robotId,
      robotNameById,
      userId,
      userNameById,
      parsedStatuses,
      startedAtFrom,
      startedAtTo,
      hasBothDateBounds,
    ]
  );

  const episodes = data?.episodes ?? [];
  const pagination = data?.pagination;
  const totalPages = pagination
    ? Math.ceil(pagination.count / pagination.limit)
    : 1;

  const hasTaskSelected = !!taskId;

  useEffect(() => {
    if (!hasTaskSelected && taskVersionId) {
      setTaskVersionId(null);
      return;
    }

    if (hasTaskSelected && taskVersionId) {
      if (isTaskVersionsLoading || isTaskVersionsError || !taskVersions) {
        return;
      }

      const hasSelectedVersion = taskVersions.some(
        (version) => version.id === taskVersionId
      );
      if (!hasSelectedVersion) {
        setTaskVersionId(null);
      }
    }
  }, [
    hasTaskSelected,
    taskVersionId,
    taskVersions,
    isTaskVersionsLoading,
    isTaskVersionsError,
    setTaskVersionId,
  ]);

  // Sorting state: URL <-> TanStack Table
  const sorting: SortingState = sortBy
    ? [{ id: sortBy, desc: sortOrder === "desc" }]
    : [];

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

    const first = newSorting[0]!;
    setSortBy(first.id as (typeof validSortBy)[number]);
    setSortOrder(first.desc ? "desc" : "asc");
    setPage(1);
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
          {t("episodesPage.title")}
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-2">
          {t("episodesPage.subtitle")}
        </p>
      </div>

      {/* Filters and Actions */}
      <div className="flex items-center justify-between gap-4">
        <div className="flex flex-wrap items-center gap-4">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
              {t("episodesPage.task")}:
            </span>
            <SearchableSelect
              value={taskId || ""}
              onValueChange={(value) => {
                setTaskId(value === "" ? null : value);
                onTaskSelectChange(value);
                setPage(1);
              }}
              options={[
                { value: "", label: t("episodesPage.allTasks") },
                ...taskSearchOptions,
              ]}
              onSearch={onTaskSearch}
              isLoading={taskSearchLoading}
              selectedLabel={taskId ? taskSelectedLabel : undefined}
              placeholder={t("episodesPage.allTasks")}
              disabled={isLoading}
              className="min-w-50"
            />
          </div>

          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
              {t("episodesPage.status")}:
            </span>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <button
                  type="button"
                  disabled={isLoading}
                  className="flex h-10 min-w-40 items-center justify-between rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white focus:outline-none focus:ring-2 focus:ring-gray-950 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 dark:border-gray-800 dark:bg-gray-950 dark:ring-offset-gray-950 dark:focus:ring-gray-300"
                >
                  <span className="truncate text-gray-500 dark:text-gray-400">
                    {(statuses || []).length === 0
                      ? t("episodesPage.allStatuses")
                      : t("common.selectedCount", {
                          count: (statuses || []).length,
                        })}
                  </span>
                  <ChevronDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                {episodeStatusOptions.map((s) => (
                  <DropdownMenuCheckboxItem
                    key={s}
                    checked={(statuses || []).includes(s)}
                    onSelect={(e) => e.preventDefault()}
                    onCheckedChange={(checked) => {
                      const currentStatuses = statuses || [];
                      setStatuses(
                        checked
                          ? [...currentStatuses, s]
                          : currentStatuses.filter((v) => v !== s)
                      );
                      setPage(1);
                    }}
                  >
                    {getEpisodeCollectionStatusLabel(s)}
                  </DropdownMenuCheckboxItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
              {t("episodesPage.robot")}:
            </span>
            <SearchableSelect
              value={robotId || ""}
              onValueChange={(value) => {
                setRobotId(value === "" ? null : value);
                onRobotSelectChange(value);
                setPage(1);
              }}
              options={[
                { value: "", label: t("episodesPage.allRobots") },
                ...robotSearchOptions,
              ]}
              onSearch={onRobotSearch}
              isLoading={robotSearchLoading}
              selectedLabel={robotId ? robotSelectedLabel : undefined}
              placeholder={t("episodesPage.allRobots")}
              disabled={isLoading}
              className="min-w-45"
            />
          </div>

          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
              {t("episodesPage.taskVersion")}:
            </span>
            <SearchableSelect
              value={taskVersionId || ""}
              onValueChange={(value) => {
                setTaskVersionId(value === "" ? null : value);
                setPage(1);
              }}
              options={[
                {
                  value: "",
                  label: hasTaskSelected
                    ? t("episodesPage.allTaskVersions")
                    : t("episodesPage.selectTaskFirst"),
                },
                ...(taskVersions || []).map((v) => ({
                  value: v.id,
                  label: v.version,
                })),
              ]}
              placeholder={
                hasTaskSelected
                  ? t("episodesPage.allTaskVersions")
                  : t("episodesPage.selectTaskFirst")
              }
              disabled={isLoading || !hasTaskSelected}
              className="min-w-55"
            />
          </div>

          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
              {t("episodesPage.user")}:
            </span>
            <SearchableSelect
              value={userId || ""}
              onValueChange={(value) => {
                setUserId(value === "" ? null : value);
                onUserSelectChange(value);
                setPage(1);
              }}
              options={[
                { value: "", label: t("episodesPage.allUsers") },
                ...userSearchOptions,
              ]}
              onSearch={onUserSearch}
              isLoading={userSearchLoading}
              selectedLabel={userId ? userSelectedLabel : undefined}
              placeholder={t("episodesPage.allUsers")}
              disabled={isLoading}
              className="min-w-45"
            />
          </div>

          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
              {t("episodesPage.dateRange")}:
            </span>
            <DateRangePicker
              value={dateRange}
              onChange={(range: DateRange) => {
                setStartedAtFrom(range.from);
                setStartedAtTo(range.to);
                setPage(1);
              }}
              onClear={() => {
                setStartedAtFrom(null);
                setStartedAtTo(null);
                setPage(1);
              }}
              disabled={isLoading}
            />
          </div>
        </div>

        <div className="flex items-center gap-2">
          <ExportMenu initialFilters={exportInitialFilters} />
          {canCreate && <CreateEpisodeDialog />}
        </div>
      </div>

      {/* Episodes Table */}
      {error ? (
        <div className="rounded-lg border bg-white dark:bg-gray-800 dark:border-gray-700 p-8 text-center text-red-600 dark:text-red-400">
          {t("episodesPage.errorLoadingEpisodes", { message: error.message })}
        </div>
      ) : (
        <EpisodeDataTable
          columns={columns}
          data={episodes}
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

export function EpisodeListPage() {
  const { t } = useTranslation();

  return (
    <Suspense
      fallback={
        <div className="p-8 text-center text-gray-600 dark:text-gray-400">
          {t("common.loading")}
        </div>
      }
    >
      <EpisodesContent />
    </Suspense>
  );
}
