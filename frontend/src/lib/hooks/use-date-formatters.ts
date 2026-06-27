/**
 * Date Formatter Hooks
 *
 * React hooks for formatting dates with proper i18n support.
 * These hooks use useTranslation to ensure components re-render on language change.
 */

import { formatDistanceToNow } from "date-fns";
import { enUS, ja } from "date-fns/locale";
import { useCallback } from "react";
import { useTranslation } from "react-i18next";

// Each hook wraps its returned formatter in useCallback so consumers get a
// referentially stable function across renders. Without this, putting the
// formatter in a useMemo/useCallback/useEffect dep array (e.g. TanStack
// Table's columns memo) thrashes on every parent render — the memo
// recomputes, cell function identities change, and React remounts every
// cell. That tore down local state in any nested component (notably the
// row's EditRobotDialog snapping shut on every render).

/**
 * Hook to format relative time (e.g., "2 days ago")
 */
export function useFormatRelativeTime() {
  const { t, i18n } = useTranslation();
  const lang = i18n.language ?? "en";

  return useCallback(
    (dateString: string | null | undefined): string => {
      if (!dateString) return t("dateUtils.notAvailable");

      try {
        const date = new Date(dateString);
        const now = new Date();
        const diffInMs = now.getTime() - date.getTime();
        const diffInSeconds = Math.floor(diffInMs / 1000);
        const diffInMinutes = Math.floor(diffInSeconds / 60);
        const diffInHours = Math.floor(diffInMinutes / 60);
        const diffInDays = Math.floor(diffInHours / 24);

        // Less than 1 minute
        if (diffInSeconds < 60) {
          return t("dateUtils.justNow");
        }

        // Less than 1 hour
        if (diffInMinutes < 60) {
          return t("dateUtils.minutesAgo", { count: diffInMinutes });
        }

        // Less than 24 hours
        if (diffInHours < 24) {
          return t("dateUtils.hoursAgo", { count: diffInHours });
        }

        // Less than 7 days
        if (diffInDays < 7) {
          return t("dateUtils.daysAgo", { count: diffInDays });
        }

        // Older than 7 days, show formatted date
        return new Intl.DateTimeFormat(lang === "ja" ? "ja-JP" : "en-US", {
          year: "numeric",
          month: "short",
          day: "numeric",
        }).format(date);
      } catch {
        return t("dateUtils.invalidDate");
      }
    },
    [t, lang]
  );
}

/**
 * Hook to format distance time (e.g., "2 minutes ago")
 */
export function useFormatDistanceTime() {
  const { t, i18n } = useTranslation();
  const lang = i18n.language ?? "en";

  return useCallback(
    (date: string | Date | null | undefined): string => {
      if (!date) return t("dateUtils.notAvailable");

      try {
        const dateObj = typeof date === "string" ? new Date(date) : date;
        return formatDistanceToNow(dateObj, {
          addSuffix: true,
          locale: lang === "ja" ? ja : enUS,
        });
      } catch {
        return t("dateUtils.invalidDate");
      }
    },
    [t, lang]
  );
}

/**
 * Hook to format absolute time (e.g., "2024-01-15 14:30:00")
 */
export function useFormatAbsoluteTime() {
  const { t, i18n } = useTranslation();
  const lang = i18n.language ?? "en";

  return useCallback(
    (date: string | Date | null | undefined): string => {
      if (!date) return t("dateUtils.notAvailable");

      try {
        const dateObj = typeof date === "string" ? new Date(date) : date;
        return new Intl.DateTimeFormat(lang === "ja" ? "ja-JP" : "en-US", {
          year: "numeric",
          month: "2-digit",
          day: "2-digit",
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
          hour12: false,
        }).format(dateObj);
      } catch {
        return t("dateUtils.invalidDate");
      }
    },
    [t, lang]
  );
}

/**
 * Hook to format date (e.g., "Jan 15, 2024")
 */
export function useFormatDate() {
  const { t, i18n } = useTranslation();
  const lang = i18n.language ?? "en";

  return useCallback(
    (dateString: string | null | undefined): string => {
      if (!dateString) return t("dateUtils.notAvailable");

      try {
        const date = new Date(dateString);
        return new Intl.DateTimeFormat(lang === "ja" ? "ja-JP" : "en-US", {
          year: "numeric",
          month: "short",
          day: "numeric",
        }).format(date);
      } catch {
        return t("dateUtils.invalidDate");
      }
    },
    [t, lang]
  );
}
