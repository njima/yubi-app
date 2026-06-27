import { useCallback, useState } from "react";

import { useDebounce } from "./use-debounce";

/**
 * Shared state management for server-side search dropdowns.
 * Handles search input debouncing and selected label tracking.
 */
export function useSearchState() {
  const [search, setSearch] = useState("");
  const [selectedLabel, setSelectedLabel] = useState<string>();
  const debouncedSearch = useDebounce(search, 300);

  const onSearch = useCallback((query: string) => {
    setSearch(query);
  }, []);

  return { debouncedSearch, selectedLabel, setSelectedLabel, onSearch };
}
