"use client";

import { cn } from "@/lib/utils";

const categoryColors: Record<string, string> = {
  Application: "border-l-blue-400 dark:border-l-blue-500",
  "Basic Skill": "border-l-emerald-400 dark:border-l-emerald-500",
};

interface TagCategoryBadgeProps {
  categoryTypeName: string;
  name: string;
}

export function TagCategoryBadge({
  categoryTypeName,
  name,
}: TagCategoryBadgeProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded border border-gray-200 dark:border-gray-700",
        "border-l-2 px-1.5 py-0.5 text-xs text-gray-600 dark:text-gray-400",
        categoryColors[categoryTypeName] ?? "border-l-gray-400"
      )}
    >
      {name}
    </span>
  );
}
