"use client";

import { Bot, Clock, Database } from "lucide-react";
import { Suspense, useState } from "react";
import { useTranslation } from "react-i18next";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  DateRangePicker,
  type DateRange,
} from "@/components/ui/date-range-picker";

import {
  FleetStatsPanel,
  useFleetSummaryQuery,
  computeFleetTotals,
  useFleetStatsQuery,
  computeStatsTotals,
} from "@/features/robots";

// --- Helpers ---

function formatNumber(n: number): string {
  return n.toLocaleString();
}

function defaultRange(): DateRange {
  const now = new Date();
  const firstOfMonth = new Date(now.getFullYear(), now.getMonth(), 1);
  const lastOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0);
  const pad = (n: number) => String(n).padStart(2, "0");
  const fmt = (d: Date) =>
    `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
  return { from: fmt(firstOfMonth), to: fmt(lastOfMonth) };
}

// --- Component ---

function DashboardContent() {
  const { t } = useTranslation();
  const [range, setRange] = useState<DateRange>(defaultRange);
  const { data: summaryData } = useFleetSummaryQuery();
  const { data: statsData } = useFleetStatsQuery(range.from, range.to);

  const fleet = summaryData ? computeFleetTotals(summaryData) : null;
  const stats = statsData ? computeStatsTotals(statsData) : null;

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
          {t("dashboard.title")}
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-2">
          {t("dashboard.subtitle")}
        </p>
      </div>

      {/* Date Range Filter */}
      <DateRangePicker value={range} onChange={setRange} />

      {/* Summary Cards */}
      <div className="grid gap-3 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle>{t("dashboard.fleetStatus")}</CardTitle>
            <Bot className="h-7 w-7 shrink-0 text-gray-400 dark:text-gray-500" />
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              {fleet ? fleet.operational : "—"}
              <span className="text-lg text-gray-400 dark:text-gray-500">
                {" "}
                / {fleet ? fleet.total : "—"}
              </span>
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              {t("dashboard.robotsOperational")}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle>{t("dashboard.totalRobotUptime")}</CardTitle>
            <Clock className="h-7 w-7 shrink-0 text-gray-400 dark:text-gray-500" />
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              {stats?.totalUptime != null
                ? formatNumber(stats.totalUptime)
                : "—"}
              {stats?.totalUptime != null && (
                <span className="text-sm font-normal text-gray-400 dark:text-gray-500">
                  {" "}
                  {t("dashboard.hoursUnit")}
                </span>
              )}
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              {t("dashboard.cumulativeRobotHours")}
              {stats?.totalUptimeRate != null && (
                <Badge className="ml-2 border-transparent bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400">
                  {t("dashboard.uptimeRate", {
                    rate: (stats.totalUptimeRate * 100).toFixed(1),
                  })}
                </Badge>
              )}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle>{t("dashboard.totalDataCollected")}</CardTitle>
            <Database className="h-7 w-7 shrink-0 text-gray-400 dark:text-gray-500" />
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              {stats ? formatNumber(stats.totalCollection) : "—"}
              <span className="text-sm font-normal text-gray-400 dark:text-gray-500">
                {" "}
                {t("dashboard.hoursUnit")}
              </span>
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              {t("dashboard.cumulativeCollectionHours")}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Fleet Stats Table + Chart */}
      <FleetStatsPanel from={range.from} to={range.to} />
    </div>
  );
}

export function DashboardPage() {
  const { t } = useTranslation();

  return (
    <Suspense
      fallback={
        <div className="p-8 text-center text-gray-600 dark:text-gray-400">
          {t("common.loading")}
        </div>
      }
    >
      <DashboardContent />
    </Suspense>
  );
}
