"use client";

import {
  AlertTriangle,
  ArrowLeft,
  CheckCircle,
  ChevronRight,
  Copy,
  History,
  Pencil,
  Save,
  Target,
} from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { toDateTimeLocalValue } from "@/shared/lib/date-utils";
import { APPROVAL_STATUS } from "@/shared/lib/status-constants";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent } from "@/shared/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/shared/ui/tabs";

import { usePermission } from "@/features/users";

import { useApproveTaskVersionMutation } from "../../hooks/use-approve-task-version-mutation";
import { useTaskVersionsQuery } from "../../hooks/use-task-versions-query";
import { useTaskQuery } from "../../hooks/use-tasks-query";
import { useUpdateTaskVersionParametersMutation } from "../../hooks/use-update-task-version-parameters-mutation";
import { EditTaskDialog } from "../edit-task-dialog";
import { EditTaskVersionTargetsDialog } from "../edit-task-version-targets-dialog";
import { SaveAsNewVersionDialog } from "../save-as-new-version-dialog";
import { SubTaskList } from "../subtask-list";
import {
  TaskVersionParametersEditor,
  type ParameterDefinition,
} from "../task-version-parameters-editor";
import { TeachMeBizCard } from "../teach-me-biz-card";
import { VersionHistoryDialog } from "../version-history-dialog";
import { TaskDescriptionCard } from "./task-description-card";
import { TaskInfoCard } from "./task-info-card";

interface TaskDetailPageProps {
  taskId: string;
  selectedVersionId?: string;
}

function VersionWarningBanner({
  taskId,
  versionString,
  t,
}: {
  taskId: string;
  versionString: string;
  t: (key: string, options?: Record<string, unknown>) => string;
}) {
  return (
    <div className="flex items-center justify-between p-4 rounded-lg bg-amber-50 border border-amber-200 dark:bg-amber-950 dark:border-amber-800">
      <div className="flex items-center gap-3">
        <AlertTriangle className="h-5 w-5 text-amber-600 dark:text-amber-400" />
        <div>
          <p className="font-medium text-amber-800 dark:text-amber-200">
            {t("taskDetail.viewingHistoricalVersion", {
              version: versionString,
            })}
          </p>
          <p className="text-sm text-amber-700 dark:text-amber-300">
            {t("taskDetail.notLatestVersion")}
          </p>
        </div>
      </div>
      <Link href={`/tasks/${taskId}`}>
        <Button variant="outline" size="sm">
          {t("taskDetail.viewLatest")}
        </Button>
      </Link>
    </div>
  );
}

