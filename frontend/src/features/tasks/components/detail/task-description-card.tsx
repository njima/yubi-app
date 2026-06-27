import { useTranslation } from "react-i18next";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface TaskDescriptionCardProps {
  description: string | null | undefined;
}

export function TaskDescriptionCard({ description }: TaskDescriptionCardProps) {
  const { t } = useTranslation();

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("taskInfo.description")}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-gray-700 dark:text-gray-300 whitespace-pre-wrap">
          {description || t("taskInfo.noDescription")}
        </p>
      </CardContent>
    </Card>
  );
}
