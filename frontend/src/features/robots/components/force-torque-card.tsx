"use client";

import { Gauge } from "lucide-react";
import { useTranslation } from "react-i18next";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function ForceTorqueCard() {
  const { t } = useTranslation();
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-base font-medium flex items-center gap-2">
          <Gauge className="h-4 w-4" />
          {t("forceTorqueCard.title")}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex flex-col items-center justify-center py-6 text-center">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t("common.comingSoon")}
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
