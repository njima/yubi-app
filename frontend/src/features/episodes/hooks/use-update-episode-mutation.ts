"use client";

/**
 * Episode Update Mutation Hook
 * TanStack Query mutation for updating episodes
 */

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { episodesQueryKeys } from "./use-episodes-query";

type EpisodeUpdate = z.infer<typeof schemas.EpisodeUpdate>;
type Episode = z.infer<typeof schemas.Episode>;

/**
 * Hook to update an existing episode
 *
 * @example
 * ```tsx
 * const { mutate, isPending } = useUpdateEpisodeMutation();
 *
 * mutate(
 *   { episodeId: 1, data: { title: "Updated Title" } },
 *   {
 *     onSuccess: () => console.log("Episode updated!"),
 *     onError: (error) => console.error(error),
 *   }
 * );
 * ```
 */
export function useUpdateEpisodeMutation() {
  const queryClient = useQueryClient();

  return useMutation<
    Episode,
    Error,
    { episodeId: string; data: EpisodeUpdate }
  >({
    mutationFn: async ({ episodeId, data }) => {
      const response = await fetch(`/web/api/episodes/${episodeId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      if (!response.ok) {
        throw new Error(`Failed to update episode: ${response.statusText}`);
      }
      return response.json();
    },
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({
        queryKey: episodesQueryKeys.lists(),
      });
      queryClient.invalidateQueries({
        queryKey: episodesQueryKeys.detail(variables.episodeId),
      });
      toast.success("Episode updated successfully");
    },
    onError: (error) => {
      toast.error("Failed to update episode", {
        description: error.message || "An unexpected error occurred",
      });
    },
  });
}
