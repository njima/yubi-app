"use client";

import { OctagonX, RotateCcw } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

interface EStopButtonProps {
  robotName?: string;
  onEStopTriggered?: () => void;
  onEStopReleased?: () => void;
}

export function EStopButton({
  robotName,
  onEStopTriggered,
  onEStopReleased,
}: EStopButtonProps) {
  const { t } = useTranslation();
  const [isActive, setIsActive] = useState(false);
  const [dialogOpen, setDialogOpen] = useState(false);

  const handleTriggerEStop = () => {
    setIsActive(true);
    setDialogOpen(false);
    toast.error(t("eStop.activated"), {
      description: robotName
        ? t("eStop.robotStoppedNamed", { name: robotName })
        : t("eStop.robotStopped"),
      duration: 5000,
    });
    onEStopTriggered?.();
  };

  const handleReleaseEStop = () => {
    setIsActive(false);
    toast.success(t("eStop.released"), {
      description: t("eStop.robotOperable"),
      duration: 3000,
    });
    onEStopReleased?.();
  };

  if (isActive) {
    return (
      <Dialog>
        <DialogTrigger asChild>
          <Button
            variant="destructive"
            size="sm"
            className="bg-red-600 hover:bg-red-700 animate-pulse"
          >
            <OctagonX className="h-4 w-4 mr-1" />
            {t("eStop.active")}
          </Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("eStop.releaseTitle")}</DialogTitle>
            <DialogDescription>
              {t("eStop.releaseDescription")}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline">{t("dialog.cancel")}</Button>
            </DialogClose>
            <Button onClick={handleReleaseEStop}>
              <RotateCcw className="h-4 w-4 mr-1" />
              {t("eStop.releaseButton")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
      <DialogTrigger asChild>
        <Button variant="destructive" size="sm">
          <OctagonX className="h-4 w-4 mr-1" />
          {t("eStop.trigger")}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="text-red-600">
            {t("eStop.triggerTitle")}
          </DialogTitle>
          <DialogDescription>
            {robotName
              ? t("eStop.triggerDescriptionNamed", { name: robotName })
              : t("eStop.triggerDescription")}
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline">{t("dialog.cancel")}</Button>
          </DialogClose>
          <Button variant="destructive" onClick={handleTriggerEStop}>
            <OctagonX className="h-4 w-4 mr-1" />
            {t("eStop.triggerButton")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
