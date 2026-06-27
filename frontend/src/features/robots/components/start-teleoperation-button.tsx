"use client";

import { Play } from "lucide-react";
import Link from "next/link";
import { useTranslation } from "react-i18next";

import { ROBOT_STATUS } from "@/shared/lib/status-constants";

import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

import type { Robot } from "../schemas/robot";

interface StartTeleoperationButtonProps {
  robot: Robot;
  /** Hide when robot is out of scope (default: true = always visible) */
  inScope?: boolean;
}

export function StartTeleoperationButton({
  robot,
  inScope = true,
}: StartTeleoperationButtonProps) {
  const { t } = useTranslation();
  const isOnline = robot.status === ROBOT_STATUS.ONLINE;

  if (!isOnline || !inScope) {
    return null;
  }

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Link href={`/robots/${robot.id}/teleoperation`}>
          <Button
            variant="ghost"
            size="sm"
            className="text-green-600 hover:text-green-700 hover:bg-green-50 dark:text-green-400 dark:hover:text-green-300 dark:hover:bg-green-950"
          >
            <Play className="h-4 w-4" />
          </Button>
        </Link>
      </TooltipTrigger>
      <TooltipContent>{t("episodeDetail.teleoperate")}</TooltipContent>
    </Tooltip>
  );
}
