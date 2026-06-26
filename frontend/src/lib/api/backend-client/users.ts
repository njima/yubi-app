import { z } from "zod";

import { fetchBackend } from "./core";
import { schemas } from "../generated/api";

import type {
  MeUpdateRequest,
  OrganizationResponse,
  UserCreateRequest,
  UserListResponse,
  UserResponse,
  UserRoleUpdateRequest,
  UserUpdateRequest,
} from "./types";

// =============================================================================
// User Import API
// =============================================================================

export type UserImportValidationResponse = z.infer<
  typeof schemas.UserImportValidationResponse
>;
export type UserImportValidRow = z.infer<typeof schemas.UserImportValidRow>;
export type UserImportRowError = z.infer<typeof schemas.UserImportRowError>;
export type UserImportResponse = z.infer<typeof schemas.UserImportResponse>;

export async function validateUserImport(
  csvContent: string
): Promise<UserImportValidationResponse> {
  return fetchBackend<UserImportValidationResponse>(
    "/api/users/import/validate",
    {
      method: "POST",
      body: JSON.stringify({ csv_content: csvContent }),
    }
  );
}

export async function importUsers(
  csvContent: string
): Promise<UserImportResponse> {
  return fetchBackend<UserImportResponse>("/api/users/import", {
    method: "POST",
    body: JSON.stringify({ csv_content: csvContent }),
  });
}

// =============================================================================
// Users API
// =============================================================================

export async function fetchUsers(
  params: {
    page?: number;
    limit?: number;
    organization_id?: string;
    location_id?: string;
    site_id?: string;
    search?: string;
    sort_by?: string;
    sort_order?: string;
  } = {}
): Promise<UserListResponse> {
  const query = new URLSearchParams();
  if (params.page !== undefined) query.set("page", String(params.page));
  if (params.limit !== undefined) query.set("limit", String(params.limit));
  if (params.organization_id)
    query.set("organization_id", params.organization_id);
  if (params.location_id) query.set("location_id", params.location_id);
  if (params.site_id) query.set("site_id", params.site_id);
  if (params.search) query.set("search", params.search);
  if (params.sort_by) query.set("sort_by", params.sort_by);
  if (params.sort_order) query.set("sort_order", params.sort_order);
  const qs = query.toString() ? `?${query.toString()}` : "";
  return fetchBackend<UserListResponse>(`/api/users${qs}`);
}

export async function updateUser(
  userId: string,
  data: UserUpdateRequest
): Promise<UserResponse> {
  return fetchBackend<UserResponse>(`/api/users/${userId}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function fetchUser(userId: string): Promise<UserResponse> {
  return fetchBackend<UserResponse>(`/api/users/${userId}`);
}

export async function fetchMe(): Promise<UserResponse> {
  return fetchBackend<UserResponse>("/api/me");
}

export async function updateMe(data: MeUpdateRequest): Promise<UserResponse> {
  return fetchBackend<UserResponse>("/api/me", {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function updateUserRole(
  userId: string,
  data: UserRoleUpdateRequest
): Promise<UserResponse> {
  return fetchBackend<UserResponse>(`/api/users/${userId}/role`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function updateUserLocations(
  userId: string,
  locationIds: string[]
): Promise<UserResponse> {
  return fetchBackend<UserResponse>(`/api/users/${userId}/locations`, {
    method: "PUT",
    body: JSON.stringify({ location_ids: locationIds }),
  });
}

// =============================================================================
// Organizations API
// =============================================================================

export async function fetchOrganization(
  organizationId: string
): Promise<OrganizationResponse> {
  return fetchBackend<OrganizationResponse>(
    `/api/organizations/${organizationId}`
  );
}

export async function createUser(
  data: UserCreateRequest
): Promise<UserResponse> {
  return fetchBackend<UserResponse>("/api/users", {
    method: "POST",
    body: JSON.stringify(data),
  });
}
