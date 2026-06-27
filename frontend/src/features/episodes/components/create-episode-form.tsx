"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { format } from "date-fns";
import { useCallback, useEffect, useMemo } from "react";
import { useForm, useWatch } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";
import { APPROVAL_STATUS } from "@/lib/status/constants";
import type {
  TaskPriorityValue,
  TaskStatusValue,
} from "@/lib/status/constants";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  SearchableSelect,
  type SearchableSelectOption,
} from "@/components/ui/searchable-select";

import { useLocationsQuery } from "@/features/locations";
import { useRobotSearchOptions, useRobotQuery } from "@/features/robots";
import {
  useTaskSearchOptions,
  useTaskVersionsQuery,
  TaskPriorityBadge,
  TaskStatusBadge,
} from "@/features/tasks";
import { useMeQuery, useUserSearchOptions } from "@/features/users";

import { useCreateEpisodeMutation } from "../hooks/use-create-episode-mutation";
import { useCreateEpisodesBulkMutation } from "../hooks/use-create-episodes-bulk-mutation";
const createEpisodeFormSchema = schemas.EpisodeCreate.extend({
  count: z.coerce.number().int().min(1).max(100),
});
type CreateEpisodeFormData = z.infer<typeof createEpisodeFormSchema>;

interface CreateEpisodeFormProps {
  onSuccess?: () => void;
  onCancel?: () => void;
}

const emptyFormValues = {
  organization_id: "",
  robot_id: "",
  location_id: "",
  task_id: "",
  task_version_id: "",
  recorded_by: "",
  count: 1,
} satisfies CreateEpisodeFormData;

