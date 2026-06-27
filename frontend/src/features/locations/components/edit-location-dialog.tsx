"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { Button } from "@/components/ui/button";
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

import { useUpdateLocationMutation } from "../hooks/use-update-location-mutation";

type Location = z.infer<typeof schemas.Location>;
type LocationUpdate = z.infer<typeof schemas.LocationUpdate>;

interface EditLocationDialogProps {
  location: Location;
  children?: React.ReactNode;
}

export function EditLocationDialog({
  location,
  children,
}: EditLocationDialogProps) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const { mutate, isPending } = useUpdateLocationMutation();

  const form = useForm<LocationUpdate>({
    resolver: zodResolver(schemas.LocationUpdate),
    defaultValues: {
      name: location.name,
    },
  });

  useEffect(() => {
    if (open) {
      form.reset({ name: location.name });
    }
  }, [open, location.name, form]);

  const onSubmit = (data: LocationUpdate) => {
    mutate(
      { locationId: location.id, data },
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
          <DialogTitle>{t("editLocationDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("editLocationDialog.description")}
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("formLabel.nameRequired")}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t("editLocationDialog.namePlaceholder")}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
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
