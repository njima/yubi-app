"use client";

import { useTranslation } from "react-i18next";

import { THRESHOLD_COLORS } from "@/lib/chart-colors";

import { useFleetSummaryQuery } from "../hooks/use-fleet-summary-query";

import type { FleetStatusCount } from "../schemas/fleet";

// --- Donut helpers ---

const DONUT_RADIUS = 15.5;
const DONUT_CIRCUMFERENCE = 2 * Math.PI * DONUT_RADIUS;

type DonutCell = FleetStatusCount | undefined;

function getStrokeColor(cell: DonutCell): string {
  if (!cell || cell.total === 0) return THRESHOLD_COLORS.track;
  const ratio = cell.operational / cell.total;
  if (ratio >= 0.8) return THRESHOLD_COLORS.good;
  if (ratio >= 0.5) return THRESHOLD_COLORS.warning;
  return THRESHOLD_COLORS.danger;
}

function getStrokeDasharray(cell: DonutCell): string {
  if (!cell || cell.total === 0) return `0 ${DONUT_CIRCUMFERENCE}`;
  const filled = (cell.operational / cell.total) * DONUT_CIRCUMFERENCE;
  return `${filled} ${DONUT_CIRCUMFERENCE - filled}`;
}

// --- Donut component ---

function DonutChart({ cell, label }: { cell: DonutCell; label: string }) {
  return (
    <div className="flex flex-col items-center gap-1">
      <svg width={56} height={56} viewBox="0 0 40 40">
        <circle
          cx={20}
          cy={20}
          r={DONUT_RADIUS}
          fill="none"
          stroke="currentColor"
          strokeWidth={4}
          className="text-gray-200 dark:text-gray-600"
        />
        <circle
          cx={20}
          cy={20}
          r={DONUT_RADIUS}
          fill="none"
          stroke={getStrokeColor(cell)}
          strokeWidth={4}
          strokeLinecap="round"
          strokeDasharray={getStrokeDasharray(cell)}
          transform="rotate(-90 20 20)"
        />
        <text
          x={20}
          y={22}
          textAnchor="middle"
          fontSize={cell && cell.total >= 10 ? 8 : 9}
          fontWeight={700}
          className="fill-gray-800 dark:fill-gray-100"
        >
          {cell ? `${cell.operational}/${cell.total}` : "\u2014"}
        </text>
      </svg>
      <span className="text-[11px] font-medium text-gray-500 dark:text-gray-400">
        {label}
      </span>
    </div>
  );
}

// --- Component ---

export function FleetSummaryGrid() {
  const { t } = useTranslation();
  const { data: sites, isLoading, error } = useFleetSummaryQuery();

  if (isLoading) {
    return (
      <div className="space-y-3">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
          {t("fleetSummary.title")}
        </h2>
        <div className="text-sm text-gray-500">{t("fleetSummary.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="space-y-3">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
          {t("fleetSummary.title")}
        </h2>
        <div className="text-sm text-red-600 dark:text-red-400">
          {t("fleetSummary.error", { message: error.message })}
        </div>
      </div>
    );
  }

  if (!sites || sites.length === 0) {
    return (
      <div className="space-y-3">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
          {t("fleetSummary.title")}
        </h2>
        <div className="text-sm text-gray-500">{t("fleetSummary.empty")}</div>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
        {t("fleetSummary.title")}
      </h2>

      <div className="grid gap-3 grid-cols-1 md:grid-cols-2 xl:grid-cols-3">
        {sites.map((site) => {
          const activeRobotTypes = Object.entries(site.robot_types).filter(
            ([, data]) => data.leader != null || data.follower != null
          );
          if (activeRobotTypes.length === 0) return null;

          return (
            <div
              key={site.site}
              className="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4"
            >
              <div className="flex items-center gap-2 mb-4">
                <h3 className="text-sm font-bold text-gray-900 dark:text-gray-100">
                  {site.site}
                </h3>
                <span className="text-[11px] text-gray-500 dark:text-gray-400 bg-gray-200 dark:bg-gray-700 px-2 py-0.5 rounded-full">
                  {t("fleetSummary.robotTypes", {
                    count: activeRobotTypes.length,
                  })}
                </span>
              </div>
              <div className="grid grid-cols-2 gap-3 max-h-48 overflow-y-auto">
                {activeRobotTypes.map(([robotType, data]) => (
                  <div
                    key={robotType}
                    className="flex flex-col items-center gap-1"
                  >
                    <span className="text-xs font-semibold text-gray-700 dark:text-gray-300">
                      {robotType}
                    </span>
                    <div className="flex gap-3">
                      <DonutChart
                        cell={data.leader}
                        label={t("fleetSummary.leader")}
                      />
                      <DonutChart
                        cell={data.follower}
                        label={t("fleetSummary.follower")}
                      />
                    </div>
                  </div>
                ))}
              </div>
            </div>
          );
        })}
      </div>

      {/* Legend */}
      <div className="flex flex-wrap items-center gap-4 text-xs text-gray-500 dark:text-gray-400">
        <div className="flex items-center gap-1">
          <span className="inline-block h-3 w-3 rounded-full bg-blue-500" />
          <span>&ge; 80%</span>
        </div>
        <div className="flex items-center gap-1">
          <span className="inline-block h-3 w-3 rounded-full bg-amber-500" />
          <span>&ge; 50%</span>
        </div>
        <div className="flex items-center gap-1">
          <span className="inline-block h-3 w-3 rounded-full bg-red-500" />
          <span>&lt; 50%</span>
        </div>
      </div>
    </div>
  );
}
