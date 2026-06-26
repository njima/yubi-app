import { fetchBackend } from "./core";

import type {
  ApiKeyCreateRequest,
  ApiKeyCreateResponse,
  ApiKeyListResponse,
  ApiKeyResponse,
} from "./types";

// =============================================================================
// API Keys API
// =============================================================================

export async function fetchApiKeys(params?: {
  page?: number;
  limit?: number;
  robot_id?: string;
  user_id?: string;
  include_revoked?: boolean;
}): Promise<ApiKeyListResponse> {
  const searchParams = new URLSearchParams();
  if (params?.page !== undefined) {
    searchParams.append("page", String(params.page));
  }
  if (params?.limit !== undefined) {
    searchParams.append("limit", String(params.limit));
  }
  if (params?.robot_id) {
    searchParams.append("robot_id", params.robot_id);
  }
  if (params?.user_id) {
    searchParams.append("user_id", params.user_id);
  }
  if (params?.include_revoked !== undefined) {
    searchParams.append("include_revoked", String(params.include_revoked));
  }
  const query = searchParams.toString();
  return fetchBackend<ApiKeyListResponse>(
    `/api/api-keys${query ? `?${query}` : ""}`
  );
}

export async function createApiKey(
  data: ApiKeyCreateRequest
): Promise<ApiKeyCreateResponse> {
  return fetchBackend<ApiKeyCreateResponse>("/api/api-keys", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function fetchApiKey(apiKeyId: string): Promise<ApiKeyResponse> {
  return fetchBackend<ApiKeyResponse>(`/api/api-keys/${apiKeyId}`);
}

// TODO: updateApiKey (PATCH /api/api-keys/{id}) — backend is implemented but frontend
// UI is not yet built. Use revoke + recreate as a workaround.

export async function revokeApiKey(apiKeyId: string): Promise<void> {
  await fetchBackend<void>(`/api/api-keys/${apiKeyId}`, {
    method: "DELETE",
  });
}
