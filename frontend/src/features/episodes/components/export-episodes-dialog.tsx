"use client";

import { Download } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { useEpisodeCollectionStatusLabel } from "@/lib/hooks/use-status-labels";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { SearchableSelect } from "@/components/ui/searchable-select";

import { useRobotSearchOptions } from "@/features/robots";
import { useTaskSearchOptions, useTaskVersionsQuery } from "@/features/tasks";
import { useUserSearchOptions } from "@/features/users";

import { episodeStatusOptions } from "../constants";
import { useExportEpisodesMutation } from "../hooks/use-export-episodes";

export type ExportEpisodesInitialFilters = {
  taskId?: string;
  taskLabel?: string;
  taskVersionId?: string;
  taskVersionLabel?: string;
  robotId?: string;
  robotLabel?: string;
  userId?: string;
  userLabel?: string;
  statuses?: number[];
  startedAtFrom?: string;
  startedAtTo?: string;
};

type ExportEpisodesDialogProps = {
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  showTrigger?: boolean;
  initialFilters?: ExportEpisodesInitialFilters;
};

export function ExportEpisodesDialog({
  open: controlledOpen,
  onOpenChange: controlledOnOpenChange,
  showTrigger = true,
  initialFilters,
}: ExportEpisodesDialogProps = {}) {
  const { t } = useTranslation();
  const getEpisodeCollectionStatusLabel = useEpisodeCollectionStatusLabel();
  const [internalOpen, setInternalOpen] = useState(false);
  const open = controlledOpen ?? internalOpen;
  const setOpen = controlledOnOpenChange ?? setInternalOpen;

  const [taskId, setTaskId] = useState(initialFilters?.taskId ?? "");
  const [taskVersionId, setTaskVersionId] = useState(
    initialFilters?.taskVersionId ?? ""
  );
  const [robotId, setRobotId] = useState(initialFilters?.robotId ?? "");
  const [userId, setUserId] = useState(initialFilters?.userId ?? "");
  const [selectedStatuses, setSelectedStatuses] = useState<number[]>(
    initialFilters?.statuses ?? []
  );
  const [startedAtFrom, setStartedAtFrom] = useState(
    initialFilters?.startedAtFrom ?? ""
  );
  const [startedAtTo, setStartedAtTo] = useState(
    initialFilters?.startedAtTo ?? ""
  );
  const [exportError, setExportError] = useState<string | null>(null);

  const {
    options: taskSearchOptions,
    isLoading: taskSearchLoading,
    onSearch: onTaskSearch,
    selectedLabel: taskSelectedLabel,
    onValueChange: onTaskSelectChange,
  } = useTaskSearchOptions();

  const { data: taskVersions } = useTaskVersionsQuery(taskId, {
    enabled: !!taskId,
  });

  const {
    options: robotSearchOptions,
    isLoading: robotSearchLoading,
    onSearch: onRobotSearch,
    selectedLabel: robotSelectedLabel,
    onValueChange: onRobotSelectChange,
  } = useRobotSearchOptions();

  const {
    options: userSearchOptions,
    isLoading: userSearchLoading,
    onSearch: onUserSearch,
    selectedLabel: userSelectedLabel,
    onValueChange: onUserSelectChange,
  } = useUserSearchOptions();

  const exportMutation = useExportEpisodesMutation();

  const handleExport = () => {
    if (!!startedAtFrom !== !!startedAtTo) {
      setExportError(t("exportEpisodesDialog.dateRangeBothRequired"));
      return;
    }
    if (startedAtFrom && startedAtTo && startedAtFrom > startedAtTo) {
      setExportError(t("exportEpisodesDialog.dateRangeError"));
      return;
    }
    exportMutation.mutate(
      {
        task_id: taskId || undefined,
        task_version_id: taskVersionId || undefined,
        robot_id: robotId || undefined,
        user_id: userId || undefined,
        status: selectedStatuses.length > 0 ? selectedStatuses : undefined,
        started_at_from: startedAtFrom || undefined,
        started_at_to: startedAtTo || undefined,
      },
      {
        onSuccess: () => {
          toast.success(t("exportEpisodesDialog.title"));
          setOpen(false);
        },
        onError: (err) => {
          setExportError(err.message || "Export failed");
        },
      }
    );
  };

  const toggleStatus = (val: number) => {
    setSelectedStatuses((prev) =>
      prev.includes(val) ? prev.filter((v) => v !== val) : [...prev, val]
    );
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      {showTrigger && (
        <DialogTrigger asChild>
          <Button variant="outline" size="sm">
            <Download className="mr-2 h-4 w-4" />
            {t("dialog.export")}
          </Button>
        </DialogTrigger>
      )}
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{t("exportEpisodesDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("exportEpisodesDialog.description")}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-2">
          {/* Task */}
          <div className="space-y-1.5">
            <Label className="text-sm font-medium">
              {t("exportEpisodesDialog.task")}
            </Label>
            <SearchableSelect
              value={taskId}
              onValueChange={(value) => {
                setTaskId(value);
                setTaskVersionId("");
                onTaskSelectChange(value);
              }}
              options={[
                { value: "", label: t("exportEpisodesDialog.allTasks") },
                ...taskSearchOptions,
              ]}
              onSearch={onTaskSearch}
              isLoading={taskSearchLoading}
              selectedLabel={
                taskId
                  ? (taskSelectedLabel ?? initialFilters?.taskLabel)
                  : undefined
              }
              placeholder="All Tasks"
            />
          </div>

          {/* Task Version */}
          <div className="space-y-1.5">
            <Label className="text-sm font-medium">
              {t("exportEpisodesDialog.taskVersion")}
            </Label>
            <SearchableSelect
              value={taskVersionId}
              onValueChange={setTaskVersionId}
              options={[
                {
                  value: "",
                  label: taskId
                    ? t("exportEpisodesDialog.allVersions")
                    : t("exportEpisodesDialog.selectTaskFirst"),
                },
                ...(taskVersions ?? []).map((v) => ({
                  value: v.id,
                  label: v.version,
                })),
              ]}
              selectedLabel={
                taskVersionId && taskVersionId === initialFilters?.taskVersionId
                  ? initialFilters?.taskVersionLabel
                  : undefined
              }
              placeholder={
                taskId
                  ? t("exportEpisodesDialog.allVersions")
                  : t("exportEpisodesDialog.selectTaskFirst")
              }
              disabled={!taskId}
            />
          </div>

          {/* Robot */}
          <div className="space-y-1.5">
            <Label className="text-sm font-medium">
              {t("exportEpisodesDialog.robot")}
            </Label>
            <SearchableSelect
              value={robotId}
              onValueChange={(value) => {
                setRobotId(value);
                onRobotSelectChange(value);
              }}
              options={[
                { value: "", label: t("exportEpisodesDialog.allRobots") },
                ...robotSearchOptions,
              ]}
              onSearch={onRobotSearch}
              isLoading={robotSearchLoading}
              selectedLabel={
                robotId
                  ? (robotSelectedLabel ?? initialFilters?.robotLabel)
                  : undefined
              }
              placeholder={t("exportEpisodesDialog.allRobots")}
            />
          </div>

          {/* User */}
          <div className="space-y-1.5">
            <Label className="text-sm font-medium">
              {t("exportEpisodesDialog.user")}
            </Label>
            <SearchableSelect
              value={userId}
              onValueChange={(value) => {
                setUserId(value);
                onUserSelectChange(value);
              }}
              options={[
                { value: "", label: t("exportEpisodesDialog.allUsers") },
                ...userSearchOptions,
              ]}
              onSearch={onUserSearch}
              isLoading={userSearchLoading}
              selectedLabel={
                userId
                  ? (userSelectedLabel ?? initialFilters?.userLabel)
                  : undefined
              }
              placeholder={t("exportEpisodesDialog.allUsers")}
            />
          </div>

          {/* Status */}
          <div className="space-y-2">
            <Label className="text-sm font-medium">
              {t("exportEpisodesDialog.status")}
            </Label>
            <div className="flex flex-wrap gap-3">
              {episodeStatusOptions.map((s) => (
                <div key={s} className="flex items-center gap-1.5">
                  <Checkbox
                    id={`ep-status-${s}`}
                    checked={selectedStatuses.includes(s)}
                    onCheckedChange={() => toggleStatus(s)}
                  />
                  <label
                    htmlFor={`ep-status-${s}`}
                    className="text-sm cursor-pointer"
                  >
                    {getEpisodeCollectionStatusLabel(s)}
                  </label>
                </div>
              ))}
            </div>
          </div>

          {/* Date range */}
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label className="text-sm font-medium">
                {t("exportEpisodesDialog.startedAtFrom")}
              </Label>
              <input
                type="date"
                value={startedAtFrom}
                onChange={(e) => setStartedAtFrom(e.target.value)}
                className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-950 dark:border-gray-800 dark:bg-gray-950"
              />
            </div>
            <div className="space-y-1.5">
              <Label className="text-sm font-medium">
                {t("exportEpisodesDialog.startedAtTo")}
              </Label>
              <input
                type="date"
                value={startedAtTo}
                onChange={(e) => setStartedAtTo(e.target.value)}
                className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-950 dark:border-gray-800 dark:bg-gray-950"
              />
            </div>
          </div>
        </div>

        {exportError && (
          <p className="text-sm text-red-600 dark:text-red-400">
            {exportError}
          </p>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => setOpen(false)}>
            {t("dialog.cancel")}
          </Button>
          <Button
            onClick={() => {
              setExportError(null);
              handleExport();
            }}
            disabled={exportMutation.isPending}
          >
            <Download className="mr-2 h-4 w-4" />
            {exportMutation.isPending
              ? t("dialog.exporting")
              : t("dialog.downloadCsv")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
