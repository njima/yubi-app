"use client";

import { ChevronDown, ListTodo, Clock, Film } from "lucide-react";
import { useTranslation } from "react-i18next";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
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

import { useRobotsQuery } from "@/features/robots";

import { useTaskAvailableTagsQuery } from "../hooks/use-task-available-tags-query";
import { useTaskSummaryQuery } from "../hooks/use-task-summary-query";
import { useTaskCategoryTypesQuery } from "../hooks/use-task-tags-query";
import { secondsToHoursMinutes } from "../lib/duration";

interface TaskSummaryCardProps {
  categoryTypeId: string | null;
  setCategoryTypeId: (v: string | null) => void;
  tagIds: string[];
  setTagIds: (v: string[]) => void;
  robotTypes: string[];
  setRobotTypes: (v: string[]) => void;
  dateRange: DateRange;
  setDateRange: (range: DateRange) => void;
  interval: string;
  onIntervalChange: (interval: string) => void;
}

export function TaskSummaryCard({
  categoryTypeId,
  setCategoryTypeId,
  tagIds,
  setTagIds,
  robotTypes,
  setRobotTypes,
  dateRange,
  setDateRange,
  interval,
  onIntervalChange,
}: TaskSummaryCardProps) {
  const { t } = useTranslation();
  const { data: categoryTypes = [] } = useTaskCategoryTypesQuery();
  const { data: robotsData } = useRobotsQuery({ limit: 1000 });
  const robots = robotsData?.robots ?? [];
  const uniqueRobotTypes = [
    ...new Set(robots.map((r) => r.robot_type).filter((m): m is string => !!m)),
  ].sort();

  const { data: availableTags = [] } = useTaskAvailableTagsQuery({
    robot_type: robotTypes.length > 0 ? robotTypes : undefined,
    category_type_id: categoryTypeId ?? undefined,
  });

  const { data: summary } = useTaskSummaryQuery({
    robot_type: robotTypes.length > 0 ? robotTypes : undefined,
    category_type_id: categoryTypeId ?? undefined,
    tag_id: tagIds.length > 0 ? tagIds : undefined,
    from: dateRange.from,
    to: dateRange.to,
  });

  const targetHours = summary
    ? secondsToHoursMinutes(summary.target_duration_seconds)
    : null;

  const intervalOptions = [
    { value: "1week", label: t("taskSummaryCard.oneWeek") },
    { value: "2week", label: t("taskSummaryCard.twoWeeks") },
    { value: "month", label: t("taskSummaryCard.month") },
  ];

  return (
    <div className="space-y-4">
      {/* Section Header + Filters */}
      <div className="flex flex-wrap items-center justify-between gap-4">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
          {t("taskSummaryCard.title")}
        </h2>
        <div className="flex flex-wrap items-center gap-3">
          <DateRangePicker value={dateRange} onChange={setDateRange} />

          <SearchableSelect
            value={interval}
            onValueChange={onIntervalChange}
            options={intervalOptions}
            className="w-[120px]"
          />

          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
              {t("tasksPage.robotType")}:
            </span>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <button
                  type="button"
                  className="flex h-10 min-w-40 items-center justify-between rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white focus:outline-none focus:ring-2 focus:ring-gray-950 focus:ring-offset-2 dark:border-gray-800 dark:bg-gray-950 dark:ring-offset-gray-950 dark:focus:ring-gray-300"
                >
                  <span className="truncate text-gray-500 dark:text-gray-400">
                    {robotTypes.length === 0
                      ? t("tasksPage.allRobotTypes")
                      : t("common.selectedCount", { count: robotTypes.length })}
                  </span>
                  <ChevronDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                {uniqueRobotTypes.map((robotType) => (
                  <DropdownMenuCheckboxItem
                    key={robotType}
                    checked={robotTypes.includes(robotType)}
                    onSelect={(e) => e.preventDefault()}
                    onCheckedChange={(checked) => {
                      setRobotTypes(
                        checked
                          ? [...robotTypes, robotType]
                          : robotTypes.filter((m) => m !== robotType)
                      );
                      setTagIds([]);
                    }}
                  >
                    {robotType}
                  </DropdownMenuCheckboxItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
              {t("taskSummaryCard.category")}:
            </span>
            <SearchableSelect
              value={categoryTypeId ?? ""}
              onValueChange={(v) => {
                setCategoryTypeId(v === "" ? null : v);
                setTagIds([]);
              }}
              options={[
                { value: "", label: t("taskSummaryCard.all") },
                ...categoryTypes.map((ct) => ({
                  value: ct.id,
                  label: ct.name,
                })),
              ]}
              placeholder={t("taskSummaryCard.all")}
              className="min-w-[140px]"
            />
          </div>

          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
              {t("taskSummaryCard.tags")}:
            </span>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <button
                  type="button"
                  className="flex h-10 min-w-40 items-center justify-between rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white focus:outline-none focus:ring-2 focus:ring-gray-950 focus:ring-offset-2 dark:border-gray-800 dark:bg-gray-950 dark:ring-offset-gray-950 dark:focus:ring-gray-300"
                >
                  <span className="truncate text-gray-500 dark:text-gray-400">
                    {tagIds.length === 0
                      ? t("taskSummaryCard.allTags")
                      : t("common.selectedCount", { count: tagIds.length })}
                  </span>
                  <ChevronDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                {availableTags.map((tag) => (
                  <DropdownMenuCheckboxItem
                    key={tag.id}
                    checked={tagIds.includes(tag.id)}
                    onSelect={(e) => e.preventDefault()}
                    onCheckedChange={(checked) => {
                      setTagIds(
                        checked
                          ? [...tagIds, tag.id]
                          : tagIds.filter((id) => id !== tag.id)
                      );
                    }}
                  >
                    {tag.name}
                  </DropdownMenuCheckboxItem>
                ))}
                {availableTags.length === 0 && (
                  <div className="px-2 py-1.5 text-sm text-gray-500">
                    {t("taskSummaryCard.noTagsAvailable")}
                  </div>
                )}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>

      {/* Stat Cards */}
      <div className="grid grid-cols-3 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle>{t("taskSummaryCard.totalTasks")}</CardTitle>
            <ListTodo className="h-7 w-7 shrink-0 text-gray-400 dark:text-gray-500" />
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              {summary?.total_tasks ?? "-"}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle>{t("taskSummaryCard.targetHours")}</CardTitle>
            <Clock className="h-7 w-7 shrink-0 text-gray-400 dark:text-gray-500" />
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              {targetHours
                ? t("taskSummaryCard.targetHoursValue", {
                    hours: targetHours.hours,
                    minutes: targetHours.minutes,
                  })
                : "-"}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle>{t("taskSummaryCard.targetEpisodes")}</CardTitle>
            <Film className="h-7 w-7 shrink-0 text-gray-400 dark:text-gray-500" />
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              {summary?.target_episode_count?.toLocaleString() ?? "-"}
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
