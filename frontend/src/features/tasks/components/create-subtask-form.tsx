"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

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
import { Textarea } from "@/components/ui/textarea";

import { useCreateSubTaskMutation } from "@/features/tasks/hooks/use-create-subtask-mutation";
import { useMeQuery } from "@/features/users";

const createSubTaskFormSchema = z.object({
  organization_id: z.string(),
  task_id: z.string(),
  task_version_id: z.string(),
  name: z.string().min(1, "Name is required"),
  description: z.string().optional(),
  target_duration_seconds: z.coerce.number().int().min(0).optional(),
});

type CreateSubTaskFormValues = z.infer<typeof createSubTaskFormSchema>;
type SubTaskCreateInput = z.infer<typeof schemas.SubTaskCreate>;

interface CreateSubTaskFormProps {
  taskId: string;
  taskVersionId: string;
  onSuccess?: () => void;
  onCancel?: () => void;
}

export function CreateSubTaskForm({
  taskId,
  taskVersionId,
  onSuccess,
  onCancel,
}: CreateSubTaskFormProps) {
  const { t } = useTranslation();
  const { mutate, isPending } = useCreateSubTaskMutation();
  const { data: meData } = useMeQuery();

  const form = useForm<CreateSubTaskFormValues>({
    resolver: zodResolver(createSubTaskFormSchema),
    defaultValues: {
      organization_id: meData?.organization_id ?? "",
      name: "",
      description: "",
      task_id: taskId,
      task_version_id: taskVersionId,
      target_duration_seconds: undefined,
    },
  });

  useEffect(() => {
    if (meData?.organization_id) {
      form.setValue("organization_id", meData.organization_id);
    }
  }, [meData?.organization_id, form]);

  const onSubmit = (values: CreateSubTaskFormValues) => {
    const data: SubTaskCreateInput = {
      organization_id: values.organization_id,
      task_id: values.task_id,
      task_version_id: values.task_version_id,
      name: values.name,
      description: values.description,
      target_duration_seconds:
        values.target_duration_seconds && values.target_duration_seconds > 0
          ? values.target_duration_seconds
          : undefined,
    };

    mutate(data, {
      onSuccess: () => {
        form.reset();
        onSuccess?.();
      },
    });
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("subtaskForm.name")}</FormLabel>
              <FormControl>
                <Input
                  placeholder={t("subtaskForm.namePlaceholder")}
                  {...field}
                />
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
              <FormLabel>{t("subtaskForm.descriptionOptional")}</FormLabel>
              <FormControl>
                <Textarea
                  placeholder={t("subtaskForm.descriptionPlaceholder")}
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
          name="target_duration_seconds"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Target Duration in seconds (Optional)</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  min={0}
                  placeholder="0"
                  className="w-32"
                  {...field}
                  value={field.value ?? ""}
                  onChange={(e) =>
                    field.onChange(
                      e.target.value === "" ? undefined : Number(e.target.value)
                    )
                  }
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
          <Button type="submit" disabled={isPending}>
            {isPending ? t("dialog.creating") : t("subtaskForm.createSubtask")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