export function CreateEpisodeForm({
  onSuccess,
  onCancel,
}: CreateEpisodeFormProps) {
  const { t } = useTranslation();
  const { mutate: createSingleEpisode, isPending: isSingleCreating } =
    useCreateEpisodeMutation();
  const { mutate: createBulkEpisodes, isPending: isBulkCreating } =
    useCreateEpisodesBulkMutation();

  // Fetch master data for dropdowns (async search)
  const { data: meData } = useMeQuery();
  const {
    options: robotOptions,
    isLoading: robotsLoading,
    onSearch: onRobotSearch,
    selectedLabel: robotLabel,
    onValueChange: onRobotValueChange,
  } = useRobotSearchOptions();
  const { data: locationsData, isLoading: locationsLoading } =
    useLocationsQuery({ limit: 1000 });
  const locations = locationsData?.locations;
  const {
    isLoading: tasksLoading,
    onSearch: onTaskSearch,
    selectedLabel: taskLabel,
    onValueChange: onTaskValueChange,
    tasks,
  } = useTaskSearchOptions({
    sort_by: "recommended",
    has_approved_version: true,
  });

  const { enrichedTaskOptions, taskMetaMap } = useMemo(() => {
    const options =
      tasks?.map((t) => ({
        value: t.id,
        label: t.name,
      })) ?? [];
    const metaMap = new Map(
      tasks?.map((t) => [
        t.id,
        {
          priority: t.priority as TaskPriorityValue,
          deadline: t.deadline,
          status: t.status as TaskStatusValue,
        },
      ]) ?? []
    );
    return { enrichedTaskOptions: options, taskMetaMap: metaMap };
  }, [tasks]);

  const renderTaskOption = useCallback(
    (option: SearchableSelectOption) => {
      const meta = taskMetaMap.get(option.value);
      return (
        <span className="flex w-full items-center justify-between gap-2">
          <span className="truncate">{option.label}</span>
          <span className="flex shrink-0 items-center gap-3">
            {meta != null && <TaskStatusBadge status={meta.status} />}
            {meta != null && <TaskPriorityBadge priority={meta.priority} />}
            {meta?.deadline && (
              <span className="text-xs text-gray-500 dark:text-gray-400">
                {format(new Date(meta.deadline), "MMM d")}
              </span>
            )}
          </span>
        </span>
      );
    },
    [taskMetaMap]
  );

  const {
    options: userOptions,
    isLoading: usersLoading,
    onSearch: onUserSearch,
    selectedLabel: userLabel,
    onValueChange: onUserValueChange,
  } = useUserSearchOptions();

  const form = useForm<CreateEpisodeFormData>({
    resolver: zodResolver(createEpisodeFormSchema),
    defaultValues: {
      ...emptyFormValues,
      organization_id: meData?.organization_id ?? "",
    },
  });

  const selectedRobotId = useWatch({ control: form.control, name: "robot_id" });
  const { data: selectedRobot } = useRobotQuery(selectedRobotId, {
    enabled: !!selectedRobotId,
  });
  const selectedTaskId = useWatch({ control: form.control, name: "task_id" });
  const { data: taskVersions } = useTaskVersionsQuery(selectedTaskId, {
    enabled: !!selectedTaskId,
  });
  const approvedVersions = taskVersions?.filter(
    (v) => v.approval_status === APPROVAL_STATUS.APPROVED
  );
  const currentVersion = approvedVersions?.find((v) => v.is_current);

  const approvedVersionParams = useMemo(
    () => currentVersion?.parameters ?? [],
    [currentVersion]
  );

  // Initialize parameter values to empty strings when parameters load,
  // so z.record(z.string()) validation doesn't fail on undefined values.
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

  useEffect(() => {
    if (meData?.organization_id) {
      form.setValue("organization_id", meData.organization_id);
    }
  }, [meData?.organization_id, form]);

  useEffect(() => {
    form.setValue("task_version_id", "");
  }, [selectedTaskId, form]);

  // Auto-set location from selected robot
  useEffect(() => {
    form.setValue("location_id", selectedRobot?.location_id ?? "");
  }, [selectedRobot, form]);

  // Auto-select the top recommended task
  useEffect(() => {
    const firstOption = enrichedTaskOptions[0];
    if (!tasksLoading && firstOption && !selectedTaskId) {
      form.setValue("task_id", firstOption.value);
      onTaskValueChange(firstOption.value);
    }
  }, [
    enrichedTaskOptions,
    tasksLoading,
    selectedTaskId,
    form,
    onTaskValueChange,
  ]);

  const selectedLocationName = useMemo(() => {
    if (!selectedRobot?.location_id) return "";
    return (
      locations?.find((l) => l.id === selectedRobot.location_id)?.name ?? ""
    );
  }, [selectedRobot, locations]);

  const onSubmit = (data: CreateEpisodeFormData) => {
    const { count, parameter_values: rawPV, ...episodeData } = data;

    // Build parameter_values, omitting keys left as "Random"
    const parameterValues: Record<string, string> = {};
    if (rawPV) {
      for (const [k, v] of Object.entries(rawPV)) {
        if (v) parameterValues[k] = v;
      }
    }

    const basePayload = {
      ...episodeData,
      recorded_by: episodeData.recorded_by || undefined,
      task_version_id: episodeData.task_version_id || undefined,
      parameter_values:
        Object.keys(parameterValues).length > 0 ? parameterValues : undefined,
    };

    const resetValues = {
      ...emptyFormValues,
      organization_id: meData?.organization_id ?? "",
    };

    if (count > 1) {
      createBulkEpisodes(
        {
          ...basePayload,
          count,
        },
        {
          onSuccess: () => {
            form.reset(resetValues);
            onSuccess?.();
          },
        }
      );
      return;
    }

    createSingleEpisode(basePayload, {
      onSuccess: () => {
        form.reset(resetValues);
        onSuccess?.();
      },
      onError: (error) => {
        const taskName = taskLabel;
        toast.error(
          t("createEpisodeForm.failedToCreate", {
            suffix: taskName
              ? t("createEpisodeForm.forTask", { taskName })
              : "",
          }),
          { description: error.message }
        );
      },
    });
  };

  const isLoadingMasterData =
    robotsLoading || locationsLoading || tasksLoading || usersLoading;
  const isPending = isSingleCreating || isBulkCreating;

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        {/* Robot Selection */}
        <FormField
          control={form.control}
          name="robot_id"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("episodesPage.robot")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value}
                  onValueChange={(v) => {
                    field.onChange(v);
                    onRobotValueChange(v);
                  }}
                  options={robotOptions}
                  onSearch={onRobotSearch}
                  isLoading={robotsLoading}
                  selectedLabel={robotLabel}
                  placeholder={t("createEpisodeForm.selectRobot")}
                  disabled={isLoadingMasterData}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Location (auto-set from robot) */}
        <FormField
          control={form.control}
          name="location_id"
          render={() => (
            <FormItem>
              <FormLabel>{t("usersPage.location")}</FormLabel>
              <FormControl>
                <Input value={selectedLocationName} disabled readOnly />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Task Selection */}
        <FormField
          control={form.control}
          name="task_id"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("episodesPage.task")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value}
                  onValueChange={(v) => {
                    field.onChange(v);
                    onTaskValueChange(v);
                  }}
                  options={enrichedTaskOptions}
                  onSearch={onTaskSearch}
                  isLoading={tasksLoading}
                  selectedLabel={taskLabel}
                  placeholder={t("createEpisodeForm.selectTask")}
                  disabled={isLoadingMasterData}
                  renderOption={renderTaskOption}
                  renderSelected={renderTaskOption}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Task Version Selection */}
        <FormField
          control={form.control}
          name="task_version_id"
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                {t("createEpisodeForm.taskVersionOptional")}
              </FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value ?? ""}
                  onValueChange={field.onChange}
                  options={
                    approvedVersions?.map((v) => ({
                      value: v.id,
                      label: `${v.version}${v.is_current ? t("createEpisodeForm.latestSuffix") : ""}`,
                    })) ?? []
                  }
                  placeholder={
                    selectedTaskId
                      ? (currentVersion?.version ??
                        t("createEpisodeForm.latestApprovedVersion"))
                      : t("episodesPage.selectTaskFirst")
                  }
                  disabled={isLoadingMasterData || !selectedTaskId}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Parameter Selections */}
        {approvedVersionParams.length > 0 && (
          <div className="space-y-3">
            <p className="text-sm font-medium">
              {t("parametersEditor.parameters")}
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
                            label: t("createEpisodeForm.random"),
                          },
                          ...param.values.map((v) => ({
                            value: v,
                            label: v,
                          })),
                        ]}
                        placeholder={t("createEpisodeForm.selectValue")}
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
              <FormLabel>{t("createEpisodeForm.recordedByOptional")}</FormLabel>
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
                  placeholder={t("createEpisodeForm.selectUser")}
                  disabled={isLoadingMasterData}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="count"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("createEpisodeForm.count")}</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  min={1}
                  max={100}
                  step={1}
                  value={field.value}
                  onChange={field.onChange}
                  disabled={isPending}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

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
          <Button type="submit" disabled={isPending || isLoadingMasterData}>
            {isPending
              ? t("dialog.creating")
              : t("createEpisodeForm.createEpisode")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
