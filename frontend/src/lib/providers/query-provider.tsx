"use client";

/**
 * TanStack Query Provider
 * Global configuration for React Query
 */

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

import type { ReactNode } from "react";

/**
 * Create QueryClient with default options (2025 best practices)
 */
function makeQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        // Data is considered fresh for 1 minute
        staleTime: 60 * 1000,
        // Cached data is garbage collected after 5 minutes
        gcTime: 5 * 60 * 1000,
        // Retry failed requests once
        retry: 1,
        // Don't refetch on window focus (can be enabled per-query if needed)
        refetchOnWindowFocus: false,
        // Refetch on network reconnect
        refetchOnReconnect: true,
      },
      mutations: {
        // Don't retry mutations
        retry: 0,
      },
    },
  });
}

let browserQueryClient: QueryClient | undefined = undefined;

function getQueryClient() {
  if (typeof window === "undefined") {
    // Server: always create a new query client
    return makeQueryClient();
  } else {
    // Browser: reuse existing client if available
    if (!browserQueryClient) browserQueryClient = makeQueryClient();
    return browserQueryClient;
  }
}

interface QueryProviderProps {
  children: ReactNode;
}

export function QueryProvider({ children }: QueryProviderProps) {
  const queryClient = getQueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}
