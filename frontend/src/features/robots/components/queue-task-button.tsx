"use client";

import { ListPlus } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { ROBOT_STATUS } from "@/shared/lib/status-constants";

import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

import { QueueTaskDialog } from "./queue-task-dialog";

import type { Robot } from "../schemas/robot";

interface QueueTaskButtonProps {
  robot: Robot;
}

export function QueueTaskButton({ robot }: QueueTaskButtonProps) {
  const { t } = useTranslation();
  const [dialogOpen, setDialogOpen] = useState(false);
  const isOnline = robot.status === ROBOT_STATUS.ONLINE;

  if (!isOnline) {
    return null;
  }

  return (
    <>
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setDialogOpen(true)}
            className="text-blue-600 hover:text-blue-700 hover:bg-blue-50 dark:text-blue-400 dark:hover:text-blue-300 dark:hover:bg-blue-950"
          >
            <ListPlus className="h-4 w-4" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>{t("dialog.queueTask")}</TooltipContent>
      </Tooltip>

      <QueueTaskDialog
        robot={robot}
        open={dialogOpen}
        onOpenChange={setDialogOpen}
      />
    </>
  );
}
