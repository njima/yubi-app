"use client";

import { useMutation } from "@tanstack/react-query";

export type ExportOperatorYieldFilter = {
  date_from: string;
  date_to: string;
  location_id?: string;
  task_id?: string;
  user_id?: string;
};

function filenameFromContentDisposition(header: string | null): string {
  if (!header) return "operator_yield_export.csv";
  const match = /filename="?([^";]+)"?/i.exec(header);
  return match?.[1] ?? "operator_yield_export.csv";
}

export function useExportOperatorYieldMutation() {
  return useMutation<void, Error, ExportOperatorYieldFilter>({
    mutationFn: async (filter: ExportOperatorYieldFilter) => {
      const params = new URLSearchParams();
      params.set("date_from", filter.date_from);
      params.set("date_to", filter.date_to);
      if (filter.location_id) params.set("location_id", filter.location_id);
      if (filter.task_id) params.set("task_id", filter.task_id);
      if (filter.user_id) params.set("user_id", filter.user_id);

      const response = await fetch(
        `/web/api/reports/operator-yield/export?${params.toString()}`
      );
      if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw new Error(
          error.message || `Export failed: ${response.statusText}`
        );
      }

      const blob = await response.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = filenameFromContentDisposition(
        response.headers.get("Content-Disposition")
      );
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    },
  });
}
