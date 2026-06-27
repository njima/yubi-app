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

import { EditEpisodeForm } from "./edit-episode-form";

interface EditEpisodeDialogProps {
  episodeId: string;
  onSuccess?: () => void;
  children?: React.ReactNode;
}

export function EditEpisodeDialog({
  episodeId,
  onSuccess,
  children,
}: EditEpisodeDialogProps) {
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
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("editEpisodeDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("editEpisodeDialog.description", {
              id: episodeId.substring(0, 8),
            })}
          </DialogDescription>
        </DialogHeader>
        <EditEpisodeForm
          episodeId={episodeId}
          onSuccess={handleSuccess}
          onCancel={handleCancel}
        />
      </DialogContent>
    </Dialog>
  );
}
