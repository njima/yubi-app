import { cva } from "class-variance-authority";

import { cn } from "@/lib/utils";

interface ConsecutiveFaultDaysBadgeProps {
  days?: number | null;
  className?: string;
}

function getFaultBadgeVariant(days: number): "blue" | "yellow" | "red" {
  if (days >= 7) {
    return "red";
  }
  if (days >= 3) {
    return "yellow";
  }
  return "blue";
}

const faultDaysBadgeVariants = cva(
  "inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium border",
  {
    variants: {
      variant: {
        blue: "bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-950 dark:text-blue-400 dark:border-blue-800",
        yellow:
          "bg-yellow-50 text-yellow-800 border-yellow-200 dark:bg-yellow-950 dark:text-yellow-400 dark:border-yellow-800",
        red: "bg-red-50 text-red-700 border-red-200 dark:bg-red-950 dark:text-red-400 dark:border-red-800",
      },
    },
  }
);

const dotVariants = cva("h-1.5 w-1.5 rounded-full", {
  variants: {
    variant: {
      blue: "bg-blue-500",
      yellow: "bg-yellow-500",
      red: "bg-red-500",
    },
  },
});

export function ConsecutiveFaultDaysBadge({
  days,
  className,
}: ConsecutiveFaultDaysBadgeProps) {
  if (days == null) {
    return null;
  }

  const variant = getFaultBadgeVariant(days);

  return (
    <span
      className={cn(faultDaysBadgeVariants({ variant }), className)}
      role="status"
      aria-label={`Consecutive fault days: ${days}`}
    >
      <span className={cn(dotVariants({ variant }))} aria-hidden="true" />
      {days <= 1 ? `${days} day` : `${days} days`}
    </span>
  );
}
