"use client";

import {
  CheckCircle,
  ChevronDown,
  ChevronRight,
  Clock,
  ExternalLink,
  History,
} from "lucide-react";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { APPROVAL_STATUS } from "@/shared/lib/status-constants";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { useSubTasksByVersionQuery } from "../hooks/use-subtasks-by-version-query";
import { useTaskVersionsQuery } from "../hooks/use-task-versions-query";

import type { TaskVersion } from "../schemas";

interface VersionHistoryDialogProps {
  taskId: string;
  children?: React.ReactNode;
}

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

interface VersionItemProps {
  version: TaskVersion;
  taskId: string;
  isExpanded: boolean;
  onToggleExpand: () => void;
  onClose: () => void;
}

function VersionItem({
  version,
  taskId,
  isExpanded,
  onToggleExpand,
  onClose,
}: VersionItemProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const { data: subtasks, isLoading } = useSubTasksByVersionQuery(
    isExpanded ? version.id : undefined
  );

  const handleViewVersion = () => {
    onClose();
    if (version.is_current) {
      router.push(`/tasks/${taskId}`);
    } else {
      router.push(`/tasks/${taskId}?version=${version.id}`);
    }
  };

  return (
    <div
      className={`rounded-lg border ${
        version.is_current
          ? "bg-blue-50 border-blue-200 dark:bg-blue-950 dark:border-blue-800"
          : "bg-gray-50 border-gray-200 dark:bg-gray-800 dark:border-gray-700"
      }`}
    >
      {/* Header Row - Clickable for expansion */}
      <div
        className="flex items-center justify-between p-3 cursor-pointer"
        onClick={onToggleExpand}
        role="button"
        aria-expanded={isExpanded}
      >
        <div className="flex items-center gap-3">
          {isExpanded ? (
            <ChevronDown className="h-4 w-4 text-gray-500" />
          ) : (
            <ChevronRight className="h-4 w-4 text-gray-500" />
          )}
          <Badge variant={version.is_current ? "default" : "outline"}>
            {version.version}
          </Badge>
          {version.display_name && (
            <span
              className="text-sm text-gray-700 dark:text-gray-300 truncate max-w-[200px]"
              title={version.display_name}
            >
              {version.display_name}
            </span>
          )}
          {version.is_current && (
            <Badge variant="secondary" className="text-xs">
              {t("versionHistoryDialog.current")}
            </Badge>
          )}
          {version.approval_status === APPROVAL_STATUS.APPROVED ? (
            <Badge className="text-xs bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300 gap-1">
              <CheckCircle className="h-3 w-3" />
              Approved
            </Badge>
          ) : (
            <Badge variant="outline" className="text-xs">
              Draft
            </Badge>
          )}
        </div>
        <div className="flex items-center gap-1 text-sm text-gray-500 dark:text-gray-400">
          <Clock className="h-3 w-3" />
          {formatDate(version.created_at)}
        </div>
      </div>

      {/* Expanded Content */}
      {isExpanded && (
        <div className="px-3 pb-3 border-t border-gray-200 dark:border-gray-700">
          <div className="pt-3 space-y-3">
            {/* Parameters */}
            {version.parameters && version.parameters.length > 0 && (
              <div className="text-sm text-gray-600 dark:text-gray-400">
                <p className="font-medium mb-1">
                  {t("versionHistoryDialog.parameters")}
                </p>
                <div className="space-y-1">
                  {version.parameters.map((p) => (
                    <div key={p.key} className="flex items-center gap-2">
                      <code className="text-xs bg-gray-100 dark:bg-gray-700 px-1.5 py-0.5 rounded">
                        {"{" + p.key + "}"}
                      </code>
                      <span className="text-xs text-gray-500">
                        {p.values.join(", ")}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Subtasks Preview */}
            <div className="text-sm text-gray-600 dark:text-gray-400">
              <p className="font-medium mb-2">
                {t("versionHistoryDialog.subtasksLabel")}
              </p>
              {isLoading ? (
                <div className="space-y-1">
                  {[1, 2, 3].map((i) => (
                    <div
                      key={i}
                      className="h-5 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"
                    />
                  ))}
                </div>
              ) : subtasks && subtasks.length > 0 ? (
                <ul className="space-y-1 pl-4 list-disc max-h-40 overflow-y-auto">
                  {subtasks.map((subtask) => (
                    <li
                      key={subtask.id}
                      className="text-gray-700 dark:text-gray-300"
                    >
                      {subtask.name}
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="text-gray-500 italic">
                  {t("versionHistoryDialog.noSubtasks")}
                </p>
              )}
            </div>

            {/* View Version Button */}
            <Button
              variant="outline"
              size="sm"
              className="w-full gap-2"
              onClick={(e) => {
                e.stopPropagation();
                handleViewVersion();
              }}
            >
              <ExternalLink className="h-4 w-4" />
              {version.is_current
                ? t("versionHistoryDialog.viewCurrentVersion")
                : t("versionHistoryDialog.viewThisVersion")}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}

function VersionListSkeleton() {
  return (
    <div className="space-y-2 animate-pulse">
      {[1, 2, 3].map((i) => (
        <div
          key={i}
          className="h-14 w-full bg-gray-200 dark:bg-gray-700 rounded-lg"
        />
      ))}
    </div>
  );
}

export function VersionHistoryDialog({
  taskId,
  children,
}: VersionHistoryDialogProps) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const [expandedVersionId, setExpandedVersionId] = useState<string | null>(
    null
  );
  const {
    data: versions,
    isLoading,
    error,
  } = useTaskVersionsQuery(taskId, {
    enabled: open,
  });

  const handleToggleExpand = (versionId: string) => {
    setExpandedVersionId(expandedVersionId === versionId ? null : versionId);
  };

  return (
    <Dialog
      open={open}
      onOpenChange={(newOpen) => {
        setOpen(newOpen);
        if (!newOpen) {
          setExpandedVersionId(null);
        }
      }}
    >
      <DialogTrigger asChild>
        {children ?? (
          <Button variant="outline" size="sm" className="gap-2">
            <History className="h-4 w-4" />
            {t("versionHistoryDialog.title")}
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>{t("versionHistoryDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("versionHistoryDialog.description")}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-2 max-h-[60vh] overflow-y-auto">
          {isLoading && <VersionListSkeleton />}

          {error && (
            <div className="text-center py-4 text-red-600 dark:text-red-400">
              {t("versionHistoryDialog.failedToLoad", {
                message: error.message,
              })}
            </div>
          )}

          {!isLoading && !error && versions?.length === 0 && (
            <div className="text-center py-4 text-gray-500 dark:text-gray-400">
              {t("versionHistoryDialog.noVersions")}
            </div>
          )}

          {!isLoading &&
            !error &&
            versions?.map((version) => (
              <VersionItem
                key={version.id}
                version={version}
                taskId={taskId}
                isExpanded={expandedVersionId === version.id}
                onToggleExpand={() => handleToggleExpand(version.id)}
                onClose={() => setOpen(false)}
              />
            ))}
        </div>
      </DialogContent>
    </Dialog>
  );
}
