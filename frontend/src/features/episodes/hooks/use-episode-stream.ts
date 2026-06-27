"use client";

import { useCallback, useMemo } from "react";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";
import { EPISODE_COLLECTION_STATUS } from "@/lib/status/constants";

import { useSSEStream } from "./use-sse-stream";

type Episode = z.infer<typeof schemas.Episode>;

const TERMINAL_STATUSES: ReadonlySet<number> = new Set([
  EPISODE_COLLECTION_STATUS.COMPLETED,
  EPISODE_COLLECTION_STATUS.CANCEL,
]);

interface UseEpisodeStreamResult {
  data: Episode | null;
  isConnected: boolean;
  error: string | null;
}

export function useEpisodeStream(
  episodeId: string,
  enabled: boolean = true
): UseEpisodeStreamResult {
  const url = useMemo(
    () => `/web/api/episodes/${episodeId}/stream`,
    [episodeId]
  );

  const parse = useCallback((data: string): Episode | null => {
    const parsed = JSON.parse(data);
    return schemas.Episode.parse(parsed);
  }, []);

  const shouldClose = useCallback(
    (episode: Episode) => TERMINAL_STATUSES.has(episode.status),
    []
  );

  return useSSEStream<Episode>({
    url,
    enabled: enabled && !!episodeId,
    label: "EpisodeStream",
    parse,
    shouldClose,
  });
}
