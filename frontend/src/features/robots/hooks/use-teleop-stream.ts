"use client";

import { useCallback, useState } from "react";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";
import { useManagedEventSource } from "@/lib/hooks/use-managed-event-source";

import type { RobotStatusStreamDetail } from "../schemas/robot";

type Episode = z.infer<typeof schemas.Episode>;

export interface TeleopTaskMeta {
  id: string;
  name: string;
  description?: string;
  manual_url: string;
  version: string;
}

interface UseTeleopStreamResult {
  status: RobotStatusStreamDetail | null;
  episode: Episode | null;
  task: TeleopTaskMeta | null;
  isConnected: boolean;
  error: string | null;
}

/**
 * Subscribes to the combined teleop SSE endpoint and exposes three
 * independent state slices populated from named events:
 *
 *   event: robot_status → status
 *   event: episode      → episode
 *   event: task         → task
 */
export function useTeleopStream(
  robotId: string,
  enabled: boolean = true
): UseTeleopStreamResult {
  const [status, setStatus] = useState<RobotStatusStreamDetail | null>(null);
  const [episode, setEpisode] = useState<Episode | null>(null);
  const [task, setTask] = useState<TeleopTaskMeta | null>(null);

  const url = `/web/api/robots/${robotId}/teleop/stream`;

  const { isConnected, error } = useManagedEventSource(
    url,
    enabled && !!robotId,
    useCallback((es: EventSource) => {
      es.addEventListener("robot_status", (e: MessageEvent) => {
        try {
          const parsed = JSON.parse(e.data);
          if (parsed && parsed.robot_type && parsed.status) {
            setStatus(parsed.status as RobotStatusStreamDetail);
          } else {
            setStatus(null);
          }
        } catch {
          // ignore malformed
        }
      });

      es.addEventListener("episode", (e: MessageEvent) => {
        try {
          if (e.data === "null") {
            setEpisode(null);
            return;
          }
          const parsed = JSON.parse(e.data);
          setEpisode(schemas.Episode.parse(parsed));
        } catch (err) {
          console.warn("[TeleopStream] failed to parse episode:", err);
        }
      });

      es.addEventListener("task", (e: MessageEvent) => {
        try {
          const parsed = JSON.parse(e.data) as {
            id: string;
            name: string;
            description?: string | null;
            manual_url: string;
            version: string;
          };
          setTask({
            id: parsed.id,
            name: parsed.name,
            description: parsed.description ?? undefined,
            manual_url: parsed.manual_url,
            version: parsed.version,
          });
        } catch (err) {
          console.warn("[TeleopStream] failed to parse task:", err);
        }
      });
    }, [])
  );

  return { status, episode, task, isConnected, error };
}
