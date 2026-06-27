"use client";

import { useEpisodeCollectionStatusLabel } from "@/lib/hooks/use-status-labels";
import {
  EPISODE_COLLECTION_STATUS,
  type EpisodeCollectionStatusValue,
} from "@/lib/status/constants";
import { EPISODE_COLLECTION_STATUS_DISPLAY } from "@/lib/status/display";

import { Badge } from "@/components/ui/badge";

interface EpisodeStatusBadgeProps {
  status?: number;
}

export function EpisodeStatusBadge({ status }: EpisodeStatusBadgeProps) {
  const getStatusLabel = useEpisodeCollectionStatusLabel();

  if (status == null) {
    return <span className="text-sm text-gray-500">-</span>;
  }

  const statusCode = status as EpisodeCollectionStatusValue;
  const display =
    EPISODE_COLLECTION_STATUS_DISPLAY[statusCode] ??
    EPISODE_COLLECTION_STATUS_DISPLAY[EPISODE_COLLECTION_STATUS.READY];

  return (
    <Badge variant="outline" className={display.className}>
      {getStatusLabel(statusCode)}
    </Badge>
  );
}
