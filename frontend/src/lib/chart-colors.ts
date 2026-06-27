/**
 * Chart Color Palette
 *
 * Centralized color definitions for SVG charts (hex codes).
 * Tailwind CSS class colors (for badges, etc.) are managed separately via cva variants.
 */

/** Series colors for line/bar/area charts */
export const CHART_PALETTE: string[] = [
  "#3b82f6", // blue
  "#ef4444", // red
  "#8b5cf6", // purple
  "#06b6d4", // cyan
  "#f59e0b", // amber
  "#10b981", // emerald
];

/** Threshold-based colors for donut charts, efficiency badges, etc. */
export const THRESHOLD_COLORS = {
  good: "#3b82f6", // blue
  warning: "#f59e0b", // amber
  danger: "#ef4444", // red
  track: "#e5e7eb", // light gray
  muted: "#6b7280", // gray
} as const;

/** Assign colors from CHART_PALETTE to series data */
export function assignColors<T>(series: T[]): (T & { color: string })[] {
  return series.map((s, i) => ({
    ...s,
    color: CHART_PALETTE[i % CHART_PALETTE.length] ?? "#6b7280",
  }));
}
