"use client";

import { useTranslation } from "react-i18next";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { StartTeleoperationForm } from "./start-teleoperation-form";

import type { Robot } from "../schemas/robot";

interface StartTeleoperationDialogProps {
  robot: Robot;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function StartTeleoperationDialog({
  robot,
  open,
  onOpenChange,
}: StartTeleoperationDialogProps) {
  const { t } = useTranslation();
  const handleSuccess = () => {
    onOpenChange(false);
  };

  const handleCancel = () => {
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{t("startTeleoperationDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("startTeleoperationDialog.description", {
              robotName: robot.name,
            })}
          </DialogDescription>
        </DialogHeader>
        <StartTeleoperationForm
          robot={robot}
          onSuccess={handleSuccess}
          onCancel={handleCancel}
        />
      </DialogContent>
    </Dialog>
  );
}
