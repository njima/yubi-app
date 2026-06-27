"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
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

import { useUpdateDisplayNameMutation } from "../hooks/use-update-display-name-mutation";

const formSchema = z.object({
  display_name: z
    .string()
    .trim()
    .min(1, "Display name is required")
    .max(60, "Display name must be 60 characters or fewer"),
});

type FormValues = z.infer<typeof formSchema>;

interface EditDisplayNameDialogProps {
  currentDisplayName: string;
  children: React.ReactNode;
}

export function EditDisplayNameDialog({
  currentDisplayName,
  children,
}: EditDisplayNameDialogProps) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const { mutateAsync, isPending } = useUpdateDisplayNameMutation();

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: { display_name: currentDisplayName },
  });

  const handleSubmit = async (values: FormValues) => {
    if (values.display_name === currentDisplayName) {
      setOpen(false);
      return;
    }
    try {
      await mutateAsync({ displayName: values.display_name });
      setOpen(false);
    } catch {
      // toast handled in mutation onError
    }
  };

  return (
    <Dialog
      open={open}
      onOpenChange={(next) => {
        if (next) {
          form.reset({ display_name: currentDisplayName });
        }
        setOpen(next);
      }}
    >
      <DialogTrigger asChild>{children}</DialogTrigger>

      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("editDisplayNameDialog.title")}</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
            <FormField
              control={form.control}
              name="display_name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {t("editDisplayNameDialog.displayName")}
                  </FormLabel>
                  <FormControl>
                    <Input maxLength={60} autoFocus {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
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
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
