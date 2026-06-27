"use client";

import { AlertCircle, Video } from "lucide-react";
import { useTranslation } from "react-i18next";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { ForceTorqueCard } from "../../force-torque-card";
import { JointTemperaturesCard } from "../../joint-temperatures-card";
import { RobotCameraViewer } from "../../robot-camera-viewer";
import { RobotStatusCard } from "../../robot-status-card";

import type { Robot, RobotStatusStreamDetail } from "../../../schemas/robot";

interface LiveViewTabProps {
  robot: Robot;
  realtimeStatus?: RobotStatusStreamDetail | null;
}

export function LiveViewTab({ robot, realtimeStatus }: LiveViewTabProps) {
  const { t } = useTranslation();
  const isOperating = !!robot.active_user_id;

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
      {/* Left side - Live Streams (2/3 width on large screens) */}
      <div className="lg:col-span-2">
        <Card className="h-full">
          <CardHeader className="pb-2">
            <CardTitle className="text-base font-medium flex items-center gap-2">
              <Video className="h-5 w-5" />
              {t("robotLiveView.liveStreams")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            {isOperating ? (
              <RobotCameraViewer
                robotConfig={robot.robot_config ?? undefined}
                robotName={robot.name}
              />
            ) : (
              <div className="flex flex-col items-center justify-center py-16 text-center">
                <AlertCircle className="h-12 w-12 text-gray-400 dark:text-gray-500 mb-4" />
                <p className="text-gray-600 dark:text-gray-400 font-medium">
                  {t("robotLiveView.robotNotOperating")}
                </p>
                <p className="text-gray-500 dark:text-gray-500 text-sm mt-1">
                  {t("robotLiveView.streamsAvailable")}
                </p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Right side - Status cards (1/3 width on large screens) */}
      <div className="space-y-4">
        <RobotStatusCard robot={robot} realtimeStatus={realtimeStatus} />
        <ForceTorqueCard />
        <JointTemperaturesCard />
      </div>
    </div>
  );
}
