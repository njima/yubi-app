"use client";

import { useQuery, type UseQueryOptions } from "@tanstack/react-query";
import { z } from "zod";

import type { CurrentUserResponse } from "@/lib/api/backend-client/types";
import { schemas } from "@/lib/api/generated/api";

type MeResponse = z.infer<typeof schemas.MeResponse>;

function currentUserResponse(me: MeResponse): CurrentUserResponse {
  return {
    ...me,
    role: me.active_role,
    organization_id: me.active_organization_id,
    organization_name: me.active_organization_name,
  };
}

export const meQueryKeys = {
  all: ["me"] as const,
  detail: () => [...meQueryKeys.all, "detail"] as const,
};

export function useMeQuery(
  options?: Omit<
    UseQueryOptions<CurrentUserResponse, Error>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<CurrentUserResponse, Error>({
    queryKey: meQueryKeys.detail(),
    queryFn: async () => {
      const response = await fetch("/web/api/me");
      if (!response.ok) {
        throw new Error(`Failed to fetch current user: ${response.statusText}`);
      }
      const data = await response.json();
      return currentUserResponse(schemas.MeResponse.parse(data));
    },
    ...options,
  });
}
