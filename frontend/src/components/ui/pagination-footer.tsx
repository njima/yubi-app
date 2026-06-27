"use client";

import { ChevronLeft, ChevronRight } from "lucide-react";
import { useTranslation } from "react-i18next";

import { PAGE_SIZE_OPTIONS } from "@/lib/pagination";

import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface PaginationFooterProps {
  page: number;
  totalPages: number;
  totalCount: number;
  onPageChange: (page: number) => void;
  itemLabel?: string | null;
  limit?: number;
  onLimitChange?: (limit: number) => void;
}

export function PaginationFooter({
  page,
  totalPages,
  totalCount,
  onPageChange,
  itemLabel,
  limit,
  onLimitChange,
}: PaginationFooterProps) {
  const { t } = useTranslation();
  const displayLabel = itemLabel ?? t("pagination.items");
  return (
    <div className="border-t px-6 py-4 flex items-center justify-between text-sm text-gray-600 dark:text-gray-400 dark:border-gray-700">
      <span>
        {t("pagination.total", { count: totalCount, label: displayLabel })}
      </span>
      <div className="flex items-center gap-4">
        {limit != null && onLimitChange && (
          <div className="flex items-center gap-2">
            <span className="text-sm">{t("pagination.rows")}</span>
            <Select
              value={String(limit)}
              onValueChange={(value) => onLimitChange(Number(value))}
            >
              <SelectTrigger className="h-8 w-[70px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {PAGE_SIZE_OPTIONS.map((option) => (
                  <SelectItem key={option} value={String(option)}>
                    {option}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )}
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onPageChange(page - 1)}
            disabled={page <= 1}
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <span>{t("pagination.page", { page, totalPages })}</span>
          <Button
            variant="outline"
            size="sm"
            onClick={() => onPageChange(page + 1)}
            disabled={page >= totalPages}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
