"use client";

import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  ComposedChart,
  Bar,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  type TooltipContentProps,
} from "recharts";

import { CHART_PALETTE } from "@/lib/chart-colors";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { type DateRange } from "@/components/ui/date-range-picker";
import { SearchableSelect } from "@/components/ui/searchable-select";

import { useTaskCompletionTrendQuery } from "../hooks/use-task-completion-trend-query";
import { secondsToHoursMinutes } from "../lib/duration";

type Metric = "duration" | "tasks" | "episodes";

interface TaskCompletionTrendProps {
  categoryTypeId: string | null;
  tagIds: string[];
  robotTypes: string[];
  dateRange: DateRange;
  interval: string;
}

function formatPeriodLabel(
  start?: string,
  end?: string,
  overdueLabel?: string,
  locale?: string
): string {
  if (!start) return overdueLabel ?? "Overdue";
  const s = new Date(start);
  const e = end ? new Date(end) : null;
  const isJa = locale?.startsWith("ja");
  const fmt = (d: Date) =>
    isJa
      ? `${String(d.getFullYear()).slice(2)}/${String(d.getMonth() + 1)}/${String(d.getDate())}`
      : `${d.getMonth() + 1}/${d.getDate()}/${String(d.getFullYear()).slice(2)}`;
  if (e) {
    const MS_PER_DAY = 24 * 60 * 60 * 1000;
    const lastDay = new Date(e.getTime() - MS_PER_DAY);
    return `${fmt(s)}~${fmt(lastDay)}`;
  }
  return fmt(s);
}

function formatDuration(
  seconds: number,
  hoursLabel: string,
  minutesLabel: string
): string {
  const { hours, minutes } = secondsToHoursMinutes(seconds);
  return `${hours}${hoursLabel} ${minutes}${minutesLabel}`;
}

function getMetricValue(
  group: {
    target_duration: number;
    actual_duration: number;
    target_tasks: number;
    actual_tasks: number;
    target_episodes: number;
    actual_episodes: number;
  },
  metric: Metric,
  type: "target" | "actual"
): number {
  if (type === "target") {
    return metric === "duration"
      ? group.target_duration
      : metric === "tasks"
        ? group.target_tasks
        : group.target_episodes;
  }
  return metric === "duration"
    ? group.actual_duration
    : metric === "tasks"
      ? group.actual_tasks
      : group.actual_episodes;
}

