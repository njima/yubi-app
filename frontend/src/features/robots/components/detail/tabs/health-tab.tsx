"use client";

import { useTranslation } from "react-i18next";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { cn } from "@/shared/lib/utils";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { GateStatusBadge } from "../../gate-status-badge";

import type { RobotStatusStreamDetail } from "../../../schemas/robot";

type GateGroupStatus = z.infer<typeof schemas.GateGroupStatus>;

interface HealthTabProps {
  realtimeStatus?: RobotStatusStreamDetail | null;
}

export function HealthTab({ realtimeStatus }: HealthTabProps) {
  const { t } = useTranslation();
  const gate = realtimeStatus?.gate_conditions;

  if (!gate) {
    return (
      <div className="text-center py-12 text-gray-500 dark:text-gray-400">
        {t("robotHealth.noGateData")}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Overall Gate Level */}
      <Card>
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between">
            <CardTitle className="text-base font-medium">
              {t("robotHealth.overallGateLevel")}
            </CardTitle>
            <GateStatusBadge gateConditions={gate} />
          </div>
        </CardHeader>
      </Card>

      {/* Groups */}
      {Object.entries(gate.groups).map(([name, group]) => (
        <GroupCard key={name} name={name} group={group} t={t} />
      ))}
    </div>
  );
}

function GroupCard({
  name,
  group,
  t,
}: {
  name: string;
  group: GateGroupStatus;
  t: (key: string) => string;
}) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="text-base font-medium">{name}</CardTitle>
          <div className="flex items-center gap-3 text-sm">
            <LevelBadge level={group.level} t={t} />
            <span className="text-gray-500 dark:text-gray-400">
              {group.settled
                ? t("robotHealth.settled")
                : t("robotHealth.unsettled")}
            </span>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {group.conditions.map((c) => (
            <div
              key={c.name}
              className={cn(
                "flex items-center justify-between rounded-md border px-3 py-2",
                c.passed
                  ? "border-green-200 bg-green-50/50 dark:border-green-800 dark:bg-green-950/30"
                  : "border-yellow-200 bg-yellow-50/50 dark:border-yellow-800 dark:bg-yellow-950/30"
              )}
            >
              <div className="flex items-center gap-2">
                <span
                  className={cn(
                    "h-2 w-2 rounded-full",
                    c.passed ? "bg-green-500" : "bg-yellow-500"
                  )}
                />
                <span className="text-sm font-medium">{c.name}</span>
              </div>
              <div className="flex items-center gap-3 text-sm text-gray-500 dark:text-gray-400">
                <span>{c.reason}</span>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

function LevelBadge({
  level,
  t,
}: {
  level: number;
  t: (key: string) => string;
}) {
  const config =
    level === 0
      ? {
          label: t("robotHealth.open"),
          className:
            "text-green-700 bg-green-100 dark:text-green-400 dark:bg-green-900/30",
        }
      : level === 1
        ? {
            label: t("robotHealth.block"),
            className:
              "text-yellow-700 bg-yellow-100 dark:text-yellow-400 dark:bg-yellow-900/30",
          }
        : {
            label: t("robotHealth.hardStop"),
            className:
              "text-red-700 bg-red-100 dark:text-red-400 dark:bg-red-900/30",
          };

  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-2 py-0.5 text-xs font-semibold",
        config.className
      )}
    >
      {config.label}
    </span>
  );
}
