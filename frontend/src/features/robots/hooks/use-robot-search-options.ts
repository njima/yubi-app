"use client";

import { useCallback, useMemo } from "react";

import { useSearchState } from "@/lib/hooks/use-search-state";

import { useRobotsQuery } from "./use-robots-query";

import type { SearchableSelectOption } from "@/components/ui/searchable-select";

export function useRobotSearchOptions({
  enabled = true,
}: { enabled?: boolean } = {}) {
  const { debouncedSearch, selectedLabel, setSelectedLabel, onSearch } =
    useSearchState();

  const { data, isLoading } = useRobotsQuery(
    {
      search: debouncedSearch || undefined,
      limit: 20,
    },
    { enabled }
  );

  const options: SearchableSelectOption[] = useMemo(
    () => data?.robots?.map((r) => ({ value: r.id, label: r.name })) ?? [],
    [data]
  );

  const onValueChange = useCallback(
    (value: string) => {
      const label = options.find((o) => o.value === value)?.label;
      if (label) setSelectedLabel(label);
    },
    [options, setSelectedLabel]
  );

  return { options, isLoading, onSearch, selectedLabel, onValueChange };
}
