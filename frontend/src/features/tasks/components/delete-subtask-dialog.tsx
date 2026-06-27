"use client";

import { ReactNode, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { useDeleteSubTaskMutation } from "@/features/tasks/hooks/use-delete-subtask-mutation";

interface DeleteSubTaskDialogProps {
  subtaskId: string;
  taskId: string;
  name: string;
  children: ReactNode;
}

export function DeleteSubTaskDialog({
  subtaskId,
  taskId,
  name,
  children,
}: DeleteSubTaskDialogProps) {
  const { t } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);
  const { mutate, isPending } = useDeleteSubTaskMutation();

  const handleDelete = () => {
    mutate(
      { subtaskId, taskId },
      {
        onSuccess: () => {
          setIsOpen(false);
        },
      }
    );
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("deleteSubtaskDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("deleteSubtaskDialog.description")}
          </DialogDescription>
        </DialogHeader>

        <div className="rounded-md bg-gray-100 p-4 dark:bg-gray-800">
          <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
            {name}
          </p>
          <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
            ID: {subtaskId}
          </p>
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => setIsOpen(false)}
            disabled={isPending}
          >
            {t("dialog.cancel")}
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={isPending}
          >
            {isPending ? t("dialog.deleting") : t("deleteSubtaskDialog.button")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
