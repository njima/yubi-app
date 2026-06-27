"use client";

import { Check, Copy } from "lucide-react";
import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";

interface RawKeyDisplayDialogProps {
  rawKey: string | null;
  onClose: () => void;
}

export function RawKeyDisplayDialog({
  rawKey,
  onClose,
}: RawKeyDisplayDialogProps) {
  const { t } = useTranslation();
  const [confirmed, setConfirmed] = useState(false);
  const [copied, setCopied] = useState(false);
  const keyRef = useRef<HTMLInputElement | null>(null);

  const handleCopy = async () => {
    if (!rawKey) return;
    // Prefer the async Clipboard API. Some browsers/contexts (insecure http,
    // iframes without clipboard permission) reject it; in that case select the
    // input below so the operator can copy it manually with the keyboard.
    if (typeof navigator !== "undefined" && navigator.clipboard) {
      try {
        await navigator.clipboard.writeText(rawKey);
        setCopied(true);
        toast.success(t("rawKeyDisplayDialog.copied"));
        return;
      } catch {
        // fall through to the manual fallback
      }
    }
    keyRef.current?.focus();
    keyRef.current?.select();
    toast.error(t("rawKeyDisplayDialog.copyFailed"));
  };

  const handleOpenChange = (open: boolean) => {
    if (!open && confirmed) {
      // Reset before fully closing so a subsequent issue starts fresh.
      setConfirmed(false);
      setCopied(false);
      onClose();
    }
  };

  return (
    <Dialog open={rawKey !== null} onOpenChange={handleOpenChange}>
      <DialogContent
        className="sm:max-w-lg"
        onPointerDownOutside={(e) => e.preventDefault()}
        onEscapeKeyDown={(e) => e.preventDefault()}
      >
        <DialogHeader>
          <DialogTitle>{t("rawKeyDisplayDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("rawKeyDisplayDialog.warning")}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div className="flex items-center gap-2 rounded-md border bg-gray-50 dark:bg-gray-900 p-3">
            <input
              ref={keyRef}
              readOnly
              value={rawKey ?? ""}
              onFocus={(e) => e.currentTarget.select()}
              aria-label={t("rawKeyDisplayDialog.title")}
              className="flex-1 break-all bg-transparent font-mono text-sm outline-none"
            />
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={handleCopy}
              className="shrink-0 gap-1"
            >
              {copied ? (
                <Check className="h-4 w-4" />
              ) : (
                <Copy className="h-4 w-4" />
              )}
              {copied
                ? t("rawKeyDisplayDialog.copied")
                : t("rawKeyDisplayDialog.copyButton")}
            </Button>
          </div>

          <div className="flex items-center gap-2">
            <Checkbox
              id="api-key-saved"
              checked={confirmed}
              onCheckedChange={(value) => setConfirmed(value === true)}
            />
            <Label
              htmlFor="api-key-saved"
              className="text-sm font-normal cursor-pointer"
            >
              {t("rawKeyDisplayDialog.checkboxConfirm")}
            </Label>
          </div>
        </div>

        <DialogFooter>
          <Button
            type="button"
            disabled={!confirmed}
            onClick={() => {
              setConfirmed(false);
              setCopied(false);
              onClose();
            }}
          >
            {t("rawKeyDisplayDialog.close")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
