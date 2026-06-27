"use client";

import { useCallback, useState } from "react";

import { useManagedEventSource } from "@/lib/hooks/use-managed-event-source";

import type {
  RobotStatusStreamDetail,
  RobotStatusStreamResponse,
} from "../schemas/robot";

interface UseRobotStatusStreamResult {
  data: RobotStatusStreamDetail | null;
  isConnected: boolean;
  error: string | null;
}

export function useRobotStatusStream(
  robotId: string,
  enabled: boolean = true
): UseRobotStatusStreamResult {
  const [data, setData] = useState<RobotStatusStreamDetail | null>(null);

  const url = `/web/api/robots/${robotId}/status/stream`;

  const { isConnected, error } = useManagedEventSource(
    url,
    enabled && !!robotId,
    useCallback((es: EventSource) => {
      es.onmessage = (event: MessageEvent) => {
        try {
          const parsed = JSON.parse(event.data) as RobotStatusStreamResponse;
          // When robot_type is empty, Redis has no status data — treat as null
          if (parsed.robot_type && parsed.status) {
            setData(parsed.status as RobotStatusStreamDetail);
          } else {
            setData(null);
          }
        } catch {
          // Ignore malformed messages
        }
      };
    }, [])
  );

  return { data, isConnected, error };
}
