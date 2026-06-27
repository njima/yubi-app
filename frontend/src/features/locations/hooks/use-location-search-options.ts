"use client";

import { useCallback, useMemo } from "react";

import { useSearchState } from "@/lib/hooks/use-search-state";

import { useLocationsQuery } from "./use-locations-query";

import type { SearchableSelectOption } from "@/components/ui/searchable-select";

export function useLocationSearchOptions(params?: { site_id?: string }) {
  const { debouncedSearch, selectedLabel, setSelectedLabel, onSearch } =
    useSearchState();

  const { data, isLoading } = useLocationsQuery({
    search: debouncedSearch || undefined,
    limit: 20,
    site_id: params?.site_id,
  });

  const options: SearchableSelectOption[] = useMemo(
    () => data?.locations?.map((l) => ({ value: l.id, label: l.name })) ?? [],
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
