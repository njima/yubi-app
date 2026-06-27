/**
 * DateTime Picker Component
 * Allows input of date and time in ISO8601 format
 */

"use client";

import { format } from "date-fns";
import { Calendar as CalendarIcon } from "lucide-react";
import * as React from "react";

import { cn } from "@/lib/utils";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

interface DateTimePickerProps {
  value?: string; // ISO8601 string
  onChange?: (value: string) => void;
  className?: string;
}

export function DateTimePicker({
  value,
  onChange,
  className,
}: DateTimePickerProps) {
  const [dateInput, setDateInput] = React.useState("");
  const [timeInput, setTimeInput] = React.useState("");

  // Parse ISO8601 string to date and time inputs
  React.useEffect(() => {
    if (value) {
      try {
        const date = new Date(value);
        // Format as YYYY-MM-DD
        const dateStr = format(date, "yyyy-MM-dd");
        // Format as HH:mm
        const timeStr = format(date, "HH:mm");
        setDateInput(dateStr);
        setTimeInput(timeStr);
      } catch {
        // Invalid date, reset
        setDateInput("");
        setTimeInput("");
      }
    }
  }, [value]);

  const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newDate = e.target.value;
    setDateInput(newDate);
    updateISO8601(newDate, timeInput);
  };

  const handleTimeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newTime = e.target.value;
    setTimeInput(newTime);
    updateISO8601(dateInput, newTime);
  };

  const updateISO8601 = (date: string, time: string) => {
    if (!date) {
      onChange?.("");
      return;
    }

    // Default time to 00:00 if not provided
    const finalTime = time || "00:00";

    try {
      // Combine date and time into ISO8601 format
      const dateTime = new Date(`${date}T${finalTime}`);
      if (!isNaN(dateTime.getTime())) {
        onChange?.(dateTime.toISOString());
      }
    } catch {
      // Invalid date/time combination
      onChange?.("");
    }
  };

  const handleClear = () => {
    setDateInput("");
    setTimeInput("");
    onChange?.("");
  };

  return (
    <div className={cn("flex gap-2", className)}>
      <div className="flex-1 relative">
        <Input
          type="date"
          value={dateInput}
          onChange={handleDateChange}
          placeholder="YYYY-MM-DD"
          className="pr-10"
        />
        <CalendarIcon className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400 pointer-events-none" />
      </div>
      <Input
        type="time"
        value={timeInput}
        onChange={handleTimeChange}
        placeholder="HH:mm"
        className="w-[120px]"
      />
      {value && (
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={handleClear}
          className="px-3"
        >
          Clear
        </Button>
      )}
    </div>
  );
}
