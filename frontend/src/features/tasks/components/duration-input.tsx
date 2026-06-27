"use client";

import {
  type Control,
  type FieldPath,
  type FieldValues,
  useController,
} from "react-hook-form";
import { useTranslation } from "react-i18next";

import { FormItem, FormLabel } from "@/components/ui/form";
import { Input } from "@/components/ui/input";

interface DurationInputProps<T extends FieldValues> {
  control: Control<T>;
  hoursName: FieldPath<T>;
  minutesName: FieldPath<T>;
  label?: string;
}

export function DurationInput<T extends FieldValues>({
  control,
  hoursName,
  minutesName,
  label,
}: DurationInputProps<T>) {
  const { t } = useTranslation();
  const { field: hoursField } = useController({ control, name: hoursName });
  const { field: minutesField } = useController({ control, name: minutesName });
  const resolvedLabel = label ?? t("durationInput.targetDurationOptional");

  return (
    <FormItem>
      <FormLabel>{resolvedLabel}</FormLabel>
      <div className="flex items-center gap-2">
        <Input
          type="number"
          min={0}
          placeholder="0"
          className="w-20"
          value={hoursField.value ?? ""}
          onChange={(e) =>
            hoursField.onChange(
              e.target.value === "" ? undefined : Number(e.target.value)
            )
          }
        />
        <span className="text-sm text-gray-500">
          {t("durationInput.hoursShort")}
        </span>
        <Input
          type="number"
          min={0}
          max={59}
          placeholder="0"
          className="w-20"
          value={minutesField.value ?? ""}
          onChange={(e) =>
            minutesField.onChange(
              e.target.value === "" ? undefined : Number(e.target.value)
            )
          }
        />
        <span className="text-sm text-gray-500">
          {t("durationInput.minutesShort")}
        </span>
      </div>
    </FormItem>
  );
}
