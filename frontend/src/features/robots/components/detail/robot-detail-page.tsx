"use client";

import { ArrowLeft, ChevronRight } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useTranslation } from "react-i18next";

import { USER_ROLE } from "@/shared/lib/status-constants";

import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { TooltipProvider } from "@/components/ui/tooltip";

import { useMeQuery } from "@/features/users";

import { QueueTaskButton } from "../queue-task-button";
import { StartTeleoperationButton } from "../start-teleoperation-button";
import { HealthTab } from "./tabs/health-tab";
import { LiveViewTab } from "./tabs/live-view-tab";
import { LogsTab } from "./tabs/logs-tab";
import { OverviewTab } from "./tabs/overview-tab";
import { SensorsTab } from "./tabs/sensors-tab";
import { useRobotScope, isRobotInScope } from "../../hooks/use-robot-scope";
import { useRobotStatusStream } from "../../hooks/use-robot-status-stream";
import { useRobotQuery } from "../../hooks/use-robots-query";

interface RobotDetailPageProps {
  robotId: string;
}

export function RobotDetailPage({ robotId }: RobotDetailPageProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const { data: robot, isLoading, error } = useRobotQuery(robotId);
  const { data: realtimeStatus, isConnected } = useRobotStatusStream(
    robotId,
    !isLoading && !!robot
  );
  const { scopeIds } = useRobotScope();
  const { data: me } = useMeQuery();

  if (isLoading) {
    return <RobotDetailSkeleton />;
  }

  if (error) {
    return (
      <div className="p-8 text-center">
        <div className="text-red-600 dark:text-red-400 mb-4">
          {t("robotDetail.errorLoadingRobot", { message: error.message })}
        </div>
        <Button onClick={() => window.location.reload()}>
          {t("taskDetail.retry")}
        </Button>
      </div>
    );
  }

  if (!robot) {
    return (
      <div className="p-8 text-center">
        <div className="text-gray-600 dark:text-gray-400 mb-4">
          {t("robotDetail.robotNotFound")}
        </div>
        <Button onClick={() => router.push("/robots")}>
          {t("robotDetail.backToList")}
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Breadcrumb */}
      <nav className="flex items-center gap-1 text-sm text-gray-500 dark:text-gray-400">
        <Link
          href="/robots"
          className="hover:text-gray-900 dark:hover:text-gray-100"
        >
          {t("robotDetail.robots")}
        </Link>
        <ChevronRight className="h-4 w-4" />
        <span className="text-gray-900 dark:text-gray-100">{robot.name}</span>
      </nav>

      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="flex items-start gap-4">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => router.push("/robots")}
            className="mt-1"
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
              {robot.name}
            </h1>
            <p className="text-gray-500 dark:text-gray-400">
              {t("robotDetail.robotType", {
                value: robot.robot_type || t("status.unknown"),
              })}
            </p>
            <p className="text-gray-500 dark:text-gray-400 mt-1 font-mono text-sm">
              {t("robotDetail.robotId", {
                value: robot.id || t("status.unknown"),
              })}
            </p>
          </div>
        </div>
        <TooltipProvider>
          <div className="flex items-center gap-2">
            <StartTeleoperationButton
              robot={robot}
              inScope={
                me?.role === USER_ROLE.ADMIN ||
                isRobotInScope(robot.id, scopeIds)
              }
            />
            {me?.role !== undefined && me.role <= USER_ROLE.MANAGER && (
              <QueueTaskButton robot={robot} />
            )}
          </div>
        </TooltipProvider>
      </div>

      {/* Tabs */}
      <Tabs defaultValue="overview" className="w-full">
        <TabsList>
          <TabsTrigger value="overview">
            {t("robotDetail.overview")}
          </TabsTrigger>
          <TabsTrigger value="health">{t("robotDetail.health")}</TabsTrigger>
          <TabsTrigger value="live-view">
            {t("robotDetail.liveView")}
          </TabsTrigger>
          <TabsTrigger value="sensors">{t("robotDetail.sensors")}</TabsTrigger>
          <TabsTrigger value="logs">{t("robotDetail.logs")}</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="mt-6">
          <OverviewTab
            robot={robot}
            realtimeStatus={realtimeStatus}
            isConnected={isConnected}
          />
        </TabsContent>

        <TabsContent value="health" className="mt-6">
          <HealthTab realtimeStatus={realtimeStatus} />
        </TabsContent>

        <TabsContent value="live-view" className="mt-6">
          <LiveViewTab robot={robot} realtimeStatus={realtimeStatus} />
        </TabsContent>

        <TabsContent value="sensors" className="mt-6">
          <SensorsTab />
        </TabsContent>

        <TabsContent value="logs" className="mt-6">
          <LogsTab />
        </TabsContent>
      </Tabs>
    </div>
  );
}

function RobotDetailSkeleton() {
  return (
    <div className="space-y-6 animate-pulse">
      <div className="flex items-start gap-4">
        <div className="h-8 w-8 bg-gray-200 dark:bg-gray-700 rounded"></div>
        <div>
          <div className="h-9 bg-gray-200 dark:bg-gray-700 rounded w-48 mb-2"></div>
          <div className="h-5 bg-gray-200 dark:bg-gray-700 rounded w-32"></div>
        </div>
      </div>
      <div className="h-10 bg-gray-200 dark:bg-gray-700 rounded w-96"></div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="h-32 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
        <div className="h-32 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
      </div>
    </div>
  );
}
