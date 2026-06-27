import { useEffect, useState } from "react";

/**
 * Debounces a value by the given delay in milliseconds.
 * Returns the debounced value that only updates after the delay has passed
 * since the last change to the input value.
 */
export function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedValue(value), delay);
    return () => clearTimeout(timer);
  }, [value, delay]);

  return debouncedValue;
}
