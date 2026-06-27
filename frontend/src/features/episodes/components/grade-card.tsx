"use client";

import type { EpisodeGrade } from "@/lib/api/backend-client";
import { useFormatRelativeTime } from "@/lib/hooks/use-date-formatters";

import { GradeBar } from "./grade-bar";

interface GradeCardProps {
  grade: EpisodeGrade;
}

export function GradeCard({ grade }: GradeCardProps) {
  const formatRelative = useFormatRelativeTime();

  return (
    <div className="rounded-md border border-gray-200 dark:border-gray-700 p-3">
      <div className="flex items-center justify-between gap-3">
        <div className="flex items-center gap-3 min-w-0">
          <span className="font-medium truncate">{grade.user_name}</span>
          <GradeBar value={grade.grade} size="sm" />
        </div>
        <span className="text-xs text-gray-500 dark:text-gray-400 whitespace-nowrap">
          {formatRelative(grade.graded_at)}
        </span>
      </div>
      {grade.comment ? (
        <p className="mt-2 text-sm text-gray-700 dark:text-gray-300 whitespace-pre-wrap">
          {grade.comment}
        </p>
      ) : null}
    </div>
  );
}
