import { fetchBackend } from "./core";

import type {
  CollectionTrend,
  FleetSiteStats,
  FleetSiteSummary,
} from "./types";

// =============================================================================
// Fleet
// =============================================================================

export async function fetchFleetSummary(): Promise<FleetSiteSummary[]> {
  return fetchBackend<FleetSiteSummary[]>("/api/fleet/summary");
}

export async function fetchFleetStats(
  from: string,
  to: string
): Promise<FleetSiteStats[]> {
  const searchParams = new URLSearchParams();
  searchParams.append("from", from);
  searchParams.append("to", to);
  return fetchBackend<FleetSiteStats[]>(
    `/api/fleet/stats?${searchParams.toString()}`
  );
}

export async function fetchFleetCollectionTrend(
  granularity: string,
  from: string,
  to: string
): Promise<CollectionTrend> {
  const searchParams = new URLSearchParams();
  searchParams.append("granularity", granularity);
  searchParams.append("from", from);
  searchParams.append("to", to);
  return fetchBackend<CollectionTrend>(
    `/api/fleet/collection-trend?${searchParams.toString()}`
  );
}
