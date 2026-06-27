"use client";

import { useTranslation } from "react-i18next";

import { Slider } from "@/components/ui/slider";
import { Textarea } from "@/components/ui/textarea";

import { getGradeColor } from "../lib/grade-color";

const valueColorClass: Record<
  Exclude<ReturnType<typeof getGradeColor>, "none">,
  string
> = {
  green: "text-green-600 dark:text-green-400",
  yellow: "text-yellow-600 dark:text-yellow-400",
  red: "text-red-600 dark:text-red-400",
};

interface GradeSliderProps {
  value: number;
  onChange: (next: number) => void;
  comment: string;
  onCommentChange: (next: string) => void;
  disabled?: boolean;
}

export function GradeSlider({
  value,
  onChange,
  comment,
  onCommentChange,
  disabled,
}: GradeSliderProps) {
  const { t } = useTranslation();
  const color = getGradeColor(value);
  const valueClass =
    color === "none" ? "text-gray-500" : valueColorClass[color];

  return (
    <div className="space-y-3">
      <div className="flex items-center gap-4">
        <Slider
          value={[value]}
          min={0}
          max={1}
          step={0.01}
          onValueChange={(next) => onChange(next[0] ?? 0)}
          disabled={disabled}
          className="flex-1"
        />
        <div
          className={`w-16 rounded-md border border-gray-200 dark:border-gray-700 px-2 py-1 text-right text-sm font-medium tabular-nums ${valueClass}`}
        >
          {value.toFixed(2)}
        </div>
      </div>
      <Textarea
        placeholder={t("gradeSlider.commentPlaceholder")}
        value={comment}
        onChange={(e) => onCommentChange(e.target.value)}
        disabled={disabled}
        rows={3}
      />
      <p className="text-xs text-gray-500 dark:text-gray-400">
        {t("gradeSlider.help")}
      </p>
    </div>
  );
}
