"use client";

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

import { useRobotQuery } from "@/features/robots/hooks/use-robots-query";

import { EditRobotForm } from "./edit-robot-form";

interface EditRobotDialogProps {
  robotId: string;
  onSuccess?: () => void;
  children?: React.ReactNode;
}

export function EditRobotDialog({
  robotId,
  onSuccess,
  children,
}: EditRobotDialogProps) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);

  // Fetch robot data when dialog opens
  const { data: robot, isLoading } = useRobotQuery(robotId, {
    enabled: open,
  });

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
      <DialogContent className="max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{t("editRobotDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("editRobotDialog.description")}
          </DialogDescription>
        </DialogHeader>
        {isLoading ? (
          <div className="p-8 text-center text-gray-600 dark:text-gray-400">
            {t("editRobotDialog.loadingData")}
          </div>
        ) : robot ? (
          <EditRobotForm
            robotId={robotId}
            defaultValues={robot}
            onSuccess={handleSuccess}
            onCancel={handleCancel}
          />
        ) : (
          <div className="p-8 text-center text-red-600 dark:text-red-400">
            {t("editRobotDialog.failedToLoad")}
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
