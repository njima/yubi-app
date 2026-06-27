"use client";

import { Video } from "lucide-react";
import { useTranslation } from "react-i18next";

import { buildMjpegStreamUrl } from "@/lib/mjpeg";
import { EPISODE_COLLECTION_STATUS } from "@/lib/status/constants";
import { cn } from "@/lib/utils";

import {
  ErrorOverlay,
  GateOverlay,
  RecordingIndicator,
  SubtaskOverlay,
} from "./camera-overlays";
import { MjpegViewer } from "./mjpeg-viewer";

import type { Camera } from "../lib/camera-utils";

export interface CameraViewProps {
  camera?: Camera;
  host?: string;
  port?: number;
  robotName?: string;
  placeholderLabel?: string;
  /** Show REC / subtask / error overlays (layout config `overlay: true`) */
  showOverlays?: boolean;
  episodeStatus?: number;
  errorDetails?: string;
  currentSubtask?: { order_index: number; name: string } | null;
  nextSubtask?: { order_index: number; name: string } | null;
  parameterValues?: Record<string, string> | null;
  gateLevel?: number;
  /** MJPEG stream quality/resolution config */
  streamConfig?: { quality?: number; width?: number; height?: number };
}

export function CameraView({
  camera,
  host,
  port,
  robotName = "Robot",
  placeholderLabel,
  showOverlays = false,
  episodeStatus,
  errorDetails,
  currentSubtask,
  nextSubtask,
  parameterValues,
  gateLevel,
  streamConfig,
}: CameraViewProps) {
  const { t } = useTranslation();

  if (!host || !port) {
    return <CameraPlaceholder label={t("cameraView.connection")} />;
  }
  if (!camera) {
    return (
      <CameraPlaceholder label={placeholderLabel ?? t("cameraView.camera")} />
    );
  }

  const cameraLabel =
    camera.name || camera.namespace || t("cameraView.unknownCamera");
  const streamUrl = buildMjpegStreamUrl({
    host,
    port,
    namespace: camera.namespace,
    ...streamConfig,
  });

  const isRecording = episodeStatus === EPISODE_COLLECTION_STATUS.RECORDING;
  const hasError = !!errorDetails;
  const ringClass = cn(
    hasError && "ring-2 ring-yellow-500 rounded-lg",
    !hasError && isRecording && "ring-2 ring-red-500 rounded-lg"
  );

  return (
    <div className="relative">
      <MjpegViewer
        className={showOverlays ? ringClass : ""}
        streamUrl={streamUrl}
        alt={`${robotName} - ${cameraLabel}`}
      />
      {/* Camera name label */}
      <div className="absolute top-2 left-2 z-10 bg-black/60 text-white text-xs rounded px-2 py-1">
        {cameraLabel}
      </div>
      {/* Overlays only on designated camera */}
      {showOverlays && (
        <>
          {gateLevel !== undefined && <GateOverlay level={gateLevel} />}
          <div className="absolute top-2 right-2 z-20 flex flex-col items-end gap-1">
            {isRecording && <RecordingIndicator />}
            {currentSubtask && (
              <SubtaskOverlay
                currentSubtask={currentSubtask}
                nextSubtask={nextSubtask}
                parameterValues={parameterValues}
              />
            )}
          </div>
          {errorDetails && <ErrorOverlay message={errorDetails} />}
        </>
      )}
    </div>
  );
}

export function CameraPlaceholder({ label }: { label: string }) {
  const { t } = useTranslation();

  return (
    <div className="aspect-video bg-gray-100 dark:bg-gray-800 rounded-lg flex flex-col items-center justify-center">
      <Video className="h-8 w-8 text-gray-400 dark:text-gray-500" />
      <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">
        {t("cameraView.notConfigured", { label })}
      </p>
    </div>
  );
}

export function CameraSkeleton() {
  return (
    <div className="aspect-video bg-gray-200 dark:bg-gray-700 rounded-lg animate-pulse" />
  );
}
