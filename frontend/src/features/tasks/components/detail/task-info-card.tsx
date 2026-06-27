import { useTranslation } from "react-i18next";

import { APPROVAL_STATUS } from "@/shared/lib/status-constants";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { secondsToHoursMinutes } from "../../lib/duration";
import { TagCategoryBadge } from "../tag-category-badge";
import { TaskDifficultyBadge } from "../task-difficulty-badge";
import { TaskPriorityBadge } from "../task-priority-badge";
import { TaskStatusBadge } from "../task-status-badge";

import type { Task, TaskVersion } from "../../schemas";

interface TaskInfoCardProps {
  task: Task;
  displayVersion?: string;
  selectedVersion?: TaskVersion;
  approvalStatus?: number;
}

export function TaskInfoCard({
  task,
  displayVersion,
  selectedVersion,
  approvalStatus,
}: TaskInfoCardProps) {
  const { t } = useTranslation();
  const deadline = new Date(task.deadline);
  const tags = task.tags ?? [];

  const targetDuration = selectedVersion?.target_duration_seconds
    ? secondsToHoursMinutes(selectedVersion.target_duration_seconds)
    : null;

  const perEpisodeDuration =
    selectedVersion?.target_duration_per_episode_seconds
      ? secondsToHoursMinutes(
          selectedVersion.target_duration_per_episode_seconds
        )
      : null;

  const actualDuration =
    selectedVersion?.actual_duration_seconds != null
      ? secondsToHoursMinutes(selectedVersion.actual_duration_seconds)
      : null;

  const remainingSeconds =
    selectedVersion?.target_duration_seconds &&
    selectedVersion?.actual_duration_seconds != null
      ? selectedVersion.target_duration_seconds -
        selectedVersion.actual_duration_seconds
      : null;

  const remainingDuration =
    remainingSeconds != null
      ? secondsToHoursMinutes(Math.abs(remainingSeconds))
      : null;

  return (
    <div className="space-y-4">
      {/* Basic Info */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base font-medium">
            {t("taskInfo.basicInfo")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                {t("taskInfo.status")}
              </p>
              <TaskStatusBadge status={task.status} />
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                {t("taskInfo.priority")}
              </p>
              <TaskPriorityBadge priority={task.priority} />
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                {t("taskInfo.difficulty")}
              </p>
              <TaskDifficultyBadge difficulty={task.difficulty} />
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                {t("taskInfo.robotType")}
              </p>
              <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
                {task.robot_type ?? "-"}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                {t("taskInfo.categoryTags")}
              </p>
              {tags.length > 0 ? (
                <div className="flex flex-wrap gap-1">
                  {tags.map((tag) => (
                    <TagCategoryBadge
                      key={tag.id}
                      categoryTypeName={tag.category_type_name}
                      name={`${tag.category_type_name}: ${tag.name}`}
                    />
                  ))}
                </div>
              ) : (
                <p className="text-sm text-gray-400 dark:text-gray-500">-</p>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Collection Targets */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base font-medium">
            {t("taskInfo.collectionTargets")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                {t("taskInfo.targetDuration")}
              </p>
              <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                {targetDuration
                  ? t("taskInfo.durationFormat", {
                      hours: targetDuration.hours,
                      minutes: targetDuration.minutes,
                    })
                  : "-"}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                {t("taskInfo.targetEpisodes")}
              </p>
              <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                {selectedVersion?.target_episode_count ?? "-"}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                {t("taskInfo.perEpisode")}
              </p>
              <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                {perEpisodeDuration
                  ? t("taskInfo.durationFormat", {
                      hours: perEpisodeDuration.hours,
                      minutes: perEpisodeDuration.minutes,
                    })
                  : "-"}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                {t("taskInfo.deadline")}
              </p>
              <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                {Number.isNaN(deadline.getTime())
                  ? "-"
                  : deadline.toLocaleDateString()}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Collection Progress */}
      {(actualDuration != null ||
        selectedVersion?.actual_episode_count != null ||
        selectedVersion?.target_duration_seconds != null ||
        selectedVersion?.target_episode_count != null) && (
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-medium">
              {t("taskInfo.collectionProgress")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                  {t("taskInfo.duration")}
                </p>
                <p className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-1">
                  {remainingDuration != null
                    ? remainingSeconds != null && remainingSeconds <= 0
                      ? t("taskInfo.completed")
                      : t("taskInfo.remaining", {
                          value: t("taskInfo.durationFormat", {
                            hours: remainingDuration.hours,
                            minutes: remainingDuration.minutes,
                          }),
                        })
                    : actualDuration
                      ? t("taskInfo.durationFormat", {
                          hours: actualDuration.hours,
                          minutes: actualDuration.minutes,
                        })
                      : t("taskInfo.durationFormat", { hours: 0, minutes: 0 })}
                </p>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  {actualDuration
                    ? t("taskInfo.durationFormat", {
                        hours: actualDuration.hours,
                        minutes: actualDuration.minutes,
                      })
                    : t("taskInfo.durationFormat", { hours: 0, minutes: 0 })}
                  {targetDuration &&
                    ` / ${t("taskInfo.durationFormat", { hours: targetDuration.hours, minutes: targetDuration.minutes })}`}
                </p>
              </div>
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                  {t("taskInfo.episodes")}
                </p>
                <p className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-1">
                  {selectedVersion?.target_episode_count != null &&
                  selectedVersion?.actual_episode_count != null
                    ? selectedVersion.actual_episode_count >=
                      selectedVersion.target_episode_count
                      ? t("taskInfo.completed")
                      : t("taskInfo.remaining", {
                          value: String(
                            selectedVersion.target_episode_count -
                              selectedVersion.actual_episode_count
                          ),
                        })
                    : (selectedVersion?.actual_episode_count ?? 0)}
                </p>
                {selectedVersion?.target_episode_count != null &&
                  selectedVersion.target_episode_count > 0 && (
                    <div className="h-1.5 w-full rounded-full bg-gray-200 dark:bg-gray-700 mb-1">
                      <div
                        className={`h-full rounded-full ${
                          (selectedVersion?.actual_episode_count ?? 0) >=
                          selectedVersion.target_episode_count
                            ? "bg-green-500"
                            : "bg-blue-500"
                        }`}
                        style={{
                          width: `${Math.min(
                            Math.round(
                              ((selectedVersion?.actual_episode_count ?? 0) /
                                selectedVersion.target_episode_count) *
                                100
                            ),
                            100
                          )}%`,
                        }}
                      />
                    </div>
                  )}
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  {selectedVersion?.actual_episode_count ?? 0}
                  {selectedVersion?.target_episode_count != null &&
                    ` / ${selectedVersion.target_episode_count}`}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Version */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base font-medium">
            {t("taskInfo.version")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                {t("taskInfo.currentVersion")}
              </p>
              <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                {displayVersion ?? task.version ?? "-"}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                {t("taskInfo.approvalStatus")}
              </p>
              {approvalStatus === APPROVAL_STATUS.APPROVED ? (
                <Badge className="bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300">
                  {t("taskDetail.approved")}
                </Badge>
              ) : (
                <Badge variant="secondary">{t("taskDetail.draft")}</Badge>
              )}
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                {t("taskInfo.lastUpdated")}
              </p>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                {t("taskInfo.notAvailable")}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
