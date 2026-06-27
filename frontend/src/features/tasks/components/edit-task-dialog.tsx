"use client";

import { useState } from "react";
import { useTranslation } from "react-i18next";

import { type TaskStatusValue } from "@/shared/lib/status-constants";
import {
  type TaskDifficultyValue,
  type TaskPriorityValue,
} from "@/shared/lib/status-constants";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { type TaskTag } from "../schemas";
import { EditTaskForm } from "./edit-task-form";

interface EditTaskDialogProps {
  taskId: string;
  name: string;
  description?: string;
  manual_url?: string;
  priority?: TaskPriorityValue;
  difficulty?: TaskDifficultyValue;
  status?: TaskStatusValue;
  deadline: string;
  robot_type?: string;
  tags?: TaskTag[];
  onSuccess?: () => void;
  children?: React.ReactNode;
}

export function EditTaskDialog({
  taskId,
  name,
  description,
  manual_url,
  priority,
  difficulty,
  status,
  deadline,
  robot_type,
  tags,
  onSuccess,
  children,
}: EditTaskDialogProps) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);

  const handleSuccess = () => {
    setOpen(false);
    onSuccess?.();
  };

  const handleCancel = () => {
    setOpen(false);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild={!!children}>{children}</DialogTrigger>
      <DialogContent className="max-w-lg max-h-[90vh] flex flex-col">
        <DialogHeader>
          <DialogTitle>{t("editTaskDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("editTaskDialog.description")}
          </DialogDescription>
        </DialogHeader>
        <div className="flex-1 overflow-y-auto px-2">
          <EditTaskForm
            taskId={taskId}
            defaultValues={{
              name,
              description,
              manual_url,
              priority,
              difficulty,
              status,
              deadline,
              robot_type,
            }}
            initialTags={tags}
            onSuccess={handleSuccess}
            onCancel={handleCancel}
          />
        </div>
      </DialogContent>
    </Dialog>
  );
}
