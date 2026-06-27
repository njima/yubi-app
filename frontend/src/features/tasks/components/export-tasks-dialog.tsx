"use client";

import { Download } from "lucide-react";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import {
  useTaskStatusLabel,
  useTaskPriorityLabel,
} from "@/lib/hooks/use-status-labels";
import {
  TASK_STATUS,
  TASK_PRIORITY,
  TASK_DIFFICULTY,
} from "@/lib/status/constants";
import { getTaskDifficultyLabel } from "@/lib/status/utils";

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

import { useRobotsQuery } from "@/features/robots";

import { useExportTasksMutation } from "../hooks/use-export-tasks";

export function ExportTasksDialog() {
  const { t } = useTranslation();
  const getTaskStatusLabel = useTaskStatusLabel();
  const getTaskPriorityLabel = useTaskPriorityLabel();

  const statusOptions = useMemo(
    () => [
      {
        value: TASK_STATUS.PLANNING,
        label: getTaskStatusLabel(TASK_STATUS.PLANNING),
      },
      {
        value: TASK_STATUS.DOING,
        label: getTaskStatusLabel(TASK_STATUS.DOING),
      },
      {
        value: TASK_STATUS.COMPLETED,
        label: getTaskStatusLabel(TASK_STATUS.COMPLETED),
      },
      {
        value: TASK_STATUS.CANCELED,
        label: getTaskStatusLabel(TASK_STATUS.CANCELED),
      },
    ],
    [getTaskStatusLabel]
  );

  const priorityOptions = useMemo(
    () => [
      {
        value: TASK_PRIORITY.LOW,
        label: getTaskPriorityLabel(TASK_PRIORITY.LOW),
      },
      {
        value: TASK_PRIORITY.NORMAL,
        label: getTaskPriorityLabel(TASK_PRIORITY.NORMAL),
      },
      {
        value: TASK_PRIORITY.HIGH,
        label: getTaskPriorityLabel(TASK_PRIORITY.HIGH),
      },
      {
        value: TASK_PRIORITY.URGENT,
        label: getTaskPriorityLabel(TASK_PRIORITY.URGENT),
      },
    ],
    [getTaskPriorityLabel]
  );

  const difficultyOptions = useMemo(
    () => [
      {
        value: TASK_DIFFICULTY.S,
        label: getTaskDifficultyLabel(TASK_DIFFICULTY.S),
      },
      {
        value: TASK_DIFFICULTY.A,
        label: getTaskDifficultyLabel(TASK_DIFFICULTY.A),
      },
      {
        value: TASK_DIFFICULTY.B,
        label: getTaskDifficultyLabel(TASK_DIFFICULTY.B),
      },
      {
        value: TASK_DIFFICULTY.C,
        label: getTaskDifficultyLabel(TASK_DIFFICULTY.C),
      },
    ],
    []
  );

  const [open, setOpen] = useState(false);
  const [selectedStatuses, setSelectedStatuses] = useState<number[]>([]);
  const [selectedPriorities, setSelectedPriorities] = useState<number[]>([]);
  const [selectedDifficulties, setSelectedDifficulties] = useState<number[]>(
    []
  );
  const [robotType, setRobotType] = useState<string>("");
  const [exportError, setExportError] = useState<string | null>(null);

  const { data: robotsData } = useRobotsQuery();
  const uniqueRobotTypes = useMemo(
    () => [
      ...new Set(
        (robotsData?.robots ?? [])
          .map((r) => r.robot_type)
          .filter((m): m is string => !!m)
      ),
    ],
    [robotsData]
  );

  const exportMutation = useExportTasksMutation();

  function toggleValue<T>(arr: T[], val: T): T[] {
    return arr.includes(val) ? arr.filter((v) => v !== val) : [...arr, val];
  }

  const handleExport = () => {
    exportMutation.mutate(
      {
        status: selectedStatuses.length > 0 ? selectedStatuses : undefined,
        priority:
          selectedPriorities.length > 0 ? selectedPriorities : undefined,
        difficulty:
          selectedDifficulties.length > 0 ? selectedDifficulties : undefined,
        robot_type: robotType || undefined,
      },
      {
        onSuccess: () => {
          toast.success(t("exportTasksDialog.title"));
          setOpen(false);
        },
        onError: (err) => {
          setExportError(err.message || t("exportTasksDialog.exportFailed"));
        },
      }
    );
  };

  const handleOpenChange = (next: boolean) => {
    if (!next) {
      setSelectedStatuses([]);
      setSelectedPriorities([]);
      setSelectedDifficulties([]);
      setRobotType("");
      setExportError(null);
    }
    setOpen(next);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">
          <Download className="mr-2 h-4 w-4" />
          {t("dialog.export")}
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{t("exportTasksDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("exportTasksDialog.description")}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-5 py-2">
          {/* Status */}
          <div className="space-y-2">
            <Label className="text-sm font-medium">
              {t("exportTasksDialog.status")}
            </Label>
            <div className="flex flex-wrap gap-3">
              {statusOptions.map((opt) => (
                <div key={opt.value} className="flex items-center gap-1.5">
                  <Checkbox
                    id={`status-${opt.value}`}
                    checked={selectedStatuses.includes(opt.value)}
                    onCheckedChange={() =>
                      setSelectedStatuses(
                        toggleValue(selectedStatuses, opt.value)
                      )
                    }
                  />
                  <label
                    htmlFor={`status-${opt.value}`}
                    className="text-sm cursor-pointer"
                  >
                    {opt.label}
                  </label>
                </div>
              ))}
            </div>
          </div>

          {/* Priority */}
          <div className="space-y-2">
            <Label className="text-sm font-medium">
              {t("exportTasksDialog.priority")}
            </Label>
            <div className="flex flex-wrap gap-3">
              {priorityOptions.map((opt) => (
                <div key={opt.value} className="flex items-center gap-1.5">
                  <Checkbox
                    id={`priority-${opt.value}`}
                    checked={selectedPriorities.includes(opt.value)}
                    onCheckedChange={() =>
                      setSelectedPriorities(
                        toggleValue(selectedPriorities, opt.value)
                      )
                    }
                  />
                  <label
                    htmlFor={`priority-${opt.value}`}
                    className="text-sm cursor-pointer"
                  >
                    {opt.label}
                  </label>
                </div>
              ))}
            </div>
          </div>

          {/* Difficulty */}
          <div className="space-y-2">
            <Label className="text-sm font-medium">
              {t("exportTasksDialog.difficulty")}
            </Label>
            <div className="flex flex-wrap gap-3">
              {difficultyOptions.map((opt) => (
                <div key={opt.value} className="flex items-center gap-1.5">
                  <Checkbox
                    id={`difficulty-${opt.value}`}
                    checked={selectedDifficulties.includes(opt.value)}
                    onCheckedChange={() =>
                      setSelectedDifficulties(
                        toggleValue(selectedDifficulties, opt.value)
                      )
                    }
                  />
                  <label
                    htmlFor={`difficulty-${opt.value}`}
                    className="text-sm cursor-pointer"
                  >
                    {opt.label}
                  </label>
                </div>
              ))}
            </div>
          </div>

          {/* Robot Type */}
          <div className="space-y-2">
            <Label className="text-sm font-medium">
              {t("exportTasksDialog.robotType")}
            </Label>
            <SearchableSelect
              value={robotType}
              onValueChange={setRobotType}
              options={[
                { value: "", label: t("exportTasksDialog.allRobotTypes") },
                ...uniqueRobotTypes.map((m) => ({ value: m, label: m })),
              ]}
              placeholder={t("exportTasksDialog.allRobotTypes")}
            />
          </div>
        </div>

        {exportError && (
          <p className="text-sm text-red-600 dark:text-red-400">
            {exportError}
          </p>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => handleOpenChange(false)}>
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
