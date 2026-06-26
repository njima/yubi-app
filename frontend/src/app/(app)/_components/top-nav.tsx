"use client";

import {
  Bot,
  ClipboardList,
  KeyRound,
  LayoutGrid,
  Link2,
  MapPin,
  Users,
} from "lucide-react";
import { useTranslation } from "react-i18next";

import { usePermission } from "@/features/users";

import { LanguageSwitcher } from "./language-switcher";
import { NavItem } from "./nav-item";
import { UserMenu } from "./user-menu";

export function TopNav() {
  const { t } = useTranslation();
  const canViewApiKeys = usePermission("api_key:list");
  const navItems = [
    {
      href: "/dashboard",
      label: t("topNav.dashboard"),
      icon: <LayoutGrid className="h-4 w-4" />,
    },
    {
      href: "/robots",
      label: t("topNav.robots"),
      icon: <Bot className="h-4 w-4" />,
    },
    {
      href: "/tasks",
      label: t("topNav.tasks"),
      icon: <ClipboardList className="h-4 w-4" />,
    },
    {
      href: "/episodes",
      label: t("topNav.episodes"),
      icon: <Link2 className="h-4 w-4" />,
    },
    {
      href: "/users",
      label: t("topNav.users"),
      icon: <Users className="h-4 w-4" />,
    },
    {
      href: "/locations",
      label: t("topNav.locations"),
      icon: <MapPin className="h-4 w-4" />,
    },
    ...(canViewApiKeys
      ? [
          {
            href: "/api-keys",
            label: t("topNav.apiKeys"),
            icon: <KeyRound className="h-4 w-4" />,
          },
        ]
      : []),
  ];

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-white dark:bg-gray-950 dark:border-gray-800">
      <div className="flex h-16 items-center justify-between px-4">
        {/* Left: Logo + Navigation */}
        <div className="flex items-center gap-6">
          {/* Logo */}
          <div className="flex items-center gap-2">
            <Bot className="h-5 w-5 text-blue-600 dark:text-blue-400" />
            <span className="font-semibold text-gray-900 dark:text-gray-100 hidden sm:inline-block">
              {t("topNav.productName")}
            </span>
          </div>

          {/* Navigation Items */}
          <nav className="hidden md:flex items-center gap-1">
            {navItems.map((item) => (
              <NavItem
                key={item.href}
                href={item.href}
                label={item.label}
                icon={item.icon}
              />
            ))}
          </nav>
        </div>
        {/* Right: Language Switcher + User Menu */}
        <div className="flex items-center gap-4">
          <LanguageSwitcher />
          <UserMenu />
        </div>
      </div>
    </header>
  );
}
