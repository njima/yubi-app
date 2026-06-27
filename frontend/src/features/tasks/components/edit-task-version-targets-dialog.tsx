"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";

import { DurationInput } from "./duration-input";
import { useUpdateTaskVersionMutation } from "../hooks/use-update-task-version-mutation";
import { hoursMinutesToSeconds, secondsToHoursMinutes } from "../lib/duration";

const editTargetsFormSchema = z.object({
  display_name: z.string().max(100).optional(),
  duration_hours: z.coerce.number().int().min(0).optional(),
  duration_minutes: z.coerce.number().int().min(0).max(59).optional(),
  target_episode_count: z.coerce.number().int().min(1).optional(),
  per_episode_hours: z.coerce.number().int().min(0).optional(),
  per_episode_minutes: z.coerce.number().int().min(0).max(59).optional(),
});

type EditTargetsFormValues = z.infer<typeof editTargetsFormSchema>;

interface EditTaskVersionTargetsDialogProps {
  taskId: string;
  versionId: string;
  defaultValues: {
    display_name?: string | null;
    target_duration_seconds?: number | null;
    target_episode_count?: number | null;
    target_duration_per_episode_seconds?: number | null;
  };
  children?: React.ReactNode;
}

export function EditTaskVersionTargetsDialog({
  taskId,
  versionId,
  defaultValues,
  children,
}: EditTaskVersionTargetsDialogProps) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const hasExistingTarget =
    defaultValues.target_duration_seconds != null ||
    defaultValues.target_episode_count != null ||
    defaultValues.target_duration_per_episode_seconds != null;
  const [enableTarget, setEnableTarget] = useState(hasExistingTarget);
  const { mutate, isPending } = useUpdateTaskVersionMutation();

  const { hours: defaultHours, minutes: defaultMinutes } =
    defaultValues.target_duration_seconds
      ? secondsToHoursMinutes(defaultValues.target_duration_seconds)
      : { hours: undefined, minutes: undefined };

  const { hours: defaultPerEpHours, minutes: defaultPerEpMinutes } =
    defaultValues.target_duration_per_episode_seconds
      ? secondsToHoursMinutes(defaultValues.target_duration_per_episode_seconds)
      : { hours: undefined, minutes: undefined };

  const form = useForm<EditTargetsFormValues>({
    resolver: zodResolver(editTargetsFormSchema),
    defaultValues: {
      display_name: defaultValues.display_name ?? "",
      duration_hours: defaultHours,
      duration_minutes: defaultMinutes,
      target_episode_count: defaultValues.target_episode_count ?? undefined,
      per_episode_hours: defaultPerEpHours,
      per_episode_minutes: defaultPerEpMinutes,
    },
  });

  const onSubmit = (values: EditTargetsFormValues) => {
    const hours = values.duration_hours ?? 0;
    const minutes = values.duration_minutes ?? 0;
    const totalSeconds = hoursMinutesToSeconds(hours, minutes);

    const perEpHours = values.per_episode_hours ?? 0;
    const perEpMinutes = values.per_episode_minutes ?? 0;
    const perEpSeconds = hoursMinutesToSeconds(perEpHours, perEpMinutes);

    mutate(
      {
        taskId,
        versionId,
        data: {
          display_name: (values.display_name ?? "").trim(),
          target_duration_seconds: enableTarget
            ? totalSeconds > 0
              ? totalSeconds
              : undefined
            : 0,
          target_episode_count: enableTarget ? values.target_episode_count : 0,
          target_duration_per_episode_seconds: enableTarget
            ? perEpSeconds > 0
              ? perEpSeconds
              : undefined
            : 0,
        },
      },
      {
        onSuccess: () => {
          setOpen(false);
        },
      }
    );
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild={!!children}>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("editTargetsDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("editTargetsDialog.description")}
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="display_name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("saveAsNewVersionForm.displayName")}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t(
                        "saveAsNewVersionForm.displayNamePlaceholder"
                      )}
                      maxLength={100}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

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

            <div className="flex justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => setOpen(false)}
                disabled={isPending}
              >
                {t("dialog.cancel")}
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending ? t("dialog.saving") : t("dialog.save")}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
