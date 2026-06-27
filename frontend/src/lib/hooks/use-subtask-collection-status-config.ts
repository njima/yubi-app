/**
 * Subtask Collection Status Config Hook
 *
 * React hook for getting subtask collection status UI configuration.
 * Uses useTranslation to ensure components re-render on language change.
 */

import { Ban, CheckCircle, Circle, Loader2, SkipForward } from "lucide-react";
import { useTranslation } from "react-i18next";

import { SUBTASK_COLLECTION_STATUS } from "@/lib/status/constants";

/**
 * Hook to get subtask collection status configuration
 * Returns UI configuration (label, icon, variant) for each status
 */
export function useSubtaskCollectionStatusConfig() {
  const { t } = useTranslation();

  return {
    [SUBTASK_COLLECTION_STATUS.READY]: {
      label: t("status.ready"),
      icon: Circle,
      variant: "secondary" as const,
    },
    [SUBTASK_COLLECTION_STATUS.IN_PROGRESS]: {
      label: t("status.running"),
      icon: Loader2,
      variant: "default" as const,
    },
    [SUBTASK_COLLECTION_STATUS.COMPLETED]: {
      label: t("status.completed"),
      icon: CheckCircle,
      variant: "secondary" as const,
    },
    [SUBTASK_COLLECTION_STATUS.SKIPPED]: {
      label: t("status.skipped"),
      icon: SkipForward,
      variant: "outline" as const,
    },
    [SUBTASK_COLLECTION_STATUS.CANCELLED]: {
      label: t("status.cancelled"),
      icon: Ban,
      variant: "destructive" as const,
    },
  } as const;
}
