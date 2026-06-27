"use client";

import { useCallback, useState } from "react";

import { useManagedEventSource } from "@/lib/hooks/use-managed-event-source";

interface UseSSEStreamOptions<T> {
  url: string;
  enabled?: boolean;
  label: string;
  parse: (data: string) => T | null;
  shouldClose?: (value: T) => boolean;
}

interface UseSSEStreamResult<T> {
  data: T | null;
  isConnected: boolean;
  error: string | null;
}

export function useSSEStream<T>({
  url,
  enabled = true,
  label,
  parse,
  shouldClose,
}: UseSSEStreamOptions<T>): UseSSEStreamResult<T> {
  const [data, setData] = useState<T | null>(null);

  const { isConnected, error } = useManagedEventSource(
    url,
    enabled,
    useCallback(
      (es: EventSource, onTerminal: () => void) => {
        es.onmessage = (event: MessageEvent) => {
          try {
            const value = parse(event.data);
            setData(value);
            if (value !== null && shouldClose?.(value)) {
              onTerminal();
            }
          } catch (err) {
            console.warn(`[${label}] failed to parse message:`, err);
          }
        };
      },
      [parse, shouldClose, label]
    )
  );

  return { data, isConnected, error };
}
