import { fetchBackend } from "./core";

import type {
  Location,
  LocationCreate,
  LocationListResponse,
  LocationUpdate,
  SiteListResponse,
} from "./types";

// =============================================================================
// Locations API
// =============================================================================

export async function fetchLocations(params?: {
  page?: number;
  limit?: number;
  search?: string;
  site_id?: string;
  sort_by?: string;
  sort_order?: string;
}): Promise<LocationListResponse> {
  const searchParams = new URLSearchParams();
  if (params?.page !== undefined) {
    searchParams.append("page", String(params.page));
  }
  if (params?.limit !== undefined) {
    searchParams.append("limit", String(params.limit));
  }
  if (params?.search) {
    searchParams.append("search", params.search);
  }
  if (params?.site_id) {
    searchParams.append("site_id", params.site_id);
  }
  if (params?.sort_by) {
    searchParams.append("sort_by", params.sort_by);
  }
  if (params?.sort_order) {
    searchParams.append("sort_order", params.sort_order);
  }
  const query = searchParams.toString();
  return fetchBackend<LocationListResponse>(
    `/api/locations${query ? `?${query}` : ""}`
  );
}

export async function fetchLocation(locationId: string): Promise<Location> {
  return fetchBackend<Location>(`/api/locations/${locationId}`);
}

export async function fetchSites(params?: {
  page?: number;
  limit?: number;
  search?: string;
  organization_id?: string;
}): Promise<SiteListResponse> {
  const searchParams = new URLSearchParams();
  if (params?.page !== undefined) {
    searchParams.append("page", String(params.page));
  }
  if (params?.limit !== undefined) {
    searchParams.append("limit", String(params.limit));
  }
  if (params?.search) {
    searchParams.append("search", params.search);
  }
  if (params?.organization_id) {
    searchParams.append("organization_id", params.organization_id);
  }
  const query = searchParams.toString();
  return fetchBackend<SiteListResponse>(
    `/api/sites${query ? `?${query}` : ""}`
  );
}

export async function createLocation(data: LocationCreate): Promise<Location> {
  return fetchBackend<Location>("/api/locations", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function updateLocation(
  locationId: string,
  data: LocationUpdate
): Promise<Location> {
  return fetchBackend<Location>(`/api/locations/${locationId}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function deleteLocation(locationId: string): Promise<void> {
  await fetchBackend<void>(`/api/locations/${locationId}`, {
    method: "DELETE",
  });
}
