import { ExternalLink } from "lucide-react";
import { useTranslation } from "react-i18next";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface TeachMeBizCardProps {
  manualUrl?: string;
}

export function TeachMeBizCard({ manualUrl }: TeachMeBizCardProps) {
  const { t } = useTranslation();
  const isSafeUrl = manualUrl?.startsWith("https://") ?? false;
  // Ensure the iframe uses the /embed URL for TeachMe Biz
  const embedUrl = manualUrl?.endsWith("/embed")
    ? manualUrl
    : `${manualUrl}/embed`;
  // Use /slideshow URL for opening in a new tab
  const viewUrl = manualUrl?.replace(/\/embed$/, "") + "/slideshow";

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-base font-medium">
          {t("teachMeBizCard.title")}
        </CardTitle>
      </CardHeader>
      <CardContent>
        {manualUrl && isSafeUrl ? (
          <div className="space-y-2">
            <div
              className="relative w-full"
              style={{ paddingBottom: "56.25%" }}
            >
              <iframe
                src={embedUrl}
                sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
                className="absolute inset-0 w-full h-full rounded border border-gray-200 dark:border-gray-700"
                title="TeachMeBiz Manual"
              />
            </div>
            <a
              href={viewUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1 text-xs text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
            >
              <ExternalLink className="h-3 w-3" />
              {t("teachMeBizCard.openInNewTab")}
            </a>
          </div>
        ) : (
          <div className="h-64 flex items-center justify-center text-sm text-gray-400">
            {t("teachMeBizCard.noManualUrl")}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
