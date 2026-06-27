"use client";

import { useCallback, useMemo } from "react";

import { useSearchState } from "@/shared/hooks/use-search-state";

import { useTasksQuery } from "./use-tasks-query";

import type { SearchableSelectOption } from "@/components/ui/searchable-select";

export function useTaskSearchOptions(params?: {
  has_approved_version?: boolean;
  sort_by?: string;
  sort_order?: string;
  status?: number[];
}) {
  const { debouncedSearch, selectedLabel, setSelectedLabel, onSearch } =
    useSearchState();

  const { data, isLoading } = useTasksQuery({
    search: debouncedSearch || undefined,
    limit: 20,
    has_approved_version: params?.has_approved_version,
    sort_by: params?.sort_by,
    sort_order: params?.sort_order,
    status: params?.status,
  });

  const options: SearchableSelectOption[] = useMemo(
    () => data?.tasks?.map((t) => ({ value: t.id, label: t.name })) ?? [],
    [data]
  );

  const onValueChange = useCallback(
    (value: string) => {
      const label = options.find((o) => o.value === value)?.label;
      if (label) setSelectedLabel(label);
    },
    [options, setSelectedLabel]
  );

  return {
    options,
    isLoading,
    onSearch,
    selectedLabel,
    onValueChange,
    tasks: data?.tasks,
  };
}
