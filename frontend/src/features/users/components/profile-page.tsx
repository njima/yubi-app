"use client";

import { Mail, Shield, Calendar, User, Pencil } from "lucide-react";
import { useTranslation } from "react-i18next";

import { useUserRoleLabel } from "@/lib/hooks/use-status-labels";
import type { UserRoleValue } from "@/lib/status/constants";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { EditDisplayNameDialog } from "./edit-display-name-dialog";
import { useMeQuery } from "../hooks/use-me-query";

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString("ja-JP", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}

export function ProfilePage() {
  const { t } = useTranslation();
  const getRoleLabel = useUserRoleLabel();
  const { data: user, isLoading, error } = useMeQuery();

  if (isLoading) {
    return (
      <div className="max-w-2xl mx-auto">
        <Card>
          <CardContent className="py-8">
            <p className="text-center text-gray-500">{t("common.loading")}</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-2xl mx-auto">
        <Card>
          <CardContent className="py-8">
            <p className="text-center text-red-600">
              {t("profilePage.failedToLoadProfile", { message: error.message })}
            </p>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  const initials = user.display_name
    .split(" ")
    .map((n: string) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);

  return (
    <div className="max-w-2xl mx-auto">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <User className="h-5 w-5" />
            {t("profilePage.title")}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Avatar and Name */}
          <div className="flex items-center gap-4">
            <Avatar className="h-20 w-20">
              <AvatarFallback className="bg-blue-600 text-white text-xl">
                {initials}
              </AvatarFallback>
            </Avatar>
            <div>
              <div className="flex items-center gap-2">
                <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                  {user.display_name}
                </h2>
                <EditDisplayNameDialog currentDisplayName={user.display_name}>
                  <Button
                    size="sm"
                    variant="ghost"
                    aria-label={t("profilePage.editDisplayName")}
                  >
                    <Pencil className="h-4 w-4" />
                  </Button>
                </EditDisplayNameDialog>
              </div>
              <p className="text-gray-500 dark:text-gray-400">
                {t("episodeDetail.id")}: {user.user_id}
              </p>
            </div>
          </div>

          {/* User Details */}
          <div className="grid gap-4">
            <div className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
              <Mail className="h-5 w-5 text-gray-500" />
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  {t("profilePage.email")}
                </p>
                <p className="font-medium text-gray-900 dark:text-gray-100">
                  {user.email}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
              <Shield className="h-5 w-5 text-gray-500" />
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  {t("profilePage.role")}
                </p>
                <p className="font-medium text-gray-900 dark:text-gray-100">
                  {getRoleLabel(user.role as UserRoleValue)}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
              <Calendar className="h-5 w-5 text-gray-500" />
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  {t("profilePage.createdAt")}
                </p>
                <p className="font-medium text-gray-900 dark:text-gray-100">
                  {formatDate(user.created_at)}
                </p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
