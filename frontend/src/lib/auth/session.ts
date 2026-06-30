import "server-only";
import { cookies } from "next/headers";

import { auth } from "@/auth";

export interface UserSession {
  username: string;
  email: string;
  userId: string;
  activeOrganizationId?: string;
  activeOrganizationName?: string;
}

const COOKIE_NAME = "active_user_id";
export const ACTIVE_ORG_COOKIE_NAME = "active_organization_id";

/**
 * Get the active user ID from cookie, falling back to the DEFAULT_USER_ID env var.
 */
export async function getUserId(): Promise<string> {
  const session = await auth();
  if (session?.user?.id) {
    return session.user.id;
  }

  const cookieStore = await cookies();
  const cookieUserId = cookieStore.get(COOKIE_NAME)?.value;
  if (cookieUserId) {
    return cookieUserId;
  }

  const envUserId = process.env.DEFAULT_USER_ID;
  if (!envUserId) {
    throw new Error("DEFAULT_USER_ID environment variable is required");
  }
  return envUserId;
}

export async function getActiveOrganizationId(): Promise<string | undefined> {
  const session = await auth();
  if (session?.user?.activeOrganizationId) {
    return session.user.activeOrganizationId;
  }

  const cookieStore = await cookies();
  return cookieStore.get(ACTIVE_ORG_COOKIE_NAME)?.value;
}

export async function getUserSession(): Promise<UserSession | null> {
  const session = await auth();
  if (session?.user?.id) {
    return {
      username: session.user.name ?? "User",
      email: session.user.email ?? "",
      userId: session.user.id,
      activeOrganizationId: session.user.activeOrganizationId,
      activeOrganizationName: session.user.activeOrganizationName,
    };
  }

  let userId: string;
  try {
    userId = await getUserId();
  } catch {
    return null;
  }
  return {
    username: "User",
    email: "",
    userId,
  };
}

export async function requireAuth(): Promise<UserSession> {
  const userSession = await getUserSession();

  if (!userSession) {
    const { redirect } = await import("next/navigation");
    return redirect("/login");
  }

  return userSession;
}
