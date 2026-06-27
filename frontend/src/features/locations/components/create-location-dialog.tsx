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
import { SearchableSelect } from "@/components/ui/searchable-select";

import { useSiteSearchOptions } from "@/features/sites";
import { useMeQuery } from "@/features/users";

import { useCreateLocationMutation } from "../hooks/use-create-location-mutation";

const formSchema = schemas.LocationCreate.extend({
  name: z.string().min(1, "Name is required"),
  site_id: z.string().min(1, "Site is required"),
});

type LocationCreateForm = z.infer<typeof formSchema>;

export function CreateLocationDialog() {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const { mutate, isPending } = useCreateLocationMutation();
  const { data: me } = useMeQuery();
  const {
    options: siteOptions,
    isLoading: sitesLoading,
    onSearch: onSiteSearch,
    selectedLabel: siteLabel,
    onValueChange: onSiteValueChange,
  } = useSiteSearchOptions();

  const form = useForm<LocationCreateForm>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      organization_id: me?.organization_id ?? "",
      site_id: "",
    },
  });

  useEffect(() => {
    if (me?.organization_id) {
      form.setValue("organization_id", me.organization_id);
    }
  }, [me?.organization_id, form]);

  const onSubmit = (data: LocationCreateForm) => {
    mutate(data, {
      onSuccess: () => {
        form.reset();
        setOpen(false);
      },
    });
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>{t("createLocationDialog.trigger")}</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("createLocationDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("createLocationDialog.description")}
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
                      placeholder={t("createLocationDialog.namePlaceholder")}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="site_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Site *</FormLabel>
                  <FormControl>
                    <SearchableSelect
                      value={field.value}
                      onValueChange={(v) => {
                        field.onChange(v);
                        onSiteValueChange(v);
                      }}
                      options={siteOptions}
                      onSearch={onSiteSearch}
                      isLoading={sitesLoading}
                      selectedLabel={field.value ? siteLabel : undefined}
                      placeholder="Select site"
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
                {isPending
                  ? t("dialog.creating")
                  : t("createLocationDialog.createButton")}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
