"use client";

import { useQuery, type UseQueryOptions } from "@tanstack/react-query";

import { assignColors } from "@/lib/chart-colors";

import type {
  ColoredCollectionTrend,
  FleetSiteStats,
  Granularity,
  TrendSeries,
} from "../schemas/fleet";

export const fleetStatsQueryKeys = {
  all: ["fleetStats"] as const,
  lists: () => [...fleetStatsQueryKeys.all, "list"] as const,
  list: (from?: string, to?: string) =>
    [...fleetStatsQueryKeys.lists(), { from, to }] as const,
  trends: () => [...fleetStatsQueryKeys.all, "trend"] as const,
  trend: (granularity: Granularity, from?: string, to?: string) =>
    [...fleetStatsQueryKeys.trends(), granularity, { from, to }] as const,
};

/** Ensure date string is in RFC3339 format for the backend */
function toRFC3339(date: string): string {
  if (date.includes("T")) return date;
  return `${date}T00:00:00Z`;
}

/**
 * Convert inclusive end date to exclusive for the API.
 * User selects "to 3/18" meaning "include 3/18",
 * but API uses exclusive to (period_start < to).
 * So we add 1 day: "2026-03-18" → "2026-03-19T00:00:00Z"
 */
function toExclusiveEnd(date: string): string {
  const d = new Date(toRFC3339(date));
  d.setUTCDate(d.getUTCDate() + 1);
  return d.toISOString().replace(".000Z", "Z");
}

/**
 * Generate all labels for a date range based on granularity.
 * Ensures the chart shows every time slot, even if there's no data.
 */
function generateAllLabels(
  from: string,
  to: string,
  granularity: Granularity
): string[] {
  const labels: string[] = [];
  const fromDate = new Date(toRFC3339(from));
  const toDate = new Date(toRFC3339(to));

  const current = new Date(fromDate);

  while (current < toDate) {
    if (granularity === "hourly") {
      labels.push(current.toISOString().replace(".000Z", "Z"));
      current.setUTCHours(current.getUTCHours() + 1);
    } else if (granularity === "daily") {
      const y = current.getUTCFullYear();
      const m = String(current.getUTCMonth() + 1).padStart(2, "0");
      const d = String(current.getUTCDate()).padStart(2, "0");
      labels.push(`${y}-${m}-${d}`);
      current.setUTCDate(current.getUTCDate() + 1);
    } else {
      const y = current.getUTCFullYear();
      const m = String(current.getUTCMonth() + 1).padStart(2, "0");
      labels.push(`${y}-${m}`);
      current.setUTCMonth(current.getUTCMonth() + 1);
    }
  }

  return labels;
}

/**
 * Fill series data to match allLabels, inserting 0 for missing labels.
 */
function fillSeries(
  series: TrendSeries[],
  apiLabels: string[],
  allLabels: string[]
): TrendSeries[] {
  const apiLabelIndex = new Map<string, number>();
  apiLabels.forEach((label, i) => {
    apiLabelIndex.set(label, i);
  });

  return series.map((s) => ({
    label: s.label,
    data: allLabels.map((label) => {
      const idx = apiLabelIndex.get(label);
      return idx != null ? (s.data[idx] ?? 0) : 0;
    }),
  }));
}

/**
 * Fetch fleet statistics data (table: uptime & collection by site/robot type).
 */
export function useFleetStatsQuery(
  from?: string,
  to?: string,
  options?: Omit<
    UseQueryOptions<FleetSiteStats[], Error>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery({
    queryKey: fleetStatsQueryKeys.list(from, to),
    queryFn: async (): Promise<FleetSiteStats[]> => {
      if (!from || !to) return [];

      const searchParams = new URLSearchParams();
      searchParams.append("from", toRFC3339(from));
      searchParams.append("to", toExclusiveEnd(to));

      const response = await fetch(
        `/web/api/fleet/stats?${searchParams.toString()}`
      );
      if (!response.ok) {
        throw new Error(`Failed to fetch fleet stats: ${response.statusText}`);
      }
      return response.json();
    },
    enabled: !!from && !!to,
    ...options,
  });
}

/**
 * Fetch collection trend data (chart: by location & by model).
 * Generates all labels for the date range and fills missing data with 0.
 * Assigns colors from CHART_PALETTE to each series.
 */
export function useCollectionTrendQuery(
  granularity: Granularity = "monthly",
  from?: string,
  to?: string,
  options?: Omit<
    UseQueryOptions<ColoredCollectionTrend, Error>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery({
    queryKey: fleetStatsQueryKeys.trend(granularity, from, to),
    queryFn: async (): Promise<ColoredCollectionTrend> => {
      if (!from || !to) {
        return { labels: [], by_site: [], by_robot_type: [] };
      }

      const searchParams = new URLSearchParams();
      searchParams.append("granularity", granularity);
      searchParams.append("from", toRFC3339(from));
      searchParams.append("to", toExclusiveEnd(to));

      const response = await fetch(
        `/web/api/fleet/collection-trend?${searchParams.toString()}`
      );
      if (!response.ok) {
        throw new Error(
          `Failed to fetch collection trend: ${response.statusText}`
        );
      }
      const data = await response.json();

      const apiLabels: string[] = data.labels ?? [];
      const allLabels = generateAllLabels(
        from,
        toExclusiveEnd(to),
        granularity
      );

      return {
        labels: allLabels,
        by_site: assignColors(
          fillSeries(data.by_site ?? [], apiLabels, allLabels)
        ),
        by_robot_type: assignColors(
          fillSeries(data.by_robot_type ?? [], apiLabels, allLabels)
        ),
      };
    },
    enabled: !!from && !!to,
    ...options,
  });
}
