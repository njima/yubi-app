import "server-only";

import { getActiveOrganizationId, getUserId } from "@/lib/auth/session";
import { clearActiveUser } from "@/lib/auth/switch-user";

const BACKEND_API_URL = process.env.BACKEND_API_URL || "http://backend:8000";

/**
 * Custom error class for Backend API errors
 */
export class BackendApiError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message);
    this.name = "BackendApiError";
  }
}

/**
 * Make a request to the Backend API and return the raw Response.
 * Use this for non-JSON responses such as CSV downloads.
 */
export async function fetchBackendRaw(
  path: string,
  options: RequestInit = {}
): Promise<Response> {
  const userId = await getUserId();
  const activeOrganizationId = await getActiveOrganizationId();
  const url = `${BACKEND_API_URL}${path}`;
  const authHeaders: Record<string, string> = {
    "X-User-ID": userId,
  };
  if (activeOrganizationId) {
    authHeaders["X-Organization-ID"] = activeOrganizationId;
  }

  const response = await fetch(url, {
    ...options,
    headers: {
      ...authHeaders,
      ...options.headers,
    },
  });

  if (!response.ok) {
    // If the user in the cookie was deleted, clear the cookie and redirect to force reload with DEFAULT_USER_ID
    if (response.status === 401) {
      await clearActiveUser();
      const { redirect } = await import("next/navigation");
      redirect("/web");
    }
    const error = await response.json().catch(() => ({}));
    throw new BackendApiError(
      response.status,
      error.message || response.statusText
    );
  }

  return response;
}

/**
 * Make a request to the Backend API
 */
export async function fetchBackend<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const userId = await getUserId();
  const activeOrganizationId = await getActiveOrganizationId();
  const url = `${BACKEND_API_URL}${path}`;
  const authHeaders: Record<string, string> = {
    "X-User-ID": userId,
  };
  if (activeOrganizationId) {
    authHeaders["X-Organization-ID"] = activeOrganizationId;
  }

  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...authHeaders,
      ...options.headers,
    },
  });

  if (!response.ok) {
    // If the user in the cookie was deleted, clear the cookie and redirect to force reload with DEFAULT_USER_ID
    if (response.status === 401) {
      await clearActiveUser();
      const { redirect } = await import("next/navigation");
      redirect("/web");
    }
    const error = await response.json().catch(() => ({}));
    throw new BackendApiError(
      response.status,
      error.message || response.statusText
    );
  }

  // Handle empty responses (204 No Content, DELETE, etc.)
  const text = await response.text();
  if (!text) {
    return undefined as T;
  }

  return JSON.parse(text);
}
