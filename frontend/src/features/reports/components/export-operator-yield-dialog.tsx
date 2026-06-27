"use client";

import { Download } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
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

import { useLocationsQuery } from "@/features/locations";
import { useTasksQuery } from "@/features/tasks";
import { useUsersQuery } from "@/features/users";

import { useExportOperatorYieldMutation } from "../hooks/use-export-operator-yield";

function defaultDateFromAndTo(): { from: string; to: string } {
  const today = new Date();
  const lastMonth = new Date(today);
  lastMonth.setDate(today.getDate() - 30);
  const pad = (n: number) => String(n).padStart(2, "0");
  const fmt = (d: Date) =>
    `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
  return { from: fmt(lastMonth), to: fmt(today) };
}

type ExportOperatorYieldDialogProps = {
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  showTrigger?: boolean;
};

export function ExportOperatorYieldDialog({
  open: controlledOpen,
  onOpenChange: controlledOnOpenChange,
  showTrigger = true,
}: ExportOperatorYieldDialogProps = {}) {
  const { t } = useTranslation();
  const [internalOpen, setInternalOpen] = useState(false);
  const open = controlledOpen ?? internalOpen;
  const setOpen = controlledOnOpenChange ?? setInternalOpen;

  const initial = defaultDateFromAndTo();
  const [dateFrom, setDateFrom] = useState(initial.from);
  const [dateTo, setDateTo] = useState(initial.to);
  const [taskId, setTaskId] = useState("");
  const [locationId, setLocationId] = useState("");
  const [userId, setUserId] = useState("");
  const [exportError, setExportError] = useState<string | null>(null);

  // Bumped above the API default so dropdowns show every option in single-tenant
  // sized data sets. Switch to async search if a tenant exceeds 1000 of any kind.
  const { data: tasksData } = useTasksQuery({ limit: 1000 });
  const tasks = tasksData?.tasks ?? [];

  const { data: locationsData } = useLocationsQuery({ limit: 1000 });
  const locations = locationsData?.locations ?? [];

  const { data: usersData } = useUsersQuery({ limit: 1000 });
  const users = usersData?.users ?? [];

  const exportMutation = useExportOperatorYieldMutation();

  const handleExport = () => {
    if (!dateFrom || !dateTo) {
      setExportError(t("exportOperatorYieldDialog.dateRequiredError"));
      return;
    }
    if (dateFrom > dateTo) {
      setExportError(t("exportOperatorYieldDialog.dateRangeError"));
      return;
    }
    exportMutation.mutate(
      {
        date_from: dateFrom,
        date_to: dateTo,
        task_id: taskId || undefined,
        location_id: locationId || undefined,
        user_id: userId || undefined,
      },
      {
        onSuccess: () => {
          toast.success(t("exportOperatorYieldDialog.successToast"));
          setOpen(false);
        },
        onError: (err) => {
          setExportError(
            err.message || t("exportOperatorYieldDialog.exportFailed")
          );
        },
      }
    );
  };

  const handleOpenChange = (next: boolean) => {
    if (!next) {
      const reset = defaultDateFromAndTo();
      setDateFrom(reset.from);
      setDateTo(reset.to);
      setTaskId("");
      setLocationId("");
      setUserId("");
      setExportError(null);
    }
    setOpen(next);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      {showTrigger && (
        <DialogTrigger asChild>
          <Button variant="outline" size="sm">
            <Download className="mr-2 h-4 w-4" />
            {t("exportOperatorYieldDialog.triggerLabel")}
          </Button>
        </DialogTrigger>
      )}
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{t("exportOperatorYieldDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("exportOperatorYieldDialog.description")}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-2">
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label className="text-sm font-medium">
                {t("exportOperatorYieldDialog.dateFrom")}
              </Label>
              <input
                type="date"
                value={dateFrom}
                required
                onChange={(e) => setDateFrom(e.target.value)}
                className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-950 dark:border-gray-800 dark:bg-gray-950"
              />
            </div>
            <div className="space-y-1.5">
              <Label className="text-sm font-medium">
                {t("exportOperatorYieldDialog.dateTo")}
              </Label>
              <input
                type="date"
                value={dateTo}
                required
                onChange={(e) => setDateTo(e.target.value)}
                className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-950 dark:border-gray-800 dark:bg-gray-950"
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label className="text-sm font-medium">
              {t("exportOperatorYieldDialog.task")}
            </Label>
            <SearchableSelect
              value={taskId}
              onValueChange={setTaskId}
              options={[
                { value: "", label: t("exportOperatorYieldDialog.allTasks") },
                ...tasks.map((task) => ({ value: task.id, label: task.name })),
              ]}
              placeholder={t("exportOperatorYieldDialog.allTasks")}
            />
          </div>

          <div className="space-y-1.5">
            <Label className="text-sm font-medium">
              {t("exportOperatorYieldDialog.location")}
            </Label>
            <SearchableSelect
              value={locationId}
              onValueChange={setLocationId}
              options={[
                {
                  value: "",
                  label: t("exportOperatorYieldDialog.allLocations"),
                },
                ...locations.map((loc) => ({ value: loc.id, label: loc.name })),
              ]}
              placeholder={t("exportOperatorYieldDialog.allLocations")}
            />
          </div>

          <div className="space-y-1.5">
            <Label className="text-sm font-medium">
              {t("exportOperatorYieldDialog.user")}
            </Label>
            <SearchableSelect
              value={userId}
              onValueChange={setUserId}
              options={[
                { value: "", label: t("exportOperatorYieldDialog.allUsers") },
                ...users.map((u) => ({
                  value: u.user_id,
                  label: u.display_name,
                })),
              ]}
              placeholder={t("exportOperatorYieldDialog.allUsers")}
            />
          </div>

          <p className="text-xs text-gray-500 dark:text-gray-400">
            {t("exportOperatorYieldDialog.cleansingNotice")}
          </p>
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
