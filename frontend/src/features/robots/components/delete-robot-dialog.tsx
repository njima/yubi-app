"use client";

import { useState } from "react";
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

import { useDeleteRobotMutation } from "@/features/robots/hooks/use-delete-robot-mutation";

interface DeleteRobotDialogProps {
  robotId: string;
  name: string;
  onSuccess?: () => void;
  children?: React.ReactNode;
}

export function DeleteRobotDialog({
  robotId,
  name,
  onSuccess,
  children,
}: DeleteRobotDialogProps) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const { mutate, isPending } = useDeleteRobotMutation();

  const handleDelete = () => {
    mutate(
      { robotId },
      {
        onSuccess: () => {
          setOpen(false);
          onSuccess?.();
        },
      }
    );
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild={!!children}>
        {children || <Button variant="destructive">Delete</Button>}
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("deleteRobotDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("deleteRobotDialog.description", { name })}
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => setOpen(false)}
            disabled={isPending}
          >
            {t("dialog.cancel")}
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={handleDelete}
            disabled={isPending}
          >
            {isPending ? t("dialog.deleting") : t("dialog.delete")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
