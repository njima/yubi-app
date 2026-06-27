"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useRef, useState } from "react";
import { useMemo } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { z } from "zod";

import {
  useTaskStatusLabel,
  useTaskPriorityLabel,
} from "@/lib/hooks/use-status-labels";
import {
  TASK_DIFFICULTY,
  TASK_PRIORITY,
  TASK_STATUS,
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
import { useCreateSubTaskMutation } from "@/features/tasks/hooks/use-create-subtask-mutation";
import { useCreateTaskMutation } from "@/features/tasks/hooks/use-create-task-mutation";
import { hoursMinutesToSeconds } from "@/features/tasks/lib/duration";
import { useMeQuery } from "@/features/users";

import {
  taskCreateSchema as baseTaskCreateSchema,
  type TaskTag,
} from "../schemas";
import {
  PendingSubTaskInput,
  type PendingSubTask,
} from "./pending-subtask-input";
import { TagSelector } from "./tag-selector";

function buildCreateTaskSchema(t: (key: string) => string) {
  return baseTaskCreateSchema.extend({
    manual_url: z
      .string()
      .min(1, t("validation.manualUrlRequired"))
      .regex(/^https:\/\//, t("validation.manualUrlHttps")),
    deadline: z.string().min(1, t("validation.deadlineRequired")),
  });
}

type TaskCreateInput = z.infer<ReturnType<typeof buildCreateTaskSchema>>;

interface CreateTaskFormProps {
  onSuccess?: () => void;
  onCancel?: () => void;
}

export function CreateTaskForm({ onSuccess, onCancel }: CreateTaskFormProps) {
  const { t } = useTranslation();
  const getStatusLabel = useTaskStatusLabel();
  const getPriorityLabel = useTaskPriorityLabel();
  const createTaskSchema = useMemo(() => buildCreateTaskSchema(t), [t]);
  const { mutateAsync: createTask, isPending: isCreatingTask } =
    useCreateTaskMutation();
  const { mutateAsync: createSubtask, isPending: isCreatingSubtask } =
    useCreateSubTaskMutation();
  const { data: meData } = useMeQuery();
  const { data: robotsData } = useRobotsQuery({ limit: 1000 });
  const robots = robotsData?.robots;
  const uniqueRobotTypes = [
    ...new Set(robots?.map((r) => r.robot_type).filter(Boolean) ?? []),
  ] as string[];

  const [pendingSubtasks, setPendingSubtasks] = useState<PendingSubTask[]>([]);
  const subtaskCounter = useRef(0);
  const [selectedTags, setSelectedTags] = useState<TaskTag[]>([]);

  const isPending = isCreatingTask || isCreatingSubtask;

  const form = useForm<TaskCreateInput>({
    resolver: zodResolver(createTaskSchema),
    defaultValues: {
      organization_id: meData?.organization_id ?? "",
      name: "",
      description: "",
      manual_url: "",
      priority: TASK_PRIORITY.NORMAL,
      difficulty: TASK_DIFFICULTY.B,
      status: TASK_STATUS.PLANNING,
      deadline: "",
    },
  });

  useEffect(() => {
    if (meData?.organization_id) {
      form.setValue("organization_id", meData.organization_id);
    }
  }, [meData?.organization_id, form]);

  const handleAddSubtask = () => {
    setPendingSubtasks((prev) => [
      ...prev,
      {
        tempId: String(subtaskCounter.current++),
        name: "",
        description: "",
        durationHours: undefined,
        durationMinutes: undefined,
      },
    ]);
  };

  const handleRemoveSubtask = (tempId: string) => {
    setPendingSubtasks((prev) => prev.filter((st) => st.tempId !== tempId));
  };

  const handleUpdateSubtask = (
    tempId: string,
    field: "name" | "description",
    value: string
  ) => {
    setPendingSubtasks((prev) =>
      prev.map((st) => (st.tempId === tempId ? { ...st, [field]: value } : st))
    );
  };

  const handleUpdateSubtaskDuration = (
    tempId: string,
    hours: number | undefined,
    minutes: number | undefined
  ) => {
    setPendingSubtasks((prev) =>
      prev.map((st) =>
        st.tempId === tempId
          ? { ...st, durationHours: hours, durationMinutes: minutes }
          : st
      )
    );
  };

  const onSubmit = async (data: TaskCreateInput) => {
    try {
      // 1. Create task
      const newTask = await createTask({
        ...data,
        deadline: new Date(data.deadline).toISOString(),
        robot_type: data.robot_type || undefined,
        tag_ids: selectedTags.map((t) => t.id),
      });

      // 2. Create subtasks (filter out empty names)
      const validSubtasks = pendingSubtasks.filter(
        (st) => st.name.trim() !== ""
      );

      if (validSubtasks.length > 0) {
        // Fetch the initial version created with the task
        const versionsRes = await fetch(
          `/web/api/tasks/${newTask.id}/versions`
        );
        if (!versionsRes.ok) {
          toast.error(t("createTaskForm.failedToRetrieveVersion"));
          return;
        }
        const versions = await versionsRes.json();
        const initialVersionId = versions[0]?.id;
        if (!initialVersionId) {
          toast.error(t("createTaskForm.failedToRetrieveVersion"));
          return;
        }

        let failureCount = 0;
        for (const st of validSubtasks) {
          try {
            const hours = st.durationHours ?? 0;
            const minutes = st.durationMinutes ?? 0;
            const totalSeconds = hoursMinutesToSeconds(hours, minutes);
            await createSubtask({
              organization_id: data.organization_id,
              task_id: newTask.id,
              task_version_id: initialVersionId,
              name: st.name,
              description: st.description || undefined,
              target_duration_seconds:
                totalSeconds > 0 ? totalSeconds : undefined,
            });
          } catch {
            failureCount++;
          }
        }
        if (failureCount > 0) {
          toast.warning(
            t("createTaskForm.subtasksFailedWarning", { count: failureCount })
          );
        }
      }

      form.reset();
      setPendingSubtasks([]);
      setSelectedTags([]);
      onSuccess?.();
    } catch {
      // Task creation failed - error is handled by mutation's onError
    }
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
              <FormLabel>{t("taskForm.descriptionOptional")}</FormLabel>
              <FormControl>
                <Textarea
                  placeholder={t("taskForm.descriptionPlaceholder")}
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
          name="manual_url"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("taskForm.manualUrl")}</FormLabel>
              <FormControl>
                <Input placeholder="https://..." {...field} />
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
                  options={[
                    {
                      value: String(TASK_STATUS.PLANNING),
                      label: getStatusLabel(TASK_STATUS.PLANNING),
                    },
                    {
                      value: String(TASK_STATUS.DOING),
                      label: getStatusLabel(TASK_STATUS.DOING),
                    },
                    {
                      value: String(TASK_STATUS.COMPLETED),
                      label: getStatusLabel(TASK_STATUS.COMPLETED),
                    },
                    {
                      value: String(TASK_STATUS.CANCELED),
                      label: getStatusLabel(TASK_STATUS.CANCELED),
                    },
                  ]}
                  placeholder={t("taskForm.selectStatus")}
                />
              </FormControl>
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
                <Input type="datetime-local" {...field} />
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

        <TagSelector selectedTags={selectedTags} onChange={setSelectedTags} />

        <Separator />

        <PendingSubTaskInput
          subtasks={pendingSubtasks}
          onAdd={handleAddSubtask}
          onRemove={handleRemoveSubtask}
          onUpdate={handleUpdateSubtask}
          onUpdateDuration={handleUpdateSubtaskDuration}
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
          <Button type="submit" disabled={isPending}>
            {isPending ? t("dialog.creating") : t("taskForm.createTask")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
