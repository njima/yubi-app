"use client";

import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { GripVertical, Pencil, Trash2 } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";

import { useSubTasksByVersionQuery } from "@/features/tasks/hooks/use-subtasks-by-version-query";
import { useSubTasksQuery } from "@/features/tasks/hooks/use-subtasks-query";
import { secondsToHoursMinutes } from "@/features/tasks/lib/duration";

import { CreateSubTaskDialog } from "./create-subtask-dialog";
import { DeleteSubTaskDialog } from "./delete-subtask-dialog";
import { EditSubTaskDialog } from "./edit-subtask-dialog";
import { useReorderSubTasksMutation } from "../hooks/use-reorder-subtasks-mutation";

interface SubTaskItem {
  id: string;
  name: string;
  description?: string | null;
  target_duration_seconds?: number | null;
}

interface SortableSubTaskItemProps {
  subtask: SubTaskItem;
  taskId: string;
  isReadOnly: boolean;
  isDragDisabled: boolean;
}

function SortableSubTaskItem({
  subtask,
  taskId,
  isReadOnly,
  isDragDisabled,
}: SortableSubTaskItemProps) {
  const { t } = useTranslation();
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: subtask.id, disabled: isDragDisabled });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`flex items-center justify-between rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4 ${
        isDragging ? "opacity-50 shadow-lg" : ""
      }`}
    >
      <div className="flex items-center gap-3 flex-1">
        {!isReadOnly && (
          <button
            className="cursor-grab touch-none text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            {...attributes}
            {...listeners}
          >
            <GripVertical className="h-5 w-5" />
          </button>
        )}
        <div>
          <p className="font-medium text-gray-900 dark:text-gray-100">
            {subtask.name}
          </p>
          <div className="flex items-center gap-4 mt-1">
            <p className="text-xs text-gray-500 dark:text-gray-400">
              ID: {subtask.id}
            </p>
            {subtask.target_duration_seconds != null &&
              (() => {
                const { hours, minutes } = secondsToHoursMinutes(
                  subtask.target_duration_seconds
                );
                return (
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    {t("subtaskList.collectionTimeTarget", {
                      hours,
                      minutes,
                    })}
                  </p>
                );
              })()}
          </div>
        </div>
      </div>

      {!isReadOnly && (
        <div className="flex gap-2">
          <EditSubTaskDialog
            subtaskId={subtask.id}
            taskId={taskId}
            name={subtask.name}
            description={subtask.description ?? undefined}
            target_duration_seconds={subtask.target_duration_seconds}
          >
            <Button variant="ghost" size="sm">
              <Pencil className="h-4 w-4" />
            </Button>
          </EditSubTaskDialog>

          <DeleteSubTaskDialog
            subtaskId={subtask.id}
            taskId={taskId}
            name={subtask.name}
          >
            <Button variant="ghost" size="sm">
              <Trash2 className="h-4 w-4 text-red-600 dark:text-red-400" />
            </Button>
          </DeleteSubTaskDialog>
        </div>
      )}
    </div>
  );
}

interface SubTaskListProps {
  taskId: string;
  taskVersionId?: string;
  isReadOnly?: boolean;
}

export function SubTaskList({
  taskId,
  taskVersionId,
  isReadOnly = false,
}: SubTaskListProps) {
  const { t } = useTranslation();
  const versionQuery = useSubTasksByVersionQuery(taskVersionId);
  const taskQuery = useSubTasksQuery(taskVersionId ? undefined : taskId);

  const {
    data: subtasks,
    isLoading,
    error,
  } = taskVersionId ? versionQuery : taskQuery;

  const [localOrder, setLocalOrder] = useState<SubTaskItem[] | null>(null);
  const reorderMutation = useReorderSubTasksMutation();

  // Use local order for optimistic updates, fall back to server data
  const orderedSubtasks = localOrder ?? subtasks ?? [];

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (!over || active.id === over.id || !taskVersionId) return;

    const oldIndex = orderedSubtasks.findIndex((s) => s.id === active.id);
    const newIndex = orderedSubtasks.findIndex((s) => s.id === over.id);

    const newOrder = arrayMove(orderedSubtasks, oldIndex, newIndex);
    setLocalOrder(newOrder);

    reorderMutation.mutate(
      {
        taskVersionId,
        subtaskIds: newOrder.map((s) => s.id),
        taskId,
      },
      {
        onSettled: () => setLocalOrder(null),
      }
    );
  }

  if (isLoading) {
    return (
      <div className="p-4 text-center text-gray-600 dark:text-gray-400">
        {t("subtaskList.loading")}
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4 text-center text-red-600 dark:text-red-400">
        {t("subtaskList.error", { message: error.message })}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
          {t("subtaskList.title", { count: orderedSubtasks.length })}
        </h3>
        {!isReadOnly && taskVersionId && (
          <CreateSubTaskDialog taskId={taskId} taskVersionId={taskVersionId} />
        )}
      </div>

      {orderedSubtasks.length === 0 ? (
        <div className="rounded-lg border border-dashed border-gray-300 dark:border-gray-700 p-8 text-center">
          <p className="text-gray-600 dark:text-gray-400">
            {isReadOnly
              ? t("subtaskList.emptyVersion")
              : t("subtaskList.empty")}
          </p>
        </div>
      ) : (
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragEnd={handleDragEnd}
        >
          <SortableContext
            items={orderedSubtasks.map((s) => s.id)}
            strategy={verticalListSortingStrategy}
          >
            <div className="space-y-2">
              {orderedSubtasks.map((subtask) => (
                <SortableSubTaskItem
                  key={subtask.id}
                  subtask={subtask}
                  taskId={taskId}
                  isReadOnly={isReadOnly}
                  isDragDisabled={isReadOnly || !taskVersionId}
                />
              ))}
            </div>
          </SortableContext>
        </DndContext>
      )}
    </div>
  );
}
