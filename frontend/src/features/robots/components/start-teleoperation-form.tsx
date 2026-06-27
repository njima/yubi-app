"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useEffect, useMemo } from "react";
import { useForm, useWatch } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { APPROVAL_STATUS } from "@/shared/lib/status-constants";

import { Button } from "@/components/ui/button";
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

const startTeleoperationFormSchema = z.object({
  task_id: z.string().min(1, "Task is required"),
  recorded_by: z.string().optional(),
  parameter_values: z.record(z.string()).optional(),
});

type StartTeleoperationFormValues = z.infer<
  typeof startTeleoperationFormSchema
>;

interface StartTeleoperationFormProps {
  robot: Robot;
  onSuccess?: () => void;
  onCancel?: () => void;
}

export function StartTeleoperationForm({
  robot,
  onSuccess,
  onCancel,
}: StartTeleoperationFormProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const { mutate: createEpisode, isPending } = useCreateEpisodeMutation();
  const {
    options: taskOptions,
    isLoading: tasksLoading,
    onSearch: onTaskSearch,
    selectedLabel: taskLabel,
    onValueChange: onTaskValueChange,
  } = useTaskSearchOptions({ has_approved_version: true });
  const {
    options: userOptions,
    isLoading: usersLoading,
    onSearch: onUserSearch,
    selectedLabel: userLabel,
    onValueChange: onUserValueChange,
  } = useUserSearchOptions();
  const { data: meData } = useMeQuery();

  const form = useForm<StartTeleoperationFormValues>({
    resolver: zodResolver(startTeleoperationFormSchema),
    defaultValues: {
      task_id: "",
      recorded_by: "",
    },
  });

  useEffect(() => {
    if (meData?.user_id) {
      form.setValue("recorded_by", meData.user_id);
    }
  }, [meData?.user_id, form]);

  // Fetch task version parameters for the selected task
  const selectedTaskId = useWatch({ control: form.control, name: "task_id" });
  const { data: taskVersions } = useTaskVersionsQuery(selectedTaskId || "", {
    enabled: !!selectedTaskId,
  });
  const currentVersion = taskVersions?.find(
    (v) => v.approval_status === APPROVAL_STATUS.APPROVED && v.is_current
  );
  const approvedVersionParams = useMemo(
    () => currentVersion?.parameters ?? [],
    [currentVersion]
  );

  // Initialize parameter values when parameters load
  useEffect(() => {
    if (approvedVersionParams.length > 0) {
      const defaults: Record<string, string> = {};
      for (const param of approvedVersionParams) {
        defaults[param.key] = "";
      }
      form.setValue("parameter_values", defaults);
    } else {
      form.setValue("parameter_values", undefined);
    }
  }, [approvedVersionParams, form]);

  const onSubmit = (data: StartTeleoperationFormValues) => {
    // Build parameter_values, omitting keys left as "Random"
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
          onSuccess?.();
          router.push(`/robots/${robot.id}/teleoperation`);
        },
      }
    );
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        {/* Task Selection */}
        <FormField
          control={form.control}
          name="task_id"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("teleoperationForm.taskRequired")}</FormLabel>
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
                  placeholder={t("teleoperationForm.selectTask")}
                  disabled={tasksLoading}
                />
              </FormControl>
              <FormDescription>
                {t("teleoperationForm.taskDescription")}
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Parameter Selections */}
        {approvedVersionParams.length > 0 && (
          <div className="space-y-3">
            <p className="text-sm font-medium">
              {t("teleoperationForm.parameters")}
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
                            label: t("teleoperationForm.random"),
                          },
                          ...param.values.map((v) => ({
                            value: v,
                            label: v,
                          })),
                        ]}
                        placeholder="Select value"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            ))}
          </div>
        )}

        {/* Recorded By Selection */}
        <FormField
          control={form.control}
          name="recorded_by"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("teleoperationForm.recordedBy")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value ?? ""}
                  onValueChange={(v) => {
                    field.onChange(v);
                    onUserValueChange(v);
                  }}
                  options={userOptions}
                  onSearch={onUserSearch}
                  isLoading={usersLoading}
                  selectedLabel={userLabel}
                  placeholder={t("teleoperationForm.selectUser")}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Action Buttons */}
        <div className="flex justify-end gap-2 pt-4">
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
            {isPending
              ? t("teleoperationForm.starting")
              : t("teleop.startTeleoperation")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