export function TaskCompletionTrend({
  categoryTypeId,
  tagIds,
  robotTypes,
  dateRange,
  interval,
}: TaskCompletionTrendProps) {
  const { t, i18n } = useTranslation();
  const language = i18n.language;
  const [metric, setMetric] = useState<Metric>("duration");
  const [groupBy, setGroupBy] = useState("category");

  const metricOptions = useMemo(
    () => [
      { value: "duration", label: t("taskCompletionTrend.hours") },
      { value: "tasks", label: t("taskCompletionTrend.tasks") },
      { value: "episodes", label: t("taskCompletionTrend.episodes") },
    ],
    [t]
  );

  const groupByOptions = useMemo(
    () => [
      { value: "category", label: t("taskCompletionTrend.category") },
      { value: "status", label: t("taskCompletionTrend.status") },
    ],
    [t]
  );

  const { data: trend, isLoading } = useTaskCompletionTrendQuery({
    group_by: groupBy,
    robot_type: robotTypes.length > 0 ? robotTypes : undefined,
    category_type_id: categoryTypeId ?? undefined,
    tag_id: tagIds.length > 0 ? tagIds : undefined,
    from: dateRange.from,
    to: dateRange.to,
    interval,
  });

  const allLabels = useMemo(() => {
    if (!trend) return [];
    const labelSet = new Set<string>();
    for (const period of trend.periods) {
      for (const group of period.groups) {
        labelSet.add(group.label);
      }
    }
    return Array.from(labelSet).sort();
  }, [trend]);

  const overdueLabel = t("taskCompletionTrend.overdue");

  const chartData = useMemo(() => {
    if (!trend) return [];

    const regular: Record<string, string | number>[] = [];
    let overdueEntry: Record<string, string | number> | null = null;

    for (const period of trend.periods) {
      const isOverdue = !period.start;
      const entry: Record<string, string | number> = {
        name: isOverdue
          ? overdueLabel
          : formatPeriodLabel(period.start, period.end, overdueLabel, language),
      };

      let targetTotal = 0;
      for (const group of period.groups) {
        entry[`actual_${group.label}`] = getMetricValue(
          group,
          metric,
          "actual"
        );
        targetTotal += getMetricValue(group, metric, "target");
      }

      if (isOverdue) {
        // No target line for overdue — leave target_total undefined to break the line
        overdueEntry = entry;
      } else {
        entry.target_total = targetTotal;
        regular.push(entry);
      }
    }

    // Overdue at the end (right side of chart)
    if (overdueEntry) {
      regular.push(overdueEntry);
    }

    return regular;
  }, [trend, metric, overdueLabel, language]);

  const yAxisLabel =
    metric === "duration"
      ? t("taskCompletionTrend.hours")
      : metric === "tasks"
        ? t("taskCompletionTrend.tasks")
        : t("taskCompletionTrend.episodes");

  const formatYValue = (value: number): string => {
    if (metric === "duration") {
      return `${Math.round(value / 3600)}`;
    }
    return String(value);
  };

  const renderTooltip = ({ active, payload, label }: TooltipContentProps) => {
    if (!active || !payload || payload.length === 0) return null;
    const nonZero = payload.filter(
      (p) => typeof p.value === "number" && p.value > 0
    );
    if (nonZero.length === 0) return null;

    return (
      <div className="rounded-md border bg-white dark:bg-gray-800 dark:border-gray-700 p-2 shadow-sm text-xs max-w-64">
        <p className="font-semibold mb-1">{label}</p>
        {nonZero.map((p) => {
          const dataKey = String(p.dataKey ?? "");
          const value = Number(p.value ?? 0);
          const isTarget = dataKey === "target_total";
          const name = isTarget
            ? t("taskCompletionTrend.targetTotal")
            : dataKey.replace(/^actual_/, "");
          const val =
            metric === "duration"
              ? formatDuration(
                  value,
                  t("durationInput.hoursShort"),
                  t("durationInput.minutesShort")
                )
              : String(value);
          return (
            <div key={dataKey} className="flex items-center gap-2">
              <div
                className="w-2 h-2 rounded-full shrink-0"
                style={{
                  backgroundColor: isTarget
                    ? "#6b7280"
                    : (p.color ?? p.fill ?? "#6b7280"),
                }}
              />
              <span className="text-gray-600 dark:text-gray-300 truncate">
                {name}
              </span>
              <span className="font-medium ml-auto shrink-0">{val}</span>
            </div>
          );
        })}
      </div>
    );
  };

  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <CardTitle>{t("taskCompletionTrend.title")}</CardTitle>
          <div className="flex items-center gap-2">
            <SearchableSelect
              value={metric}
              onValueChange={(v) => setMetric(v as Metric)}
              options={metricOptions}
              className="w-[110px]"
            />
            <SearchableSelect
              value={groupBy}
              onValueChange={setGroupBy}
              options={groupByOptions}
              className="w-[110px]"
            />
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="h-[280px] flex items-center justify-center text-gray-500">
            {t("common.loading")}
          </div>
        ) : chartData.length === 0 ? (
          <div className="h-[280px] flex items-center justify-center text-gray-500">
            {t("taskCompletionTrend.noData")}
          </div>
        ) : (
          <>
            <ResponsiveContainer width="100%" height={280}>
              <ComposedChart data={chartData} barSize={56} barCategoryGap="15%">
                <XAxis
                  dataKey="name"
                  tick={{ fontSize: 10 }}
                  angle={-30}
                  textAnchor="end"
                  height={50}
                />
                <YAxis
                  tickFormatter={formatYValue}
                  tick={{ fontSize: 11 }}
                  width={45}
                  label={{
                    value: yAxisLabel,
                    angle: -90,
                    position: "insideLeft",
                    style: { fontSize: 11 },
                  }}
                />
                <Tooltip content={renderTooltip} />
                {/* Actual bars (stacked by category/status) */}
                {allLabels.map((label, i) => (
                  <Bar
                    key={`actual_${label}`}
                    dataKey={`actual_${label}`}
                    stackId="actual"
                    fill={CHART_PALETTE[i % CHART_PALETTE.length]}
                    radius={
                      i === allLabels.length - 1 ? [2, 2, 0, 0] : undefined
                    }
                  />
                ))}
                {/* Target line (dashed) */}
                <Line
                  dataKey="target_total"
                  stroke="#6b7280"
                  strokeWidth={2}
                  strokeDasharray="6 3"
                  dot={{ r: 3, fill: "#6b7280" }}
                  type="monotone"
                />
              </ComposedChart>
            </ResponsiveContainer>
            {/* Custom Legend */}
            <div className="flex flex-wrap items-center gap-x-4 gap-y-1 mt-1 text-xs text-gray-500 dark:text-gray-400">
              {allLabels.map((label, i) => (
                <div key={label} className="flex items-center gap-1.5">
                  <div
                    className="w-2.5 h-2.5 rounded-sm"
                    style={{
                      backgroundColor: CHART_PALETTE[i % CHART_PALETTE.length],
                    }}
                  />
                  <span>{label}</span>
                </div>
              ))}
              <div className="flex items-center gap-1.5">
                <svg width="16" height="2" className="shrink-0">
                  <line
                    x1="0"
                    y1="1"
                    x2="16"
                    y2="1"
                    stroke="#6b7280"
                    strokeWidth="2"
                    strokeDasharray="4 2"
                  />
                </svg>
                <span>{t("taskCompletionTrend.target")}</span>
              </div>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}
