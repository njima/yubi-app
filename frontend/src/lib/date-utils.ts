/**
 * Date utility functions
 */
import i18n from "./i18n";

function lang(): string {
  return i18n.language ?? "en";
}

/**
 * Format a date string to a full date and time (e.g., "Jan 15, 2024, 10:30 AM")
 */
export function formatDateTime(dateString: string | null | undefined): string {
  if (!dateString) return i18n.t("dateUtils.notAvailable");

  try {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat(lang(), {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "numeric",
      minute: "2-digit",
      hour12: true,
    }).format(date);
  } catch {
    return i18n.t("dateUtils.invalidDate");
  }
}

/**
 * Format uptime from seconds to human-readable format (e.g., "1h 2m 3s")
 * Unlike formatDuration which takes two timestamps, this takes raw seconds
 */
export function formatUptime(seconds: number): string {
  if (seconds <= 0) return i18n.t("dateUtils.uptimeZero");

  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;

  const parts: string[] = [];
  if (days > 0) parts.push(i18n.t("dateUtils.uptimeDay", { count: days }));
  if (hours > 0) parts.push(i18n.t("dateUtils.uptimeHour", { count: hours }));
  if (minutes > 0)
    parts.push(i18n.t("dateUtils.uptimeMinute", { count: minutes }));
  if (secs > 0 && parts.length < 3)
    parts.push(i18n.t("dateUtils.uptimeSecond", { count: secs }));

  return parts.length > 0 ? parts.join(" ") : i18n.t("dateUtils.uptimeZero");
}

/**
 * Format a time string to time only (e.g., "10:30:15 AM")
 */
export function formatTime(dateString: string | null | undefined): string {
  if (!dateString) return i18n.t("dateUtils.notAvailable");

  try {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat(lang(), {
      hour: "numeric",
      minute: "2-digit",
      second: "2-digit",
      hour12: true,
    }).format(date);
  } catch {
    return i18n.t("dateUtils.invalidDate");
  }
}

/**
 * Format the duration between two timestamps (e.g., "2m 30s")
 * Returns "-" if either timestamp is missing
 */
export function formatDuration(
  startedAt: string | null | undefined,
  finishedAt: string | null | undefined
): string {
  if (!startedAt || !finishedAt) return "-";

  try {
    const start = new Date(startedAt);
    const end = new Date(finishedAt);
    const diffInSeconds = Math.floor((end.getTime() - start.getTime()) / 1000);

    if (diffInSeconds < 0) return "-";

    const minutes = Math.floor(diffInSeconds / 60);
    const seconds = diffInSeconds % 60;

    if (minutes === 0) return i18n.t("dateUtils.durationSecond", { seconds });
    return i18n.t("dateUtils.durationMinuteSecond", { minutes, seconds });
  } catch {
    return "-";
  }
}

/**
 * Convert an ISO date string to a value suitable for <input type="datetime-local" />.
 * Returns "" if the input is not a valid date.
 */
export function toDateTimeLocalValue(dateText: string): string {
  const d = new Date(dateText);
  if (Number.isNaN(d.getTime())) {
    return "";
  }
  const tzOffsetMs = d.getTimezoneOffset() * 60 * 1000;
  return new Date(d.getTime() - tzOffsetMs).toISOString().slice(0, 16);
}
