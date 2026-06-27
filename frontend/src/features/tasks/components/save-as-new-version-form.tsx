"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm, useWatch } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { SearchableSelect } from "@/components/ui/searchable-select";

import { DurationInput } from "./duration-input";
import { useCreateTaskVersionMutation } from "../hooks/use-create-task-version-mutation";
import { useSubTasksByVersionQuery } from "../hooks/use-subtasks-by-version-query";
import { useTaskVersionsQuery } from "../hooks/use-task-versions-query";
import { hoursMinutesToSeconds } from "../lib/duration";

import type { TaskVersion } from "../schemas";

const saveAsNewVersionFormSchema = z.object({
  version: z.string().min(1).max(50),
  display_name: z.string().max(100).optional(),
  base_task_version_id: z.string().min(1),
  duration_hours: z.coerce.number().int().min(0).optional(),
  duration_minutes: z.coerce.number().int().min(0).max(59).optional(),
  target_episode_count: z.coerce.number().int().min(1).optional(),
  per_episode_hours: z.coerce.number().int().min(0).optional(),
  per_episode_minutes: z.coerce.number().int().min(0).max(59).optional(),
});

type SaveAsNewVersionFormValues = z.infer<typeof saveAsNewVersionFormSchema>;

interface SaveAsNewVersionFormProps {
  taskId: string;
  versions: TaskVersion[];
  defaultBaseVersionId: string;
  onSuccess?: (newVersionId: string) => void;
  onCancel?: () => void;
}

