/**
 * MJPEG Viewer Component
 * Displays MJPEG stream from ROS2 web_video_server with loading and error states
 */

"use client";

import { AlertCircle, RefreshCw, Video } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";

interface MjpegViewerProps {
  /** MJPEG stream URL */
  streamUrl: string;
  /** Alternative text for the image */
  alt?: string;
  /** Additional CSS classes */
  className?: string;
}

type LoadingState = "loading" | "loaded" | "error";

export function MjpegViewer({
  streamUrl,
  alt = "Camera stream",
  className = "",
}: MjpegViewerProps) {
  const { t } = useTranslation();
  const [loadingState, setLoadingState] = useState<LoadingState>("loading");
  const [retryKey, setRetryKey] = useState(0);

  const handleLoad = () => {
    setLoadingState("loaded");
  };

  const handleError = () => {
    setLoadingState("error");
  };

  const handleRetry = () => {
    setLoadingState("loading");
    setRetryKey((prev) => prev + 1);
  };

  const imageUrl = `${streamUrl}&_retry=${retryKey}`;

  return (
    <div className={`relative w-full aspect-video ${className}`}>
      {/* Loading State */}
      {loadingState === "loading" && (
        <div className="absolute inset-0 flex items-center justify-center bg-gray-100 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
          <div className="flex flex-col items-center gap-3 text-gray-600 dark:text-gray-400">
            <Video className="h-12 w-12 animate-pulse" />
            <p className="text-sm font-medium">{t("mjpegViewer.connecting")}</p>
          </div>
        </div>
      )}

      {/* Error State */}
      {loadingState === "error" && (
        <div className="absolute inset-0 flex items-center justify-center bg-gray-100 dark:bg-gray-800 rounded-lg border border-red-300 dark:border-red-700">
          <div className="flex flex-col items-center gap-3 text-red-600 dark:text-red-400 p-6 text-center">
            <AlertCircle className="h-12 w-12" />
            <div>
              <p className="text-sm font-medium mb-1">
                {t("mjpegViewer.error")}
              </p>
              <p className="text-xs text-gray-600 dark:text-gray-400">
                {t("mjpegViewer.errorHint")}
              </p>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={handleRetry}
              className="mt-2"
            >
              <RefreshCw className="h-4 w-4 mr-2" />
              {t("mjpegViewer.reload")}
            </Button>
          </div>
        </div>
      )}

      {/* MJPEG Stream Image */}
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={imageUrl}
        alt={alt}
        onLoad={handleLoad}
        onError={handleError}
        className={`w-full h-auto rounded-lg border border-gray-200 dark:border-gray-700 ${
          loadingState === "loaded" ? "block" : "hidden"
        }`}
        style={{
          maxHeight: "600px",
          objectFit: "contain",
        }}
      />
    </div>
  );
}