export function TaskDetailPage({
  taskId,
  selectedVersionId,
}: TaskDetailPageProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const { data: task, isLoading, error } = useTaskQuery(taskId);
  const { data: versions } = useTaskVersionsQuery(taskId);

  const currentVersion = versions?.find((v) => v.is_current);
  const latestVersion = versions?.[0];
  const selectedVersion = selectedVersionId
    ? versions?.find((v) => v.id === selectedVersionId)
    : (currentVersion ?? latestVersion);

  const effectiveVersionId =
    selectedVersionId ?? currentVersion?.id ?? latestVersion?.id;

  const isViewingOldVersion =
    !!selectedVersionId && !selectedVersion?.is_current;
  const isApproved =
    selectedVersion?.approval_status === APPROVAL_STATUS.APPROVED;
  const isDraft = selectedVersion?.approval_status === APPROVAL_STATUS.DRAFT;
  const isReadOnly = isApproved;
  const approveMutation = useApproveTaskVersionMutation(taskId);
  const updateParametersMutation =
    useUpdateTaskVersionParametersMutation(taskId);
  const [approveDialogOpen, setApproveDialogOpen] = useState(false);
  const [previousVersionId, setPreviousVersionId] =
    useState(effectiveVersionId);
  const [editingParameters, setEditingParameters] = useState<
    ParameterDefinition[] | null
  >(null);

  if (previousVersionId !== effectiveVersionId) {
    setPreviousVersionId(effectiveVersionId);
    setEditingParameters(null);
  }

  const canEdit = usePermission("task:update");
  const canManageSubtasks = usePermission("subtask:create");

  if (isLoading) {
    return <TaskDetailSkeleton />;
  }

  if (error) {
    return (
      <div className="p-8 text-center">
        <div className="text-red-600 dark:text-red-400 mb-4">
          {t("taskDetail.errorLoadingTask", { message: error.message })}
        </div>
        <Button onClick={() => window.location.reload()}>
          {t("taskDetail.retry")}
        </Button>
      </div>
    );
  }

  if (!task) {
    return (
      <div className="p-8 text-center">
        <div className="text-gray-600 dark:text-gray-400 mb-4">
          {t("taskDetail.taskNotFound")}
        </div>
        <Button onClick={() => router.push("/tasks")}>
          {t("taskDetail.backToList")}
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Version Warning Banner */}
      {isViewingOldVersion && selectedVersion && (
        <VersionWarningBanner
          taskId={taskId}
          versionString={selectedVersion.version}
          t={t}
        />
      )}

      {/* Breadcrumb */}
      <nav className="flex items-center gap-1 text-sm text-gray-500 dark:text-gray-400">
        <Link
          href="/tasks"
          className="hover:text-gray-900 dark:hover:text-gray-100"
        >
          {t("topNav.tasks")}
        </Link>
        <ChevronRight className="h-4 w-4" />
        <span className="text-gray-900 dark:text-gray-100">{task.name}</span>
        {selectedVersion && (
          <>
            <ChevronRight className="h-4 w-4" />
            <Badge variant="outline">{selectedVersion.version}</Badge>
            {isApproved ? (
              <Badge className="bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300 gap-1">
                <CheckCircle className="h-3 w-3" />
                {t("taskDetail.approved")}
              </Badge>
            ) : (
              <Badge variant="secondary">{t("taskDetail.draft")}</Badge>
            )}
          </>
        )}
      </nav>

      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="flex items-start gap-4">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => router.push("/tasks")}
            className="mt-1"
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
              {task.name}
            </h1>
            <p className="text-gray-500 dark:text-gray-400 mt-1 text-sm">
              <span className="font-medium">{t("taskDetail.taskId")}:</span>{" "}
              <span className="font-mono">{task.id}</span>
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {/* Version History */}
          <VersionHistoryDialog taskId={task.id}>
            <Button variant="outline" size="sm" className="gap-2">
              <History className="h-4 w-4" />
              {t("taskDetail.versionHistory")}
            </Button>
          </VersionHistoryDialog>

          {/* Approve */}
          {isDraft && selectedVersion && (
            <>
              <Button
                size="sm"
                className="gap-2"
                disabled={approveMutation.isPending}
                onClick={() => setApproveDialogOpen(true)}
              >
                <CheckCircle className="h-4 w-4" />
                {approveMutation.isPending
                  ? t("taskDetail.approving")
                  : t("taskDetail.approve")}
              </Button>
              <Dialog
                open={approveDialogOpen}
                onOpenChange={setApproveDialogOpen}
              >
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>{t("taskDetail.approveVersion")}</DialogTitle>
                    <DialogDescription>
                      {t("taskDetail.approveVersionConfirm", {
                        version: selectedVersion.version,
                      })}
                    </DialogDescription>
                  </DialogHeader>
                  <DialogFooter>
                    <Button
                      variant="outline"
                      onClick={() => setApproveDialogOpen(false)}
                    >
                      {t("dialog.cancel")}
                    </Button>
                    <Button
                      disabled={approveMutation.isPending}
                      onClick={() => {
                        approveMutation.mutate(selectedVersion.id);
                        setApproveDialogOpen(false);
                      }}
                    >
                      <CheckCircle className="h-4 w-4 mr-1" />
                      {t("taskDetail.approve")}
                    </Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </>
          )}

          {/* Save as New Version */}
          {canEdit && (
            <SaveAsNewVersionDialog
              taskId={task.id}
              versions={versions ?? []}
              defaultBaseVersionId={effectiveVersionId ?? ""}
            >
              <Button variant="outline" size="sm" className="gap-2">
                <Copy className="h-4 w-4" />
                {t("taskDetail.saveAsNewVersion")}
              </Button>
            </SaveAsNewVersionDialog>
          )}

          {/* Edit Version Targets */}
          {canEdit && isDraft && selectedVersion && effectiveVersionId && (
            <EditTaskVersionTargetsDialog
              taskId={task.id}
              versionId={effectiveVersionId}
              defaultValues={{
                display_name: selectedVersion.display_name,
                target_duration_seconds:
                  selectedVersion.target_duration_seconds,
                target_episode_count: selectedVersion.target_episode_count,
                target_duration_per_episode_seconds:
                  selectedVersion.target_duration_per_episode_seconds,
              }}
            >
              <Button variant="outline" size="sm" className="gap-2">
                <Target className="h-4 w-4" />
                {t("taskDetail.editTargets")}
              </Button>
            </EditTaskVersionTargetsDialog>
          )}

          {/* Edit Task */}
          {canEdit && (
            <EditTaskDialog
              taskId={task.id}
              name={task.name}
              description={task.description}
              manual_url={task.manual_url}
              priority={task.priority}
              difficulty={task.difficulty}
              status={task.status}
              deadline={toDateTimeLocalValue(task.deadline)}
              robot_type={task.robot_type ?? undefined}
              tags={task.tags}
            >
              <Button variant="outline" size="sm">
                <Pencil className="mr-2 h-4 w-4" />
                {t("taskDetail.editTask")}
              </Button>
            </EditTaskDialog>
          )}
        </div>
      </div>

      {/* Tabs */}
      <Tabs defaultValue="overview" className="w-full">
        <TabsList>
          <TabsTrigger value="overview">{t("taskDetail.overview")}</TabsTrigger>
          <TabsTrigger value="subtasks">{t("taskDetail.subtasks")}</TabsTrigger>
          <TabsTrigger value="details">{t("taskDetail.details")}</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="mt-6">
          <TaskInfoCard
            task={task}
            displayVersion={selectedVersion?.version}
            selectedVersion={selectedVersion}
            approvalStatus={selectedVersion?.approval_status}
          />
        </TabsContent>

        <TabsContent value="subtasks" className="mt-6">
          <Card>
            <CardContent className="pt-6">
              <SubTaskList
                taskId={taskId}
                taskVersionId={effectiveVersionId}
                isReadOnly={isReadOnly || !canManageSubtasks}
              />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="details" className="mt-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <TaskDescriptionCard description={task.description} />
            <TeachMeBizCard manualUrl={task.manual_url} />
          </div>

          {/* Parameters */}
          {selectedVersion &&
            (isDraft && canEdit ? (
              <Card className="mt-6">
                <CardContent className="pt-6">
                  <TaskVersionParametersEditor
                    value={
                      editingParameters ?? selectedVersion.parameters ?? []
                    }
                    onChange={setEditingParameters}
                  />
                  {editingParameters !== null && (
                    <div className="flex justify-end gap-2 mt-4">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setEditingParameters(null)}
                      >
                        {t("taskDetail.cancel")}
                      </Button>
                      <Button
                        size="sm"
                        className="gap-1"
                        disabled={updateParametersMutation.isPending}
                        onClick={() => {
                          updateParametersMutation.mutate(
                            {
                              versionId: selectedVersion.id,
                              parameters: editingParameters,
                            },
                            {
                              onSuccess: () => setEditingParameters(null),
                            }
                          );
                        }}
                      >
                        <Save className="h-3 w-3" />
                        {updateParametersMutation.isPending
                          ? t("taskDetail.saving")
                          : t("taskDetail.saveParameters")}
                      </Button>
                    </div>
                  )}
                </CardContent>
              </Card>
            ) : (
              selectedVersion.parameters &&
              selectedVersion.parameters.length > 0 && (
                <Card className="mt-6">
                  <CardContent className="pt-6">
                    <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-3">
                      {t("taskDetail.parameters")}
                    </h3>
                    <div className="space-y-2">
                      {selectedVersion.parameters.map((p) => (
                        <div key={p.key} className="flex items-start gap-3">
                          <code className="text-sm bg-gray-100 dark:bg-gray-800 px-2 py-0.5 rounded font-mono">
                            {"{" + p.key + "}"}
                          </code>
                          <div className="flex flex-wrap gap-1">
                            {p.values.map((v) => (
                              <Badge
                                key={v}
                                variant="secondary"
                                className="text-xs"
                              >
                                {v}
                              </Badge>
                            ))}
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              )
            ))}
        </TabsContent>
      </Tabs>
    </div>
  );
}

function TaskDetailSkeleton() {
  return (
    <div className="space-y-6 animate-pulse">
      {/* Breadcrumb skeleton */}
      <div className="h-5 bg-gray-200 dark:bg-gray-700 rounded w-32"></div>

      {/* Header skeleton */}
      <div className="flex items-start justify-between">
        <div className="flex items-start gap-4">
          <div className="h-8 w-8 bg-gray-200 dark:bg-gray-700 rounded"></div>
          <div>
            <div className="h-9 bg-gray-200 dark:bg-gray-700 rounded w-48 mb-2"></div>
            <div className="h-5 bg-gray-200 dark:bg-gray-700 rounded w-64"></div>
          </div>
        </div>
        <div className="flex gap-2">
          <div className="h-9 bg-gray-200 dark:bg-gray-700 rounded w-36"></div>
          <div className="h-9 bg-gray-200 dark:bg-gray-700 rounded w-24"></div>
        </div>
      </div>

      {/* Info cards skeleton */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="h-24 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
        <div className="h-24 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
      </div>

      {/* Description skeleton */}
      <div className="h-32 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>

      {/* Subtasks skeleton */}
      <div className="h-48 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
    </div>
  );
}
