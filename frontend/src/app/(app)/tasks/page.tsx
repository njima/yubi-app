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
import { Suspense, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import {
  useTaskStatusLabel,
  useTaskPriorityLabel,
} from "@/shared/hooks/use-status-labels";
import { DEFAULT_PAGE_SIZE } from "@/shared/lib/pagination";
import {
  TASK_STATUS,
  TASK_PRIORITY,
  TASK_DIFFICULTY,
} from "@/shared/lib/status-constants";
import {
  parseTaskStatus,
  parseTaskPriority,
  parseTaskDifficulty,
  getTaskDifficultyLabel,
} from "@/shared/lib/status-utils";

import { type DateRange } from "@/components/ui/date-range-picker";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { SearchableSelect } from "@/components/ui/searchable-select";

import { useRobotsQuery } from "@/features/robots";
import {
  useTasksQuery,
  CreateTaskDialog,
  ImportTasksDialog,
  ExportTasksDialog,
  TaskCompletionTrend,
  TaskSummaryCard,
  TaskDataTable,
  getTaskColumns,
} from "@/features/tasks";
import { usePermission } from "@/features/users";

const validSortBy = [
  "name",
  "robot_type",
  "priority",
  "difficulty",
  "target_duration_seconds",
  "status",
] as const;
const validSortOrder = ["asc", "desc"] as const;

const statusOptions = [
  TASK_STATUS.PLANNING,
  TASK_STATUS.DOING,
  TASK_STATUS.COMPLETED,
  TASK_STATUS.CANCELED,
];

const priorityOptions = [
  TASK_PRIORITY.LOW,
  TASK_PRIORITY.NORMAL,
  TASK_PRIORITY.HIGH,
  TASK_PRIORITY.URGENT,
];

const difficultyOptions = [
  TASK_DIFFICULTY.S,
  TASK_DIFFICULTY.A,
  TASK_DIFFICULTY.B,
  TASK_DIFFICULTY.C,
];

function defaultSummaryDateRange(): DateRange {
  const now = new Date();
  const from = new Date(now.getFullYear(), now.getMonth() - 2, now.getDate());
  const to = new Date(now.getFullYear(), now.getMonth() + 2, now.getDate());
  const pad = (n: number) => String(n).padStart(2, "0");
  const fmt = (d: Date) =>
    `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
  return { from: fmt(from), to: fmt(to) };
}

function TasksContent() {
  const { t } = useTranslation();
  const getTaskStatusLabel = useTaskStatusLabel();
  const getTaskPriorityLabel = useTaskPriorityLabel();
  // Summary filters (independent from task list filters)
  const [summaryCategoryTypeId, setSummaryCategoryTypeId] = useState<
    string | null
  >(null);
  const [summaryTagIds, setSummaryTagIds] = useState<string[]>([]);
  const [summaryRobotTypes, setSummaryRobotTypes] = useState<string[]>([]);
  const [summaryInterval, setSummaryInterval] = useState("2week");
  const [summaryDateRange, setSummaryDateRange] = useState<DateRange>(
    defaultSummaryDateRange
  );

  const [sortBy, setSortBy] = useQueryState(
    "sort_by",
    parseAsStringEnum([...validSortBy])
  );
  const [sortOrder, setSortOrder] = useQueryState(
    "sort_order",
    parseAsStringEnum([...validSortOrder])
  );
  const [statuses, setStatuses] = useQueryState(
    "status",
    parseAsArrayOf(parseAsInteger).withDefault([])
  );
  const [priorities, setPriorities] = useQueryState(
    "priority",
    parseAsArrayOf(parseAsInteger).withDefault([])
  );
  const [difficulties, setDifficulties] = useQueryState(
    "difficulty",
    parseAsArrayOf(parseAsInteger).withDefault([])
  );
  const [robotType, setRobotType] = useQueryState("robot_type", parseAsString);
  const [page, setPage] = useQueryState("page", parseAsInteger.withDefault(1));
  const [limit, setLimit] = useQueryState(
    "limit",
    parseAsInteger.withDefault(DEFAULT_PAGE_SIZE)
  );

  const parsedStatuses = statuses
    .map((s) => parseTaskStatus(String(s)))
    .filter((s): s is NonNullable<typeof s> => s !== undefined);
  const parsedPriorities = priorities
    .map((p) => parseTaskPriority(String(p)))
    .filter((p): p is NonNullable<typeof p> => p !== undefined);
  const parsedDifficulties = difficulties
    .map((d) => parseTaskDifficulty(String(d)))
    .filter((d): d is NonNullable<typeof d> => d !== undefined);

  const {
    data: tasksData,
    isLoading,
    error,
  } = useTasksQuery({
    sort_by: sortBy ?? undefined,
    sort_order: sortOrder ?? undefined,
    status: parsedStatuses.length > 0 ? parsedStatuses : undefined,
    priority: parsedPriorities.length > 0 ? parsedPriorities : undefined,
    difficulty: parsedDifficulties.length > 0 ? parsedDifficulties : undefined,
    robot_type: robotType ?? undefined,
    page,
    limit,
  });
  const tasks = tasksData?.tasks ?? [];
  const pagination = tasksData?.pagination;
  const totalPages = pagination
    ? Math.ceil(pagination.count / pagination.limit)
    : 1;

  // Robot model dynamic list
  const { data: robotsData } = useRobotsQuery({ limit: 1000 });
  const robots = robotsData?.robots ?? [];
  const uniqueRobotTypes = [
    ...new Set(robots.map((r) => r.robot_type).filter((m): m is string => !!m)),
  ].sort();

  const canCreate = usePermission("task:create");
  const canEdit = usePermission("task:update");

  const columns = useMemo(() => getTaskColumns({ canEdit, t }), [canEdit, t]);

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
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
          {t("tasksPage.title")}
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-2">
          {t("tasksPage.subtitle")}
        </p>
      </div>

      {/* Task Summary */}
      <TaskSummaryCard
        categoryTypeId={summaryCategoryTypeId}
        setCategoryTypeId={setSummaryCategoryTypeId}
        tagIds={summaryTagIds}
        setTagIds={setSummaryTagIds}
        robotTypes={summaryRobotTypes}
        setRobotTypes={setSummaryRobotTypes}
        dateRange={summaryDateRange}
        setDateRange={setSummaryDateRange}
        interval={summaryInterval}
        onIntervalChange={setSummaryInterval}
      />

      {/* Completion Trend */}
      <TaskCompletionTrend
        categoryTypeId={summaryCategoryTypeId}
        tagIds={summaryTagIds}
        robotTypes={summaryRobotTypes}
        dateRange={summaryDateRange}
        interval={summaryInterval}
      />

      {/* Task List section */}
      <div className="border-t border-gray-200 dark:border-gray-700 pt-6">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
          {t("tasksPage.taskList")}
        </h2>

        {/* Filters and Actions */}
        <div className="flex items-center justify-between gap-4">
          <div className="flex flex-wrap items-center gap-4">
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
                {t("tasksPage.status")}:
              </span>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <button
                    type="button"
                    disabled={isLoading}
                    className="flex h-10 min-w-40 items-center justify-between rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white focus:outline-none focus:ring-2 focus:ring-gray-950 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 dark:border-gray-800 dark:bg-gray-950 dark:ring-offset-gray-950 dark:focus:ring-gray-300"
                  >
                    <span className="truncate text-gray-500 dark:text-gray-400">
                      {statuses.length === 0
                        ? t("tasksPage.allStatuses")
                        : t("common.selectedCount", { count: statuses.length })}
                    </span>
                    <ChevronDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  {statusOptions.map((s) => (
                    <DropdownMenuCheckboxItem
                      key={s}
                      checked={statuses.includes(s)}
                      onSelect={(e) => e.preventDefault()}
                      onCheckedChange={(checked) => {
                        setStatuses(
                          checked
                            ? [...statuses, s]
                            : statuses.filter((v) => v !== s)
                        );
                        setPage(1);
                      }}
                    >
                      {getTaskStatusLabel(s)}
                    </DropdownMenuCheckboxItem>
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            </div>

            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
                {t("tasksPage.priority")}:
              </span>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <button
                    type="button"
                    disabled={isLoading}
                    className="flex h-10 min-w-40 items-center justify-between rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white focus:outline-none focus:ring-2 focus:ring-gray-950 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 dark:border-gray-800 dark:bg-gray-950 dark:ring-offset-gray-950 dark:focus:ring-gray-300"
                  >
                    <span className="truncate text-gray-500 dark:text-gray-400">
                      {priorities.length === 0
                        ? t("tasksPage.allPriorities")
                        : t("common.selectedCount", {
                            count: priorities.length,
                          })}
                    </span>
                    <ChevronDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  {priorityOptions.map((p) => (
                    <DropdownMenuCheckboxItem
                      key={p}
                      checked={priorities.includes(p)}
                      onSelect={(e) => e.preventDefault()}
                      onCheckedChange={(checked) => {
                        setPriorities(
                          checked
                            ? [...priorities, p]
                            : priorities.filter((v) => v !== p)
                        );
                        setPage(1);
                      }}
                    >
                      {getTaskPriorityLabel(p)}
                    </DropdownMenuCheckboxItem>
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            </div>

            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
                {t("tasksPage.difficulty")}:
              </span>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <button
                    type="button"
                    disabled={isLoading}
                    className="flex h-10 min-w-40 items-center justify-between rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white focus:outline-none focus:ring-2 focus:ring-gray-950 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 dark:border-gray-800 dark:bg-gray-950 dark:ring-offset-gray-950 dark:focus:ring-gray-300"
                  >
                    <span className="truncate text-gray-500 dark:text-gray-400">
                      {difficulties.length === 0
                        ? t("tasksPage.allDifficulties")
                        : t("common.selectedCount", {
                            count: difficulties.length,
                          })}
                    </span>
                    <ChevronDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  {difficultyOptions.map((d) => (
                    <DropdownMenuCheckboxItem
                      key={d}
                      checked={difficulties.includes(d)}
                      onSelect={(e) => e.preventDefault()}
                      onCheckedChange={(checked) => {
                        setDifficulties(
                          checked
                            ? [...difficulties, d]
                            : difficulties.filter((v) => v !== d)
                        );
                        setPage(1);
                      }}
                    >
                      {getTaskDifficultyLabel(d)}
                    </DropdownMenuCheckboxItem>
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            </div>

            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
                {t("tasksPage.robotType")}:
              </span>
              <SearchableSelect
                value={robotType ?? ""}
                onValueChange={(value) => {
                  setRobotType(value === "" ? null : value);
                  setPage(1);
                }}
                options={[
                  { value: "", label: t("tasksPage.allRobotTypes") },
                  ...uniqueRobotTypes.map((m) => ({ value: m, label: m })),
                ]}
                placeholder={t("tasksPage.allRobotTypes")}
                disabled={isLoading}
                className="min-w-40"
              />
            </div>
          </div>

          <div className="flex items-center gap-2">
            <ExportTasksDialog />
            {canCreate && (
              <>
                <ImportTasksDialog />
                <CreateTaskDialog />
              </>
            )}
          </div>
        </div>

        {/* Tasks Table */}
        <div className="mt-4">
          {error ? (
            <div className="rounded-lg border bg-white dark:bg-gray-800 dark:border-gray-700 p-8 text-center text-red-600 dark:text-red-400">
              {t("tasksPage.errorLoadingTasks", { message: error.message })}
            </div>
          ) : (
            <TaskDataTable
              columns={columns}
              data={tasks}
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
      </div>
    </div>
  );
}

export default function TasksPage() {
  const { t } = useTranslation();

  return (
    <Suspense
      fallback={
        <div className="p-8 text-center text-gray-600 dark:text-gray-400">
          {t("common.loading")}
        </div>
      }
    >
      <TasksContent />
    </Suspense>
  );
}
