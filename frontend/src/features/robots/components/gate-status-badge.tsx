import { cva } from "class-variance-authority";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";
import { cn } from "@/lib/utils";

import { Badge } from "@/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

type GateConditionStatus = z.infer<typeof schemas.GateConditionStatus>;
type GateGroupStatus = z.infer<typeof schemas.GateGroupStatus>;

const badgeVariants = cva(
  "inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-sm font-medium border",
  {
    variants: {
      level: {
        open: "text-green-600 bg-green-50 border-green-200 dark:text-green-400 dark:bg-green-950 dark:border-green-800",
        block:
          "text-yellow-600 bg-yellow-50 border-yellow-200 dark:text-yellow-400 dark:bg-yellow-950 dark:border-yellow-800",
        stop: "text-red-600 bg-red-50 border-red-200 dark:text-red-400 dark:bg-red-950 dark:border-red-800",
      },
    },
    defaultVariants: {
      level: "open",
    },
  }
);

const badgeVariantsInRobotList = cva("", {
  variants: {
    level: {
      open: "bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300",
      block:
        "bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300",
      stop: "bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300",
    },
  },
});

const dotVariants = cva("h-1.5 w-1.5 rounded-full", {
  variants: {
    level: {
      open: "bg-green-500 animate-pulse",
      block: "bg-yellow-500",
      stop: "bg-red-500",
    },
  },
});

const LEVEL_MAP = {
  0: { variant: "open", label: "Open" },
  1: { variant: "block", label: "Block" },
  2: { variant: "stop", label: "Hard Stop" },
} as const;

function levelVariant(level: number) {
  return LEVEL_MAP[level as keyof typeof LEVEL_MAP] ?? LEVEL_MAP[0];
}

interface GateStatusBadgeProps {
  gateConditions?: GateConditionStatus | null;
  className?: string;
}

export function GateStatusBadge({
  gateConditions,
  className,
}: GateStatusBadgeProps) {
  if (!gateConditions) {
    return null;
  }

  const { variant, label } = levelVariant(gateConditions.gate_level);

  return (
    <span
      className={cn(badgeVariants({ level: variant }), className)}
      aria-label={`Gate status: ${label}`}
      role="status"
    >
      <span
        className={cn(dotVariants({ level: variant }))}
        aria-hidden="true"
      />
      {label}
    </span>
  );
}

// Gate status badge used in the robot list view
export function GateStatusBadgeInRobotList({
  gateConditions,
  className,
}: GateStatusBadgeProps) {
  if (!gateConditions) {
    return null;
  }

  const { variant, label } = levelVariant(gateConditions.gate_level);

  return (
    <Badge
      variant="outline"
      className={cn(badgeVariantsInRobotList({ level: variant }), className)}
    >
      {label}
    </Badge>
  );
}

interface GateGroupGridProps {
  gateConditions?: GateConditionStatus | null;
}

export function GateGroupGrid({ gateConditions }: GateGroupGridProps) {
  if (!gateConditions || Object.keys(gateConditions.groups).length === 0) {
    return null;
  }

  return (
    <div className="grid grid-cols-2 gap-2">
      {Object.entries(gateConditions.groups).map(([name, group]) => (
        <GateGroupCell key={name} name={name} group={group} />
      ))}
    </div>
  );
}

export function GateGroupCell({
  name,
  group,
  className,
}: {
  name: string;
  group: GateGroupStatus;
  className?: string;
}) {
  const { variant, label } = levelVariant(group.level);

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <div
            className={cn(
              "flex items-center gap-2 rounded-md border px-2.5 py-1.5 text-xs cursor-default",
              group.level === 0 && "border-green-200 dark:border-green-800",
              group.level === 1 && "border-yellow-200 dark:border-yellow-800",
              group.level >= 2 && "border-red-200 dark:border-red-800",
              className
            )}
          >
            <span
              className={cn(dotVariants({ level: variant }))}
              aria-hidden="true"
            />
            <span className="font-medium truncate">{name}</span>
            <span className="ml-auto text-gray-500 dark:text-gray-400">
              {label}
            </span>
          </div>
        </TooltipTrigger>
        <TooltipContent side="bottom" className="max-w-xs">
          <div className="space-y-1.5">
            <div className="font-semibold text-sm">
              {name}{" "}
              <span className="font-normal text-gray-500">
                {group.settled ? "(settled)" : "(unsettled)"}
              </span>
            </div>
            {group.conditions.map((c) => (
              <div key={c.name} className="flex items-start gap-2 text-xs">
                <span
                  className={cn(
                    "mt-1 h-1.5 w-1.5 shrink-0 rounded-full",
                    c.passed
                      ? "bg-green-500"
                      : c.escalation >= 2
                        ? "bg-red-500"
                        : "bg-yellow-500"
                  )}
                />
                <div>
                  <span className="font-medium">{c.name}</span>
                  <span className="ml-1 text-gray-500">{c.reason}</span>
                </div>
              </div>
            ))}
          </div>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
}
