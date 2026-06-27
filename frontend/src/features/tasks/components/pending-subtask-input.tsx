"use client";

import { Plus, X } from "lucide-react";
import { useEffect, useRef } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";

export interface PendingSubTask {
  tempId: string;
  name: string;
  description: string;
  durationHours: number | undefined;
  durationMinutes: number | undefined;
}

interface PendingSubTaskInputProps {
  subtasks: PendingSubTask[];
  onAdd: () => void;
  onRemove: (tempId: string) => void;
  onUpdate: (
    tempId: string,
    field: "name" | "description",
    value: string
  ) => void;
  onUpdateDuration: (
    tempId: string,
    hours: number | undefined,
    minutes: number | undefined
  ) => void;
}

export function PendingSubTaskInput({
  subtasks,
  onAdd,
  onRemove,
  onUpdate,
  onUpdateDuration,
}: PendingSubTaskInputProps) {
  const { t } = useTranslation();
  const listRef = useRef<HTMLDivElement>(null);
  const prevCountRef = useRef(subtasks.length);

  useEffect(() => {
    if (subtasks.length > prevCountRef.current && listRef.current) {
      listRef.current.scrollTo({
        top: listRef.current.scrollHeight,
        behavior: "smooth",
      });
    }
    prevCountRef.current = subtasks.length;
  }, [subtasks.length]);

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <p className="text-sm font-medium text-gray-700 dark:text-gray-300">
          {t("pendingSubtaskInput.subtasksOptional")}
        </p>
        <Button type="button" variant="outline" size="sm" onClick={onAdd}>
          <Plus className="h-4 w-4 mr-1" />
          {t("createSubtaskDialog.trigger")}
        </Button>
      </div>

      {subtasks.length === 0 ? (
        <div className="rounded-lg border border-dashed border-gray-300 dark:border-gray-700 p-4 text-center">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t("pendingSubtaskInput.empty")}
          </p>
        </div>
      ) : (
        <div ref={listRef} className="space-y-3 max-h-64 overflow-y-auto">
          {subtasks.map((subtask, index) => (
            <div
              key={subtask.tempId}
              className="rounded-lg border border-gray-200 dark:border-gray-700 p-3 space-y-2"
            >
              <div className="flex items-center justify-between">
                <span className="text-xs font-medium text-gray-500 dark:text-gray-400">
                  {t("pendingSubtaskInput.subtaskNumber", {
                    index: index + 1,
                  })}
                </span>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => onRemove(subtask.tempId)}
                  className="h-6 w-6 p-0"
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>
              <Input
                placeholder={t("subtaskForm.namePlaceholder")}
                value={subtask.name}
                onChange={(e) =>
                  onUpdate(subtask.tempId, "name", e.target.value)
                }
              />
              <Textarea
                placeholder={t("subtaskForm.descriptionOptional")}
                value={subtask.description}
                onChange={(e) =>
                  onUpdate(subtask.tempId, "description", e.target.value)
                }
                className="resize-none"
                rows={2}
              />
              <div className="flex items-center gap-2">
                <span className="text-xs text-gray-500 dark:text-gray-400 shrink-0">
                  {t("pendingSubtaskInput.targetDuration")}
                </span>
                <Input
                  type="number"
                  min={0}
                  placeholder="0"
                  className="w-16 h-7 text-sm"
                  value={subtask.durationHours ?? ""}
                  onChange={(e) =>
                    onUpdateDuration(
                      subtask.tempId,
                      e.target.value === ""
                        ? undefined
                        : Number(e.target.value),
                      subtask.durationMinutes
                    )
                  }
                />
                <span className="text-xs text-gray-500">
                  {t("durationInput.hoursShort")}
                </span>
                <Input
                  type="number"
                  min={0}
                  max={59}
                  placeholder="0"
                  className="w-16 h-7 text-sm"
                  value={subtask.durationMinutes ?? ""}
                  onChange={(e) =>
                    onUpdateDuration(
                      subtask.tempId,
                      subtask.durationHours,
                      e.target.value === "" ? undefined : Number(e.target.value)
                    )
                  }
                />
                <span className="text-xs text-gray-500">
                  {t("durationInput.minutesShort")}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
