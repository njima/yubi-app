import { fetchBackend } from "./core";

import type {
  Robot,
  RobotCreate,
  RobotListResponse,
  RobotUpdate,
} from "./types";

// =============================================================================
// Robots API
// =============================================================================

export interface RobotFilterParams {
  site_id?: string;
  location_id?: string;
  status?: number;
  robot_type?: string;
  page?: number;
  limit?: number;
  search?: string;
  sort_by?: string;
  sort_order?: string;
}

export async function fetchRobots(
  params?: RobotFilterParams
): Promise<RobotListResponse> {
  const searchParams = new URLSearchParams();
  if (params?.site_id) {
    searchParams.append("site_id", params.site_id);
  }
  if (params?.location_id) {
    searchParams.append("location_id", params.location_id);
  }
  if (params?.status !== undefined) {
    searchParams.append("status", String(params.status));
  }
  if (params?.robot_type) {
    searchParams.append("robot_type", params.robot_type);
  }
  if (params?.page !== undefined) {
    searchParams.append("page", String(params.page));
  }
  if (params?.limit !== undefined) {
    searchParams.append("limit", String(params.limit));
  }
  if (params?.search) {
    searchParams.append("search", params.search);
  }
  if (params?.sort_by) {
    searchParams.append("sort_by", params.sort_by);
  }
  if (params?.sort_order) {
    searchParams.append("sort_order", params.sort_order);
  }
  const query = searchParams.toString();
  return fetchBackend<RobotListResponse>(
    `/api/robots${query ? `?${query}` : ""}`
  );
}

export async function fetchRobot(robotId: string): Promise<Robot> {
  return fetchBackend<Robot>(`/api/robots/${robotId}`);
}

export async function createRobot(data: RobotCreate): Promise<Robot> {
  return fetchBackend<Robot>("/api/robots", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function updateRobot(
  robotId: string,
  data: RobotUpdate
): Promise<Robot> {
  return fetchBackend<Robot>(`/api/robots/${robotId}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function deleteRobot(robotId: string): Promise<void> {
  await fetchBackend<void>(`/api/robots/${robotId}`, {
    method: "DELETE",
  });
}

export async function fetchRobotTypes(params?: {
  site_id?: string;
  location_id?: string;
  status?: number;
}): Promise<string[]> {
  const searchParams = new URLSearchParams();
  if (params?.site_id) searchParams.append("site_id", params.site_id);
  if (params?.location_id)
    searchParams.append("location_id", params.location_id);
  if (params?.status !== undefined)
    searchParams.append("status", String(params.status));
  const query = searchParams.toString();
  return fetchBackend<string[]>(`/api/robot-types${query ? `?${query}` : ""}`);
}

// --- Robot Operator ---

export async function setRobotOperator(
  robotId: string,
  data: { user_id: string; display_name: string; organization_name: string }
): Promise<void> {
  await fetchBackend<void>(`/api/robots/${robotId}/operator`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
}

export async function clearRobotOperator(robotId: string): Promise<void> {
  await fetchBackend<void>(`/api/robots/${robotId}/operator`, {
    method: "DELETE",
  });
}
