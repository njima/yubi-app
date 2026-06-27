/**
 * Episodes Query Hooks
 * Provides TanStack Query hooks for fetching episodes data
 */

"use client";

import { useQuery } from "@tanstack/react-query";
import { z } from "zod";

import { fetchAndParse } from "@/lib/api/client-fetch";
import { schemas } from "@/lib/api/generated/api";
import { withQueryString } from "@/lib/api/query-string";
import type { EpisodeCollectionStatusValue } from "@/lib/status/constants";

import type { UseQueryOptions } from "@tanstack/react-query";

type Episode = z.infer<typeof schemas.Episode>;
type EpisodeListResponse = z.infer<typeof schemas.EpisodeListResponse>;

/**
 * Query parameters for listing episodes
 */
interface ListEpisodesParams {
  task_id?: string;
  task_version_id?: string;
  robot_id?: string;
  user_id?: string;
  status?: EpisodeCollectionStatusValue[];
  started_at_from?: string;
  started_at_to?: string;
  page?: number;
  limit?: number;
  sort_by?: string;
  sort_order?: string;
}

/**
 * Query keys for episodes
 * Hierarchical structure for efficient cache invalidation
 */
export const episodesQueryKeys = {
  all: ["episodes"] as const,
  lists: () => [...episodesQueryKeys.all, "list"] as const,
  list: (params: ListEpisodesParams = {}) =>
    [...episodesQueryKeys.lists(), params] as const,
  details: () => [...episodesQueryKeys.all, "detail"] as const,
  detail: (id: string) => [...episodesQueryKeys.details(), id] as const,
  next: () => [...episodesQueryKeys.all, "next"] as const,
};

/**
 * Hook to fetch episodes list
 * @param params - Query parameters (task_id, status)
 * @param options - TanStack Query options
 */
export function useEpisodesQuery(
  params: ListEpisodesParams = {},
  options?: Omit<
    UseQueryOptions<EpisodeListResponse, Error>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<EpisodeListResponse, Error>({
    queryKey: episodesQueryKeys.list(params),
    queryFn: async () => {
      return fetchAndParse(
        withQueryString("/web/api/episodes", params),
        schemas.EpisodeListResponse,
        "Failed to fetch episodes"
      );
    },
    ...options,
  });
}

/**
 * Hook to fetch a single episode by ID
 * @param episodeId - Episode ID
 * @param options - TanStack Query options
 */
export function useEpisodeQuery(
  episodeId: string,
  options?: Omit<UseQueryOptions<Episode, Error>, "queryKey" | "queryFn">
) {
  return useQuery<Episode, Error>({
    queryKey: episodesQueryKeys.detail(episodeId),
    queryFn: async () => {
      return fetchAndParse(
        `/web/api/episodes/${episodeId}`,
        schemas.Episode,
        "Failed to fetch episode"
      );
    },
    ...options,
  });
}

// TODO: Implement useNextEpisodeQuery when API endpoint is available
// type GetNextEpisodeResponse = Awaited<
//   ReturnType<typeof apiClient.getNextEpisode>
// >;

// /**
//  * Hook to fetch the next episode to process
//  * @param options - TanStack Query options
//  */
// export function useNextEpisodeQuery(
//   options?: Omit<
//     UseQueryOptions<GetNextEpisodeResponse, ZodiosError>,
//     "queryKey" | "queryFn"
//   >
// ) {
//   return useQuery<GetNextEpisodeResponse, ZodiosError>({
//     queryKey: episodesQueryKeys.next(),
//     queryFn: async () => {
//       const nextEpisode = await apiClient.getNextEpisode();
//       return nextEpisode;
//     },
//     ...options,
//   });
// }
