"use client";

import { ArrowLeft, ChevronRight, Pencil, Radio } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useTranslation } from "react-i18next";

import { formatDateTime } from "@/lib/date-utils";
import { EPISODE_COLLECTION_STATUS } from "@/lib/status/constants";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { useLocationsQuery } from "@/features/locations";
import {
  useRobotsQuery,
  useRobotScope,
  isRobotInScope,
} from "@/features/robots";
import { useTaskVersionsQuery } from "@/features/tasks";
import { usePermission } from "@/features/users";
import { useUsersQuery } from "@/features/users";

import { useEpisodeStream } from "../../hooks/use-episode-stream";
import { EditEpisodeDialog } from "../edit-episode-dialog";
import { EpisodeStatusBadge } from "../episode-status-badge";
import { EpisodeGradeSection } from "./episode-grade-section";
import { SubtaskListCard } from "./subtask-list-card";
import { RecordingTab } from "./tabs/recording-tab";
import { StatsTab } from "./tabs/stats-tab";

interface EpisodeDetailPageProps {
  episodeId: string;
}

export function EpisodeDetailPage({ episodeId }: EpisodeDetailPageProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const { data: episode, error } = useEpisodeStream(episodeId, true);
  const isLoading = !episode && !error;

  const { data: usersData } = useUsersQuery({ limit: 1000 });
  const users = usersData?.users;
  const { data: robotsData } = useRobotsQuery({ limit: 1000 });
  const robots = robotsData?.robots;
  const { data: locationsData } = useLocationsQuery({ limit: 1000 });
  const locations = locationsData?.locations;
  const { scopeIds } = useRobotScope();
  const { data: taskVersions } = useTaskVersionsQuery(episode?.task_id ?? "", {
    enabled: !!episode?.task_id,
  });

  const userNameById = new Map(
    (users || []).map((user) => [user.user_id, user.display_name])
  );
  const robotNameById = new Map(
    (robots || []).map((robot) => [robot.id, robot.name])
  );
  const locationNameById = new Map(
    (locations || []).map((location) => [location.id, location.name])
  );

  const canEdit = usePermission("episode:update");

  if (isLoading) {
    return <EpisodeDetailSkeleton />;
  }

  if (!episode) {
    return (
      <div className="p-8 text-center">
        <div className="text-gray-600 dark:text-gray-400 mb-4">
          {t("episodeDetail.episodeNotFound")}
        </div>
        <Button onClick={() => router.push("/episodes")}>
          {t("episodeDetail.backToList")}
        </Button>
      </div>
    );
  }

  const shortId = episode.id.substring(0, 8);
  const taskName = episode.task_name ?? episode.task_id ?? "-";
  const robotName = episode.robot_id
    ? (robotNameById.get(episode.robot_id) ?? episode.robot_id)
    : "-";
  const userName = episode.user_id
    ? (userNameById.get(episode.user_id) ?? episode.user_id)
    : "-";
  const locationName = episode.location_id
    ? (locationNameById.get(episode.location_id) ?? episode.location_id)
    : "-";
  const recordedByName = episode.recorded_by
    ? (userNameById.get(episode.recorded_by) ?? episode.recorded_by)
    : "-";
  const taskVersion = taskVersions?.find(
    (v) => v.id === episode.task_version_id
  );
  const taskVersionDisplay =
    episode.task_version_display_name ??
    (taskVersion
      ? `${taskName} ${taskVersion.version}`
      : (episode.task_version_id ?? "-"));

  return (
    <div className="space-y-6">
      {/* Breadcrumb */}
      <nav className="flex items-center gap-1 text-sm text-gray-500 dark:text-gray-400">
        <Link
          href="/episodes"
          className="hover:text-gray-900 dark:hover:text-gray-100"
        >
          {t("episodeDetail.episodes")}
        </Link>
        <ChevronRight className="h-4 w-4" />
        <span className="text-gray-900 dark:text-gray-100">{shortId}</span>
      </nav>

      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="flex items-start gap-4">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => router.push("/episodes")}
            className="mt-1"
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <div className="flex items-center gap-3">
              <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
                {t("episodeDetail.episode", { id: shortId })}
              </h1>
              <EpisodeStatusBadge status={episode.status} />
            </div>
            <p className="text-gray-500 dark:text-gray-400 mt-1 font-mono text-sm">
              {t("episodeDetail.id")}: {episode.id}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {episode.status === EPISODE_COLLECTION_STATUS.RECORDING &&
            isRobotInScope(episode.robot_id, scopeIds) && (
              <Button variant="default" size="sm" asChild>
                <Link href={`/robots/${episode.robot_id}/teleoperation`}>
                  <Radio className="mr-2 h-4 w-4" />
                  {t("episodeDetail.teleoperate")}
                </Link>
              </Button>
            )}
          {canEdit && (
            <EditEpisodeDialog episodeId={episode.id}>
              <Button variant="outline" size="sm">
                <Pencil className="mr-2 h-4 w-4" />
                {t("episodeDetail.edit")}
              </Button>
            </EditEpisodeDialog>
          )}
        </div>
      </div>

      {/* Info Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">
              {t("episodeDetail.task")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <Link
              href={`/tasks/${episode.task_id}`}
              className="text-lg font-semibold text-gray-900 dark:text-gray-100 hover:underline"
            >
              {taskName}
            </Link>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">
              {t("episodeDetail.taskVersion")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {taskVersionDisplay}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">
              {t("episodeDetail.robot")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {robotName}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">
              {t("episodeDetail.user")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {userName}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">
              {t("episodeDetail.location")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {locationName}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">
              {t("episodeDetail.recordedBy")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {recordedByName}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Time Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">
              {t("episodeDetail.startedAt")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {episode.started_at ? formatDateTime(episode.started_at) : "-"}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">
              {t("episodeDetail.endedAt")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {episode.ended_at ? formatDateTime(episode.ended_at) : "-"}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Parameter Values */}
      {episode.parameter_values &&
        Object.keys(episode.parameter_values).length > 0 && (
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">
                {t("episodeDetail.parameterValues")}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex flex-wrap gap-2">
                {Object.entries(episode.parameter_values).map(
                  ([key, value]) => (
                    <Badge key={key} variant="secondary" className="text-sm">
                      {key}: <strong className="ml-1">{value}</strong>
                    </Badge>
                  )
                )}
              </div>
            </CardContent>
          </Card>
        )}

      {/* Subtasks */}
      {episode.subtasks && episode.subtasks.length > 0 && (
        <SubtaskListCard
          subtasks={episode.subtasks}
          parameterValues={episode.parameter_values}
        />
      )}

      {/* Grades */}
      <EpisodeGradeSection
        episodeId={episode.id}
        averageGrade={episode.average_grade}
        gradeCount={episode.grade_count}
      />

      {/* Recording */}
      <RecordingTab
        episodeId={episode.id}
        startedAt={episode.started_at}
        subtasks={episode.subtasks ?? []}
      />

      {/* Statistics */}
      <StatsTab episodeId={episode.id} />

      {/* Error Details */}
      {episode.error_details && (
        <Card className="border-red-200 dark:border-red-800">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-red-600 dark:text-red-400">
              {t("episodeDetail.errorDetails")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <pre className="p-3 bg-red-50 dark:bg-red-950 rounded-md text-sm text-red-700 dark:text-red-300 overflow-x-auto whitespace-pre-wrap">
              <code>{episode.error_details}</code>
            </pre>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

function EpisodeDetailSkeleton() {
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
          <div className="h-9 bg-gray-200 dark:bg-gray-700 rounded w-20"></div>
          <div className="h-9 bg-gray-200 dark:bg-gray-700 rounded w-20"></div>
        </div>
      </div>

      {/* Info cards skeleton */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="h-24 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
        <div className="h-24 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
        <div className="h-24 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
        <div className="h-24 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
      </div>

      {/* Time cards skeleton */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="h-24 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
        <div className="h-24 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
      </div>
    </div>
  );
}
