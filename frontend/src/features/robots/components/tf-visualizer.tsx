"use client";

import { useTranslation } from "react-i18next";

import type { LayoutContext } from "@/features/robots/lib/teleop-layout-registry";

interface TfVisualizerProps {
  ctx: LayoutContext;
}

export function TfVisualizer({ ctx }: TfVisualizerProps) {
  const { t } = useTranslation();
  const host = ctx.host;
  const rosbridgePort = ctx.rosbridgePort ?? 9090;

  if (!host) {
    return (
      <div className="flex items-center justify-center h-full rounded-lg border bg-white dark:bg-gray-800 dark:border-gray-700 p-6">
        <p className="text-sm text-gray-500 dark:text-gray-400">
          {t("tfVisualizer.hostNotConfigured")}
        </p>
      </div>
    );
  }

  const wsUrl = `ws://${host}:${rosbridgePort}`;
  const src = `/web/tf-visualizer?wsUrl=${encodeURIComponent(wsUrl)}`;

  return (
    <iframe
      src={src}
      className="w-full h-full min-h-75 rounded-lg border-0"
      title={t("tfVisualizer.title")}
    />
  );
}
