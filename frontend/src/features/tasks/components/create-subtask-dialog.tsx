"use client";

import { Plus } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { CreateSubTaskForm } from "./create-subtask-form";

interface CreateSubTaskDialogProps {
  taskId: string;
  taskVersionId: string;
}

export function CreateSubTaskDialog({
  taskId,
  taskVersionId,
}: CreateSubTaskDialogProps) {
  const { t } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button size="sm">
          <Plus className="mr-2 h-4 w-4" />
          {t("createSubtaskDialog.trigger")}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("createSubtaskDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("createSubtaskDialog.description")}
          </DialogDescription>
        </DialogHeader>
        <CreateSubTaskForm
          taskId={taskId}
          taskVersionId={taskVersionId}
          onSuccess={() => setIsOpen(false)}
          onCancel={() => setIsOpen(false)}
        />
      </DialogContent>
    </Dialog>
  );
}
