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

import { useDeleteLocationMutation } from "../hooks/use-delete-location-mutation";

interface DeleteLocationDialogProps {
  locationId: string;
  name: string;
  children?: React.ReactNode;
}

export function DeleteLocationDialog({
  locationId,
  name,
  children,
}: DeleteLocationDialogProps) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const { mutate, isPending } = useDeleteLocationMutation();

  const handleDelete = () => {
    mutate(
      { locationId },
      {
        onSuccess: () => {
          setOpen(false);
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
          <DialogTitle>{t("deleteLocationDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("deleteLocationDialog.description", { name })}
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
