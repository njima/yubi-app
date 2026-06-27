"use client";

import { Check, ChevronDown, Plus, Search, X } from "lucide-react";
import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";

import { Badge } from "@/components/ui/badge";
import { Label } from "@/components/ui/label";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { SearchableSelect } from "@/components/ui/searchable-select";

import { useCreateTaskTagMutation } from "../hooks/use-create-task-tag-mutation";
import {
  useTaskCategoryTypesQuery,
  useTaskTagsQuery,
} from "../hooks/use-task-tags-query";
import { type TaskTag } from "../schemas";

interface TagSelectorProps {
  selectedTags: TaskTag[];
  onChange: (tags: TaskTag[]) => void;
}

export function TagSelector({ selectedTags, onChange }: TagSelectorProps) {
  const { t } = useTranslation();
  const [selectedCategoryId, setSelectedCategoryId] = useState<string>("");
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  const { data: categoryTypes = [] } = useTaskCategoryTypesQuery();
  const { data: allTags = [] } = useTaskTagsQuery();
  const { mutateAsync: createTag, isPending: isCreating } =
    useCreateTaskTagMutation();

  const selectedTagIds = new Set(selectedTags.map((t) => t.id));

  const tagsInCategory = allTags.filter(
    (t) =>
      t.category_type_id === selectedCategoryId &&
      t.name.toLowerCase().includes(search.toLowerCase())
  );

  const exactMatch = tagsInCategory.some(
    (t) => t.name.toLowerCase() === search.trim().toLowerCase()
  );

  const handleSelectCategory = (categoryId: string) => {
    setSelectedCategoryId(categoryId);
    setSearch("");
    setOpen(true);
  };

  const handleToggle = (tag: TaskTag) => {
    if (selectedTagIds.has(tag.id)) {
      onChange(selectedTags.filter((t) => t.id !== tag.id));
    } else {
      onChange([...selectedTags, tag]);
    }
  };

  const handleRemove = (tagId: string, e: React.MouseEvent) => {
    e.stopPropagation();
    onChange(selectedTags.filter((t) => t.id !== tagId));
  };

  const handleCreate = async () => {
    if (!search.trim() || !selectedCategoryId) return;
    const tag = await createTag({
      name: search.trim(),
      category_type_id: selectedCategoryId,
    });
    onChange([...selectedTags, tag]);
    setSearch("");
  };

  return (
    <div className="space-y-3">
      <Label>{t("tagSelector.tags")}</Label>

      {/* Selected tags */}
      {selectedTags.length > 0 && (
        <div className="flex flex-wrap gap-1.5">
          {selectedTags.map((tag) => (
            <Badge
              key={tag.id}
              variant="secondary"
              className="gap-1 pr-1 shrink-0"
            >
              <span className="text-xs text-muted-foreground">
                {tag.category_type_name}:
              </span>
              {tag.name}
              <button
                type="button"
                onClick={(e) => handleRemove(tag.id, e)}
                className="ml-0.5 rounded-full hover:bg-gray-300 dark:hover:bg-gray-600 p-0.5"
              >
                <X className="h-3 w-3" />
              </button>
            </Badge>
          ))}
        </div>
      )}

      <div className="flex gap-2">
        {/* Category dropdown */}
        <SearchableSelect
          value={selectedCategoryId}
          onValueChange={handleSelectCategory}
          options={categoryTypes.map((ct) => ({
            value: ct.id,
            label: ct.name,
          }))}
          placeholder={t("tagSelector.selectCategory")}
          className="w-44"
        />

        {/* Tag search & dropdown */}
        {selectedCategoryId && (
          <Popover
            open={open}
            onOpenChange={(nextOpen) => {
              setOpen(nextOpen);
              if (!nextOpen) setSearch("");
            }}
          >
            <PopoverTrigger asChild>
              <button
                type="button"
                onClick={() => setOpen(true)}
                className="h-10 w-full flex-1 cursor-text rounded-md border border-input bg-background px-3 flex items-center gap-1.5 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
              >
                <span className="flex-1 text-sm text-left text-muted-foreground">
                  {t("tagSelector.searchTags")}
                </span>
                <ChevronDown className="h-4 w-4 text-muted-foreground shrink-0" />
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
                  onChange={(e) => setSearch(e.target.value)}
                  placeholder={t("tagSelector.searchTags")}
                  className="flex h-9 w-full bg-transparent py-2 text-sm outline-none placeholder:text-gray-500 dark:placeholder:text-gray-400"
                />
              </div>
              <div className="max-h-52 overflow-y-auto py-1">
                {tagsInCategory.length === 0 && !search.trim() && (
                  <p className="px-3 py-2 text-sm text-muted-foreground text-center">
                    {t("tagSelector.noTagsFound")}
                  </p>
                )}

                {tagsInCategory.map((tag) => {
                  const selected = selectedTagIds.has(tag.id);
                  return (
                    <button
                      key={tag.id}
                      type="button"
                      onClick={() => handleToggle(tag)}
                      className="w-full flex items-center gap-2 px-3 py-1.5 text-sm hover:bg-accent text-left"
                    >
                      <span
                        className={`h-4 w-4 shrink-0 flex items-center justify-center rounded border transition-colors ${
                          selected
                            ? "bg-primary border-primary"
                            : "border-input"
                        }`}
                      >
                        {selected && (
                          <Check className="h-3 w-3 text-primary-foreground" />
                        )}
                      </span>
                      {tag.name}
                    </button>
                  );
                })}

                {search.trim() && !exactMatch && (
                  <button
                    type="button"
                    disabled={isCreating}
                    onClick={handleCreate}
                    className="w-full flex items-center gap-2 px-3 py-1.5 text-sm hover:bg-accent text-left border-t border-input disabled:opacity-50"
                  >
                    <Plus className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
                    <span className="text-muted-foreground">
                      {t("tagSelector.createTag", { name: search.trim() })}
                    </span>
                  </button>
                )}
              </div>
            </PopoverContent>
          </Popover>
        )}
      </div>
    </div>
  );
}
