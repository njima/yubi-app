"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";
import { useEpisodeCollectionStatusLabel } from "@/lib/hooks/use-status-labels";
import { EPISODE_COLLECTION_STATUS } from "@/lib/status/constants";

import { Button } from "@/components/ui/button";
import { DateTimePicker } from "@/components/ui/datetime-picker";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { SearchableSelect } from "@/components/ui/searchable-select";
import { Textarea } from "@/components/ui/textarea";

import { useUserSearchOptions } from "@/features/users";

import { useEpisodeQuery } from "../hooks/use-episodes-query";
import { useUpdateEpisodeMutation } from "../hooks/use-update-episode-mutation";

const episodeUpdateSchema = schemas.EpisodeUpdate;

type Episode = z.infer<typeof schemas.Episode>;
type EpisodeUpdate = z.infer<typeof schemas.EpisodeUpdate>;

export interface EditEpisodeFormProps {
  episodeId: string;
  onSuccess?: () => void;
  onCancel?: () => void;
}

export function EditEpisodeForm({
  episodeId,
  onSuccess,
  onCancel,
}: EditEpisodeFormProps) {
  const { t } = useTranslation();
  const { data: episode, isLoading } = useEpisodeQuery(episodeId);

  if (isLoading || !episode) {
    return (
      <div className="py-8 text-center text-gray-600 dark:text-gray-400">
        {t("editEpisodeForm.loadingEpisodeData")}
      </div>
    );
  }

  return (
    <EditEpisodeFormInner
      episode={episode}
      episodeId={episodeId}
      onSuccess={onSuccess}
      onCancel={onCancel}
    />
  );
}

function EditEpisodeFormInner({
  episode,
  episodeId,
  onSuccess,
  onCancel,
}: {
  episode: Episode;
  episodeId: string;
  onSuccess?: () => void;
  onCancel?: () => void;
}) {
  const { t } = useTranslation();
  const getStatusLabel = useEpisodeCollectionStatusLabel();
  const { mutate, isPending } = useUpdateEpisodeMutation();
  const {
    options: userOptions,
    isLoading: usersLoading,
    onSearch: onUserSearch,
    selectedLabel: userLabel,
    onValueChange: onUserValueChange,
  } = useUserSearchOptions();

  const form = useForm<EpisodeUpdate>({
    resolver: zodResolver(episodeUpdateSchema),
    defaultValues: {
      status: episode.status,
      start_time: episode.started_at ?? undefined,
      end_time: episode.ended_at ?? undefined,
      error_details: episode.error_details ?? "",
      recorded_by: episode.recorded_by ?? "",
    },
  });

  const onSubmit = (data: EpisodeUpdate) => {
    // Clean up data: remove empty strings
    const cleanedData: Partial<EpisodeUpdate> = {};

    if (data.status !== undefined) {
      cleanedData.status = data.status;
    }
    if (data.start_time) {
      cleanedData.start_time = data.start_time;
    }
    if (data.end_time) {
      cleanedData.end_time = data.end_time;
    }
    if (data.error_details) {
      cleanedData.error_details = data.error_details;
    }
    if (data.recorded_by) {
      cleanedData.recorded_by = data.recorded_by;
    }

    // Data is already validated by React Hook Form with zodResolver
    mutate(
      { episodeId, data: cleanedData },
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
        {/* Status Selection */}
        <FormField
          control={form.control}
          name="status"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("episodesPage.status")}</FormLabel>
              <FormControl>
                <SearchableSelect
                  value={field.value?.toString() ?? ""}
                  onValueChange={(value) => {
                    const parsed = parseInt(value, 10);
                    field.onChange(
                      parsed === EPISODE_COLLECTION_STATUS.READY ||
                        parsed === EPISODE_COLLECTION_STATUS.RECORDING ||
                        parsed === EPISODE_COLLECTION_STATUS.CANCEL ||
                        parsed === EPISODE_COLLECTION_STATUS.COMPLETED
                        ? parsed
                        : undefined
                    );
                  }}
                  options={[
                    {
                      value: EPISODE_COLLECTION_STATUS.READY.toString(),
                      label: getStatusLabel(EPISODE_COLLECTION_STATUS.READY),
                    },
                    {
                      value: EPISODE_COLLECTION_STATUS.RECORDING.toString(),
                      label: getStatusLabel(
                        EPISODE_COLLECTION_STATUS.RECORDING
                      ),
                    },
                    {
                      value: EPISODE_COLLECTION_STATUS.CANCEL.toString(),
                      label: getStatusLabel(EPISODE_COLLECTION_STATUS.CANCEL),
                    },
                    {
                      value: EPISODE_COLLECTION_STATUS.COMPLETED.toString(),
                      label: getStatusLabel(
                        EPISODE_COLLECTION_STATUS.COMPLETED
                      ),
                    },
                  ]}
                  placeholder={t("editEpisodeForm.selectStatus")}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Start Time */}
        <FormField
          control={form.control}
          name="start_time"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("editEpisodeForm.startTimeOptional")}</FormLabel>
              <FormControl>
                <DateTimePicker value={field.value} onChange={field.onChange} />
              </FormControl>
              <FormDescription>
                {t("editEpisodeForm.iso8601WithTimezone")}
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* End Time */}
        <FormField
          control={form.control}
          name="end_time"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("editEpisodeForm.endTimeOptional")}</FormLabel>
              <FormControl>
                <DateTimePicker value={field.value} onChange={field.onChange} />
              </FormControl>
              <FormDescription>
                {t("editEpisodeForm.iso8601WithTimezone")}
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

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
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Error Details */}
        <FormField
          control={form.control}
          name="error_details"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("editEpisodeForm.errorDetailsOptional")}</FormLabel>
              <FormControl>
                <Textarea
                  placeholder={t("editEpisodeForm.errorDetailsPlaceholder")}
                  {...field}
                  value={field.value ?? ""}
                  rows={3}
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
            {isPending
              ? t("editEpisodeForm.updating")
              : t("editEpisodeForm.updateEpisode")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
