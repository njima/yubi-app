"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";

import { schemas } from "@/lib/api/generated/api";

import { useRobotTypeLabel } from "@/shared/hooks/use-status-labels";
import { ROBOT_TYPE } from "@/shared/lib/status-constants";
import { Button } from "@/shared/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/ui/form";
import { Input } from "@/shared/ui/input";
import { SearchableSelect } from "@/shared/ui/searchable-select";
import { Textarea } from "@/shared/ui/textarea";

import { useLocationSearchOptions } from "@/features/locations";
import { useCreateRobotMutation } from "@/features/robots/hooks/use-create-robot-mutation";
import { useSiteSearchOptions } from "@/features/sites";
import { useMeQuery } from "@/features/users";

import type { z } from "zod";

type RobotCreate = z.infer<typeof schemas.RobotCreate>;

interface CreateRobotFormProps {
  onSuccess?: () => void;
  onCancel?: () => void;
}

export function CreateRobotForm({ onSuccess, onCancel }: CreateRobotFormProps) {
  const { t } = useTranslation();
  const getRobotTypeLabel = useRobotTypeLabel();
  const { mutate, isPending } = useCreateRobotMutation();

  // Fetch sites and locations for dropdowns
  const { data: meData } = useMeQuery();
  const [selectedSiteId, setSelectedSiteId] = useState<string>("");
  const {
    options: siteOptions,
    isLoading: sitesLoading,
    onSearch: onSiteSearch,
    selectedLabel: siteLabel,
    onValueChange: onSiteValueChange,
  } = useSiteSearchOptions();
  const {
    options: locationOptions,
    isLoading: locationsLoading,
    onSearch: onLocationSearch,
    selectedLabel: locationLabel,
    onValueChange: onLocationValueChange,
  } = useLocationSearchOptions({ site_id: selectedSiteId || undefined });

  const [robotConfigError, setRobotConfigError] = useState<string>("");

  const form = useForm<RobotCreate>({
    resolver: zodResolver(schemas.RobotCreate),
    defaultValues: {
      name: "",
      organization_id: meData?.organization_id ?? undefined,
      location_id: undefined,
      robot_type: ROBOT_TYPE.YUBI,
      robot_config: undefined,
    },
  });

  useEffect(() => {
    if (meData?.organization_id) {
      form.setValue("organization_id", meData.organization_id);
    }
  }, [meData?.organization_id, form]);

  const onSubmit = (data: RobotCreate) => {
    const robotConfigValue = form.getValues("robot_config") as unknown;
    if (typeof robotConfigValue === "string") {
      if (robotConfigValue === "") {
        setRobotConfigError("");
      } else {
        try {
          JSON.parse(robotConfigValue);
          setRobotConfigError("");
        } catch {
          setRobotConfigError(t("robotForm.invalidJson"));
          return;
        }
      }
    }

    mutate(data, {
      onSuccess: () => {
        form.reset();
        setRobotConfigError("");
        onSuccess?.();
      },
    });
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        {/* Name */}
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("robotForm.nameRequired")}</FormLabel>
              <FormControl>
                <Input
                  placeholder={t("robotForm.namePlaceholderCreate")}
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Site Filter */}
        <FormItem>
          <FormLabel>{t("robotForm.siteFilter")}</FormLabel>
          <FormControl>
            <SearchableSelect
              value={selectedSiteId}
              onValueChange={(v) => {
                setSelectedSiteId(v);
                onSiteValueChange(v);
                form.setValue("location_id", "");
              }}
              options={[
                { value: "", label: t("robotForm.allSites") },
                ...siteOptions,
              ]}
              onSearch={onSiteSearch}
              isLoading={sitesLoading}
              selectedLabel={selectedSiteId ? siteLabel : undefined}
              placeholder={t("robotForm.allSites")}
            />
          </FormControl>
        </FormItem>

        {/* Location */}
        <FormField
          control={form.control}
          name="location_id"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("robotForm.location")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value}
                  onValueChange={(v) => {
                    field.onChange(v);
                    onLocationValueChange(v);
                  }}
                  options={locationOptions}
                  onSearch={onLocationSearch}
                  isLoading={locationsLoading}
                  selectedLabel={locationLabel}
                  placeholder={t("robotForm.selectLocation")}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Robot Type */}
        <FormField
          control={form.control}
          name="robot_type"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("robotForm.robotType")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value ?? ""}
                  onValueChange={field.onChange}
                  options={[
                    {
                      value: ROBOT_TYPE.YUBI,
                      label: getRobotTypeLabel(ROBOT_TYPE.YUBI),
                    },
                    {
                      value: ROBOT_TYPE.YUBI_PORTABLE,
                      label: getRobotTypeLabel(ROBOT_TYPE.YUBI_PORTABLE),
                    },
                  ]}
                  placeholder={t("robotForm.selectRobotType")}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Leader Status */}
        <FormField
          control={form.control}
          name="leader_status"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("robotForm.leaderStatus")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value?.toString() ?? ""}
                  onValueChange={(value) => field.onChange(parseInt(value, 10))}
                  options={[
                    { value: "0", label: t("status.ready") },
                    { value: "1", label: t("status.faulted") },
                    { value: "2", label: t("status.maintenance") },
                  ]}
                  placeholder={t("robotForm.notSetNoLeader")}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Robot Config */}
        <FormField
          control={form.control}
          name="robot_config"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("robotForm.robotConfig")}</FormLabel>
              <FormControl>
                <Textarea
                  placeholder='{"host": "192.168.1.101", "port": 9090, "cameras": [{"namespace": "camera_0", "name": "Front"}]}'
                  className="font-mono text-sm"
                  rows={4}
                  {...field}
                  value={
                    typeof field.value === "string"
                      ? field.value
                      : field.value
                        ? JSON.stringify(field.value, null, 2)
                        : ""
                  }
                  onChange={(e) => {
                    const value = e.target.value;
                    if (value.trim() === "") {
                      field.onChange(undefined);
                      setRobotConfigError("");
                    } else {
                      try {
                        const parsed = JSON.parse(value);
                        field.onChange(parsed);
                        setRobotConfigError("");
                      } catch {
                        field.onChange(value);
                        setRobotConfigError("Invalid JSON format");
                      }
                    }
                  }}
                />
              </FormControl>
              {robotConfigError && (
                <p className="text-sm font-medium text-destructive">
                  {robotConfigError}
                </p>
              )}
              <FormDescription>
                {t("robotForm.robotConfigDescription")}
              </FormDescription>
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
          <Button type="submit" disabled={isPending || !!robotConfigError}>
            {isPending ? t("dialog.creating") : t("robotForm.createRobot")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
