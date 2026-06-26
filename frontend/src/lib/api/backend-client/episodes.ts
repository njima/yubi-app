import { BackendApiError, fetchBackend } from "./core";

import type {
  Episode,
  EpisodeBulkCreateRequest,
  EpisodeCreate,
  EpisodeGrade,
  EpisodeGradeListResponse,
  EpisodeGradeUpdate,
  EpisodeListResponse,
  EpisodeRecordingsResponse,
  EpisodeStatsResponse,
  EpisodeUpdate,
} from "./types";

// =============================================================================
// Episodes API
// =============================================================================

export interface EpisodeFilterParams {
  task_id?: string;
  task_version_id?: string;
  robot_id?: string;
  user_id?: string;
  status?: number[];
  started_at_from?: string;
  started_at_to?: string;
  page?: number;
  limit?: number;
  sort_by?: string;
  sort_order?: string;
}

export async function fetchEpisodes(
  params: EpisodeFilterParams = {}
): Promise<EpisodeListResponse> {
  const searchParams = new URLSearchParams();
  if (params.task_id) searchParams.set("task_id", params.task_id);
  if (params.task_version_id)
    searchParams.set("task_version_id", params.task_version_id);
  if (params.robot_id) searchParams.set("robot_id", params.robot_id);
  if (params.user_id) searchParams.set("user_id", params.user_id);
  params.status?.forEach((s) => searchParams.append("status", String(s)));
  if (params.started_at_from)
    searchParams.set("started_at_from", params.started_at_from);
  if (params.started_at_to)
    searchParams.set("started_at_to", params.started_at_to);
  if (params.page !== undefined) searchParams.set("page", String(params.page));
  if (params.limit !== undefined)
    searchParams.set("limit", String(params.limit));
  if (params.sort_by) searchParams.set("sort_by", params.sort_by);
  if (params.sort_order) searchParams.set("sort_order", params.sort_order);
  const query = searchParams.toString();
  return fetchBackend<EpisodeListResponse>(
    `/api/episodes${query ? `?${query}` : ""}`
  );
}

export async function fetchEpisode(episodeId: string): Promise<Episode> {
  return fetchBackend<Episode>(`/api/episodes/${episodeId}`);
}

export async function createEpisode(data: EpisodeCreate): Promise<Episode> {
  return fetchBackend<Episode>("/api/episodes", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function createEpisodesBulk(
  data: EpisodeBulkCreateRequest
): Promise<Episode[]> {
  return fetchBackend<Episode[]>("/api/episodes/bulk", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function updateEpisode(
  episodeId: string,
  data: EpisodeUpdate
): Promise<Episode> {
  return fetchBackend<Episode>(`/api/episodes/${episodeId}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

// fetchMyEpisodeGrade returns null on 404 ("not graded yet"); other errors throw.
export async function fetchMyEpisodeGrade(
  episodeId: string
): Promise<EpisodeGrade | null> {
  try {
    return await fetchBackend<EpisodeGrade>(
      `/api/episodes/${episodeId}/grades/me`
    );
  } catch (error) {
    if (error instanceof BackendApiError && error.status === 404) {
      return null;
    }
    throw error;
  }
}

export async function updateMyEpisodeGrade(
  episodeId: string,
  data: EpisodeGradeUpdate
): Promise<EpisodeGrade> {
  return fetchBackend<EpisodeGrade>(`/api/episodes/${episodeId}/grades/me`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function listEpisodeGrades(
  episodeId: string,
  params: { page?: number; limit?: number } = {}
): Promise<EpisodeGradeListResponse> {
  const searchParams = new URLSearchParams();
  if (params.page !== undefined) searchParams.set("page", String(params.page));
  if (params.limit !== undefined)
    searchParams.set("limit", String(params.limit));
  const query = searchParams.toString();
  return fetchBackend<EpisodeGradeListResponse>(
    `/api/episodes/${episodeId}/grades${query ? `?${query}` : ""}`
  );
}

export async function fetchEpisodeRecordings(
  episodeId: string
): Promise<EpisodeRecordingsResponse> {
  return fetchBackend<EpisodeRecordingsResponse>(
    `/api/episodes/${episodeId}/recordings`
  );
}

export async function fetchEpisodeStats(
  episodeId: string
): Promise<EpisodeStatsResponse> {
  return fetchBackend<EpisodeStatsResponse>(`/api/episodes/${episodeId}/stats`);
}