export function SaveAsNewVersionForm({
  taskId,
  versions,
  defaultBaseVersionId,
  onSuccess,
  onCancel,
}: SaveAsNewVersionFormProps) {
  const { t } = useTranslation();
  const { mutate, isPending } = useCreateTaskVersionMutation(taskId);

  const { data: allVersions } = useTaskVersionsQuery(taskId);

  const form = useForm<SaveAsNewVersionFormValues>({
    resolver: zodResolver(saveAsNewVersionFormSchema),
    defaultValues: {
      version: "",
      display_name: "",
      base_task_version_id: defaultBaseVersionId,
      duration_hours: undefined,
      duration_minutes: undefined,
      target_episode_count: undefined,
      per_episode_hours: undefined,
      per_episode_minutes: undefined,
    },
  });

  const selectedBaseVersionId = useWatch({
    control: form.control,
    name: "base_task_version_id",
  });
  const { data: subtasks, isLoading: subtasksLoading } =
    useSubTasksByVersionQuery(selectedBaseVersionId);

  const baseVersion = allVersions?.find((v) => v.id === selectedBaseVersionId);
  const baseParameters = baseVersion?.parameters ?? [];

  const onSubmit = (values: SaveAsNewVersionFormValues) => {
    const hours = values.duration_hours ?? 0;
    const minutes = values.duration_minutes ?? 0;
    const totalSeconds = hoursMinutesToSeconds(hours, minutes);

    const perEpHours = values.per_episode_hours ?? 0;
    const perEpMinutes = values.per_episode_minutes ?? 0;
    const perEpSeconds = hoursMinutesToSeconds(perEpHours, perEpMinutes);

    mutate(
      {
        version: values.version,
        display_name: values.display_name?.trim() || undefined,
        base_task_version_id: values.base_task_version_id,
        target_duration_seconds: totalSeconds > 0 ? totalSeconds : undefined,
        target_episode_count: values.target_episode_count,
        target_duration_per_episode_seconds:
          perEpSeconds > 0 ? perEpSeconds : undefined,
      },
      {
        onSuccess: (newVersion) => {
          onSuccess?.(newVersion.id);
        },
      }
    );
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        {/* Base Version Selection */}
        <FormField
          control={form.control}
          name="base_task_version_id"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("saveAsNewVersionForm.baseVersion")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value}
                  onValueChange={field.onChange}
                  options={versions.map((v) => ({
                    value: v.id,
                    label: `${v.version}${v.is_current ? t("saveAsNewVersionForm.currentSuffix") : ""}`,
                  }))}
                  placeholder={t("saveAsNewVersionForm.selectBaseVersion")}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* New Version Input */}
        <FormField
          control={form.control}
          name="version"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("saveAsNewVersionForm.newVersion")}</FormLabel>
              <FormControl>
                <Input placeholder="e.g., v2.0.0" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Display Name (optional) */}
        <FormField
          control={form.control}
          name="display_name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("saveAsNewVersionForm.displayName")}</FormLabel>
              <FormControl>
                <Input
                  placeholder={t("saveAsNewVersionForm.displayNamePlaceholder")}
                  maxLength={100}
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Target Settings */}
        <TargetFields form={form} />

        {/* Parameters (read-only, copied from base version) */}
        {baseParameters.length > 0 && (
          <div className="space-y-2">
            <p className="text-sm font-medium">
              {t("saveAsNewVersionForm.parametersInherited")}
            </p>
            <div className="space-y-2">
              {baseParameters.map((p) => (
                <div key={p.key} className="flex items-start gap-3">
                  <code className="text-sm bg-gray-100 dark:bg-gray-800 px-2 py-0.5 rounded font-mono">
                    {"{" + p.key + "}"}
                  </code>
                  <div className="flex flex-wrap gap-1">
                    {p.values.map((v) => (
                      <Badge key={v} variant="secondary" className="text-xs">
                        {v}
                      </Badge>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Subtasks Preview */}
        <div className="text-sm text-gray-600 dark:text-gray-400">
          <p className="font-medium mb-2">
            {t("saveAsNewVersionForm.subtasksToCopy")}
          </p>
          {subtasksLoading ? (
            <div className="space-y-1">
              {[1, 2, 3].map((i) => (
                <div
                  key={i}
                  className="h-5 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"
                />
              ))}
            </div>
          ) : subtasks && subtasks.length > 0 ? (
            <ul className="space-y-1 pl-4 list-disc max-h-40 overflow-y-auto">
              {subtasks.map((subtask) => (
                <li
                  key={subtask.id}
                  className="text-gray-700 dark:text-gray-300"
                >
                  {subtask.name}
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-gray-500 italic">
              {t("versionHistoryDialog.noSubtasks")}
            </p>
          )}
        </div>

        <div className="flex justify-end gap-2">
          {onCancel && (
            <Button
              type="button"
              variant="outline"
              onClick={onCancel}
              disabled={isPending}
            >
              {t("dialog.cancel")}
            </Button>
          )}
          <Button type="submit" disabled={isPending}>
            {isPending ? t("dialog.creating") : t("dialog.create")}
          </Button>
        </div>
      </form>
    </Form>
  );
}

function TargetFields({
  form,
}: {
  form: ReturnType<typeof useForm<SaveAsNewVersionFormValues>>;
}) {
  const { t } = useTranslation();
  const [enableTarget, setEnableTarget] = useState(false);

  return (
    <div className="space-y-3">
      <div className="flex items-center gap-2">
        <Checkbox
          checked={enableTarget}
          onCheckedChange={(checked) => {
            setEnableTarget(!!checked);
            if (!checked) {
              form.setValue("duration_hours", undefined);
              form.setValue("duration_minutes", undefined);
              form.setValue("target_episode_count", undefined);
              form.setValue("per_episode_hours", undefined);
              form.setValue("per_episode_minutes", undefined);
            }
          }}
        />
        <span className="text-sm font-medium">
          {t("editTargetsDialog.setCollectionTargets")}
        </span>
      </div>

      {enableTarget && (
        <>
          <DurationInput
            control={form.control}
            hoursName="duration_hours"
            minutesName="duration_minutes"
          />
          <DurationInput
            control={form.control}
            hoursName="per_episode_hours"
            minutesName="per_episode_minutes"
            label={t("editTargetsDialog.targetDurationPerEpisode")}
          />
          <FormField
            control={form.control}
            name="target_episode_count"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  {t("editTargetsDialog.targetEpisodeCount")}
                </FormLabel>
                <FormControl>
                  <Input
                    type="number"
                    min={1}
                    placeholder="e.g., 100"
                    className="w-32"
                    {...field}
                    value={field.value ?? ""}
                    onChange={(e) =>
                      field.onChange(
                        e.target.value === ""
                          ? undefined
                          : Number(e.target.value)
                      )
                    }
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </>
      )}
    </div>
  );
}
