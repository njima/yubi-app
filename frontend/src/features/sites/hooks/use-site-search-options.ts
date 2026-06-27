"use client";

import { useCallback, useMemo } from "react";

import { useSearchState } from "@/lib/hooks/use-search-state";

import { useSitesQuery } from "./use-sites-query";

import type { SearchableSelectOption } from "@/components/ui/searchable-select";

export function useSiteSearchOptions() {
  const { debouncedSearch, selectedLabel, setSelectedLabel, onSearch } =
    useSearchState();

  const { data, isLoading } = useSitesQuery({
    search: debouncedSearch || undefined,
    limit: 20,
  });

  const options: SearchableSelectOption[] = useMemo(
    () => data?.sites?.map((s) => ({ value: s.id, label: s.name })) ?? [],
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
