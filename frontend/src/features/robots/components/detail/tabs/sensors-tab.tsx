"use client";

import { Gauge } from "lucide-react";
import { useTranslation } from "react-i18next";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function SensorsTab() {
  const { t } = useTranslation();

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-base font-medium flex items-center gap-2">
          <Gauge className="h-5 w-5" />
          {t("robotDetail.sensors")}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex flex-col items-center justify-center py-16 text-center">
          <p className="text-gray-400 dark:text-gray-500">
            {t("common.comingSoon")}
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
