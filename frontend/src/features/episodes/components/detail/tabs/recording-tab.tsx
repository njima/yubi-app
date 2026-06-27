"use client";

import {
  Fragment,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

import { useEpisodeRecordingsQuery } from "../../../hooks/use-episode-recordings-query";

type EpisodeSubTask = z.infer<typeof schemas.EpisodeSubTask>;

interface RecordingTabProps {
  episodeId: string;
  startedAt?: string | null;
  subtasks?: EpisodeSubTask[];
}

interface SubtaskGroup {
  id: string;
  subtaskName: string;
  executions: { id: string; startSec: number; endSec: number }[];
}

interface SelectedCell {
  groupIdx: number;
  execIdx: number;
}

function toElapsedSec(timestamp: string, base: string): number {
  return (new Date(timestamp).getTime() - new Date(base).getTime()) / 1000;
}

function buildSubtaskGroups(
  subtasks: EpisodeSubTask[],
  startedAt: string
): SubtaskGroup[] {
  const groups: SubtaskGroup[] = [];
  for (const subtask of subtasks) {
    if (!subtask.executions?.length) continue;
    const executions = subtask.executions
      .filter((exec) => exec.started_at)
      .map((exec) => ({
        id: exec.id,
        startSec: toElapsedSec(exec.started_at!, startedAt),
        endSec: exec.finished_at
          ? toElapsedSec(exec.finished_at, startedAt)
          : toElapsedSec(exec.started_at!, startedAt),
      }));
    if (executions.length > 0) {
      groups.push({ id: subtask.id, subtaskName: subtask.name, executions });
    }
  }
  return groups;
}

export function RecordingTab({
  episodeId,
  startedAt,
  subtasks,
}: RecordingTabProps) {
  const { t } = useTranslation();
  const { data, isLoading, error } = useEpisodeRecordingsQuery(episodeId);

  const videoRefs = useRef<Map<string, HTMLVideoElement>>(new Map());
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [selectedCell, setSelectedCell] = useState<SelectedCell | null>(null);

  const features = useMemo(
    () => (data ? Object.keys(data.recordings).sort() : []),
    [data]
  );
  const subtaskGroups = useMemo(
    () =>
      startedAt && subtasks ? buildSubtaskGroups(subtasks, startedAt) : [],
    [startedAt, subtasks]
  );

  const handleTimeUpdate = useCallback((e: Event) => {
    const video = e.target as HTMLVideoElement;
    setCurrentTime(video.currentTime);
  }, []);

  const handleLoadedMetadata = useCallback((e: Event) => {
    const video = e.target as HTMLVideoElement;
    setDuration((prev) => Math.max(prev, video.duration));
  }, []);

  const handleEnded = useCallback(() => {}, []);

  useEffect(() => {
    const firstKey = features[0];
    if (!firstKey) return;
    const video = videoRefs.current.get(firstKey);
    if (!video) return;
    video.addEventListener("timeupdate", handleTimeUpdate);
    video.addEventListener("loadedmetadata", handleLoadedMetadata);
    video.addEventListener("ended", handleEnded);
    return () => {
      video.removeEventListener("timeupdate", handleTimeUpdate);
      video.removeEventListener("loadedmetadata", handleLoadedMetadata);
      video.removeEventListener("ended", handleEnded);
    };
  }, [features, handleTimeUpdate, handleLoadedMetadata, handleEnded]);

  const syncAllVideos = useCallback((time: number) => {
    videoRefs.current.forEach((video) => {
      video.currentTime = time;
    });
  }, []);

  const handlePlay = useCallback(() => {
    videoRefs.current.forEach((video) => video.play());
  }, []);

  const handlePause = useCallback(() => {
    videoRefs.current.forEach((video) => video.pause());
  }, []);

  const handleScrubberChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const time = parseFloat(e.target.value);
      setCurrentTime(time);
      syncAllVideos(time);
    },
    [syncAllVideos]
  );

  const handleRowClick = useCallback(
    (startSec: number, groupIdx: number, execIdx: number) => {
      setSelectedCell({ groupIdx, execIdx });
      videoRefs.current.forEach((video) => {
        video.pause();
        video.currentTime = startSec;
      });
      setCurrentTime(startSec);
    },
    []
  );

  const hasRecordings = !isLoading && !error && data && features.length > 0;

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("episodeRecording.title")}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {isLoading && (
          <p className="text-muted-foreground text-sm">
            {t("episodeRecording.loading")}
          </p>
        )}
        {!isLoading && (error || !hasRecordings) && (
          <p className="text-muted-foreground text-sm">
            {error
              ? t("episodeRecording.loadFailed")
              : t("episodeRecording.empty")}
          </p>
        )}
        {hasRecordings && (
          <>
            <div className="flex gap-4 overflow-x-auto">
              {features.map((feature) => (
                <div
                  key={feature}
                  className="flex min-w-0 flex-1 flex-col gap-1"
                >
                  <span className="text-muted-foreground truncate text-xs">
                    {feature}
                  </span>
                  <video
                    ref={(el) => {
                      if (el) {
                        videoRefs.current.set(feature, el);
                      } else {
                        videoRefs.current.delete(feature);
                      }
                    }}
                    src={data!.recordings[feature]}
                    className="w-full rounded bg-black"
                    preload="metadata"
                    playsInline
                  />
                </div>
              ))}
            </div>

            <div className="space-y-1">
              <input
                type="range"
                min={0}
                max={duration || 1}
                step={0.05}
                value={currentTime}
                onChange={handleScrubberChange}
                className="w-full cursor-pointer accent-blue-500"
              />
              <div className="text-muted-foreground flex justify-between text-xs">
                <span>{currentTime.toFixed(1)} s</span>
                <span>{duration.toFixed(1)} s</span>
              </div>
            </div>

            <div className="flex gap-2">
              <Button variant="default" size="sm" onClick={handlePlay}>
                {t("episodeRecording.play")}
              </Button>
              <Button variant="outline" size="sm" onClick={handlePause}>
                {t("episodeRecording.pause")}
              </Button>
            </div>
          </>
        )}

        {subtaskGroups.length > 0 && (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t("episodeRecording.startSec")}</TableHead>
                <TableHead>{t("episodeRecording.endSec")}</TableHead>
                <TableHead>{t("episodeRecording.task")}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {subtaskGroups.map((group, groupIdx) => (
                <Fragment key={group.id}>
                  <TableRow>
                    <TableCell
                      colSpan={3}
                      className="bg-muted/40 py-1 text-xs font-semibold uppercase tracking-wider"
                    >
                      {group.subtaskName}
                    </TableCell>
                  </TableRow>
                  {group.executions.map((exec, execIdx) => {
                    const isSelected =
                      selectedCell?.groupIdx === groupIdx &&
                      selectedCell?.execIdx === execIdx;
                    return (
                      <TableRow
                        key={exec.id}
                        className={`cursor-pointer hover:bg-muted/50 ${isSelected ? "bg-muted" : ""}`}
                        onClick={() =>
                          handleRowClick(exec.startSec, groupIdx, execIdx)
                        }
                      >
                        <TableCell className="font-mono text-sm">
                          {exec.startSec.toFixed(2)}
                        </TableCell>
                        <TableCell className="font-mono text-sm">
                          {exec.endSec.toFixed(2)}
                        </TableCell>
                        <TableCell className="text-sm">
                          #{execIdx + 1} {group.subtaskName}
                        </TableCell>
                      </TableRow>
                    );
                  })}
                </Fragment>
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>
    </Card>
  );
}
