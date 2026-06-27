"use client";

import { useCallback, useMemo } from "react";

import { useSearchState } from "@/shared/hooks/use-search-state";

import { useUsersQuery } from "./use-users-query";

import type { SearchableSelectOption } from "@/components/ui/searchable-select";

export function useUserSearchOptions() {
  const { debouncedSearch, selectedLabel, setSelectedLabel, onSearch } =
    useSearchState();

  const { data, isLoading } = useUsersQuery({
    search: debouncedSearch || undefined,
    limit: 20,
  });

  const options: SearchableSelectOption[] = useMemo(
    () =>
      data?.users?.map((u) => ({
        value: u.user_id,
        label: u.display_name,
      })) ?? [],
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
