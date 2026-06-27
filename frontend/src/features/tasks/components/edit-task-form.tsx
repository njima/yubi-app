"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import {
  useTaskStatusLabel,
  useTaskPriorityLabel,
} from "@/lib/hooks/use-status-labels";
import {
  TASK_DIFFICULTY,
  type TaskDifficultyValue,
  TASK_PRIORITY,
  type TaskPriorityValue,
  TASK_STATUS,
  type TaskStatusValue,
} from "@/lib/status/constants";
import { getTaskDifficultyLabel } from "@/lib/status/utils";

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
import { SearchableSelect } from "@/components/ui/searchable-select";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";

import { useRobotsQuery } from "@/features/robots";
import { useUpdateTaskMutation } from "@/features/tasks/hooks/use-update-task-mutation";

import {
  taskUpdateSchema as baseTaskUpdateSchema,
  type TaskTag,
} from "../schemas";
import { TagSelector } from "./tag-selector";

const editTaskSchema = baseTaskUpdateSchema.extend({
  manual_url: z
    .string()
    .min(1, "Manual URL is required")
    .regex(/^https:\/\//, "Manual URL must start with https://"),
  deadline: z.string().min(1, "Deadline is required"),
});

type TaskUpdateInput = z.infer<typeof editTaskSchema>;

interface EditTaskFormProps {
  taskId: string;
  defaultValues: {
    name?: string;
    description?: string;
    manual_url?: string;
    priority?: TaskPriorityValue;
    difficulty?: TaskDifficultyValue;
    status?: TaskStatusValue;
    deadline?: string;
    robot_type?: string;
  };
  initialTags?: TaskTag[];
  onSuccess?: () => void;
  onCancel?: () => void;
}

export function EditTaskForm({
  taskId,
  defaultValues,
  initialTags = [],
  onSuccess,
  onCancel,
}: EditTaskFormProps) {
  const { t } = useTranslation();
  const getStatusLabel = useTaskStatusLabel();
  const getPriorityLabel = useTaskPriorityLabel();
  const { mutate, isPending } = useUpdateTaskMutation();
  const { data: robotsData } = useRobotsQuery({ limit: 1000 });
  const robots = robotsData?.robots;
  const uniqueRobotTypes = [
    ...new Set(robots?.map((r) => r.robot_type).filter(Boolean) ?? []),
  ] as string[];
  const [selectedTags, setSelectedTags] = useState<TaskTag[]>(initialTags);

  const form = useForm<TaskUpdateInput>({
    resolver: zodResolver(editTaskSchema),
    defaultValues,
  });

  const onSubmit = (data: TaskUpdateInput) => {
    mutate(
      {
        taskId,
        data: {
          ...data,
          deadline: new Date(data.deadline).toISOString(),
          robot_type: data.robot_type || undefined,
          tag_ids: selectedTags.map((t) => t.id),
        },
      },
      {
        onSuccess: () => {
          onSuccess?.();
        },
      }
    );
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("taskForm.name")}</FormLabel>
              <FormControl>
                <Input placeholder={t("taskForm.namePlaceholder")} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("taskForm.description")}</FormLabel>
              <FormControl>
                <Textarea
                  placeholder={t("taskForm.descriptionPlaceholder")}
                  className="resize-none"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="manual_url"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("taskForm.manualUrl")}</FormLabel>
              <FormControl>
                <Input
                  placeholder="https://..."
                  {...field}
                  value={field.value || ""}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="status"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("tasksPage.status")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value?.toString() ?? ""}
                  onValueChange={(value) => field.onChange(parseInt(value, 10))}
                  options={
                    field.value === TASK_STATUS.CANCELED
                      ? [
                          {
                            value: String(TASK_STATUS.CANCELED),
                            label: getStatusLabel(TASK_STATUS.CANCELED),
                          },
                          {
                            value: String(TASK_STATUS.PLANNING),
                            label: `${getStatusLabel(TASK_STATUS.PLANNING)} (${t("taskForm.uncancel")})`,
                          },
                        ]
                      : [
                          {
                            value: String(field.value ?? TASK_STATUS.PLANNING),
                            label: getStatusLabel(
                              (field.value ??
                                TASK_STATUS.PLANNING) as TaskStatusValue
                            ),
                          },
                          {
                            value: String(TASK_STATUS.CANCELED),
                            label: getStatusLabel(TASK_STATUS.CANCELED),
                          },
                        ]
                  }
                  placeholder={t("taskForm.selectStatus")}
                />
              </FormControl>
              <p className="text-xs text-gray-500 dark:text-gray-400">
                {t("taskForm.statusAutoManaged")}
              </p>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="priority"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("tasksPage.priority")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value?.toString() ?? ""}
                  onValueChange={(value) => field.onChange(parseInt(value, 10))}
                  options={[
                    {
                      value: String(TASK_PRIORITY.LOW),
                      label: getPriorityLabel(TASK_PRIORITY.LOW),
                    },
                    {
                      value: String(TASK_PRIORITY.NORMAL),
                      label: getPriorityLabel(TASK_PRIORITY.NORMAL),
                    },
                    {
                      value: String(TASK_PRIORITY.HIGH),
                      label: getPriorityLabel(TASK_PRIORITY.HIGH),
                    },
                    {
                      value: String(TASK_PRIORITY.URGENT),
                      label: getPriorityLabel(TASK_PRIORITY.URGENT),
                    },
                  ]}
                  placeholder={t("taskForm.selectPriority")}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="difficulty"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("tasksPage.difficulty")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value?.toString() ?? ""}
                  onValueChange={(value) => field.onChange(parseInt(value, 10))}
                  options={[
                    {
                      value: String(TASK_DIFFICULTY.S),
                      label: getTaskDifficultyLabel(TASK_DIFFICULTY.S),
                    },
                    {
                      value: String(TASK_DIFFICULTY.A),
                      label: getTaskDifficultyLabel(TASK_DIFFICULTY.A),
                    },
                    {
                      value: String(TASK_DIFFICULTY.B),
                      label: getTaskDifficultyLabel(TASK_DIFFICULTY.B),
                    },
                    {
                      value: String(TASK_DIFFICULTY.C),
                      label: getTaskDifficultyLabel(TASK_DIFFICULTY.C),
                    },
                  ]}
                  placeholder={t("taskForm.selectDifficulty")}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="deadline"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("taskForm.deadline")}</FormLabel>
              <FormControl>
                <Input
                  type="datetime-local"
                  {...field}
                  value={field.value || ""}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="robot_type"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("taskForm.robotTypeOptional")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value ?? ""}
                  onValueChange={(value) =>
                    field.onChange(value === "" ? undefined : value)
                  }
                  options={[
                    { value: "", label: t("taskForm.notSpecified") },
                    ...uniqueRobotTypes.map((m) => ({ value: m, label: m })),
                  ]}
                  placeholder={t("taskForm.selectRobotType")}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Separator />

        <TagSelector selectedTags={selectedTags} onChange={setSelectedTags} />

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
            {isPending ? t("dialog.saving") : t("taskForm.updateTask")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
