"use client";

import { ReactNode, useState } from "react";
import { useTranslation } from "react-i18next";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { EditSubTaskForm } from "./edit-subtask-form";

interface EditSubTaskDialogProps {
  subtaskId: string;
  taskId: string;
  name: string;
  description?: string;
  target_duration_seconds?: number | null;
  children: ReactNode;
}

export function EditSubTaskDialog({
  subtaskId,
  taskId,
  name,
  description,
  target_duration_seconds,
  children,
}: EditSubTaskDialogProps) {
  const { t } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("editSubtaskDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("editSubtaskDialog.description")}
          </DialogDescription>
        </DialogHeader>
        <EditSubTaskForm
          subtaskId={subtaskId}
          taskId={taskId}
          defaultValues={{ name, description, target_duration_seconds }}
          onSuccess={() => setIsOpen(false)}
          onCancel={() => setIsOpen(false)}
        />
      </DialogContent>
    </Dialog>
  );
}
