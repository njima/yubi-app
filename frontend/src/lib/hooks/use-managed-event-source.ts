"use client";

import { useEffect, useState } from "react";

const INITIAL_RETRY_DELAY_MS = 1000;
const MAX_RETRY_DELAY_MS = 30000;
const BACKOFF_MULTIPLIER = 2;
const HIDDEN_GRACE_MS = 30000;

interface UseManagedEventSourceResult {
  isConnected: boolean;
  error: string | null;
}

/**
 * Manages an EventSource lifecycle with:
 *   - Deferred connect when mounted in a hidden tab.
 *   - 30 s grace before closing when the tab goes hidden.
 *   - Exponential backoff on connection errors.
 *   - Cleanup on unmount and on url/enabled change.
 *
 * The consumer supplies `setupListeners(es, onTerminal)` — a callback
 * that attaches whatever event handlers it needs to the fresh
 * EventSource. `onTerminal` is an effect-local function that tears
 * down the connection; call it from a message handler when a
 * terminal state is detected (e.g. shouldClose). Consumers must wrap
 * setupListeners in `useCallback` so the effect dep is stable.
 */
export function useManagedEventSource(
  url: string,
  enabled: boolean,
  setupListeners: (es: EventSource, onTerminal: () => void) => void
): UseManagedEventSourceResult {
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!enabled) {
      return;
    }

    let eventSource: EventSource | null = null;
    let retryTimeout: ReturnType<typeof setTimeout> | null = null;
    let hideTimeout: ReturnType<typeof setTimeout> | null = null;
    const retryDelay = { current: INITIAL_RETRY_DELAY_MS };
    let isCleaned = false;

    function isVisible() {
      return typeof document === "undefined" || !document.hidden;
    }

    function closeConnection() {
      if (eventSource) {
        eventSource.close();
        eventSource = null;
      }
      if (retryTimeout) {
        clearTimeout(retryTimeout);
        retryTimeout = null;
      }
      setIsConnected(false);
    }

    function cleanup() {
      isCleaned = true;
      closeConnection();
      if (hideTimeout) {
        clearTimeout(hideTimeout);
        hideTimeout = null;
      }
    }

    function connect() {
      if (isCleaned) return;
      if (!isVisible()) return;

      if (eventSource) {
        eventSource.close();
      }

      eventSource = new EventSource(url);

      eventSource.onopen = () => {
        setIsConnected(true);
        setError(null);
        retryDelay.current = INITIAL_RETRY_DELAY_MS;
      };

      setupListeners(eventSource, cleanup);

      eventSource.onerror = () => {
        if (isCleaned) return;
        setIsConnected(false);
        setError("Connection lost, retrying...");
        if (eventSource) {
          eventSource.close();
          eventSource = null;
        }

        const delay = retryDelay.current;
        retryDelay.current = Math.min(
          delay * BACKOFF_MULTIPLIER,
          MAX_RETRY_DELAY_MS
        );

        retryTimeout = setTimeout(() => {
          connect();
        }, delay);
      };
    }

    function handleVisibilityChange() {
      if (isCleaned) return;
      if (isVisible()) {
        if (hideTimeout) {
          clearTimeout(hideTimeout);
          hideTimeout = null;
        }
        if (!eventSource && !retryTimeout) {
          retryDelay.current = INITIAL_RETRY_DELAY_MS;
          connect();
        }
      } else {
        if (hideTimeout) clearTimeout(hideTimeout);
        hideTimeout = setTimeout(() => {
          if (isCleaned) return;
          closeConnection();
        }, HIDDEN_GRACE_MS);
      }
    }

    if (typeof document !== "undefined") {
      document.addEventListener("visibilitychange", handleVisibilityChange);
    }

    connect();

    return () => {
      if (typeof document !== "undefined") {
        document.removeEventListener(
          "visibilitychange",
          handleVisibilityChange
        );
      }
      cleanup();
      setIsConnected(false);
      setError(null);
    };
  }, [url, enabled, setupListeners]);

  return { isConnected, error };
}
