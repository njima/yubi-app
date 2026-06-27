/**
 * Robot Camera Viewer Component
 * Displays robot camera streams with multi-camera support via tabs
 */

"use client";

import { ExternalLink } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { buildMjpegStreamUrl, buildMjpegViewerUrl } from "@/shared/lib/mjpeg";

import { Button } from "@/components/ui/button";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";

import { MjpegViewer } from "./mjpeg-viewer";
import { extractCameras, extractHostPort } from "../lib/camera-utils";

interface RobotCameraViewerProps {
  /** Robot configuration including host, port, and cameras */
  robotConfig?: Record<string, unknown>;
  /** Robot name for display */
  robotName?: string;
}

export function RobotCameraViewer({
  robotConfig,
  robotName = "Robot",
}: RobotCameraViewerProps) {
  const { t } = useTranslation();
  const { host, port } = extractHostPort(robotConfig);
  const cameras = extractCameras(robotConfig);
  const [selectedCamera, setSelectedCamera] = useState(
    cameras.length > 0 ? cameras[0]!.namespace : ""
  );

  // Check if host/port are configured
  const hasConnectionConfig = host && port;

  if (cameras.length === 0 || !hasConnectionConfig) {
    return (
      <div className="rounded-lg border bg-white dark:bg-gray-800 dark:border-gray-700 p-6">
        <h2 className="text-lg font-semibold mb-4">
          {t("robotCameraViewer.title")}
        </h2>
        <div className="flex items-center justify-center py-12 text-gray-500 dark:text-gray-400">
          <p className="text-sm">
            {!hasConnectionConfig
              ? t("robotCameraViewer.connectionNotConfigured")
              : t("robotCameraViewer.cameraNotConfigured")}
          </p>
        </div>
      </div>
    );
  }

  // Single camera: simple display
  if (cameras.length === 1) {
    const camera = cameras[0]!;
    const streamUrl = buildMjpegStreamUrl({
      host,
      port,
      namespace: camera.namespace,
    });
    const viewerUrl = buildMjpegViewerUrl({
      host,
      port,
      namespace: camera.namespace,
    });

    return (
      <div className="rounded-lg border bg-white dark:bg-gray-800 dark:border-gray-700 p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold">
            {camera.name || t("robotCameraViewer.cameraFootage")}
          </h2>
          <Button
            variant="outline"
            size="sm"
            onClick={() => window.open(viewerUrl, "_blank")}
          >
            <ExternalLink className="h-4 w-4 mr-2" />
            {t("robotCameraViewer.openInNewTab")}
          </Button>
        </div>
        <MjpegViewer
          streamUrl={streamUrl}
          alt={`${robotName} - ${camera.name || "Camera"}`}
        />
      </div>
    );
  }

  // Multiple cameras: tabs display
  return (
    <div className="rounded-lg border bg-white dark:bg-gray-800 dark:border-gray-700 p-6">
      <h2 className="text-lg font-semibold mb-4">
        {t("robotCameraViewer.title")}
      </h2>
      <Tabs value={selectedCamera} onValueChange={setSelectedCamera}>
        <div className="flex items-center justify-between mb-4">
          <TabsList>
            {cameras.map((camera) => (
              <TabsTrigger key={camera.namespace} value={camera.namespace}>
                {camera.name || camera.namespace}
              </TabsTrigger>
            ))}
          </TabsList>
          <Button
            variant="outline"
            size="sm"
            onClick={() =>
              window.open(
                buildMjpegViewerUrl({ host, port, namespace: selectedCamera }),
                "_blank"
              )
            }
          >
            <ExternalLink className="h-4 w-4 mr-2" />
            {t("robotCameraViewer.openInNewTab")}
          </Button>
        </div>

        {(() => {
          const camera = cameras.find((c) => c.namespace === selectedCamera);
          if (!camera) return null;
          return (
            <MjpegViewer
              streamUrl={buildMjpegStreamUrl({
                host,
                port,
                namespace: camera.namespace,
              })}
              alt={`${robotName} - ${camera.name || camera.namespace}`}
            />
          );
        })()}
      </Tabs>
    </div>
  );
}
