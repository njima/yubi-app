"use client";

import { Check, ChevronDown, Loader2, Search } from "lucide-react";
import * as React from "react";
import { useTranslation } from "react-i18next";

import { cn } from "@/lib/utils";

import { Popover, PopoverContent, PopoverTrigger } from "./popover";

export interface SearchableSelectOption {
  value: string;
  label: string;
}

export interface SearchableSelectProps {
  value: string | undefined;
  onValueChange: (value: string) => void;
  options: SearchableSelectOption[];
  placeholder?: string;
  searchPlaceholder?: string;
  disabled?: boolean;
  className?: string;
  /** When provided, enables async search mode. Called on every search input change. */
  onSearch?: (query: string) => void;
  /** Shows a loading spinner in the dropdown list. */
  isLoading?: boolean;
  /** Override the display label for the selected value (useful in async mode when the selected item may not be in the current options). */
  selectedLabel?: string;
  /** Custom render function for each option in the dropdown list. Falls back to plain label text if not provided. */
  renderOption?: (option: SearchableSelectOption) => React.ReactNode;
  /** Custom render function for the selected value in the trigger button. Falls back to displayLabel text if not provided. */
  renderSelected?: (option: SearchableSelectOption) => React.ReactNode;
}

const SearchableSelect = React.forwardRef<
  HTMLButtonElement,
  SearchableSelectProps
>(
  (
    {
      value,
      onValueChange,
      options,
      placeholder,
      searchPlaceholder,
      disabled = false,
      className,
      onSearch,
      isLoading = false,
      selectedLabel: selectedLabelProp,
      renderOption,
      renderSelected,
    },
    ref
  ) => {
    const { t } = useTranslation();
    const [open, setOpen] = React.useState(false);
    const [search, setSearch] = React.useState("");
    const inputRef = React.useRef<HTMLInputElement>(null);
    const listId = React.useId();
    const resolvedPlaceholder = placeholder ?? t("common.selectPlaceholder");
    const resolvedSearchPlaceholder =
      searchPlaceholder ?? t("common.searchPlaceholder");

    // In async mode (onSearch provided), skip client-side filtering — parent controls options.
    const filtered = React.useMemo(() => {
      if (onSearch) return options;
      if (!search) return options;
      const lower = search.toLowerCase();
      return options.filter((o) => o.label.toLowerCase().includes(lower));
    }, [options, search, onSearch]);

    const selectedOption = options.find((o) => o.value === value);
    const displayLabel = selectedLabelProp ?? selectedOption?.label;

    const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const query = e.target.value;
      setSearch(query);
      onSearch?.(query);
    };

    const handleSelect = (optionValue: string) => {
      onValueChange(optionValue);
      setOpen(false);
      setSearch("");
      onSearch?.("");
    };

    return (
      <Popover
        modal
        open={open}
        onOpenChange={(nextOpen) => {
          setOpen(nextOpen);
          if (!nextOpen) {
            setSearch("");
            onSearch?.("");
          }
        }}
      >
        <PopoverTrigger asChild>
          <button
            ref={ref}
            type="button"
            role="combobox"
            aria-expanded={open}
            aria-controls={listId}
            disabled={disabled}
            className={cn(
              "flex h-10 min-w-0 w-full items-center justify-between rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white focus:outline-none focus:ring-2 focus:ring-gray-950 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 dark:border-gray-800 dark:bg-gray-950 dark:ring-offset-gray-950 dark:focus:ring-gray-300",
              !displayLabel && "text-gray-500 dark:text-gray-400",
              className
            )}
          >
            <span className="min-w-0 flex-1 truncate">
              {renderSelected && selectedOption
                ? renderSelected(selectedOption)
                : displayLabel || resolvedPlaceholder}
            </span>
            <ChevronDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </button>
        </PopoverTrigger>
        <PopoverContent
          className="w-[var(--radix-popover-trigger-width)] p-0"
          onOpenAutoFocus={(e) => {
            e.preventDefault();
            inputRef.current?.focus();
          }}
        >
          <div className="flex items-center border-b border-gray-200 px-3 dark:border-gray-800">
            <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
            <input
              ref={inputRef}
              value={search}
              onChange={handleSearchChange}
              placeholder={resolvedSearchPlaceholder}
              className="flex h-9 w-full bg-transparent py-2 text-sm outline-none placeholder:text-gray-500 dark:placeholder:text-gray-400"
            />
          </div>
          <div
            id={listId}
            role="listbox"
            className="max-h-60 overflow-y-auto p-1"
            onWheel={(e) => e.stopPropagation()}
          >
            {isLoading ? (
              <div className="flex items-center justify-center py-4">
                <Loader2 className="h-4 w-4 animate-spin text-gray-400" />
              </div>
            ) : filtered.length === 0 ? (
              <p className="px-3 py-2 text-sm text-gray-500 dark:text-gray-400 text-center">
                {t("common.noResultsFound")}
              </p>
            ) : (
              filtered.map((option) => (
                <button
                  key={option.value}
                  type="button"
                  onClick={() => handleSelect(option.value)}
                  className={cn(
                    "relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none hover:bg-gray-100 dark:hover:bg-gray-800",
                    value === option.value && "font-medium"
                  )}
                >
                  <span className="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
                    {value === option.value && <Check className="h-4 w-4" />}
                  </span>
                  {renderOption ? renderOption(option) : option.label}
                </button>
              ))
            )}
          </div>
        </PopoverContent>
      </Popover>
    );
  }
);
SearchableSelect.displayName = "SearchableSelect";

export { SearchableSelect };
