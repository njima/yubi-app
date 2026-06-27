"use client";

import { use } from "react";

import { TeleopView } from "@/features/robots/components/teleoperation/teleop-view";

interface PageProps {
  params: Promise<{ robotId: string }>;
}

export default function TeleoperationPage({ params }: PageProps) {
  const { robotId } = use(params);
  return <TeleopView robotId={robotId} viewName="default" />;
}
