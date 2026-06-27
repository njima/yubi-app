"use client";

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

import { CreateTaskForm } from "./create-task-form";

export function CreateTaskDialog() {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);

  const handleSuccess = () => {
    setOpen(false);
  };

  const handleCancel = () => {
    setOpen(false);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>{t("createTaskDialog.trigger")}</Button>
      </DialogTrigger>
      <DialogContent className="max-w-lg max-h-[90vh] flex flex-col">
        <DialogHeader>
          <DialogTitle>{t("createTaskDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("createTaskDialog.description")}
          </DialogDescription>
        </DialogHeader>
        <div className="flex-1 overflow-y-auto px-2">
          <CreateTaskForm onSuccess={handleSuccess} onCancel={handleCancel} />
        </div>
      </DialogContent>
    </Dialog>
  );
}
