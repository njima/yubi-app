"use client";

import { AlertTriangle } from "lucide-react";
import { useTranslation } from "react-i18next";

import { ParameterizedName } from "@/features/tasks";

export function GateOverlay({ level }: { level: number }) {
  const { t } = useTranslation();

  if (level <= 0) return null;

  if (level === 1) {
    return (
      <>
        {/* Yellow gradient border — wider, more visible */}
        <div
          className="absolute inset-0 z-10 pointer-events-none rounded-lg animate-[pulse_3s_ease-in-out_infinite]"
          style={{
            boxShadow: "inset 0 0 24px 10px rgba(250,204,21,0.6)",
          }}
        />
        {/* Badge */}
        <div className="absolute top-10 left-2 z-10">
          <div className="bg-yellow-400/90 text-black text-xs font-bold px-2 py-1 rounded flex items-center gap-1">
            <AlertTriangle className="h-3 w-3" />
            {t("cameraOverlays.block")}
          </div>
        </div>
      </>
    );
  }

  return (
    <>
      {/* Red gradient border — thick, slow blink */}
      <div
        className="absolute inset-0 z-10 pointer-events-none rounded-lg animate-[pulse_2s_ease-in-out_infinite]"
        style={{
          boxShadow: "inset 0 0 24px 10px rgba(220,38,38,0.6)",
        }}
      />
      {/* Badge with warning triangle */}
      <div className="absolute top-10 left-2 z-10">
        <div className="bg-red-600/90 text-white text-xs font-bold px-2 py-1 rounded flex items-center gap-1 animate-[pulse_2s_ease-in-out_infinite]">
          <AlertTriangle className="h-3 w-3" />
          {t("cameraOverlays.hardStop")}
        </div>
      </div>
    </>
  );
}

export function RecordingIndicator() {
  const { t } = useTranslation();

  return (
    <div className="bg-red-600 text-white text-xs font-bold px-2 py-1 rounded flex items-center gap-1">
      <div className="h-2 w-2 rounded-full bg-white animate-pulse" />
      {t("cameraOverlays.rec")}
    </div>
  );
}

export function SubtaskOverlay({
  currentSubtask,
  nextSubtask,
  parameterValues,
}: {
  currentSubtask: { order_index: number; name: string };
  nextSubtask?: { order_index: number; name: string } | null;
  parameterValues?: Record<string, string> | null;
}) {
  const { t } = useTranslation();

  return (
    <div className="bg-black/70 text-white text-xs rounded px-2 py-1">
      {parameterValues && Object.keys(parameterValues).length > 0 && (
        <div className="text-yellow-300 mb-0.5">
          {Object.entries(parameterValues)
            .map(([k, v]) => `${k}: ${v}`)
            .join(" | ")}
        </div>
      )}
      <div>
        {currentSubtask.order_index + 1}.{" "}
        <ParameterizedName
          name={currentSubtask.name}
          parameterValues={parameterValues}
        />
      </div>
      {nextSubtask && (
        <div className="text-gray-400">
          {t("cameraOverlays.next")}: {nextSubtask.order_index + 1}.{" "}
          <ParameterizedName
            name={nextSubtask.name}
            parameterValues={parameterValues}
          />
        </div>
      )}
    </div>
  );
}

export function ErrorOverlay({ message }: { message: string }) {
  return (
    <div className="absolute bottom-0 left-0 right-0 z-10 bg-yellow-500/90 text-black px-3 py-2 text-sm flex items-center gap-2 rounded-b-lg">
      <AlertTriangle className="h-4 w-4 shrink-0" />
      <span className="truncate">{message}</span>
    </div>
  );
}
