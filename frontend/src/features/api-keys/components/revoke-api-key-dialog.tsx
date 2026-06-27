"use client";

import { Trash2 } from "lucide-react";
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

import { useRevokeApiKeyMutation } from "../hooks/use-revoke-api-key-mutation";

interface RevokeApiKeyDialogProps {
  apiKeyId: string;
  name: string;
  robotName: string | null | undefined;
}

export function RevokeApiKeyDialog({
  apiKeyId,
  name,
  robotName,
}: RevokeApiKeyDialogProps) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const { mutate, isPending } = useRevokeApiKeyMutation();

  const handleRevoke = () => {
    mutate({ apiKeyId }, { onSuccess: () => setOpen(false) });
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className="text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950"
        >
          <Trash2 className="h-4 w-4 mr-1" />
          {t("revokeApiKeyDialog.trigger")}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("revokeApiKeyDialog.title")}</DialogTitle>
          <DialogDescription>
            <span className="block">
              {t("revokeApiKeyDialog.description", { name })}
            </span>
            {robotName && (
              <span className="block">
                {t("revokeApiKeyDialog.robotLabel", { robot: robotName })}
              </span>
            )}
            <span className="mt-2 block text-red-600 dark:text-red-400">
              {t("revokeApiKeyDialog.warning")}
            </span>
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => setOpen(false)}
            disabled={isPending}
          >
            {t("revokeApiKeyDialog.cancel")}
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={handleRevoke}
            disabled={isPending}
          >
            {isPending
              ? t("revokeApiKeyDialog.submitting")
              : t("revokeApiKeyDialog.confirm")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
