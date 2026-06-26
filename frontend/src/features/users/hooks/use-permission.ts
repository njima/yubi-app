"use client";

import { hasPermission, type PermissionAction } from "@/shared/lib/permissions";

import { useMeQuery } from "@/features/users";

export function usePermission(action: PermissionAction): boolean {
  const { data: me } = useMeQuery();
  return hasPermission(me?.role, action);
}
