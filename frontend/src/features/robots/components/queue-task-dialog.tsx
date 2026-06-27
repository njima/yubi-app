"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useMemo } from "react";
import { useForm, useWatch } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { APPROVAL_STATUS } from "@/shared/lib/status-constants";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { SearchableSelect } from "@/components/ui/searchable-select";

import { useCreateEpisodeMutation } from "@/features/episodes";
import { useTaskSearchOptions, useTaskVersionsQuery } from "@/features/tasks";
import { useMeQuery, useUserSearchOptions } from "@/features/users";

import type { Robot } from "../schemas/robot";

const queueTaskFormSchema = z.object({
  task_id: z.string().min(1, "Task is required"),
  recorded_by: z.string().optional(),
  parameter_values: z.record(z.string()).optional(),
});

type QueueTaskFormValues = z.infer<typeof queueTaskFormSchema>;

interface QueueTaskDialogProps {
  robot: Robot;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function QueueTaskDialog({
  robot,
  open,
  onOpenChange,
}: QueueTaskDialogProps) {
  const { t } = useTranslation();
  const { mutate: createEpisode, isPending } = useCreateEpisodeMutation();
  const {
    options: taskOptions,
    isLoading: tasksLoading,
    onSearch: onTaskSearch,
    selectedLabel: taskLabel,
    onValueChange: onTaskValueChange,
  } = useTaskSearchOptions({ has_approved_version: true });
  const { data: meData } = useMeQuery();
  const {
    options: userOptions,
    isLoading: usersLoading,
    onSearch: onUserSearch,
    selectedLabel: userLabel,
    onValueChange: onUserValueChange,
  } = useUserSearchOptions();

  const form = useForm<QueueTaskFormValues>({
    resolver: zodResolver(queueTaskFormSchema),
    defaultValues: {
      task_id: "",
      recorded_by: "",
    },
  });

  // Pre-fill recorded_by with current user
  useEffect(() => {
    if (meData?.user_id && !form.getValues("recorded_by")) {
      form.setValue("recorded_by", meData.user_id);
    }
  }, [meData, form]);

  // Watch selected task to fetch its versions
  const selectedTaskId = useWatch({ control: form.control, name: "task_id" });
  const { data: taskVersions } = useTaskVersionsQuery(selectedTaskId || "", {
    enabled: !!selectedTaskId,
  });

  const currentVersion = useMemo(() => {
    if (!taskVersions) return undefined;
    return taskVersions.find(
      (v) => v.approval_status === APPROVAL_STATUS.APPROVED && v.is_current
    );
  }, [taskVersions]);

  const approvedVersionParams = useMemo(
    () => currentVersion?.parameters ?? [],
    [currentVersion]
  );

  // Initialize parameter values when parameters change
  useEffect(() => {
    if (approvedVersionParams.length > 0) {
      const defaults: Record<string, string> = {};
      for (const param of approvedVersionParams) {
        defaults[param.key] = "";
      }
      form.setValue("parameter_values", defaults);
    } else {
      form.setValue("parameter_values", undefined as never);
    }
  }, [approvedVersionParams, form]);

  const onSubmit = (data: QueueTaskFormValues) => {
    // Build parameter_values, omitting "Random" (empty) values
    const parameterValues: Record<string, string> = {};
    if (data.parameter_values) {
      for (const [k, v] of Object.entries(data.parameter_values)) {
        if (v) parameterValues[k] = v;
      }
    }

    createEpisode(
      {
        organization_id: robot.organization_id ?? "",
        location_id: robot.location_id ?? "",
        robot_id: robot.id,
        task_id: data.task_id,
        recorded_by: data.recorded_by || undefined,
        parameter_values:
          Object.keys(parameterValues).length > 0 ? parameterValues : undefined,
      },
      {
        onSuccess: () => {
          form.reset();
          onOpenChange(false);
        },
      }
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{t("queueTaskDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("queueTaskDialog.description", { robotName: robot.name })}
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            {/* Task Selection */}
            <FormField
              control={form.control}
              name="task_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("queueTaskDialog.taskLabel")}</FormLabel>
                  <FormControl>
                    <SearchableSelect
                      value={field.value}
                      onValueChange={(v) => {
                        field.onChange(v);
                        onTaskValueChange(v);
                      }}
                      options={taskOptions}
                      onSearch={onTaskSearch}
                      isLoading={tasksLoading}
                      selectedLabel={taskLabel}
                      placeholder={t("queueTaskDialog.selectTask")}
                      disabled={tasksLoading}
                    />
                  </FormControl>
                  <FormDescription>
                    {t("queueTaskDialog.taskDescription")}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Parameter Selections */}
            {approvedVersionParams.length > 0 && (
              <div className="space-y-3">
                <p className="text-sm font-medium">
                  {t("queueTaskDialog.parameters")}
                </p>
                {approvedVersionParams.map((param) => (
                  <FormField
                    key={param.key}
                    control={form.control}
                    name={`parameter_values.${param.key}` as never}
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{param.key}</FormLabel>
                        <FormControl>
                          <SearchableSelect
                            value={(field.value as string) || "__random__"}
                            onValueChange={(v) =>
                              field.onChange(v === "__random__" ? "" : v)
                            }
                            options={[
                              {
                                value: "__random__",
                                label: t("queueTaskDialog.random"),
                              },
                              ...param.values.map((v) => ({
                                value: v,
                                label: v,
                              })),
                            ]}
                            placeholder={t("queueTaskDialog.selectValue")}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                ))}
              </div>
            )}

            {/* Recorded By */}
            <FormField
              control={form.control}
              name="recorded_by"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("queueTaskDialog.recordedBy")}</FormLabel>
                  <FormControl>
                    <SearchableSelect
                      value={field.value || ""}
                      onValueChange={(v) => {
                        field.onChange(v);
                        onUserValueChange(v);
                      }}
                      options={userOptions}
                      onSearch={onUserSearch}
                      isLoading={usersLoading}
                      selectedLabel={userLabel}
                      placeholder={t("queueTaskDialog.selectUser")}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="flex justify-end gap-2 pt-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isPending}
              >
                {t("dialog.cancel")}
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending ? t("dialog.queuing") : t("dialog.queueTask")}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
