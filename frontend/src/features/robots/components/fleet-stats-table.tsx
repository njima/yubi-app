"use client";

import { Bot, Clock, Database, MapPin } from "lucide-react";
import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";

import { THRESHOLD_COLORS } from "@/shared/lib/chart-colors";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { SearchableSelect } from "@/components/ui/searchable-select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

import {
  useFleetStatsQuery,
  useCollectionTrendQuery,
} from "../hooks/use-fleet-stats-query";

import type {
  FleetSiteStats,
  ColoredCollectionTrend,
  ColoredTrendSeries,
  Granularity,
} from "../schemas/fleet";

// --- Components ---

export function FleetStatsPanel({ from, to }: { from?: string; to?: string }) {
  return (
    <div className="grid gap-4 lg:grid-cols-2">
      <FleetStatsTable from={from} to={to} />
      <MonthlyCollectionChart from={from} to={to} />
    </div>
  );
}

function FleetStatsTable({ from, to }: { from?: string; to?: string }) {
  const { t } = useTranslation();
  const { data: statsData, isLoading, error } = useFleetStatsQuery(from, to);

  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-3">
          <CardTitle>{t("fleetStats.title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-gray-500">{t("fleetStats.loading")}</div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardHeader className="pb-3">
          <CardTitle>{t("fleetStats.title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-red-600 dark:text-red-400">
            {t("fleetStats.error", { message: error.message })}
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!statsData || statsData.length === 0) {
    return (
      <Card>
        <CardHeader className="pb-3">
          <CardTitle>{t("fleetStats.title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-gray-500">{t("fleetStats.empty")}</div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className="pb-3">
        <CardTitle>{t("fleetStats.title")}</CardTitle>
        <CardDescription>{t("fleetStats.description")}</CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-35 text-center">
                <span className="inline-flex items-center gap-1.5">
                  <MapPin className="h-3.5 w-3.5" />
                  {t("fleetStats.location")}
                </span>
              </TableHead>
              <TableHead className="text-center">
                <span className="inline-flex items-center gap-1.5">
                  <Bot className="h-3.5 w-3.5" />
                  {t("fleetStats.robotType")}
                </span>
              </TableHead>
              <TableHead className="text-center">
                <span className="inline-flex items-center gap-1.5">
                  <Clock className="h-3.5 w-3.5" />
                  {t("fleetStats.robotUptime")}
                </span>
              </TableHead>
              <TableHead className="text-center">
                <span className="inline-flex items-center gap-1.5">
                  <Clock className="h-3.5 w-3.5" />
                  {t("fleetStats.uptimeRate")}
                </span>
              </TableHead>
              <TableHead className="text-center">
                <span className="inline-flex items-center gap-1.5">
                  <Database className="h-3.5 w-3.5" />
                  {t("fleetStats.collectionTime")}
                </span>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {statsData.map((site) =>
              site.robot_types.map((robotType, robotTypeIdx) => {
                const key = `${site.site}-${robotType.robot_type}`;
                const uptimeRatePct =
                  robotType.uptime_rate != null
                    ? `${(robotType.uptime_rate * 100).toFixed(1)}%`
                    : "—";

                return (
                  <TableRow key={key}>
                    {robotTypeIdx === 0 ? (
                      <TableCell
                        rowSpan={site.robot_types.length}
                        className="font-semibold text-sm text-center align-middle border-r dark:border-gray-700"
                      >
                        {site.site}
                      </TableCell>
                    ) : null}
                    <TableCell className="text-sm font-medium text-center">
                      {robotType.robot_type}
                    </TableCell>
                    <TableCell className="text-sm text-center tabular-nums">
                      {robotType.robot_uptime != null
                        ? robotType.robot_uptime.toLocaleString()
                        : "—"}
                    </TableCell>
                    <TableCell className="text-sm text-center tabular-nums">
                      {uptimeRatePct}
                    </TableCell>
                    <TableCell className="text-sm text-center tabular-nums">
                      {robotType.data_collection_time.toLocaleString()}
                    </TableCell>
                  </TableRow>
                );
              })
            )}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}

type ViewMode = "site" | "robot_type";

function MonthlyCollectionChart({ from, to }: { from?: string; to?: string }) {
  const { t } = useTranslation();
  const { data: statsData } = useFleetStatsQuery(from, to);
  const [granularity, setGranularity] = useState<Granularity>("monthly");
  const {
    data: collectionData,
    isLoading,
    error,
  } = useCollectionTrendQuery(granularity, from, to);
  const [viewMode, setViewMode] = useState<ViewMode>("site");

  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader className="pb-3">
          <CardTitle>{t("fleetStats.trendTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-gray-500">
            {t("fleetStats.trendLoading")}
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="flex flex-col">
        <CardHeader className="pb-3">
          <CardTitle>{t("fleetStats.trendTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-red-600 dark:text-red-400">
            {t("fleetStats.trendError", { message: error.message })}
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!collectionData || !statsData) {
    return (
      <Card className="flex flex-col">
        <CardHeader className="pb-3">
          <CardTitle>{t("fleetStats.trendTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-gray-500">
            {t("fleetStats.trendEmpty")}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <MonthlyCollectionChartContent
      statsData={statsData}
      collectionData={collectionData}
      viewMode={viewMode}
      onViewModeChange={setViewMode}
      granularity={granularity}
      onGranularityChange={setGranularity}
    />
  );
}

const GRANULARITY_LABELS: Record<Granularity, string> = {
  hourly: "fleetStats.hourly",
  daily: "fleetStats.daily",
  monthly: "fleetStats.monthly",
};

function MonthlyCollectionChartContent({
  statsData,
  collectionData,
  viewMode,
  onViewModeChange,
  granularity,
  onGranularityChange,
}: {
  statsData: FleetSiteStats[];
  collectionData: ColoredCollectionTrend;
  viewMode: ViewMode;
  onViewModeChange: (mode: ViewMode) => void;
  granularity: Granularity;
  onGranularityChange: (g: Granularity) => void;
}) {
  const { t } = useTranslation();
  const labels = collectionData.labels;
  const series =
    viewMode === "site" ? collectionData.by_site : collectionData.by_robot_type;
  const allValues = series.flatMap((s) => s.data);
  const yTicks = getYTicks(Math.max(...allValues));
  const yMax = yTicks[yTicks.length - 1] ?? 1;

  const PAD = { top: 16, right: 20, bottom: 52, left: 56 };
  const W = 600;
  const H = 320;
  const chartW = W - PAD.left - PAD.right;
  const chartH = H - PAD.top - PAD.bottom;

  const [activeIdx, setActiveIdx] = useState<number | null>(null);
  const svgRef = useRef<SVGSVGElement>(null);

  function toX(i: number) {
    if (labels.length <= 1) return PAD.left + chartW / 2;
    return PAD.left + (i / (labels.length - 1)) * chartW;
  }
  function toY(v: number) {
    return PAD.top + chartH - (v / yMax) * chartH;
  }

  function linePath(data: number[]) {
    return data
      .map((v, i) => `${i === 0 ? "M" : "L"}${toX(i)},${toY(v)}`)
      .join(" ");
  }
  function areaPath(data: number[]) {
    const baseline = toY(0);
    return `${linePath(data)} L${toX(data.length - 1)},${baseline} L${toX(0)},${baseline} Z`;
  }

  function handleMouseMove(e: React.MouseEvent<SVGSVGElement>) {
    const svg = svgRef.current;
    if (!svg) return;
    const rect = svg.getBoundingClientRect();
    const svgX = ((e.clientX - rect.left) / rect.width) * W;
    let closest = 0;
    let minDist = Infinity;
    for (let i = 0; i < labels.length; i++) {
      const dist = Math.abs(svgX - toX(i));
      if (dist < minDist) {
        minDist = dist;
        closest = i;
      }
    }
    if (svgX >= PAD.left - 10 && svgX <= W - PAD.right + 10) {
      setActiveIdx(closest);
    } else {
      setActiveIdx(null);
    }
  }

  function handleMouseLeave() {
    setActiveIdx(null);
  }

  const tooltipX = activeIdx !== null ? toX(activeIdx) : 0;
  const tooltipRight = activeIdx !== null && tooltipX > W / 2;

  // Build summary from statsData
  const summaryItems = buildSummaryItems(statsData, series, viewMode);

  return (
    <Card className="flex flex-col">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>{t("fleetStats.trendTitle")}</CardTitle>
            <CardDescription>
              {t("fleetStats.trendDescription", {
                granularity: t(GRANULARITY_LABELS[granularity]),
              })}
            </CardDescription>
          </div>
          <div className="flex items-center gap-2">
            <SearchableSelect
              value={granularity}
              onValueChange={(value) =>
                onGranularityChange(value as Granularity)
              }
              options={[
                { value: "hourly", label: t("fleetStats.hourly") },
                { value: "daily", label: t("fleetStats.daily") },
                { value: "monthly", label: t("fleetStats.monthly") },
              ]}
              placeholder={t("fleetStats.granularity")}
              className="h-8"
            />
            <SearchableSelect
              value={viewMode}
              onValueChange={(value) => onViewModeChange(value as ViewMode)}
              options={[
                { value: "site", label: t("fleetStats.bySite") },
                { value: "robot_type", label: t("fleetStats.byRobotType") },
              ]}
              placeholder={t("fleetStats.viewMode")}
              className="h-8"
            />
          </div>
        </div>
      </CardHeader>
      <CardContent className="p-4 pt-0 flex-1 flex flex-col">
        <div className="flex gap-6 flex-1">
          {/* Chart */}
          <div className="flex-1 min-w-0">
            <svg
              ref={svgRef}
              viewBox={`0 0 ${W} ${H}`}
              preserveAspectRatio="xMidYMid meet"
              className="w-full h-full"
              onMouseMove={handleMouseMove}
              onMouseLeave={handleMouseLeave}
            >
              {/* Grid lines & Y labels */}
              {yTicks.map((tick) => {
                const y = toY(tick);
                return (
                  <g key={tick}>
                    <line
                      x1={PAD.left}
                      x2={W - PAD.right}
                      y1={y}
                      y2={y}
                      className="stroke-gray-200 dark:stroke-gray-700"
                      strokeWidth={0.5}
                    />
                    <text
                      x={PAD.left - 8}
                      y={y + 4}
                      textAnchor="end"
                      className="fill-gray-400 dark:fill-gray-500"
                      fontSize={10}
                    >
                      {tick.toLocaleString()}
                    </text>
                  </g>
                );
              })}

              {/* Y-axis label */}
              <text
                x={12}
                y={PAD.top + chartH / 2}
                textAnchor="middle"
                className="fill-gray-400 dark:fill-gray-500"
                fontSize={10}
                transform={`rotate(-90, 12, ${PAD.top + chartH / 2})`}
              >
                {t("fleetStats.dataCollected")}
              </text>

              {/* X labels (thinned to avoid overlap) */}
              {labels.map((m, i) => {
                const maxLabels = Math.floor(chartW / 70);
                const step = Math.max(1, Math.ceil(labels.length / maxLabels));
                if (i % step !== 0 && i !== labels.length - 1) return null;
                return (
                  <text
                    key={m}
                    x={toX(i)}
                    y={H - PAD.bottom + 18}
                    textAnchor="middle"
                    className="fill-gray-500 dark:fill-gray-400"
                    fontSize={10}
                  >
                    {formatChartLabel(m, granularity)}
                  </text>
                );
              })}

              {/* X-axis label */}
              <text
                x={PAD.left + chartW / 2}
                y={H - 6}
                textAnchor="middle"
                className="fill-gray-400 dark:fill-gray-500"
                fontSize={10}
              >
                {t("fleetStats.timeUtc")}
              </text>

              {/* Area fills */}
              {series.map((s) => (
                <path
                  key={`area-${s.label}`}
                  d={areaPath(s.data)}
                  fill={s.color}
                  opacity={0.1}
                />
              ))}

              {/* Lines */}
              {series.map((s) => (
                <path
                  key={`line-${s.label}`}
                  d={linePath(s.data)}
                  fill="none"
                  stroke={s.color}
                  strokeWidth={1.5}
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              ))}

              {/* Dots - dimmed when not active */}
              {series.map((s) =>
                s.data.map((v, i) => (
                  <circle
                    key={`dot-${s.label}-${i}`}
                    cx={toX(i)}
                    cy={toY(v)}
                    r={activeIdx === i ? 3.5 : 2}
                    fill={activeIdx === i ? s.color : "white"}
                    stroke={s.color}
                    strokeWidth={1.5}
                    style={{ transition: "r 0.1s, fill 0.1s" }}
                  />
                ))
              )}

              {/* Hover: vertical line */}
              {activeIdx !== null && (
                <line
                  x1={tooltipX}
                  x2={tooltipX}
                  y1={PAD.top}
                  y2={PAD.top + chartH}
                  className="stroke-gray-400 dark:stroke-gray-500"
                  strokeWidth={0.75}
                  strokeDasharray="3 2"
                />
              )}

              {/* Hover: tooltip */}
              {activeIdx !== null && (
                <g>
                  <rect
                    x={tooltipRight ? tooltipX - 140 : tooltipX + 10}
                    y={PAD.top}
                    width={130}
                    height={18 + series.length * 18}
                    rx={6}
                    className="fill-white dark:fill-gray-800 stroke-gray-200 dark:stroke-gray-600"
                    strokeWidth={0.5}
                    filter="drop-shadow(0 1px 3px rgba(0,0,0,0.1))"
                  />
                  {/* Month label */}
                  <text
                    x={tooltipRight ? tooltipX - 132 : tooltipX + 18}
                    y={PAD.top + 14}
                    fontSize={10}
                    fontWeight={600}
                    className="fill-gray-800 dark:fill-gray-100"
                  >
                    {formatChartLabel(labels[activeIdx] ?? "", granularity)}
                  </text>
                  {/* Series values */}
                  {series.map((s, si) => (
                    <g key={s.label}>
                      <circle
                        cx={tooltipRight ? tooltipX - 132 : tooltipX + 18}
                        cy={PAD.top + 30 + si * 18}
                        r={3.5}
                        fill={s.color}
                      />
                      <text
                        x={tooltipRight ? tooltipX - 124 : tooltipX + 26}
                        y={PAD.top + 34 + si * 18}
                        fontSize={10}
                        className="fill-gray-600 dark:fill-gray-300"
                      >
                        {s.label}
                      </text>
                      <text
                        x={tooltipRight ? tooltipX - 18 : tooltipX + 132}
                        y={PAD.top + 34 + si * 18}
                        fontSize={10}
                        fontWeight={600}
                        textAnchor="end"
                        className="fill-gray-800 dark:fill-gray-100"
                      >
                        {s.data[activeIdx]?.toLocaleString()}
                      </text>
                    </g>
                  ))}
                </g>
              )}
            </svg>
          </div>

          {/* Period Summary */}
          <div className="w-44 shrink-0 border-l border-gray-200 dark:border-gray-700 pl-6 flex flex-col justify-center gap-3">
            {summaryItems.map((item) => (
              <div key={item.label}>
                <div className="flex items-center gap-1.5 mb-1">
                  <span
                    className="inline-block h-2.5 w-2.5 rounded-full"
                    style={{ backgroundColor: item.color }}
                  />
                  <span className="text-xs text-gray-500 dark:text-gray-400">
                    {item.label}
                  </span>
                </div>
                <div className="flex items-baseline gap-2">
                  <p className="text-lg font-bold text-gray-900 dark:text-gray-100 tabular-nums">
                    {item.total.toLocaleString()}
                    <span className="text-xs font-normal text-gray-400 dark:text-gray-500">
                      {" "}
                      {t("dashboard.hoursUnit")}
                    </span>
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

// --- Helpers ---

type SummaryItem = {
  label: string;
  color: string;
  total: number;
};

function buildSummaryItems(
  statsData: FleetSiteStats[],
  series: ColoredTrendSeries[],
  viewMode: ViewMode
): SummaryItem[] {
  const items: SummaryItem[] = [];
  if (viewMode === "site") {
    for (const site of statsData) {
      const total = site.robot_types.reduce(
        (sum, m) => sum + m.data_collection_time,
        0
      );
      const color =
        series.find((s) => s.label === site.site)?.color ??
        THRESHOLD_COLORS.muted;
      items.push({ label: site.site, color, total });
    }
  } else {
    const robotTypeMap = new Map<string, number>();
    for (const site of statsData) {
      for (const m of site.robot_types) {
        const prev = robotTypeMap.get(m.robot_type) ?? 0;
        robotTypeMap.set(m.robot_type, prev + m.data_collection_time);
      }
    }
    for (const [robotType, total] of robotTypeMap) {
      const color =
        series.find((s) => s.label === robotType)?.color ??
        THRESHOLD_COLORS.muted;
      items.push({ label: robotType, color, total });
    }
  }
  return items;
}

const MONTH_NAMES = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec",
];

function formatChartLabel(raw: string, granularity: Granularity): string {
  if (granularity === "hourly") {
    // "2026-03-15T09:00:00Z" → "Mar 15 09:00"
    const d = new Date(raw);
    const monthName = MONTH_NAMES[d.getUTCMonth()];
    const day = d.getUTCDate();
    const hours = String(d.getUTCHours()).padStart(2, "0");
    const mins = String(d.getUTCMinutes()).padStart(2, "0");
    return `${monthName} ${day} ${hours}:${mins}`;
  }
  if (granularity === "daily") {
    // "2026-03-15" → "Mar 15"
    const parts = raw.split("-");
    const monthIdx = parseInt(parts[1] ?? "0", 10) - 1;
    return `${MONTH_NAMES[monthIdx] ?? ""} ${parseInt(parts[2] ?? "0", 10)}`;
  }
  if (granularity === "monthly") {
    // "2026-03" → "Mar"
    const parts = raw.split("-");
    const monthIdx = parseInt(parts[1] ?? "0", 10) - 1;
    return MONTH_NAMES[monthIdx] ?? "";
  }
  return raw;
}

function getYTicks(max: number): number[] {
  if (max <= 0) return [0, 1];

  // Choose a nice step size based on the data range
  const rawStep = max / 4;
  const magnitude = Math.pow(10, Math.floor(Math.log10(rawStep)));
  const niceSteps = [1, 2, 5, 10];
  const step = niceSteps.find((s) => s * magnitude >= rawStep)! * magnitude;

  const ticks: number[] = [];
  for (let i = 0; i <= max + step; i += step) {
    ticks.push(Math.round(i * 1000) / 1000);
    if (i >= max) break;
  }
  return ticks;
}
