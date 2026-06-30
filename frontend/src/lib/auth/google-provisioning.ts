import "server-only";

import { z } from "zod";

const BACKEND_API_URL = process.env.BACKEND_API_URL || "http://backend:8000";

const googleAuthSessionResponseSchema = z.object({
  user_id: z.string(),
  email: z.string(),
  display_name: z.string(),
  avatar_url: z.string().nullable().optional(),
  active_organization_id: z.string(),
  active_organization_name: z.string(),
  active_role: z.number(),
});

export type GoogleAuthSession = z.infer<typeof googleAuthSessionResponseSchema>;

export interface GoogleProfileInput {
  googleSub: string;
  email: string;
  name: string;
  avatarUrl?: string;
}

export async function provisionGoogleUserSession(
  input: GoogleProfileInput
): Promise<GoogleAuthSession> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (process.env.AUTH_INTERNAL_API_SECRET) {
    headers["X-Internal-Auth-Secret"] = process.env.AUTH_INTERNAL_API_SECRET;
  }

  const response = await fetch(`${BACKEND_API_URL}/api/auth/google/session`, {
    method: "POST",
    headers,
    body: JSON.stringify({
      google_sub: input.googleSub,
      email: input.email,
      name: input.name,
      avatar_url: input.avatarUrl ?? "",
    }),
  });

  if (!response.ok) {
    throw new Error(`Failed to provision Google user: ${response.status}`);
  }

  return googleAuthSessionResponseSchema.parse(await response.json());
}
