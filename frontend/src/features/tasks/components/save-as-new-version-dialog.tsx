"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { SaveAsNewVersionForm } from "./save-as-new-version-form";

import type { TaskVersion } from "../schemas";

interface SaveAsNewVersionDialogProps {
  taskId: string;
  versions: TaskVersion[];
  defaultBaseVersionId: string;
  children?: React.ReactNode;
}

export function SaveAsNewVersionDialog({
  taskId,
  versions,
  defaultBaseVersionId,
  children,
}: SaveAsNewVersionDialogProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const [open, setOpen] = useState(false);

  const handleSuccess = (newVersionId: string) => {
    setOpen(false);
    router.push(`/tasks/${taskId}?version=${newVersionId}`);
  };

  const handleCancel = () => {
    setOpen(false);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild={!!children}>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("saveAsNewVersionDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("saveAsNewVersionDialog.description")}
          </DialogDescription>
        </DialogHeader>
        <SaveAsNewVersionForm
          taskId={taskId}
          versions={versions}
          defaultBaseVersionId={defaultBaseVersionId}
          onSuccess={handleSuccess}
          onCancel={handleCancel}
        />
      </DialogContent>
    </Dialog>
  );
}
