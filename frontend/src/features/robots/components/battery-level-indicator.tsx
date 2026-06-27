import { Battery, BatteryLow, BatteryWarning } from "lucide-react";

import { cn } from "@/lib/utils";

/**
 * Battery Level Indicator
 * Displays battery level with progress bar and threshold-based coloring
 *
 * Thresholds:
 * - < 10%: Critical (red, alert icon)
 * - < 30%: Warning (yellow, warning icon)
 * - >= 30%: Normal (green, battery icon)
 */

interface BatteryLevelIndicatorProps {
  level: number;
  className?: string;
  showIcon?: boolean;
}

export function BatteryLevelIndicator({
  level,
  className,
  showIcon = true,
}: BatteryLevelIndicatorProps) {
  const percentage = Math.max(0, Math.min(100, level));

  const getColorClasses = () => {
    if (percentage < 10) {
      return {
        bg: "bg-red-500",
        text: "text-red-700 dark:text-red-400",
        icon: BatteryLow,
      };
    }
    if (percentage < 30) {
      return {
        bg: "bg-yellow-500",
        text: "text-yellow-700 dark:text-yellow-400",
        icon: BatteryWarning,
      };
    }
    return {
      bg: "bg-green-500",
      text: "text-green-700 dark:text-green-400",
      icon: Battery,
    };
  };

  const { bg, text, icon: Icon } = getColorClasses();

  return (
    <div className={cn("flex items-center gap-2", className)}>
      {showIcon && <Icon className={cn("h-4 w-4", text)} aria-hidden="true" />}
      <div className="flex items-center gap-2 flex-1">
        <div className="relative h-2 w-full bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
          <div
            className={cn("h-full transition-all duration-300", bg)}
            style={{ width: `${percentage}%` }}
            role="progressbar"
            aria-valuenow={percentage}
            aria-valuemin={0}
            aria-valuemax={100}
            aria-label={`Battery level: ${percentage}%`}
          />
        </div>
        <span
          className={cn("text-xs font-medium tabular-nums min-w-[3ch]", text)}
        >
          {percentage}%
        </span>
      </div>
    </div>
  );
}
