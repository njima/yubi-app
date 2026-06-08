"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { useRobotTypeLabel } from "@/shared/hooks/use-status-labels";
import { ROBOT_STATUS, ROBOT_TYPE } from "@/shared/lib/status-constants";
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

import {
  useLocationSearchOptions,
  useLocationQuery,
} from "@/features/locations";
import { useUpdateRobotMutation } from "@/features/robots/hooks/use-update-robot-mutation";
import {
  mergeRobotConfigWithHost,
  splitHostFromRobotConfig,
} from "@/features/robots/lib/robot-config-utils";
import { useSiteSearchOptions } from "@/features/sites";

type Robot = z.infer<typeof schemas.Robot>;

function buildUpdateRobotSchema(t: (key: string) => string) {
  return schemas.RobotUpdate.extend({
    robot_type: z.enum([ROBOT_TYPE.YUBI_STATIONARY, ROBOT_TYPE.YUBI_PORTABLE], {
      errorMap: () => ({ message: t("robotForm.robotTypeRequired") }),
    }),
    host: z.string().min(1, t("validation.hostRequired")),
  });
}

type RobotUpdateInput = z.infer<ReturnType<typeof buildUpdateRobotSchema>>;

interface EditRobotFormProps {
  robotId: string;
  defaultValues: Robot;
  onSuccess?: () => void;
  onCancel?: () => void;
}

export function EditRobotForm({
  robotId,
  defaultValues,
  onSuccess,
  onCancel,
}: EditRobotFormProps) {
  const { t } = useTranslation();
  const getRobotTypeLabel = useRobotTypeLabel();
  const { mutate, isPending } = useUpdateRobotMutation();

  // Auto-select site from robot's current location
  const { data: currentLocation } = useLocationQuery(
    defaultValues.location_id ?? "",
    { enabled: !!defaultValues.location_id }
  );
  const [siteOverride, setSiteOverride] = useState<string | null>(null);
  const selectedSiteId = siteOverride ?? currentLocation?.site_id ?? "";

  // Fetch sites and locations for dropdowns
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

  // Map resolved status to the manually settable DB value.
  // Online/Offline → Ready (resolved from Redis, DB value is Ready)
  // Busy → undefined (cannot be set manually; exclude from request)
  const defaultStatus = (() => {
    const s = defaultValues.status;
    if (s === ROBOT_STATUS.ONLINE || s === ROBOT_STATUS.OFFLINE)
      return ROBOT_STATUS.READY;
    if (s === ROBOT_STATUS.BUSY) return undefined;
    return s;
  })();

  const updateRobotSchema = useMemo(() => buildUpdateRobotSchema(t), [t]);

  const { host: existingHost, advancedSettings: defaultRobotConfig } =
    splitHostFromRobotConfig(defaultValues.robot_config);

  const form = useForm<RobotUpdateInput>({
    resolver: zodResolver(updateRobotSchema),
    defaultValues: {
      name: defaultValues.name,
      organization_id: defaultValues.organization_id,
      location_id: defaultValues.location_id,
      robot_type:
        defaultValues.robot_type === ROBOT_TYPE.YUBI_STATIONARY ||
        defaultValues.robot_type === ROBOT_TYPE.YUBI_PORTABLE
          ? defaultValues.robot_type
          : undefined,
      leader_status: defaultValues.leader_status,
      status: defaultStatus,
      last_heartbeat_at: defaultValues.last_heartbeat_at,
      offline_reason: defaultValues.offline_reason,
      host: existingHost,
      robot_config: defaultRobotConfig,
    },
  });

  const onSubmit = (data: RobotUpdateInput) => {
    const { host, robot_config, ...rest } = data;
    const merged = mergeRobotConfigWithHost(robot_config, host);
    if (!merged) {
      setRobotConfigError(t("robotForm.invalidJson"));
      return;
    }

    mutate(
      { robotId, data: { ...rest, robot_config: merged } },
      {
        onSuccess: () => {
          // Textarea isn't disabled during mutate; clear any stale "Invalid JSON" error.
          setRobotConfigError("");
          onSuccess?.();
        },
      }
    );
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
              <FormLabel>{t("robotForm.name")}</FormLabel>
              <FormControl>
                <Input
                  placeholder={t("robotForm.namePlaceholder")}
                  {...field}
                  value={field.value || ""}
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
                setSiteOverride(v);
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
                      value: ROBOT_TYPE.YUBI_STATIONARY,
                      label: getRobotTypeLabel(ROBOT_TYPE.YUBI_STATIONARY),
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

        {/* Follower Status */}
        <FormField
          control={form.control}
          name="status"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("robotForm.followerStatus")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value?.toString() ?? ""}
                  onValueChange={(value) => field.onChange(parseInt(value, 10))}
                  options={[
                    {
                      value: String(ROBOT_STATUS.READY),
                      label: t("status.ready"),
                    },
                    {
                      value: String(ROBOT_STATUS.FAULTED),
                      label: t("status.faulted"),
                    },
                    {
                      value: String(ROBOT_STATUS.MAINTENANCE),
                      label: t("status.maintenance"),
                    },
                  ]}
                  placeholder={t("robotForm.selectStatus")}
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
                  placeholder={t("robotForm.notSet")}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Last Heartbeat */}
        <FormField
          control={form.control}
          name="last_heartbeat_at"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("robotForm.lastHeartbeat")}</FormLabel>
              <FormControl>
                <Input
                  type="datetime-local"
                  {...field}
                  value={
                    field.value
                      ? new Date(field.value).toISOString().slice(0, 16)
                      : ""
                  }
                  onChange={(e) => {
                    const value = e.target.value;
                    field.onChange(value ? new Date(value).toISOString() : "");
                  }}
                />
              </FormControl>
              <FormDescription>
                {t("robotForm.lastHeartbeatDescription")}
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Offline Reason */}
        <FormField
          control={form.control}
          name="offline_reason"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("robotForm.offlineReason")}</FormLabel>
              <FormControl>
                <Textarea
                  placeholder={t("robotForm.offlineReasonPlaceholder")}
                  rows={2}
                  {...field}
                  value={field.value || ""}
                />
              </FormControl>
              <FormDescription>
                {t("robotForm.offlineReasonDescription")}
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Host */}
        <FormField
          control={form.control}
          name="host"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("robotForm.host")}</FormLabel>
              <FormControl>
                <Input placeholder="192.168.1.101" {...field} />
              </FormControl>
              <FormDescription>
                {t("robotForm.hostDescription")}
              </FormDescription>
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
              <FormLabel>{t("robotForm.advancedSettings")}</FormLabel>
              <FormControl>
                <Textarea
                  placeholder='{"port": 9090, "cameras": [{"namespace": "camera_0", "name": "Front"}]}'
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
                {t("robotForm.advancedSettingsDescription")}
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
            {isPending ? t("dialog.saving") : t("robotForm.updateRobot")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
