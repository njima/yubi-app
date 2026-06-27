"use client";

import { zodResolver } from "@hookform/resolvers/zod";
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

import { useUpdateSubTaskMutation } from "@/features/tasks/hooks/use-update-subtask-mutation";

const editSubTaskFormSchema = z.object({
  name: z.string().min(1, "Name is required").optional(),
  description: z.string().optional(),
  target_duration_seconds: z.coerce.number().int().min(0).optional(),
});

type EditSubTaskFormValues = z.infer<typeof editSubTaskFormSchema>;
type SubTaskUpdateInput = z.infer<typeof schemas.SubTaskUpdate>;

interface EditSubTaskFormProps {
  subtaskId: string;
  taskId: string;
  defaultValues: {
    name?: string;
    description?: string;
    target_duration_seconds?: number | null;
  };
  onSuccess?: () => void;
  onCancel?: () => void;
}

export function EditSubTaskForm({
  subtaskId,
  taskId,
  defaultValues,
  onSuccess,
  onCancel,
}: EditSubTaskFormProps) {
  const { t } = useTranslation();
  const { mutate, isPending } = useUpdateSubTaskMutation();

  const form = useForm<EditSubTaskFormValues>({
    resolver: zodResolver(editSubTaskFormSchema),
    defaultValues: {
      name: defaultValues.name,
      description: defaultValues.description,
      target_duration_seconds:
        defaultValues.target_duration_seconds ?? undefined,
    },
  });

  const onSubmit = (values: EditSubTaskFormValues) => {
    const data: SubTaskUpdateInput = {
      name: values.name,
      description: values.description,
      target_duration_seconds:
        values.target_duration_seconds && values.target_duration_seconds > 0
          ? values.target_duration_seconds
          : undefined,
    };

    mutate(
      { subtaskId, taskId, data },
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
              <FormLabel>{t("subtaskForm.description")}</FormLabel>
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
            {isPending ? t("dialog.saving") : t("subtaskForm.updateSubtask")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
