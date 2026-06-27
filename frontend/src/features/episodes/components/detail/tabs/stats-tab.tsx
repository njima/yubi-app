"use client";

import { useTranslation } from "react-i18next";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

import { useEpisodeStatsQuery } from "../../../hooks/use-episode-stats-query";

type EpisodeFeatureStats = z.infer<typeof schemas.EpisodeFeatureStats>;

interface StatsTabProps {
  episodeId: string;
}

interface StatsRow {
  feature: string;
  ch: number;
  min: number;
  max: number;
  mean: number;
  std: number;
  count: number;
  isFirstInFeature: boolean;
  rowSpan: number;
}

function fmt(v: number): string {
  return v.toFixed(4);
}

function extractArray(value: unknown): number[] {
  if (typeof value === "number") return [value];
  if (Array.isArray(value)) {
    const flat: number[] = [];
    for (const item of value) {
      flat.push(...extractArray(item));
    }
    return flat;
  }
  return [];
}

function buildRows(feature: string, stats: EpisodeFeatureStats): StatsRow[] {
  const mins = extractArray(stats.min);
  const maxs = extractArray(stats.max);
  const means = extractArray(stats.mean);
  const stds = extractArray(stats.std);
  const len = Math.max(mins.length, maxs.length, means.length, stds.length, 1);

  return Array.from({ length: len }, (_, i) => ({
    feature,
    ch: i,
    min: mins[i] ?? 0,
    max: maxs[i] ?? 0,
    mean: means[i] ?? 0,
    std: stds[i] ?? 0,
    count: stats.count,
    isFirstInFeature: i === 0,
    rowSpan: len,
  }));
}

export function StatsTab({ episodeId }: StatsTabProps) {
  const { t } = useTranslation();
  const { data, isLoading, error } = useEpisodeStatsQuery(episodeId);

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>{t("episodeStats.title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground text-sm">
            {t("episodeStats.loading")}
          </p>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>{t("episodeStats.title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground text-sm">
            {t("episodeStats.loadFailed")}
          </p>
        </CardContent>
      </Card>
    );
  }

  if (!data || Object.keys(data.stats).length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>{t("episodeStats.title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground text-sm">
            {t("episodeStats.empty")}
          </p>
        </CardContent>
      </Card>
    );
  }

  const rows: StatsRow[] = Object.entries(data.stats)
    .sort(([a], [b]) => a.localeCompare(b))
    .flatMap(([feature, stats]) => buildRows(feature, stats));

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("episodeStats.title")}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t("episodeStats.key")}</TableHead>
                <TableHead className="w-12 text-right">
                  {t("episodeStats.ch")}
                </TableHead>
                <TableHead className="text-right">
                  {t("episodeStats.min")}
                </TableHead>
                <TableHead className="text-right">
                  {t("episodeStats.max")}
                </TableHead>
                <TableHead className="text-right">
                  {t("episodeStats.mean")}
                </TableHead>
                <TableHead className="text-right">
                  {t("episodeStats.std")}
                </TableHead>
                <TableHead className="text-right">
                  {t("episodeStats.count")}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {rows.map((row) => (
                <TableRow key={`${row.feature}-${row.ch}`}>
                  {row.isFirstInFeature && (
                    <TableCell
                      rowSpan={row.rowSpan}
                      className="align-top font-mono text-xs"
                    >
                      {row.feature}
                    </TableCell>
                  )}
                  <TableCell className="text-muted-foreground text-right font-mono text-xs">
                    {row.ch}
                  </TableCell>
                  <TableCell className="text-right font-mono text-xs">
                    {fmt(row.min)}
                  </TableCell>
                  <TableCell className="text-right font-mono text-xs">
                    {fmt(row.max)}
                  </TableCell>
                  <TableCell className="text-right font-mono text-xs">
                    {fmt(row.mean)}
                  </TableCell>
                  <TableCell className="text-right font-mono text-xs">
                    {fmt(row.std)}
                  </TableCell>
                  {row.isFirstInFeature && (
                    <TableCell
                      rowSpan={row.rowSpan}
                      className="text-right align-top font-mono text-xs"
                    >
                      {row.count}
                    </TableCell>
                  )}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}
