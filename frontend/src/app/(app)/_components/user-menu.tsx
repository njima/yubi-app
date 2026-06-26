"use client";

import { User } from "lucide-react";
import Link from "next/link";
import { useTranslation } from "react-i18next";

import { Avatar, AvatarFallback } from "@/shared/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/shared/ui/dropdown-menu";

import { useMeQuery } from "@/features/users";

import { SwitchUserDialog } from "./switch-user-dialog";

export function UserMenu() {
  const { t } = useTranslation();
  const { data: user, isLoading, error } = useMeQuery();

  if (error) {
    return null;
  }

  if (isLoading) {
    return (
      <div className="flex items-center gap-2 px-3 py-2">
        <div className="h-8 w-8 rounded-full bg-gray-200 animate-pulse" />
        <div className="hidden md:block h-4 w-24 bg-gray-200 rounded animate-pulse" />
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
    <DropdownMenu>
      <DropdownMenuTrigger className="flex items-center gap-2 rounded-md px-3 py-2 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors focus:outline-none">
        <Avatar className="h-8 w-8">
          <AvatarFallback className="bg-blue-600 text-white text-xs">
            {initials}
          </AvatarFallback>
        </Avatar>
        <span className="hidden md:inline font-medium text-sm text-gray-900 dark:text-gray-100">
          {user.display_name}
        </span>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-56">
        <DropdownMenuLabel>{t("userMenu.myAccount")}</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem asChild>
          <Link href="/profile">
            <User className="mr-2 h-4 w-4" />
            {t("userMenu.profile")}
          </Link>
        </DropdownMenuItem>
        <SwitchUserDialog />
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
