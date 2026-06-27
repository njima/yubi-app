"use client";

import { Calendar, ChevronLeft, ChevronRight, X } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { cn } from "@/lib/utils";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

// --- Types ---

export type DateRange = {
  from: string; // YYYY-MM-DD
  to: string; // YYYY-MM-DD
};

type Mode = "absolute" | "relative";

type RelativePreset = {
  label: string;
  getRange: () => DateRange;
};

// --- Helpers ---

/** Subtract months safely, clamping to end-of-month to avoid overflow (e.g. Mar 31 - 1 month = Feb 28). */
function subtractMonths(base: Date, months: number): Date {
  const d = new Date(base.getFullYear(), base.getMonth() - months, 1);
  const lastDay = new Date(d.getFullYear(), d.getMonth() + 1, 0).getDate();
  d.setDate(Math.min(base.getDate(), lastDay));
  return d;
}

function formatDate(d: Date): string {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

function parseDate(s: string): Date | null {
  const d = new Date(s + "T00:00:00");
  return isNaN(d.getTime()) ? null : d;
}

function isSameDay(a: Date, b: Date): boolean {
  return (
    a.getFullYear() === b.getFullYear() &&
    a.getMonth() === b.getMonth() &&
    a.getDate() === b.getDate()
  );
}

function isBetween(d: Date, from: Date, to: Date): boolean {
  return d > from && d < to;
}

const MONTH_KEYS = [
  "january",
  "february",
  "march",
  "april",
  "may",
  "june",
  "july",
  "august",
  "september",
  "october",
  "november",
  "december",
] as const;

const DAY_KEYS = ["sun", "mon", "tue", "wed", "thu", "fri", "sat"] as const;

function getRelativePresets(t: (key: string) => string): RelativePreset[] {
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  return [
    {
      label: t("dateRangePicker.presets.today"),
      getRange: () => ({ from: formatDate(today), to: formatDate(today) }),
    },
    {
      label: t("dateRangePicker.presets.sevenDays"),
      getRange: () => {
        const d = new Date(today);
        d.setDate(d.getDate() - 6);
        return { from: formatDate(d), to: formatDate(today) };
      },
    },
    {
      label: t("dateRangePicker.presets.monthToDate"),
      getRange: () => {
        const d = new Date(today.getFullYear(), today.getMonth(), 1);
        return { from: formatDate(d), to: formatDate(today) };
      },
    },
    {
      label: t("dateRangePicker.presets.oneMonth"),
      getRange: () => {
        const d = subtractMonths(today, 1);
        return { from: formatDate(d), to: formatDate(today) };
      },
    },
    {
      label: t("dateRangePicker.presets.threeMonths"),
      getRange: () => {
        const d = subtractMonths(today, 3);
        return { from: formatDate(d), to: formatDate(today) };
      },
    },
    {
      label: t("dateRangePicker.presets.sixMonths"),
      getRange: () => {
        const d = subtractMonths(today, 6);
        return { from: formatDate(d), to: formatDate(today) };
      },
    },
    {
      label: t("dateRangePicker.presets.yearToDate"),
      getRange: () => {
        const d = new Date(today.getFullYear(), 0, 1);
        return { from: formatDate(d), to: formatDate(today) };
      },
    },
    {
      label: t("dateRangePicker.presets.oneYear"),
      getRange: () => {
        const d = subtractMonths(today, 12);
        return { from: formatDate(d), to: formatDate(today) };
      },
    },
  ];
}

// --- Calendar Grid ---

function CalendarMonth({
  year,
  month,
  rangeFrom,
  rangeTo,
  onSelect,
}: {
  year: number;
  month: number;
  rangeFrom: Date | null;
  rangeTo: Date | null;
  onSelect: (d: Date) => void;
}) {
  const { t } = useTranslation();
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const firstDay = new Date(year, month, 1);
  const startDow = firstDay.getDay();
  const daysInMonth = new Date(year, month + 1, 0).getDate();

  const cells: (number | null)[] = [];
  for (let i = 0; i < startDow; i++) cells.push(null);
  for (let d = 1; d <= daysInMonth; d++) cells.push(d);

  return (
    <div>
      <div className="text-center text-sm font-semibold text-gray-900 dark:text-gray-100 mb-2">
        {t(`dateRangePicker.months.${MONTH_KEYS[month]}`)} {year}
      </div>
      <div className="grid grid-cols-7 gap-0">
        {DAY_KEYS.map((key) => (
          <div
            key={key}
            className="text-center text-[11px] font-medium text-gray-400 dark:text-gray-500 py-1"
          >
            {t(`dateRangePicker.days.${key}`)}
          </div>
        ))}
        {cells.map((day, i) => {
          if (day === null) {
            return <div key={`empty-${i}`} />;
          }
          const date = new Date(year, month, day);
          const isToday = isSameDay(date, today);
          const isStart = rangeFrom && isSameDay(date, rangeFrom);
          const isEnd = rangeTo && isSameDay(date, rangeTo);
          const isInRange =
            rangeFrom && rangeTo && isBetween(date, rangeFrom, rangeTo);

          const isSingle = isStart && isEnd;
          const isStartOnly = isStart && !isEnd;
          const isEndOnly = isEnd && !isStart;

          return (
            <button
              key={day}
              type="button"
              onClick={() => onSelect(date)}
              className={cn(
                "h-8 w-full text-xs transition-colors",
                isSingle
                  ? "bg-blue-600 text-white font-semibold rounded-md"
                  : isStartOnly
                    ? "bg-blue-600 text-white font-semibold rounded-l-md"
                    : isEndOnly
                      ? "bg-blue-600 text-white font-semibold rounded-r-md"
                      : isInRange
                        ? "bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300"
                        : isToday
                          ? "font-semibold text-blue-600 dark:text-blue-400 rounded-md"
                          : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md"
              )}
            >
              {day}
            </button>
          );
        })}
      </div>
    </div>
  );
}

// --- Main Component ---

export function DateRangePicker({
  value,
  onChange,
  onClear,
  disabled = false,
}: {
  value: DateRange | undefined;
  onChange: (range: DateRange) => void;
  onClear?: () => void;
  disabled?: boolean;
}) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const [mode, setMode] = useState<Mode>("absolute");
  const [draftFrom, setDraftFrom] = useState(value?.from ?? "");
  const [draftTo, setDraftTo] = useState(value?.to ?? "");
  const [selectingEnd, setSelectingEnd] = useState(false);

  // Calendar view month: show from-date month on left
  const initDate = parseDate(value?.from ?? "") ?? new Date();
  const [viewYear, setViewYear] = useState(initDate.getFullYear());
  const [viewMonth, setViewMonth] = useState(initDate.getMonth());

  const handleOpenChange = (nextOpen: boolean) => {
    if (nextOpen) {
      setDraftFrom(value?.from ?? "");
      setDraftTo(value?.to ?? "");
      setSelectingEnd(false);
      const d = parseDate(value?.from ?? "") ?? new Date();
      setViewYear(d.getFullYear());
      setViewMonth(d.getMonth());
    }
    setOpen(nextOpen);
  };

  const handleApply = () => {
    if (draftFrom && draftTo && draftFrom <= draftTo) {
      onChange({ from: draftFrom, to: draftTo });
    }
    setOpen(false);
  };

  const handleCancel = () => {
    setOpen(false);
  };

  const handleCalendarSelect = (d: Date) => {
    const ds = formatDate(d);
    if (!selectingEnd) {
      setDraftFrom(ds);
      setDraftTo("");
      setSelectingEnd(true);
    } else {
      if (ds < draftFrom) {
        setDraftFrom(ds);
        setDraftTo(draftFrom);
      } else {
        setDraftTo(ds);
      }
      setSelectingEnd(false);
    }
  };

  const handleRelativeSelect = (preset: RelativePreset) => {
    const range = preset.getRange();
    setDraftFrom(range.from);
    setDraftTo(range.to);
    // Update calendar view
    const d = parseDate(range.from);
    if (d) {
      setViewYear(d.getFullYear());
      setViewMonth(d.getMonth());
    }
  };

  const prevMonth = () => {
    if (viewMonth === 0) {
      setViewYear(viewYear - 1);
      setViewMonth(11);
    } else {
      setViewMonth(viewMonth - 1);
    }
  };

  const nextMonth = () => {
    if (viewMonth === 11) {
      setViewYear(viewYear + 1);
      setViewMonth(0);
    } else {
      setViewMonth(viewMonth + 1);
    }
  };

  const nextMonthIdx = viewMonth === 11 ? 0 : viewMonth + 1;
  const nextMonthYear = viewMonth === 11 ? viewYear + 1 : viewYear;

  const rangeFrom = parseDate(draftFrom);
  const rangeTo = parseDate(draftTo);

  const hasValue = !!(value?.from && value?.to);
  const displayText = hasValue
    ? `${value!.from} — ${value!.to}`
    : t("dateRangePicker.selectRange");

  return (
    <div className="inline-flex items-center gap-1">
      <Popover open={open} onOpenChange={handleOpenChange}>
        <PopoverTrigger asChild>
          <button
            type="button"
            disabled={disabled}
            className="inline-flex items-center gap-2 h-9 px-3 rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors disabled:cursor-not-allowed disabled:opacity-50"
          >
            <Calendar className="h-4 w-4 text-gray-400" />
            <span className="tabular-nums">{displayText}</span>
          </button>
        </PopoverTrigger>
        <PopoverContent className="w-145 p-0">
          {/* Mode tabs */}
          <div className="flex border-b border-gray-200 dark:border-gray-700">
            <button
              type="button"
              onClick={() => setMode("absolute")}
              className={cn(
                "flex-1 py-2 text-sm font-medium transition-colors",
                mode === "absolute"
                  ? "text-blue-600 dark:text-blue-400 border-b-2 border-blue-600 dark:border-blue-400"
                  : "text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300"
              )}
            >
              {t("dateRangePicker.absolute")}
            </button>
            <button
              type="button"
              onClick={() => setMode("relative")}
              className={cn(
                "flex-1 py-2 text-sm font-medium transition-colors",
                mode === "relative"
                  ? "text-blue-600 dark:text-blue-400 border-b-2 border-blue-600 dark:border-blue-400"
                  : "text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300"
              )}
            >
              {t("dateRangePicker.relative")}
            </button>
          </div>

          <div className="p-4">
            {mode === "absolute" ? (
              <>
                {/* Calendar navigation */}
                <div className="flex items-center justify-between mb-3">
                  <button
                    type="button"
                    onClick={prevMonth}
                    className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700"
                  >
                    <ChevronLeft className="h-4 w-4 text-gray-500" />
                  </button>
                  <div className="flex-1" />
                  <button
                    type="button"
                    onClick={nextMonth}
                    className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700"
                  >
                    <ChevronRight className="h-4 w-4 text-gray-500" />
                  </button>
                </div>

                {/* Two-month calendar */}
                <div className="grid grid-cols-2 gap-6">
                  <CalendarMonth
                    year={viewYear}
                    month={viewMonth}
                    rangeFrom={rangeFrom}
                    rangeTo={rangeTo}
                    onSelect={handleCalendarSelect}
                  />
                  <CalendarMonth
                    year={nextMonthYear}
                    month={nextMonthIdx}
                    rangeFrom={rangeFrom}
                    rangeTo={rangeTo}
                    onSelect={handleCalendarSelect}
                  />
                </div>

                {/* Start / End date inputs */}
                <div className="grid grid-cols-2 gap-4 mt-4">
                  <div>
                    <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
                      {t("dateRangePicker.startDate")}
                    </label>
                    <Input
                      type="date"
                      value={draftFrom}
                      onChange={(e) => setDraftFrom(e.target.value)}
                      className="text-sm"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
                      {t("dateRangePicker.endDate")}
                    </label>
                    <Input
                      type="date"
                      value={draftTo}
                      onChange={(e) => setDraftTo(e.target.value)}
                      className="text-sm"
                    />
                  </div>
                </div>
              </>
            ) : (
              /* Relative presets */
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400 mb-3">
                  {t("dateRangePicker.relativeDescription")}
                </p>
                <div className="flex flex-wrap gap-2">
                  {getRelativePresets(t).map((preset) => {
                    const range = preset.getRange();
                    const isActive =
                      draftFrom === range.from && draftTo === range.to;
                    return (
                      <button
                        key={preset.label}
                        type="button"
                        onClick={() => handleRelativeSelect(preset)}
                        className={cn(
                          "px-3 py-1.5 rounded-md text-sm transition-colors",
                          isActive
                            ? "bg-blue-600 text-white"
                            : "bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600"
                        )}
                      >
                        {preset.label}
                      </button>
                    );
                  })}
                </div>

                {/* Preview */}
                {draftFrom && draftTo && (
                  <div className="mt-4 p-3 rounded-md bg-gray-50 dark:bg-gray-700/50 text-sm text-gray-600 dark:text-gray-400">
                    {draftFrom} — {draftTo}
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Footer */}
          <div className="flex justify-end gap-2 px-4 py-3 border-t border-gray-200 dark:border-gray-700">
            <Button variant="outline" size="sm" onClick={handleCancel}>
              {t("dialog.cancel")}
            </Button>
            <Button
              size="sm"
              onClick={handleApply}
              disabled={!draftFrom || !draftTo || draftFrom > draftTo}
            >
              {t("dateRangePicker.apply")}
            </Button>
          </div>
        </PopoverContent>
      </Popover>
      {onClear && hasValue && (
        <button
          type="button"
          onClick={onClear}
          disabled={disabled}
          aria-label={t("dateRangePicker.clear")}
          className="inline-flex items-center justify-center h-9 w-9 rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors ring-offset-white focus:outline-none focus:ring-2 focus:ring-gray-950 focus:ring-offset-2 dark:ring-offset-gray-950 dark:focus:ring-gray-300 disabled:cursor-not-allowed disabled:opacity-50"
        >
          <X className="h-4 w-4" />
        </button>
      )}
    </div>
  );
}
